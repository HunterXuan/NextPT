package mod

import (
	"context"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

type sModCheaterDomain struct{}

func init() {
	service.RegisterModCheaterDomain(NewModCheaterDomain())
}

func NewModCheaterDomain() *sModCheaterDomain {
	return &sModCheaterDomain{}
}

func (s *sModCheaterDomain) Create(ctx context.Context, log entity.ModCheaterLog) error {
	_, err := dao.ModCheaterLog.Ctx(ctx).Data(log).Insert()
	return err
}

func (s *sModCheaterDomain) Resolve(ctx context.Context, id uint64, dealtBy uint64, dealtComment string, dealtAt *gtime.Time) error {
	_, err := dao.ModCheaterLog.Ctx(ctx).WherePri(id).Data(do.ModCheaterLog{
		IsDealt:      true,
		DealtBy:      dealtBy,
		DealtComment: dealtComment,
		DealtAt:      dealtAt,
	}).Update()
	return err
}

func (s *sModCheaterDomain) QueryCheaterLogs(ctx context.Context, isDealt *int, page, size int) ([]entity.ModCheaterLog, int, error) {
	m := dao.ModCheaterLog.Ctx(ctx)
	if isDealt != nil && *isDealt >= 0 {
		m = m.Where(dao.ModCheaterLog.Columns().IsDealt, *isDealt)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var records []entity.ModCheaterLog
	err = m.Page(page, size).OrderDesc(dao.ModCheaterLog.Columns().CreatedAt).Scan(&records)
	return records, total, err
}

func (s *sModCheaterDomain) DeleteCheaterLogsByTorrentId(ctx context.Context, torrentId uint64) error {
	_, err := dao.ModCheaterLog.Ctx(ctx).Where(dao.ModCheaterLog.Columns().TorrentId, torrentId).Delete()
	return err
}
