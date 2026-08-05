package iam

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/library/twofactor"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/iamin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/iamout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

const iamTwoStepRateLimitScript = `
local current = redis.call("INCR", KEYS[1])
if current == 1 then
    redis.call("EXPIRE", KEYS[1], tonumber(ARGV[1]))
end
return current
`

type sIamTwoStepUsecase struct{}

type iamTwoStepSetupCache struct {
	ChallengeHash string `json:"challengeHash"`
	Secret        string `json:"secret"`
}

func init() {
	service.RegisterIamTwoStepUsecase(NewIamTwoStepUsecase())
}

func NewIamTwoStepUsecase() *sIamTwoStepUsecase {
	return &sIamTwoStepUsecase{}
}

func (s *sIamTwoStepUsecase) Setup(ctx context.Context, actor *model.Actor, in iamin.UserTwoStepSetupInp) (*iamout.UserTwoStepSetupOut, error) {
	user, err := s.currentUser(ctx, actor)
	if err != nil {
		return nil, err
	}
	if user.TwoStepType != consts.IamTwoStepTypeDisabled {
		return nil, gerror.New(gi18n.T(ctx, "iam.two_step.already_enabled"))
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.user.old_password_invalid"))
	}

	siteName := strings.TrimSpace(g.Cfg().MustGet(ctx, "site.name", "NextPT").String())
	if siteName == "" {
		siteName = "NextPT"
	}
	setup, err := twofactor.NewTOTPSetup(siteName, user.Email)
	if err != nil {
		return nil, err
	}
	challenge, err := s.newChallenge()
	if err != nil {
		return nil, err
	}
	if err := s.storeSetup(ctx, user.Id, challenge, setup.Secret); err != nil {
		return nil, err
	}
	return &iamout.UserTwoStepSetupOut{
		Challenge:     challenge,
		QRCodeDataURL: setup.QRCodeDataURL,
		Secret:        setup.Secret,
	}, nil
}

func (s *sIamTwoStepUsecase) Confirm(ctx context.Context, actor *model.Actor, in iamin.UserTwoStepConfirmInp) (*iamout.UserTwoStepRecoveryCodesOut, error) {
	user, err := s.currentUser(ctx, actor)
	if err != nil {
		return nil, err
	}
	if user.TwoStepType != consts.IamTwoStepTypeDisabled {
		return nil, gerror.New(gi18n.T(ctx, "iam.two_step.already_enabled"))
	}

	setup, err := s.loadSetup(ctx, user.Id)
	if err != nil || subtle.ConstantTimeCompare([]byte(setup.ChallengeHash), []byte(s.digest(in.Challenge))) != 1 {
		return nil, gerror.New(gi18n.T(ctx, "iam.two_step.challenge_invalid"))
	}
	if !twofactor.VerifyTOTP(setup.Secret, in.Code, gtime.Now().Time) {
		return nil, gerror.New(gi18n.T(ctx, "iam.two_step.code_invalid"))
	}

	encryptionKey, err := s.encryptionKey(ctx)
	if err != nil {
		return nil, err
	}
	encryptedSecret, err := twofactor.EncryptSecret(setup.Secret, encryptionKey)
	if err != nil {
		return nil, err
	}
	codes, err := twofactor.NewRecoveryCodes(consts.IamTwoStepRecoveryCodeCount)
	if err != nil {
		return nil, err
	}
	if err := s.enable(ctx, user.Id, encryptedSecret, codes); err != nil {
		return nil, err
	}
	_, _ = g.Redis().Do(ctx, "DEL", service.SysCache().KeyIamTwoStepSetup(ctx, user.Id))
	s.recordAudit(ctx, actor, consts.SiteAuditOperationEnableTwoStep)
	service.IamUserUsecase().InvalidateUserCache(ctx, user.Id)
	return &iamout.UserTwoStepRecoveryCodesOut{RecoveryCodes: codes}, nil
}

func (s *sIamTwoStepUsecase) CreateRecoveryCodes(ctx context.Context, actor *model.Actor, in iamin.UserTwoStepRecoveryCodesCreateInp) (*iamout.UserTwoStepRecoveryCodesOut, error) {
	user, err := s.currentUser(ctx, actor)
	if err != nil {
		return nil, err
	}
	if _, err := s.validateUserCode(ctx, user, in.Code); err != nil {
		return nil, err
	}
	codes, err := twofactor.NewRecoveryCodes(consts.IamTwoStepRecoveryCodeCount)
	if err != nil {
		return nil, err
	}
	if err := dao.IamUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return s.replaceRecoveryCodes(ctx, user.Id, codes)
	}); err != nil {
		return nil, err
	}
	s.recordAudit(ctx, actor, consts.SiteAuditOperationRenewTwoStepCodes)
	return &iamout.UserTwoStepRecoveryCodesOut{RecoveryCodes: codes}, nil
}

func (s *sIamTwoStepUsecase) Disable(ctx context.Context, actor *model.Actor, in iamin.UserTwoStepDisableInp) error {
	user, err := s.currentUser(ctx, actor)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return gerror.New(gi18n.T(ctx, "iam.user.old_password_invalid"))
	}
	if _, err := s.validateUserCode(ctx, user, in.Code); err != nil {
		return err
	}

	err = dao.IamUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.IamUserDomain().UpdateTwoStep(ctx, user.Id, consts.IamTwoStepTypeDisabled, ""); err != nil {
			return err
		}
		return service.IamTwoStepDomain().DeleteRecoveryCodes(ctx, user.Id)
	})
	if err != nil {
		return err
	}
	s.recordAudit(ctx, actor, consts.SiteAuditOperationDisableTwoStep)
	service.IamUserUsecase().InvalidateUserCache(ctx, user.Id)
	return nil
}

func (s *sIamTwoStepUsecase) CreateLoginChallenge(ctx context.Context, userId uint64) (string, error) {
	challenge, err := s.newChallenge()
	if err != nil {
		return "", err
	}
	_, err = g.Redis().Do(
		ctx,
		"SET",
		service.SysCache().KeyIamTwoStepLoginChallenge(ctx, s.digest(challenge)),
		gconv.String(userId),
		"EX",
		int(consts.IamTwoStepLoginChallengeTTL.Seconds()),
	)
	return challenge, err
}

func (s *sIamTwoStepUsecase) VerifyLogin(ctx context.Context, in iamin.SessionTwoStepVerifyInp) (uint64, error) {
	userId, err := s.loginChallengeUserId(ctx, in.Challenge)
	if err != nil || userId == 0 {
		return 0, gerror.New(gi18n.T(ctx, "iam.two_step.challenge_invalid"))
	}
	if err := s.loginAttemptAllowed(ctx, userId); err != nil {
		return userId, err
	}
	user, err := service.IamUserDomain().GetUserById(ctx, userId)
	if err != nil {
		return userId, err
	}
	if user == nil {
		return userId, gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}
	recoveryCodeHash, err := s.validateUserCode(ctx, user, in.Code)
	if err != nil {
		return userId, err
	}
	consumed, err := s.consumeLoginChallenge(ctx, in.Challenge)
	if err != nil {
		return userId, err
	}
	if !consumed {
		return userId, gerror.New(gi18n.T(ctx, "iam.two_step.challenge_invalid"))
	}
	if recoveryCodeHash == "" {
		return userId, nil
	}
	consumed, err = service.IamTwoStepDomain().ConsumeRecoveryCode(ctx, user.Id, recoveryCodeHash)
	if err != nil {
		return userId, err
	}
	if !consumed {
		return userId, gerror.New(gi18n.T(ctx, "iam.two_step.code_invalid"))
	}
	return userId, nil
}

func (s *sIamTwoStepUsecase) enable(ctx context.Context, userId uint64, encryptedSecret string, codes []string) error {
	return dao.IamUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.IamUserDomain().UpdateTwoStep(ctx, userId, consts.IamTwoStepTypeTOTP, encryptedSecret); err != nil {
			return err
		}
		return s.replaceRecoveryCodes(ctx, userId, codes)
	})
}

func (s *sIamTwoStepUsecase) replaceRecoveryCodes(ctx context.Context, userId uint64, codes []string) error {
	hashes := make([]string, 0, len(codes))
	for _, code := range codes {
		hashes = append(hashes, twofactor.RecoveryCodeHash(code))
	}
	return service.IamTwoStepDomain().ReplaceRecoveryCodeHashes(ctx, userId, hashes)
}

func (s *sIamTwoStepUsecase) validateUserCode(ctx context.Context, user *entity.IamUser, code string) (string, error) {
	if user == nil || user.TwoStepType != consts.IamTwoStepTypeTOTP || user.TwoStepSecret == "" {
		return "", gerror.New(gi18n.T(ctx, "iam.two_step.not_enabled"))
	}
	encryptionKey, err := s.encryptionKey(ctx)
	if err != nil {
		return "", err
	}
	secret, err := twofactor.DecryptSecret(user.TwoStepSecret, encryptionKey)
	if err != nil {
		return "", gerror.New(gi18n.T(ctx, "iam.two_step.secret_invalid"))
	}
	if twofactor.VerifyTOTP(secret, code, gtime.Now().Time) {
		return "", nil
	}
	recoveryCodeHash := twofactor.RecoveryCodeHash(code)
	exists, err := service.IamTwoStepDomain().HasUnusedRecoveryCode(ctx, user.Id, recoveryCodeHash)
	if err != nil {
		return "", err
	}
	if exists {
		return recoveryCodeHash, nil
	}
	return "", gerror.New(gi18n.T(ctx, "iam.two_step.code_invalid"))
}

func (s *sIamTwoStepUsecase) currentUser(ctx context.Context, actor *model.Actor) (*entity.IamUser, error) {
	if actor == nil || actor.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	user, err := service.IamUserDomain().GetUserById(ctx, actor.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}
	return user, nil
}

func (s *sIamTwoStepUsecase) storeSetup(ctx context.Context, userId uint64, challenge string, secret string) error {
	payload, err := json.Marshal(iamTwoStepSetupCache{ChallengeHash: s.digest(challenge), Secret: secret})
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SET", service.SysCache().KeyIamTwoStepSetup(ctx, userId), string(payload), "EX", int(consts.IamTwoStepSetupTTL.Seconds()))
	return err
}

func (s *sIamTwoStepUsecase) loadSetup(ctx context.Context, userId uint64) (*iamTwoStepSetupCache, error) {
	value, err := g.Redis().Do(ctx, "GET", service.SysCache().KeyIamTwoStepSetup(ctx, userId))
	if err != nil || value == nil || value.IsNil() {
		return nil, gerror.New(gi18n.T(ctx, "iam.two_step.challenge_invalid"))
	}
	var setup iamTwoStepSetupCache
	if err := json.Unmarshal([]byte(value.String()), &setup); err != nil || setup.ChallengeHash == "" || setup.Secret == "" {
		return nil, gerror.New(gi18n.T(ctx, "iam.two_step.challenge_invalid"))
	}
	return &setup, nil
}

func (s *sIamTwoStepUsecase) loginChallengeUserId(ctx context.Context, challenge string) (uint64, error) {
	value, err := g.Redis().Do(ctx, "GET", service.SysCache().KeyIamTwoStepLoginChallenge(ctx, s.digest(challenge)))
	if err != nil || value == nil || value.IsNil() {
		return 0, err
	}
	return value.Uint64(), nil
}

func (s *sIamTwoStepUsecase) consumeLoginChallenge(ctx context.Context, challenge string) (bool, error) {
	value, err := g.Redis().Do(ctx, "DEL", service.SysCache().KeyIamTwoStepLoginChallenge(ctx, s.digest(challenge)))
	if err != nil {
		return false, err
	}
	return value.Int() == 1, nil
}

func (s *sIamTwoStepUsecase) loginAttemptAllowed(ctx context.Context, userId uint64) error {
	value, err := g.Redis().Do(
		ctx,
		"EVAL",
		iamTwoStepRateLimitScript,
		1,
		service.SysCache().KeyIamTwoStepLoginRate(ctx, userId),
		int(consts.IamTwoStepLoginRateWindow.Seconds()),
	)
	if err != nil {
		return err
	}
	if value.Int() > consts.IamTwoStepLoginRateLimit {
		return gerror.New(gi18n.T(ctx, "iam.two_step.rate_limited"))
	}
	return nil
}

func (s *sIamTwoStepUsecase) encryptionKey(ctx context.Context) (string, error) {
	key := strings.TrimSpace(g.Cfg().MustGet(ctx, "iam.twoStepEncryptionKey", "").String())
	if key == "" {
		return "", gerror.New(gi18n.T(ctx, "iam.two_step.encryption_key_missing"))
	}
	return key, nil
}

func (s *sIamTwoStepUsecase) digest(value string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(value)))
	return hex.EncodeToString(digest[:])
}

func (s *sIamTwoStepUsecase) newChallenge() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func (s *sIamTwoStepUsecase) recordAudit(ctx context.Context, actor *model.Actor, operation string) {
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeIamUser,
		TargetId:   actor.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail:     map[string]any{"operation": operation},
	})
}
