package iam

import (
	"context"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

type sIamInviteDomain struct{}

func init() {
	service.RegisterIamInviteDomain(NewIamInviteDomain())
}

func NewIamInviteDomain() *sIamInviteDomain {
	return &sIamInviteDomain{}
}

func (s *sIamInviteDomain) GetInviteByHashForUpdate(ctx context.Context, hash string) (*entity.IamInvite, error) {
	var invite *entity.IamInvite
	err := dao.IamInvite.Ctx(ctx).LockUpdate().Where(dao.IamInvite.Columns().Hash, hash).Scan(&invite)
	return invite, err
}

func (s *sIamInviteDomain) GetInviteByHash(ctx context.Context, hash string) (*entity.IamInvite, error) {
	var invite *entity.IamInvite
	err := dao.IamInvite.Ctx(ctx).Where(dao.IamInvite.Columns().Hash, hash).Scan(&invite)
	return invite, err
}

func (s *sIamInviteDomain) GetInvitesByInviterIdAndHash(ctx context.Context, inviterId uint64, hash string) (*entity.IamInvite, error) {
	var invite *entity.IamInvite
	err := dao.IamInvite.Ctx(ctx).Where(dao.IamInvite.Columns().InviterId, inviterId).Where(dao.IamInvite.Columns().Hash, hash).Scan(&invite)
	return invite, err
}

func (s *sIamInviteDomain) QueryInvitesByInviter(ctx context.Context, inviterId uint64, page, size int, status *uint) ([]entity.IamInvite, int, error) {
	columns := dao.IamInvite.Columns()
	m := dao.IamInvite.Ctx(ctx).Where(columns.InviterId, inviterId)
	if status != nil {
		m = m.Where(columns.Status, *status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.IamInvite
	err = m.Page(page, size).OrderDesc(columns.Id).Scan(&list)
	return list, total, err
}

func (s *sIamInviteDomain) ExpireInvites(ctx context.Context, now *gtime.Time) (int64, error) {
	columns := dao.IamInvite.Columns()
	result, err := dao.IamInvite.Ctx(ctx).
		WhereIn(columns.Status, []int{consts.IamInviteStatusUnused, consts.IamInviteStatusSent}).
		WhereNotNull(columns.ExpireAt).
		WhereLTE(columns.ExpireAt, now).
		Data(columns.Status, consts.IamInviteStatusExpired).
		Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *sIamInviteDomain) UpdateInviteStatus(ctx context.Context, id uint64, status uint) error {
	_, err := dao.IamInvite.Ctx(ctx).Where(dao.IamInvite.Columns().Id, id).Data(dao.IamInvite.Columns().Status, status).Update()
	return err
}

func (s *sIamInviteDomain) UpdateInvite(ctx context.Context, id uint64, data do.IamInvite) error {
	_, err := dao.IamInvite.Ctx(ctx).Where(dao.IamInvite.Columns().Id, id).Data(data).Update()
	return err
}

func (s *sIamInviteDomain) AdminCreateInvites(ctx context.Context, invites []do.IamInvite) error {
	for _, inv := range invites {
		_, err := dao.IamInvite.Ctx(ctx).Data(inv).Insert()
		if err != nil {
			return err
		}
	}
	return nil
}
