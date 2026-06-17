package iam

import (
	"context"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/service"
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

func (s *sIamInviteDomain) QueryInvitesByInviter(ctx context.Context, inviterId uint64, page, size int) ([]entity.IamInvite, int, error) {
	m := dao.IamInvite.Ctx(ctx).Where(dao.IamInvite.Columns().InviterId, inviterId)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.IamInvite
	err = m.Page(page, size).OrderDesc(dao.IamInvite.Columns().Id).Scan(&list)
	return list, total, err
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
