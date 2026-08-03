package catalog

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/catalogin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/catalogout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sCatalogRequestUsecase struct{}

const catalogRequestRewardAmountScale = 10

func init() {
	service.RegisterCatalogRequestUsecase(NewCatalogRequestUsecase())
}

func NewCatalogRequestUsecase() *sCatalogRequestUsecase {
	return &sCatalogRequestUsecase{}
}

func (s *sCatalogRequestUsecase) List(ctx context.Context, actor *model.Actor, in catalogin.RequestListInp) (*catalogout.RequestListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	requests, total, err := service.CatalogRequestDomain().QueryRequests(ctx, model.CatalogRequestListOptions{
		ActorId:     actor.Id,
		Keyword:     in.Keyword,
		RequestType: in.RequestType,
		Status:      in.Status,
		CategoryId:  in.CategoryId,
		View:        in.View,
		Page:        in.Page,
		Size:        in.Size,
	})
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.request.query_failed"))
	}
	return &catalogout.RequestListOut{
		List:  s.buildRequestListItems(ctx, requests),
		Total: total,
	}, nil
}

func (s *sCatalogRequestUsecase) Get(ctx context.Context, actor *model.Actor, in catalogin.RequestGetInp) (*catalogout.RequestDetailOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	request, err := service.CatalogRequestDomain().GetRequestById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	item := s.buildRequestListItems(ctx, []entity.CatalogRequest{*request})[0]
	targetTorrent, resultTorrent := s.requestTorrentSummaries(ctx, request)
	return &catalogout.RequestDetailOut{
		RequestListItem: item,
		Description:     request.Description,
		TargetTorrent:   targetTorrent,
		ResultTorrent:   resultTorrent,
		CancelReason:    request.CancelReason,
		Actions:         s.requestActions(ctx, actor, request),
	}, nil
}

func (s *sCatalogRequestUsecase) Create(ctx context.Context, actor *model.Actor, in catalogin.RequestCreateInp) (*catalogout.RequestCreateOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	rewardAmount, err := s.normalizeRequestReward(ctx, in.RewardAmount)
	if err != nil {
		return nil, err
	}

	now := gtime.Now()
	request := entity.CatalogRequest{
		RequestType:  in.RequestType,
		RequesterId:  actor.Id,
		Description:  strings.TrimSpace(in.Description),
		RewardAmount: rewardAmount,
		Status:       consts.CatalogRequestStatusOpen,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	switch in.RequestType {
	case consts.CatalogRequestTypeTorrent:
		title := strings.TrimSpace(in.Title)
		if title == "" {
			return nil, gerror.New(gi18n.T(ctx, "catalog.request.title_req"))
		}
		if _, err = service.CatalogCategoryDomain().GetCategoryById(ctx, in.CategoryId); err != nil {
			return nil, gerror.New(gi18n.T(ctx, "catalog.category.invalid"))
		}
		request.CategoryId = in.CategoryId
		request.Title = title
	case consts.CatalogRequestTypeReseed:
		torrent, loadErr := service.CatalogTorrentDomain().LoadViewableTorrent(ctx, actor, in.TargetTorrentId)
		if loadErr != nil {
			return nil, loadErr
		}
		if torrent.Seeders > 0 {
			return nil, gerror.New(gi18n.T(ctx, "catalog.request.reseed_has_seeders"))
		}
		request.CategoryId = torrent.CategoryId
		request.TargetTorrentId = torrent.Id
		request.Title = torrent.Name
	default:
		return nil, gerror.New(gi18n.T(ctx, "catalog.request.type_invalid"))
	}

	var requestId uint64
	err = dao.CatalogRequest.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if request.RequestType == consts.CatalogRequestTypeReseed {
			torrent, lockErr := service.CatalogTorrentDomain().GetTorrentByIdForUpdate(ctx, request.TargetTorrentId)
			if lockErr != nil {
				return lockErr
			}
			if torrent.Seeders > 0 {
				return gerror.New(gi18n.T(ctx, "catalog.request.reseed_has_seeders"))
			}
			exists, existsErr := service.CatalogRequestDomain().HasActiveReseedRequest(ctx, request.TargetTorrentId)
			if existsErr != nil {
				return existsErr
			}
			if exists {
				return gerror.New(gi18n.T(ctx, "catalog.request.reseed_exists"))
			}
		}

		id, insertErr := service.CatalogRequestDomain().InsertRequest(ctx, request)
		if insertErr != nil {
			return insertErr
		}
		requestId = id
		if permissionErr := service.IamPermissionDomain().GrantUserPermission(ctx, actor.Id, s.catalogRequestUpdatePermission(requestId), false); permissionErr != nil {
			return permissionErr
		}
		return s.debitRequestReward(ctx, actor.Id, rewardAmount, requestId, now)
	})
	if err != nil {
		return nil, err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, actor.Id)
	return &catalogout.RequestCreateOut{Id: requestId}, nil
}

func (s *sCatalogRequestUsecase) Claim(ctx context.Context, actor *model.Actor, in catalogin.RequestClaimInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	request, err := service.CatalogRequestDomain().GetRequestById(ctx, in.Id)
	if err != nil {
		return err
	}
	if request.Status != consts.CatalogRequestStatusOpen {
		return gerror.New(gi18n.T(ctx, "catalog.request.not_open"))
	}
	if request.RequesterId == actor.Id {
		return gerror.New(gi18n.T(ctx, "catalog.request.claim_self_denied"))
	}
	if request.RequestType == consts.CatalogRequestTypeTorrent {
		canPublish, permissionErr := service.IamUserUsecase().CheckPermission(ctx, actor, consts.IamPermissionCatalogTorrentCreate)
		if permissionErr != nil || !canPublish {
			return gerror.New(gi18n.T(ctx, "catalog.request.claim_publish_denied"))
		}
	}

	now := gtime.Now()
	expiresAt := now.Add(s.requestClaimDuration(request.RequestType))
	err = dao.CatalogRequest.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		claimed, claimErr := service.CatalogRequestDomain().ClaimRequest(ctx, request.Id, actor.Id, now, expiresAt)
		if claimErr != nil {
			return claimErr
		}
		if !claimed {
			return gerror.New(gi18n.T(ctx, "catalog.request.state_changed"))
		}
		return service.IamPermissionDomain().GrantUserPermission(ctx, actor.Id, s.catalogRequestUpdatePermission(request.Id), false)
	})
	if err != nil {
		return err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, actor.Id)
	s.notifyRequest(ctx, actor.Id, request.RequesterId, "site.message.catalog_request.claimed.title", "site.message.catalog_request.claimed.content", request)
	return nil
}

func (s *sCatalogRequestUsecase) Abandon(ctx context.Context, actor *model.Actor, in catalogin.RequestAbandonInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	request, err := service.CatalogRequestDomain().GetRequestById(ctx, in.Id)
	if err != nil {
		return err
	}
	if request.Status != consts.CatalogRequestStatusClaimed || request.ClaimedBy != actor.Id {
		return gerror.New(gi18n.T(ctx, "catalog.request.abandon_denied"))
	}
	err = dao.CatalogRequest.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		abandoned, abandonErr := service.CatalogRequestDomain().AbandonRequest(ctx, request.Id, actor.Id)
		if abandonErr != nil {
			return abandonErr
		}
		if !abandoned {
			return gerror.New(gi18n.T(ctx, "catalog.request.state_changed"))
		}
		return service.IamPermissionDomain().RevokeUserPermission(ctx, actor.Id, s.catalogRequestUpdatePermission(request.Id), false)
	})
	if err != nil {
		return err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, actor.Id)
	s.notifyRequest(ctx, actor.Id, request.RequesterId, "site.message.catalog_request.abandoned.title", "site.message.catalog_request.abandoned.content", request)
	return nil
}

func (s *sCatalogRequestUsecase) Submit(ctx context.Context, actor *model.Actor, in catalogin.RequestSubmitInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	request, err := service.CatalogRequestDomain().GetRequestById(ctx, in.Id)
	if err != nil {
		return err
	}
	now := gtime.Now()
	if request.Status != consts.CatalogRequestStatusClaimed || request.ClaimedBy != actor.Id {
		return gerror.New(gi18n.T(ctx, "catalog.request.submit_denied"))
	}
	if request.ClaimExpiresAt == nil || !request.ClaimExpiresAt.After(now) {
		return gerror.New(gi18n.T(ctx, "catalog.request.claim_expired"))
	}

	resultTorrentId := uint64(0)
	switch request.RequestType {
	case consts.CatalogRequestTypeTorrent:
		if in.ResultTorrentId == 0 {
			return gerror.New(gi18n.T(ctx, "catalog.request.result_torrent_req"))
		}
		if _, err = service.CatalogTorrentDomain().LoadViewableTorrent(ctx, actor, in.ResultTorrentId); err != nil {
			return err
		}
		resultTorrentId = in.ResultTorrentId
	case consts.CatalogRequestTypeReseed:
		seeding, seedingErr := s.isUserSeedingTorrent(ctx, actor.Id, request.TargetTorrentId)
		if seedingErr != nil {
			return seedingErr
		}
		if !seeding {
			return gerror.New(gi18n.T(ctx, "catalog.request.reseed_not_seeding"))
		}
	default:
		return gerror.New(gi18n.T(ctx, "catalog.request.type_invalid"))
	}

	err = dao.CatalogRequest.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		submitted, submitErr := service.CatalogRequestDomain().SubmitRequest(ctx, request.Id, actor.Id, resultTorrentId, now)
		if submitErr != nil {
			return submitErr
		}
		if !submitted {
			return gerror.New(gi18n.T(ctx, "catalog.request.state_changed"))
		}
		return service.IamPermissionDomain().RevokeUserPermission(ctx, actor.Id, s.catalogRequestUpdatePermission(request.Id), false)
	})
	if err != nil {
		return err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, actor.Id)
	s.notifyRequest(ctx, actor.Id, request.RequesterId, "site.message.catalog_request.submitted.title", "site.message.catalog_request.submitted.content", request)
	return nil
}

func (s *sCatalogRequestUsecase) Complete(ctx context.Context, actor *model.Actor, in catalogin.RequestCompleteInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	return s.complete(ctx, actor, in, false)
}

func (s *sCatalogRequestUsecase) CompleteByAdmin(ctx context.Context, actor *model.Actor, in catalogin.RequestCompleteInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	canManage, err := service.IamUserUsecase().CheckPermission(ctx, actor, consts.IamPermissionAdminCatalogRequestManage)
	if err != nil || !canManage {
		return gerror.New(gi18n.T(ctx, "catalog.request.complete_denied"))
	}
	return s.complete(ctx, actor, in, true)
}

func (s *sCatalogRequestUsecase) complete(ctx context.Context, actor *model.Actor, in catalogin.RequestCompleteInp, byAdmin bool) error {
	now := gtime.Now()
	var completedRequest *entity.CatalogRequest
	err := dao.CatalogRequest.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		request, lockErr := service.CatalogRequestDomain().GetRequestByIdForUpdate(ctx, in.Id)
		if lockErr != nil {
			return lockErr
		}
		if request.Status != consts.CatalogRequestStatusSubmitted {
			return gerror.New(gi18n.T(ctx, "catalog.request.not_submitted"))
		}
		if !byAdmin && request.RequesterId != actor.Id {
			return gerror.New(gi18n.T(ctx, "catalog.request.complete_denied"))
		}
		if request.ClaimedBy == 0 {
			return gerror.New(gi18n.T(ctx, "catalog.request.claimer_missing"))
		}
		completed, completeErr := service.CatalogRequestDomain().CompleteRequest(ctx, request.Id, now)
		if completeErr != nil {
			return completeErr
		}
		if !completed {
			return gerror.New(gi18n.T(ctx, "catalog.request.state_changed"))
		}
		if rewardErr := s.creditRequestReward(ctx, request.ClaimedBy, request.RewardAmount, request.Id, consts.EconomyBonusActionRequestReward, now); rewardErr != nil {
			return rewardErr
		}
		if permissionErr := service.IamPermissionDomain().RevokeUserPermission(ctx, request.RequesterId, s.catalogRequestUpdatePermission(request.Id), false); permissionErr != nil {
			return permissionErr
		}
		completedRequest = request
		return nil
	})
	if err != nil {
		return err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, completedRequest.RequesterId)
	s.notifyRequest(ctx, actor.Id, completedRequest.ClaimedBy, "site.message.catalog_request.completed.title", "site.message.catalog_request.completed.content", completedRequest)
	return nil
}

func (s *sCatalogRequestUsecase) Cancel(ctx context.Context, actor *model.Actor, in catalogin.RequestCancelInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	return s.cancel(ctx, actor, in, false)
}

func (s *sCatalogRequestUsecase) CancelByAdmin(ctx context.Context, actor *model.Actor, in catalogin.RequestCancelInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	canManage, err := service.IamUserUsecase().CheckPermission(ctx, actor, consts.IamPermissionAdminCatalogRequestManage)
	if err != nil || !canManage {
		return gerror.New(gi18n.T(ctx, "catalog.request.cancel_denied"))
	}
	return s.cancel(ctx, actor, in, true)
}

func (s *sCatalogRequestUsecase) cancel(ctx context.Context, actor *model.Actor, in catalogin.RequestCancelInp, byAdmin bool) error {
	now := gtime.Now()
	var cancelledRequest *entity.CatalogRequest
	err := dao.CatalogRequest.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		request, lockErr := service.CatalogRequestDomain().GetRequestByIdForUpdate(ctx, in.Id)
		if lockErr != nil {
			return lockErr
		}
		if request.Status == consts.CatalogRequestStatusCompleted || request.Status == consts.CatalogRequestStatusCancelled {
			return gerror.New(gi18n.T(ctx, "catalog.request.cancel_denied"))
		}
		if !byAdmin && (request.RequesterId != actor.Id || request.Status != consts.CatalogRequestStatusOpen) {
			return gerror.New(gi18n.T(ctx, "catalog.request.cancel_denied"))
		}
		cancelled, cancelErr := service.CatalogRequestDomain().CancelRequest(ctx, request.Id, actor.Id, in.Reason, now)
		if cancelErr != nil {
			return cancelErr
		}
		if !cancelled {
			return gerror.New(gi18n.T(ctx, "catalog.request.state_changed"))
		}
		if refundErr := s.creditRequestReward(ctx, request.RequesterId, request.RewardAmount, request.Id, consts.EconomyBonusActionRequestRefund, now); refundErr != nil {
			return refundErr
		}
		if permissionErr := service.IamPermissionDomain().RevokeUserPermission(ctx, request.RequesterId, s.catalogRequestUpdatePermission(request.Id), false); permissionErr != nil {
			return permissionErr
		}
		if request.ClaimedBy > 0 {
			if permissionErr := service.IamPermissionDomain().RevokeUserPermission(ctx, request.ClaimedBy, s.catalogRequestUpdatePermission(request.Id), false); permissionErr != nil {
				return permissionErr
			}
		}
		cancelledRequest = request
		return nil
	})
	if err != nil {
		return err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, cancelledRequest.RequesterId)
	if cancelledRequest.ClaimedBy > 0 {
		service.IamUserUsecase().InvalidateUserCache(ctx, cancelledRequest.ClaimedBy)
	}
	if cancelledRequest.RequesterId != actor.Id {
		s.notifyRequest(ctx, actor.Id, cancelledRequest.RequesterId, "site.message.catalog_request.cancelled.title", "site.message.catalog_request.cancelled.content", cancelledRequest)
	}
	if cancelledRequest.ClaimedBy > 0 {
		s.notifyRequest(ctx, actor.Id, cancelledRequest.ClaimedBy, "site.message.catalog_request.cancelled.title", "site.message.catalog_request.cancelled.content", cancelledRequest)
	}
	return nil
}

func (s *sCatalogRequestUsecase) ReleaseExpiredClaims(ctx context.Context, limit int) (int, error) {
	now := gtime.Now()
	requests, err := service.CatalogRequestDomain().QueryExpiredClaims(ctx, now, limit)
	if err != nil {
		return 0, err
	}
	releasedCount := 0
	for i := range requests {
		request := &requests[i]
		released := false
		releaseErr := dao.CatalogRequest.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			var err error
			released, err = service.CatalogRequestDomain().ReleaseExpiredClaim(ctx, request.Id, request.ClaimedBy, now)
			if err != nil || !released {
				return err
			}
			return service.IamPermissionDomain().RevokeUserPermission(ctx, request.ClaimedBy, s.catalogRequestUpdatePermission(request.Id), false)
		})
		if releaseErr != nil {
			return releasedCount, releaseErr
		}
		if !released {
			continue
		}
		releasedCount++
		service.IamUserUsecase().InvalidateUserCache(ctx, request.ClaimedBy)
		s.notifyRequest(ctx, 0, request.RequesterId, "site.message.catalog_request.expired.title", "site.message.catalog_request.expired.content", request)
		s.notifyRequest(ctx, 0, request.ClaimedBy, "site.message.catalog_request.expired.title", "site.message.catalog_request.expired.content", request)
	}
	return releasedCount, nil
}

func (s *sCatalogRequestUsecase) normalizeRequestReward(ctx context.Context, amount float64) (float64, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return 0, gerror.New(gi18n.T(ctx, "catalog.request.reward_invalid"))
	}
	amount = math.Round(amount*catalogRequestRewardAmountScale) / catalogRequestRewardAmountScale
	if amount <= 0 {
		return 0, gerror.New(gi18n.T(ctx, "catalog.request.reward_invalid"))
	}
	return amount, nil
}

func (s *sCatalogRequestUsecase) debitRequestReward(ctx context.Context, userId uint64, amount float64, requestId uint64, now *gtime.Time) error {
	if err := service.EconomyBonusDomain().DebitBonusIfEnough(ctx, userId, amount); err != nil {
		return err
	}
	balance, err := service.EconomyBonusDomain().GetUserBonus(ctx, userId)
	if err != nil {
		return err
	}
	return service.EconomyBonusDomain().InsertBonusLog(ctx, entity.EconomyBonusLog{
		UserId:       userId,
		Amount:       -amount,
		BalanceAfter: balance,
		Action:       consts.EconomyBonusActionRequestEscrow,
		TargetType:   consts.EconomyBonusTargetTypeCatalogRequest,
		TargetId:     requestId,
		CreatedAt:    now,
	})
}

func (s *sCatalogRequestUsecase) creditRequestReward(ctx context.Context, userId uint64, amount float64, requestId uint64, action string, now *gtime.Time) error {
	if err := service.EconomyBonusDomain().CreditBonus(ctx, userId, amount); err != nil {
		return err
	}
	balance, err := service.EconomyBonusDomain().GetUserBonus(ctx, userId)
	if err != nil {
		return err
	}
	return service.EconomyBonusDomain().InsertBonusLog(ctx, entity.EconomyBonusLog{
		UserId:       userId,
		Amount:       amount,
		BalanceAfter: balance,
		Action:       action,
		TargetType:   consts.EconomyBonusTargetTypeCatalogRequest,
		TargetId:     requestId,
		CreatedAt:    now,
	})
}

func (s *sCatalogRequestUsecase) requestClaimDuration(requestType uint) time.Duration {
	if requestType == consts.CatalogRequestTypeReseed {
		return consts.CatalogRequestReseedClaimDuration
	}
	return consts.CatalogRequestTorrentClaimDuration
}

func (s *sCatalogRequestUsecase) isUserSeedingTorrent(ctx context.Context, userId uint64, torrentId uint64) (bool, error) {
	peers, err := service.TrackerPeerDomain().GetUserSeedingPeers(ctx, userId)
	if err != nil {
		return false, err
	}
	for _, peer := range peers {
		if peer.TorrentId == torrentId && peer.IsSeeder {
			return true, nil
		}
	}
	return false, nil
}

func (s *sCatalogRequestUsecase) requestActions(ctx context.Context, actor *model.Actor, request *entity.CatalogRequest) model.CatalogRequestActions {
	if actor == nil || request == nil {
		return model.CatalogRequestActions{}
	}
	canRead, _ := service.IamUserUsecase().CheckPermission(ctx, actor, consts.IamPermissionCatalogRequestRead)
	canUpdate, _ := service.IamUserUsecase().CheckPermission(ctx, actor, s.catalogRequestUpdatePermission(request.Id))
	canManage, _ := service.IamUserUsecase().CheckPermission(ctx, actor, consts.IamPermissionAdminCatalogRequestManage)
	canPublish := true
	if request.RequestType == consts.CatalogRequestTypeTorrent {
		canPublish, _ = service.IamUserUsecase().CheckPermission(ctx, actor, consts.IamPermissionCatalogTorrentCreate)
	}
	now := gtime.Now()
	claimActive := request.ClaimExpiresAt != nil && request.ClaimExpiresAt.After(now)
	return model.CatalogRequestActions{
		CanClaim:    canRead && canPublish && request.Status == consts.CatalogRequestStatusOpen && request.RequesterId != actor.Id,
		CanAbandon:  canRead && request.Status == consts.CatalogRequestStatusClaimed && request.ClaimedBy == actor.Id,
		CanSubmit:   canUpdate && claimActive && request.Status == consts.CatalogRequestStatusClaimed && request.ClaimedBy == actor.Id,
		CanComplete: request.Status == consts.CatalogRequestStatusSubmitted && ((canUpdate && request.RequesterId == actor.Id) || canManage),
		CanCancel:   (canUpdate && request.Status == consts.CatalogRequestStatusOpen && request.RequesterId == actor.Id) || (canManage && request.Status < consts.CatalogRequestStatusCompleted),
	}
}

func (s *sCatalogRequestUsecase) buildRequestListItems(ctx context.Context, requests []entity.CatalogRequest) []catalogout.RequestListItem {
	userIds := make([]uint64, 0, len(requests)*2)
	for _, request := range requests {
		userIds = append(userIds, request.RequesterId)
		if request.ClaimedBy > 0 {
			userIds = append(userIds, request.ClaimedBy)
		}
	}
	userMap := s.requestUserSummaryMap(ctx, userIds)
	items := make([]catalogout.RequestListItem, 0, len(requests))
	for _, request := range requests {
		var claimer *model.IamUserSummary
		if request.ClaimedBy > 0 {
			summary := userMap[request.ClaimedBy]
			claimer = &summary
		}
		items = append(items, catalogout.RequestListItem{
			Id:              request.Id,
			RequestType:     request.RequestType,
			CategoryId:      request.CategoryId,
			TargetTorrentId: request.TargetTorrentId,
			ResultTorrentId: request.ResultTorrentId,
			Title:           request.Title,
			RewardAmount:    request.RewardAmount,
			Status:          request.Status,
			Requester:       userMap[request.RequesterId],
			Claimer:         claimer,
			ClaimExpiresAt:  s.catalogRequestTimeString(request.ClaimExpiresAt),
			SubmittedAt:     s.catalogRequestTimeString(request.SubmittedAt),
			CompletedAt:     s.catalogRequestTimeString(request.CompletedAt),
			CreatedAt:       s.catalogRequestTimeString(request.CreatedAt),
			UpdatedAt:       s.catalogRequestTimeString(request.UpdatedAt),
		})
	}
	return items
}

func (s *sCatalogRequestUsecase) requestUserSummaryMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
	result := make(map[uint64]model.IamUserSummary)
	uniqueIds := make([]uint64, 0, len(userIds))
	seen := make(map[uint64]struct{}, len(userIds))
	for _, userId := range userIds {
		if userId == 0 {
			continue
		}
		if _, exists := seen[userId]; exists {
			continue
		}
		seen[userId] = struct{}{}
		uniqueIds = append(uniqueIds, userId)
		result[userId] = model.IamUserSummary{Id: userId}
	}
	if len(uniqueIds) == 0 {
		return result
	}
	users, err := service.IamUserDomain().GetUsersByIds(ctx, uniqueIds)
	if err == nil {
		for _, user := range users {
			summary := result[user.Id]
			summary.Username = user.Username
			result[user.Id] = summary
		}
	}
	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, uniqueIds)
	if err == nil {
		for _, profile := range profiles {
			summary := result[profile.UserId]
			summary.Avatar = profile.Avatar
			result[profile.UserId] = summary
		}
	}
	return result
}

func (s *sCatalogRequestUsecase) requestTorrentSummaries(ctx context.Context, request *entity.CatalogRequest) (*model.CatalogTorrentSummary, *model.CatalogTorrentSummary) {
	ids := make([]uint64, 0, 2)
	if request.TargetTorrentId > 0 {
		ids = append(ids, request.TargetTorrentId)
	}
	if request.ResultTorrentId > 0 && request.ResultTorrentId != request.TargetTorrentId {
		ids = append(ids, request.ResultTorrentId)
	}
	torrentMap := make(map[uint64]*entity.CatalogTorrent, len(ids))
	if len(ids) > 0 {
		torrents, err := service.CatalogTorrentDomain().GetTorrentsByIds(ctx, ids)
		if err == nil {
			for _, torrent := range torrents {
				if torrent != nil {
					torrentMap[torrent.Id] = torrent
				}
			}
		}
	}
	return s.catalogRequestTorrentSummary(request.TargetTorrentId, torrentMap), s.catalogRequestTorrentSummary(request.ResultTorrentId, torrentMap)
}

func (s *sCatalogRequestUsecase) notifyRequest(ctx context.Context, actorId uint64, receiverId uint64, titleKey string, contentKey string, request *entity.CatalogRequest) {
	if request == nil {
		return
	}
	service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
		ActorId:     actorId,
		ReceiverId:  receiverId,
		TitleKey:    titleKey,
		ContentKey:  contentKey,
		ContentArgs: []any{request.Title},
		TargetType:  consts.SiteMessageTargetTypeCatalogRequest,
		TargetId:    request.Id,
	})
}

func (s *sCatalogRequestUsecase) catalogRequestTorrentSummary(id uint64, torrentMap map[uint64]*entity.CatalogTorrent) *model.CatalogTorrentSummary {
	if id == 0 {
		return nil
	}
	summary := &model.CatalogTorrentSummary{Id: id}
	if torrent := torrentMap[id]; torrent != nil {
		summary.Name = torrent.Name
		summary.Size = torrent.Size
		summary.Exist = true
	}
	return summary
}

func (s *sCatalogRequestUsecase) catalogRequestTimeString(value *gtime.Time) string {
	if value == nil {
		return ""
	}
	return value.String()
}

func (s *sCatalogRequestUsecase) catalogRequestUpdatePermission(requestId uint64) string {
	return fmt.Sprintf("update:catalog/request:%d", requestId)
}
