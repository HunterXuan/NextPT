package mod

import (
	"context"

	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/modin"
	"server/internal/model/out/modout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sModCheaterUsecase struct{}

func init() {
	service.RegisterModCheaterUsecase(NewModCheaterUsecase())
}

func NewModCheaterUsecase() *sModCheaterUsecase {
	return &sModCheaterUsecase{}
}

func (s *sModCheaterUsecase) Record(ctx context.Context, in modin.RecordCheaterLogInp) error {
	return service.ModCheaterDomain().Create(ctx, entity.ModCheaterLog{
		UserId:       in.UserId,
		TorrentId:    in.TorrentId,
		Uploaded:     in.Uploaded,
		Downloaded:   in.Downloaded,
		AnnounceTime: in.AnnounceTime,
		Seeders:      in.Seeders,
		Leechers:     in.Leechers,
		HitCount:     in.HitCount,
		Comment:      in.Comment,
		IsDealt:      false,
	})
}

func (s *sModCheaterUsecase) List(ctx context.Context, actor *model.Actor, in modin.ListCheaterLogsInp) (*modout.ListCheaterLogsOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	records, total, err := service.ModCheaterDomain().QueryCheaterLogs(ctx, in.IsDealt, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	var list []modout.CheaterLogItem
	for _, r := range records {
		list = append(list, modout.CheaterLogItem{
			Id:           r.Id,
			UserId:       r.UserId,
			TorrentId:    r.TorrentId,
			Uploaded:     r.Uploaded,
			Downloaded:   r.Downloaded,
			AnnounceTime: r.AnnounceTime,
			Seeders:      r.Seeders,
			Leechers:     r.Leechers,
			HitCount:     r.HitCount,
			DealtBy:      r.DealtBy,
			IsDealt:      r.IsDealt,
			Comment:      r.Comment,
			CreatedAt:    r.CreatedAt,
		})
	}

	return &modout.ListCheaterLogsOut{
		Page:  in.Page,
		Size:  in.Size,
		Total: total,
		List:  list,
	}, nil
}

func (s *sModCheaterUsecase) Resolve(ctx context.Context, actor *model.Actor, in modin.ResolveCheaterLogInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	return service.ModCheaterDomain().Update(ctx, in.Id, entity.ModCheaterLog{
		IsDealt: true,
		DealtBy: actor.Id,
		Comment: in.Comment,
	})
}
