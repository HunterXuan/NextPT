package iam

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

type sIamSessionUsecase struct{}

func init() {
	service.RegisterIamSessionUsecase(NewIamSessionUsecase())
}

func NewIamSessionUsecase() *sIamSessionUsecase {
	return &sIamSessionUsecase{}
}

func (s *sIamSessionUsecase) Create(ctx context.Context, in iamin.SessionCreateInp) (*iamout.SessionCreateOut, error) {
	user, err := service.IamUserDomain().GetUserByLogin(ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		s.recordLogin(ctx, 0, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonUserNotFound)
		return nil, gerror.New(gi18n.T(ctx, "iam.session.invalid_credentials"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))
	if err != nil {
		s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonInvalidPassword)
		return nil, gerror.New(gi18n.T(ctx, "iam.session.invalid_credentials"))
	}

	if err := service.IamUserUsecase().EnsureCanAuthenticate(ctx, user); err != nil {
		s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonAccountUnavailable)
		return nil, err
	}

	role, err := service.IamRoleDomain().GetRoleById(ctx, user.Role)
	if err != nil || role == nil {
		s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonRoleMissing)
		return nil, gerror.New(gi18n.T(ctx, "iam.session.role_missing"))
	}
	if user.TwoStepType != consts.IamTwoStepTypeDisabled {
		if user.TwoStepType != consts.IamTwoStepTypeTOTP {
			s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonTwoStepInvalid)
			return nil, gerror.New(gi18n.T(ctx, "iam.two_step.unsupported_type"))
		}
		challenge, err := service.IamTwoStepUsecase().CreateLoginChallenge(ctx, user.Id)
		if err != nil {
			s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonTokenCreateFailed)
			return nil, err
		}
		return &iamout.SessionCreateOut{TwoStepRequired: true, TwoStepChallenge: challenge}, nil
	}

	token, err := s.generateToken(ctx, user, role.Level, role.IsStaff)
	if err != nil {
		s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonTokenCreateFailed)
		return nil, err
	}

	s.recordSuccessfulLogin(ctx, user)
	return &iamout.SessionCreateOut{
		Token: token,
	}, nil
}

func (s *sIamSessionUsecase) VerifyTwoStep(ctx context.Context, in iamin.SessionTwoStepVerifyInp) (*iamout.SessionCreateOut, error) {
	userId, err := service.IamTwoStepUsecase().VerifyLogin(ctx, in)
	if err != nil {
		if userId != 0 {
			s.recordLogin(ctx, userId, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonTwoStepInvalid)
		}
		return nil, err
	}
	user, err := service.IamUserDomain().GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}
	if err := service.IamUserUsecase().EnsureCanAuthenticate(ctx, user); err != nil {
		s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonAccountUnavailable)
		return nil, err
	}
	role, err := service.IamRoleDomain().GetRoleById(ctx, user.Role)
	if err != nil || role == nil {
		s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonRoleMissing)
		return nil, gerror.New(gi18n.T(ctx, "iam.session.role_missing"))
	}
	token, err := s.generateToken(ctx, user, role.Level, role.IsStaff)
	if err != nil {
		s.recordLogin(ctx, user.Id, consts.IamLoginLogResultFail, consts.IamLoginLogFailReasonTokenCreateFailed)
		return nil, err
	}
	s.recordSuccessfulLogin(ctx, user)
	return &iamout.SessionCreateOut{Token: token}, nil
}

func (s *sIamSessionUsecase) generateToken(ctx context.Context, user *entity.IamUser, roleLevel int, isStaff bool) (string, error) {
	return service.IamSessionDomain().GenerateToken(ctx, gconv.String(user.Id), g.Map{
		"roleId":    user.Role,
		"roleLevel": roleLevel,
		"isStaff":   isStaff,
	})
}

func (s *sIamSessionUsecase) recordSuccessfulLogin(ctx context.Context, user *entity.IamUser) {
	if user == nil {
		return
	}
	now := gtime.Now()
	ip := s.requestIp(ctx)
	s.recordLoginWithMeta(ctx, user.Id, consts.IamLoginLogResultSuccess, "", ip, s.requestUserAgent(ctx), now)
	if err := service.IamUserDomain().UpdateLoginTrace(ctx, user.Id, now, ip); err != nil {
		g.Log().Warningf(ctx, "iam session: update login trace failed: %v", err)
	}
}

func (s *sIamSessionUsecase) recordLogin(ctx context.Context, userId uint64, result int, failReason string) {
	s.recordLoginWithMeta(ctx, userId, result, failReason, s.requestIp(ctx), s.requestUserAgent(ctx), gtime.Now())
}

func (s *sIamSessionUsecase) recordLoginWithMeta(ctx context.Context, userId uint64, result int, failReason string, ip string, userAgent string, createdAt *gtime.Time) {
	if err := service.IamLoginLogDomain().Create(ctx, do.IamLoginLog{
		UserId:     userId,
		Ip:         ip,
		UserAgent:  userAgent,
		Result:     result,
		FailReason: failReason,
		CreatedAt:  createdAt,
	}); err != nil {
		g.Log().Warningf(ctx, "iam session: record login log failed: %v", err)
	}
}

func (s *sIamSessionUsecase) requestIp(ctx context.Context) string {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	return r.GetClientIp()
}

func (s *sIamSessionUsecase) requestUserAgent(ctx context.Context) string {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	return r.Header.Get("User-Agent")
}

func (s *sIamSessionUsecase) Delete(ctx context.Context, actor *model.Actor) error {
	if actor != nil && actor.Id > 0 {
		return service.IamSessionDomain().RemoveToken(ctx, gconv.String(actor.Id))
	}
	return nil
}

func (s *sIamSessionUsecase) VerifyPasskey(ctx context.Context, passkey string) (*model.Actor, error) {
	return service.IamUserUsecase().LoadActorByPasskey(ctx, passkey)
}
