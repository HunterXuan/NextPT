package catalog

import (
	"context"
	"fmt"
	"strings"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sCatalogRequestDomain struct{}

func init() {
	service.RegisterCatalogRequestDomain(NewCatalogRequestDomain())
}

func NewCatalogRequestDomain() *sCatalogRequestDomain {
	return &sCatalogRequestDomain{}
}

func (s *sCatalogRequestDomain) InsertRequest(ctx context.Context, request entity.CatalogRequest) (uint64, error) {
	id, err := dao.CatalogRequest.Ctx(ctx).Data(request).InsertAndGetId()
	return uint64(id), err
}

func (s *sCatalogRequestDomain) GetRequestById(ctx context.Context, id uint64) (*entity.CatalogRequest, error) {
	var request entity.CatalogRequest
	err := dao.CatalogRequest.Ctx(ctx).Where(dao.CatalogRequest.Columns().Id, id).Scan(&request)
	if err != nil || request.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "catalog.request.not_found"))
	}
	return &request, nil
}

func (s *sCatalogRequestDomain) GetRequestByIdForUpdate(ctx context.Context, id uint64) (*entity.CatalogRequest, error) {
	var request entity.CatalogRequest
	err := dao.CatalogRequest.Ctx(ctx).LockUpdate().Where(dao.CatalogRequest.Columns().Id, id).Scan(&request)
	if err != nil || request.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "catalog.request.not_found"))
	}
	return &request, nil
}

func (s *sCatalogRequestDomain) QueryRequests(ctx context.Context, options model.CatalogRequestListOptions) ([]entity.CatalogRequest, int, error) {
	columns := dao.CatalogRequest.Columns()
	m := dao.CatalogRequest.Ctx(ctx)
	if keyword := strings.TrimSpace(options.Keyword); keyword != "" {
		m = m.WhereLike(columns.Title, "%"+keyword+"%")
	}
	if options.RequestType > 0 {
		m = m.Where(columns.RequestType, options.RequestType)
	}
	if options.Status != nil {
		m = m.Where(columns.Status, *options.Status)
	}
	if options.CategoryId > 0 {
		m = m.Where(columns.CategoryId, options.CategoryId)
	}
	switch options.View {
	case consts.CatalogRequestViewCreated:
		m = m.Where(columns.RequesterId, options.ActorId)
	case consts.CatalogRequestViewClaimed:
		m = m.Where(columns.ClaimedBy, options.ActorId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	page, size := s.normalizePage(options.Page, options.Size)
	var requests []entity.CatalogRequest
	err = m.Page(page, size).
		OrderAsc(columns.Status).
		OrderDesc(columns.UpdatedAt).
		OrderDesc(columns.Id).
		Scan(&requests)
	return requests, total, err
}

func (s *sCatalogRequestDomain) HasActiveReseedRequest(ctx context.Context, torrentId uint64) (bool, error) {
	columns := dao.CatalogRequest.Columns()
	count, err := dao.CatalogRequest.Ctx(ctx).
		Where(columns.RequestType, consts.CatalogRequestTypeReseed).
		Where(columns.TargetTorrentId, torrentId).
		WhereIn(columns.Status, []uint{
			consts.CatalogRequestStatusOpen,
			consts.CatalogRequestStatusClaimed,
			consts.CatalogRequestStatusSubmitted,
		}).
		Count()
	return count > 0, err
}

func (s *sCatalogRequestDomain) ClaimRequest(ctx context.Context, id uint64, userId uint64, claimedAt *gtime.Time, expiresAt *gtime.Time) (bool, error) {
	columns := dao.CatalogRequest.Columns()
	result, err := dao.CatalogRequest.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.CatalogRequestStatusOpen).
		Data(g.Map{
			columns.Status:         consts.CatalogRequestStatusClaimed,
			columns.ClaimedBy:      userId,
			columns.ClaimedAt:      claimedAt,
			columns.ClaimExpiresAt: expiresAt,
			columns.UpdatedAt:      claimedAt,
		}).Update()
	return s.affected(result, err)
}

func (s *sCatalogRequestDomain) AbandonRequest(ctx context.Context, id uint64, userId uint64) (bool, error) {
	columns := dao.CatalogRequest.Columns()
	now := gtime.Now()
	result, err := dao.CatalogRequest.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.CatalogRequestStatusClaimed).
		Where(columns.ClaimedBy, userId).
		Data(g.Map{
			columns.Status:         consts.CatalogRequestStatusOpen,
			columns.ClaimedBy:      0,
			columns.ClaimedAt:      nil,
			columns.ClaimExpiresAt: nil,
			columns.UpdatedAt:      now,
		}).Update()
	return s.affected(result, err)
}

func (s *sCatalogRequestDomain) SubmitRequest(ctx context.Context, id uint64, userId uint64, resultTorrentId uint64, submittedAt *gtime.Time) (bool, error) {
	columns := dao.CatalogRequest.Columns()
	result, err := dao.CatalogRequest.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.CatalogRequestStatusClaimed).
		Where(columns.ClaimedBy, userId).
		WhereGT(columns.ClaimExpiresAt, submittedAt).
		Data(g.Map{
			columns.Status:          consts.CatalogRequestStatusSubmitted,
			columns.ResultTorrentId: resultTorrentId,
			columns.ClaimExpiresAt:  nil,
			columns.SubmittedAt:     submittedAt,
			columns.UpdatedAt:       submittedAt,
		}).Update()
	return s.affected(result, err)
}

func (s *sCatalogRequestDomain) CompleteRequest(ctx context.Context, id uint64, completedAt *gtime.Time) (bool, error) {
	columns := dao.CatalogRequest.Columns()
	result, err := dao.CatalogRequest.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.CatalogRequestStatusSubmitted).
		Data(g.Map{
			columns.Status:      consts.CatalogRequestStatusCompleted,
			columns.CompletedAt: completedAt,
			columns.UpdatedAt:   completedAt,
		}).Update()
	return s.affected(result, err)
}

func (s *sCatalogRequestDomain) CancelRequest(ctx context.Context, id uint64, cancelledBy uint64, reason string, cancelledAt *gtime.Time) (bool, error) {
	columns := dao.CatalogRequest.Columns()
	result, err := dao.CatalogRequest.Ctx(ctx).
		Where(columns.Id, id).
		WhereNotIn(columns.Status, []uint{consts.CatalogRequestStatusCompleted, consts.CatalogRequestStatusCancelled}).
		Data(g.Map{
			columns.Status:       consts.CatalogRequestStatusCancelled,
			columns.CancelledBy:  cancelledBy,
			columns.CancelledAt:  cancelledAt,
			columns.CancelReason: strings.TrimSpace(reason),
			columns.UpdatedAt:    cancelledAt,
		}).Update()
	return s.affected(result, err)
}

func (s *sCatalogRequestDomain) QueryExpiredClaims(ctx context.Context, now *gtime.Time, limit int) ([]entity.CatalogRequest, error) {
	if limit <= 0 {
		limit = 100
	}
	columns := dao.CatalogRequest.Columns()
	var requests []entity.CatalogRequest
	err := dao.CatalogRequest.Ctx(ctx).
		Where(columns.Status, consts.CatalogRequestStatusClaimed).
		WhereLTE(columns.ClaimExpiresAt, now).
		OrderAsc(columns.ClaimExpiresAt).
		Limit(limit).
		Scan(&requests)
	return requests, err
}

func (s *sCatalogRequestDomain) ReleaseExpiredClaim(ctx context.Context, id uint64, claimedBy uint64, now *gtime.Time) (bool, error) {
	columns := dao.CatalogRequest.Columns()
	result, err := dao.CatalogRequest.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.CatalogRequestStatusClaimed).
		Where(columns.ClaimedBy, claimedBy).
		WhereLTE(columns.ClaimExpiresAt, now).
		Data(g.Map{
			columns.Status:         consts.CatalogRequestStatusOpen,
			columns.ClaimedBy:      0,
			columns.ClaimedAt:      nil,
			columns.ClaimExpiresAt: nil,
			columns.UpdatedAt:      now,
		}).Update()
	return s.affected(result, err)
}

func (s *sCatalogRequestDomain) normalizePage(page int, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

func (s *sCatalogRequestDomain) affected(result interface{ RowsAffected() (int64, error) }, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	affected, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		return false, fmt.Errorf("read affected rows: %w", affectedErr)
	}
	return affected > 0, nil
}
