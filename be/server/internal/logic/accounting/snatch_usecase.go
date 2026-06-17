package accounting

import (
	"context"
	"server/internal/model/entity"

	"server/internal/model"
	"server/internal/model/in/accountingin"
	"server/internal/model/out/accountingout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sAccountingSnatchUsecase struct{}

func init() {
	service.RegisterAccountingSnatchUsecase(NewAccountingSnatchUsecase())
}

func NewAccountingSnatchUsecase() *sAccountingSnatchUsecase {
	return &sAccountingSnatchUsecase{}
}

func (s *sAccountingSnatchUsecase) ListMySnatches(ctx context.Context, actor *model.Actor, in accountingin.SnatchListInp) (*accountingout.SnatchListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	page := in.Page
	size := in.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	snatches, total, err := service.AccountingSnatchDomain().ListSnatches(ctx, actor.Id, page, size, in.IsFinished, in.IsActive)
	if err != nil {
		return nil, err
	}

	var torrentIds []uint64
	for _, sn := range snatches {
		torrentIds = append(torrentIds, sn.TorrentId)
	}

	torrents, _ := service.CatalogTorrentDomain().GetTorrentsByIds(ctx, torrentIds)
	torrentMap := make(map[uint64]*entity.CatalogTorrent)
	for _, t := range torrents {
		torrentMap[t.Id] = t
	}

	var list []accountingout.SnatchItem
	for _, sn := range snatches {
		tName := ""
		tSize := uint64(0)
		if t, ok := torrentMap[sn.TorrentId]; ok {
			tName = t.Name
			tSize = t.Size
		}

		list = append(list, accountingout.SnatchItem{
			Id:           sn.Id,
			TorrentId:    sn.TorrentId,
			TorrentName:  tName,
			TorrentSize:  tSize,
			Uploaded:     sn.Uploaded,
			Downloaded:   sn.Downloaded,
			SeedTime:     uint64(sn.SeedTime),
			LeechTime:    uint64(sn.LeechTime),
			IsFinished:   sn.IsFinished,
			StartedAt:    sn.StartedAt,
			FinishedAt:   sn.CompletedAt,
			LastActionAt: sn.LastAction,
		})
	}

	return &accountingout.SnatchListOut{
		Page:  page,
		Size:  size,
		Total: total,
		List:  list,
	}, nil
}

func (s *sAccountingSnatchUsecase) GetMySnatch(ctx context.Context, actor *model.Actor, in accountingin.SnatchGetInp) (*accountingout.SnatchGetOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	sn, err := service.AccountingSnatchDomain().GetSnatch(ctx, actor.Id, in.TorrentId)
	if err != nil {
		return nil, err
	}
	if sn == nil {
		return nil, gerror.New(gi18n.T(ctx, "accounting.snatch.not_found"))
	}

	return &accountingout.SnatchGetOut{
		SnatchItem: accountingout.SnatchItem{
			Id:           sn.Id,
			TorrentId:    sn.TorrentId,
			Uploaded:     sn.Uploaded,
			Downloaded:   sn.Downloaded,
			SeedTime:     uint64(sn.SeedTime),
			LeechTime:    uint64(sn.LeechTime),
			IsFinished:   sn.IsFinished,
			StartedAt:    sn.StartedAt,
			FinishedAt:   sn.CompletedAt,
			LastActionAt: sn.LastAction,
		},
	}, nil
}
