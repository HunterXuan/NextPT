package accounting

import (
	"context"
	"sort"

	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/accountingin"
	"server/internal/model/out/accountingout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sAccountingPeerUsecase struct{}

func init() {
	service.RegisterAccountingPeerUsecase(NewAccountingPeerUsecase())
}

func NewAccountingPeerUsecase() *sAccountingPeerUsecase {
	return &sAccountingPeerUsecase{}
}

func (s *sAccountingPeerUsecase) ListMyPeers(ctx context.Context, actor *model.Actor, in accountingin.PeerListInp) (*accountingout.PeerListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	page := in.Page
	size := in.Size
	status := in.Status
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if status == "" {
		status = "all"
	}

	seedingPeers, err := service.TrackerPeerDomain().GetUserSeedingPeers(ctx, actor.Id)
	if err != nil {
		return nil, err
	}
	leechingPeers, err := service.TrackerPeerDomain().GetUserLeechingPeers(ctx, actor.Id)
	if err != nil {
		return nil, err
	}

	allPeers := s.filterPeers(status, seedingPeers, leechingPeers)
	sort.SliceStable(allPeers, func(i, j int) bool {
		if allPeers[i].LastAction == nil {
			return false
		}
		if allPeers[j].LastAction == nil {
			return true
		}
		return allPeers[i].LastAction.Time.After(allPeers[j].LastAction.Time)
	})

	total := len(allPeers)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}
	peers := allPeers[start:end]

	torrentIds := make([]uint64, 0, len(peers))
	for _, peer := range peers {
		torrentIds = append(torrentIds, peer.TorrentId)
	}

	torrentMap, err := s.loadTorrentMap(ctx, torrentIds)
	if err != nil {
		return nil, err
	}

	list := make([]accountingout.PeerItem, 0, len(peers))
	for _, peer := range peers {
		tName := ""
		tSize := uint64(0)
		if t, ok := torrentMap[peer.TorrentId]; ok {
			tName = t.Name
			tSize = t.Size
		}

		list = append(list, accountingout.PeerItem{
			TorrentId:    peer.TorrentId,
			TorrentName:  tName,
			TorrentSize:  tSize,
			Uploaded:     peer.Uploaded,
			Downloaded:   peer.Downloaded,
			Remaining:    peer.Remaining,
			IsSeeder:     peer.IsSeeder,
			Agent:        peer.Agent,
			StartedAt:    peer.StartedAt,
			FinishedAt:   peer.FinishedAt,
			LastActionAt: peer.LastAction,
		})
	}

	return &accountingout.PeerListOut{
		Page:          page,
		Size:          size,
		Total:         total,
		SeedingTotal:  len(seedingPeers),
		LeechingTotal: len(leechingPeers),
		List:          list,
	}, nil
}

func (s *sAccountingPeerUsecase) filterPeers(status string, seedingPeers, leechingPeers []entity.TrackerPeer) []entity.TrackerPeer {
	if status == "seeding" {
		return seedingPeers
	}
	if status == "leeching" {
		return leechingPeers
	}

	peers := make([]entity.TrackerPeer, 0, len(seedingPeers)+len(leechingPeers))
	peers = append(peers, seedingPeers...)
	peers = append(peers, leechingPeers...)
	return peers
}

func (s *sAccountingPeerUsecase) loadTorrentMap(ctx context.Context, torrentIds []uint64) (map[uint64]*entity.CatalogTorrent, error) {
	if len(torrentIds) == 0 {
		return map[uint64]*entity.CatalogTorrent{}, nil
	}

	uniqueIds := make([]uint64, 0, len(torrentIds))
	seen := make(map[uint64]struct{}, len(torrentIds))
	for _, id := range torrentIds {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}
	if len(uniqueIds) == 0 {
		return map[uint64]*entity.CatalogTorrent{}, nil
	}

	torrents, err := service.CatalogTorrentDomain().GetTorrentsByIds(ctx, uniqueIds)
	if err != nil {
		return nil, err
	}

	torrentMap := make(map[uint64]*entity.CatalogTorrent, len(torrents))
	for _, t := range torrents {
		torrentMap[t.Id] = t
	}
	return torrentMap, nil
}
