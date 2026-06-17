package accounting

import (
	"context"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/accountingin"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sAccountingSnatchDomain struct{}

func init() {
	service.RegisterAccountingSnatchDomain(NewAccountingSnatchDomain())
}

func NewAccountingSnatchDomain() *sAccountingSnatchDomain {
	return &sAccountingSnatchDomain{}
}

func (s *sAccountingSnatchDomain) RecordSnatch(ctx context.Context, in accountingin.RecordSnatchInp) (bool, error) {
	seedTime := 0
	leechTime := 0
	if in.IsSeeder {
		seedTime = in.TimeDiff
	} else {
		leechTime = in.TimeDiff
	}

	if _, err := dao.TrackerSnatch.Ctx(ctx).InsertIgnore(do.TrackerSnatch{
		TorrentId: in.TorrentId,
		UserId:    in.UserId,
		StartedAt: in.EventTime,
	}); err != nil {
		return false, err
	}

	sql := `
		UPDATE tracker_snatch SET
			ipv4 = ?,
			ipv6 = ?,
			port = ?,
			uploaded = uploaded + ?,
			downloaded = downloaded + ?,
			remaining = ?,
			seed_time = seed_time + ?,
			leech_time = leech_time + ?,
			last_action = ?
		WHERE torrent_id = ? AND user_id = ?
	`

	_, err := g.DB().Exec(ctx, sql,
		in.Ipv4, in.Ipv6, in.Port,
		in.UploadedDiff, in.DownloadedDiff, in.Remaining,
		seedTime, leechTime,
		in.EventTime,
		in.TorrentId, in.UserId,
	)
	if err != nil {
		return false, err
	}
	if !in.IsFinished {
		return false, nil
	}

	res, err := g.DB().Exec(ctx, `
		UPDATE tracker_snatch SET
			is_finished = 1,
			completed_at = IF(completed_at IS NULL, ?, completed_at)
		WHERE torrent_id = ? AND user_id = ? AND is_finished = 0
	`, in.EventTime, in.TorrentId, in.UserId)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *sAccountingSnatchDomain) ListSnatches(ctx context.Context, userId uint64, page, size int, isFinished, isActive *bool) ([]*entity.TrackerSnatch, int, error) {
	m := dao.TrackerSnatch.Ctx(ctx).Where("user_id", userId)
	if isFinished != nil {
		m = m.Where("is_finished", *isFinished)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*entity.TrackerSnatch
	err = m.Page(page, size).OrderDesc("last_action").Scan(&list)
	return list, total, err
}

func (s *sAccountingSnatchDomain) GetSnatch(ctx context.Context, userId uint64, torrentId uint64) (*entity.TrackerSnatch, error) {
	var snatch *entity.TrackerSnatch
	err := dao.TrackerSnatch.Ctx(ctx).Where("user_id", userId).Where("torrent_id", torrentId).Scan(&snatch)
	return snatch, err
}

func (s *sAccountingSnatchDomain) DeleteSnatchesByTorrentId(ctx context.Context, torrentId uint64) error {
	_, err := dao.TrackerSnatch.Ctx(ctx).Where(dao.TrackerSnatch.Columns().TorrentId, torrentId).Delete()
	return err
}
