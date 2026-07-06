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
	"github.com/gogf/gf/v2/os/gtime"
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

	userIds := make([]uint64, 0, len(records)*2)
	torrentIds := make([]uint64, 0, len(records))
	for _, r := range records {
		userIds = append(userIds, r.UserId)
		if r.DealtBy > 0 {
			userIds = append(userIds, r.DealtBy)
		}
		torrentIds = append(torrentIds, r.TorrentId)
	}
	userMap := s.loadUserSummaryMap(ctx, userIds)
	torrentMap := s.loadTorrentSummaryMap(ctx, torrentIds)

	var list []modout.CheaterLogItem
	for _, r := range records {
		list = append(list, modout.CheaterLogItem{
			Id:           r.Id,
			UserId:       r.UserId,
			User:         s.userSummary(userMap, r.UserId),
			TorrentId:    r.TorrentId,
			Torrent:      s.torrentSummary(torrentMap, r.TorrentId),
			Uploaded:     r.Uploaded,
			Downloaded:   r.Downloaded,
			AnnounceTime: r.AnnounceTime,
			Seeders:      r.Seeders,
			Leechers:     r.Leechers,
			HitCount:     r.HitCount,
			DealtBy:      r.DealtBy,
			DealtUser:    s.userSummary(userMap, r.DealtBy),
			IsDealt:      r.IsDealt,
			Comment:      r.Comment,
			DealtComment: r.DealtComment,
			DealtAt:      r.DealtAt,
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

	return service.ModCheaterDomain().Resolve(ctx, in.Id, actor.Id, in.Comment, gtime.Now())
}

func (s *sModCheaterUsecase) loadUserSummaryMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
	userMap := make(map[uint64]model.IamUserSummary)
	uniqueIds := s.uniqueUint64s(userIds)
	for _, id := range uniqueIds {
		userMap[id] = model.IamUserSummary{Id: id}
	}
	if len(uniqueIds) == 0 {
		return userMap
	}

	users, err := service.IamUserDomain().GetUsersByIds(ctx, uniqueIds)
	if err == nil {
		for _, user := range users {
			summary := userMap[user.Id]
			summary.Id = user.Id
			summary.Username = user.Username
			userMap[user.Id] = summary
		}
	}

	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, uniqueIds)
	if err == nil {
		for _, profile := range profiles {
			summary := userMap[profile.UserId]
			summary.Id = profile.UserId
			summary.Avatar = profile.Avatar
			userMap[profile.UserId] = summary
		}
	}
	return userMap
}

func (s *sModCheaterUsecase) loadTorrentSummaryMap(ctx context.Context, torrentIds []uint64) map[uint64]model.CatalogTorrentSummary {
	torrentMap := make(map[uint64]model.CatalogTorrentSummary)
	uniqueIds := s.uniqueUint64s(torrentIds)
	for _, id := range uniqueIds {
		torrentMap[id] = model.CatalogTorrentSummary{Id: id}
	}
	if len(uniqueIds) == 0 {
		return torrentMap
	}

	torrents, err := service.CatalogTorrentDomain().GetTorrentsByIds(ctx, uniqueIds)
	if err != nil {
		return torrentMap
	}
	for _, torrent := range torrents {
		if torrent == nil {
			continue
		}
		torrentMap[torrent.Id] = model.CatalogTorrentSummary{
			Id:    torrent.Id,
			Name:  torrent.Name,
			Size:  torrent.Size,
			Exist: true,
		}
	}
	return torrentMap
}

func (s *sModCheaterUsecase) uniqueUint64s(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	uniqueIds := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}
	return uniqueIds
}

func (s *sModCheaterUsecase) userSummary(userMap map[uint64]model.IamUserSummary, userId uint64) model.IamUserSummary {
	if userId == 0 {
		return model.IamUserSummary{}
	}
	if summary, ok := userMap[userId]; ok {
		return summary
	}
	return model.IamUserSummary{Id: userId}
}

func (s *sModCheaterUsecase) torrentSummary(torrentMap map[uint64]model.CatalogTorrentSummary, torrentId uint64) model.CatalogTorrentSummary {
	if torrentId == 0 {
		return model.CatalogTorrentSummary{}
	}
	if summary, ok := torrentMap[torrentId]; ok {
		return summary
	}
	return model.CatalogTorrentSummary{Id: torrentId}
}
