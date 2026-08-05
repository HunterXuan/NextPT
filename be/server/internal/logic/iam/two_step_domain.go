package iam

import (
	"context"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sIamTwoStepDomain struct{}

func init() {
	service.RegisterIamTwoStepDomain(NewIamTwoStepDomain())
}

func NewIamTwoStepDomain() *sIamTwoStepDomain {
	return &sIamTwoStepDomain{}
}

func (s *sIamTwoStepDomain) ReplaceRecoveryCodeHashes(ctx context.Context, userId uint64, hashes []string) error {
	if err := s.DeleteRecoveryCodes(ctx, userId); err != nil {
		return err
	}
	if len(hashes) == 0 {
		return nil
	}

	items := make([]do.IamUserRecoveryCode, 0, len(hashes))
	for _, hash := range hashes {
		items = append(items, do.IamUserRecoveryCode{
			UserId:    userId,
			CodeHash:  hash,
			CreatedAt: gtime.Now(),
		})
	}
	_, err := dao.IamUserRecoveryCode.Ctx(ctx).Data(items).Insert()
	return err
}

func (s *sIamTwoStepDomain) HasUnusedRecoveryCode(ctx context.Context, userId uint64, codeHash string) (bool, error) {
	columns := dao.IamUserRecoveryCode.Columns()
	count, err := dao.IamUserRecoveryCode.Ctx(ctx).
		Where(columns.UserId, userId).
		Where(columns.CodeHash, codeHash).
		WhereNull(columns.UsedAt).
		Count()
	return count > 0, err
}

func (s *sIamTwoStepDomain) ConsumeRecoveryCode(ctx context.Context, userId uint64, codeHash string) (bool, error) {
	columns := dao.IamUserRecoveryCode.Columns()
	result, err := dao.IamUserRecoveryCode.Ctx(ctx).
		Where(columns.UserId, userId).
		Where(columns.CodeHash, codeHash).
		WhereNull(columns.UsedAt).
		Data(g.Map{columns.UsedAt: gtime.Now()}).
		Update()
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected == 1, err
}

func (s *sIamTwoStepDomain) DeleteRecoveryCodes(ctx context.Context, userId uint64) error {
	_, err := dao.IamUserRecoveryCode.Ctx(ctx).
		Where(dao.IamUserRecoveryCode.Columns().UserId, userId).
		Delete()
	return err
}
