package mod

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

type sModUserDomain struct{}

func init() {
	service.RegisterModUserDomain(NewModUserDomain())
}

func NewModUserDomain() *sModUserDomain {
	return &sModUserDomain{}
}

func (s *sModUserDomain) Create(ctx context.Context, mod entity.ModUserLog) (uint64, error) {
	result, err := dao.ModUserLog.Ctx(ctx).Data(mod).Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

func (s *sModUserDomain) GetById(ctx context.Context, id uint64) (*entity.ModUserLog, error) {
	var mod entity.ModUserLog
	err := dao.ModUserLog.Ctx(ctx).WherePri(id).Scan(&mod)
	if err != nil {
		return nil, err
	}
	return &mod, nil
}

func (s *sModUserDomain) Update(ctx context.Context, id uint64, data interface{}) error {
	_, err := dao.ModUserLog.Ctx(ctx).WherePri(id).Data(data).Update()
	return err
}

func (s *sModUserDomain) QueryExpiredActiveMods(ctx context.Context, now *gtime.Time, limit int) ([]entity.ModUserLog, error) {
	if limit <= 0 {
		limit = 500
	}

	columns := dao.ModUserLog.Columns()
	var records []entity.ModUserLog
	err := dao.ModUserLog.Ctx(ctx).
		Where(columns.IsActive, true).
		WhereNotNull(columns.ExpireAt).
		WhereLTE(columns.ExpireAt, now).
		OrderAsc(columns.Id).
		Limit(limit).
		Scan(&records)
	return records, err
}

func (s *sModUserDomain) HasActiveMod(ctx context.Context, userId uint64, modTypes []int) (bool, error) {
	if userId == 0 || len(modTypes) == 0 {
		return false, nil
	}

	columns := dao.ModUserLog.Columns()
	m := dao.ModUserLog.Ctx(ctx)
	count, err := m.
		Where(columns.UserId, userId).
		Where(columns.IsActive, true).
		WhereIn(columns.ModType, modTypes).
		Where(m.Builder().WhereNull(columns.ExpireAt).WhereOrGT(columns.ExpireAt, gtime.Now())).
		Count()
	return count > 0, err
}

func (s *sModUserDomain) QueryUserLogs(ctx context.Context, userId uint64, page, size int) ([]entity.ModUserLog, int, error) {
	m := dao.ModUserLog.Ctx(ctx).Where(dao.ModUserLog.Columns().UserId, userId)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var records []entity.ModUserLog
	err = m.Page(page, size).OrderDesc("id").Scan(&records)
	return records, total, err
}
