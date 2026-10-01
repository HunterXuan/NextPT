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

	completedTransition := false
	if in.IsSeeder {
		// Check the last reported remaining size before replacing it with this announce.
		res, err := g.DB().Exec(ctx, `
			UPDATE tracker_snatch SET completed_at = ?
			WHERE torrent_id = ? AND user_id = ? AND remaining > 0 AND completed_at IS NULL
		`, in.EventTime, in.TorrentId, in.UserId)
		if err != nil {
			return false, err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return false, err
		}
		completedTransition = affected > 0
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
			is_finished = IF(?, 1, is_finished),
			last_action = ?
		WHERE torrent_id = ? AND user_id = ?
	`

	_, err := g.DB().Exec(ctx, sql,
		in.Ipv4, in.Ipv6, in.Port,
		in.UploadedDiff, in.DownloadedDiff, in.Remaining,
		seedTime, leechTime, in.IsSeeder,
		in.EventTime,
		in.TorrentId, in.UserId,
	)
	if err != nil {
		return false, err
	}
	return completedTransition, nil
}

func (s *sAccountingSnatchDomain) ListCompletions(ctx context.Context, torrentId uint64, page, size int) ([]*entity.TrackerSnatch, int, error) {
	m := dao.TrackerSnatch.Ctx(ctx).
		Where(dao.TrackerSnatch.Columns().TorrentId, torrentId).
		WhereNotNull(dao.TrackerSnatch.Columns().CompletedAt)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.TrackerSnatch
	err = m.Page(page, size).OrderDesc(dao.TrackerSnatch.Columns().CompletedAt).OrderDesc(dao.TrackerSnatch.Columns().Id).Scan(&list)
	return list, total, err
}

func (s *sAccountingSnatchDomain) ListSnatches(ctx context.Context, userId uint64, page, size int, isFinished *bool) ([]*entity.TrackerSnatch, int, error) {
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
