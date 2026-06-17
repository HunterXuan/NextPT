package iam

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
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
		return nil, gerror.New(gi18n.T(ctx, "iam.session.invalid_credentials"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))
	if err != nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.session.invalid_credentials"))
	}

	if err := service.IamUserUsecase().EnsureCanAuthenticate(ctx, user); err != nil {
		return nil, err
	}

	role, err := service.IamRoleDomain().GetRoleById(ctx, user.Role)
	if err != nil || role == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.session.role_missing"))
	}

	token, err := service.IamSessionDomain().GenerateToken(ctx, gconv.String(user.Id), g.Map{
		"roleId":    user.Role,
		"roleLevel": role.Level,
		"isStaff":   role.IsStaff,
	})
	if err != nil {
		return nil, err
	}

	return &iamout.SessionCreateOut{
		Token: token,
	}, nil
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
