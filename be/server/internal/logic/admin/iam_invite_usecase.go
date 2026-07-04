package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/in/adminin"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sAdminIamInviteUsecase struct{}

const (
	inviteHashAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	inviteHashLength   = 32
)

func NewAdminIamInviteUsecase() *sAdminIamInviteUsecase {
	return &sAdminIamInviteUsecase{}
}

func init() {
	service.RegisterAdminIamInviteUsecase(NewAdminIamInviteUsecase())
}

func (s *sAdminIamInviteUsecase) Grant(ctx context.Context, actor *model.Actor, in adminin.IamInviteGrantInp) error {
	var list []do.IamInvite
	isTemporary := in.IsTemp || in.ExpireAt != nil
	for i := 0; i < in.Amount; i++ {
		list = append(list, do.IamInvite{
			InviterId:   0,
			Hash:        s.generateInviteHash(),
			Status:      0,
			IsTemporary: isTemporary,
			ExpireAt:    in.ExpireAt,
		})
	}

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return service.IamInviteDomain().AdminCreateInvites(ctx, list)
	})
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.invite.grant_failed"))
	}

	return nil
}

func (s *sAdminIamInviteUsecase) generateInviteHash() string {
	return grand.Str(inviteHashAlphabet, inviteHashLength)
}
