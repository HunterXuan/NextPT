package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) RequestList(ctx context.Context, req *v1.RequestListReq) (res *v1.RequestListRes, err error) {
	out, err := service.CatalogRequestUsecase().List(ctx, contexts.GetActor(ctx), req.RequestListInp)
	if err != nil {
		return nil, err
	}
	return &v1.RequestListRes{RequestListOut: *out}, nil
}

func (c *ControllerV1) RequestGet(ctx context.Context, req *v1.RequestGetReq) (res *v1.RequestGetRes, err error) {
	out, err := service.CatalogRequestUsecase().Get(ctx, contexts.GetActor(ctx), req.RequestGetInp)
	if err != nil {
		return nil, err
	}
	return &v1.RequestGetRes{RequestDetailOut: *out}, nil
}

func (c *ControllerV1) RequestCreate(ctx context.Context, req *v1.RequestCreateReq) (res *v1.RequestCreateRes, err error) {
	out, err := service.CatalogRequestUsecase().Create(ctx, contexts.GetActor(ctx), req.RequestCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.RequestCreateRes{RequestCreateOut: *out}, nil
}

func (c *ControllerV1) RequestClaim(ctx context.Context, req *v1.RequestClaimReq) (res *v1.RequestClaimRes, err error) {
	err = service.CatalogRequestUsecase().Claim(ctx, contexts.GetActor(ctx), req.RequestClaimInp)
	return
}

func (c *ControllerV1) RequestAbandon(ctx context.Context, req *v1.RequestAbandonReq) (res *v1.RequestAbandonRes, err error) {
	err = service.CatalogRequestUsecase().Abandon(ctx, contexts.GetActor(ctx), req.RequestAbandonInp)
	return
}

func (c *ControllerV1) RequestSubmit(ctx context.Context, req *v1.RequestSubmitReq) (res *v1.RequestSubmitRes, err error) {
	err = service.CatalogRequestUsecase().Submit(ctx, contexts.GetActor(ctx), req.RequestSubmitInp)
	return
}

func (c *ControllerV1) RequestComplete(ctx context.Context, req *v1.RequestCompleteReq) (res *v1.RequestCompleteRes, err error) {
	err = service.CatalogRequestUsecase().Complete(ctx, contexts.GetActor(ctx), req.RequestCompleteInp)
	return
}

func (c *ControllerV1) RequestCancel(ctx context.Context, req *v1.RequestCancelReq) (res *v1.RequestCancelRes, err error) {
	err = service.CatalogRequestUsecase().Cancel(ctx, contexts.GetActor(ctx), req.RequestCancelInp)
	return
}

func (c *ControllerV1) RequestCommentCreate(ctx context.Context, req *v1.RequestCommentCreateReq) (res *v1.RequestCommentCreateRes, err error) {
	out, err := service.CatalogCommentUsecase().CreateForRequest(ctx, contexts.GetActor(ctx), req.RequestCommentCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.RequestCommentCreateRes{CommentCreateOut: *out}, nil
}

func (c *ControllerV1) RequestCommentList(ctx context.Context, req *v1.RequestCommentListReq) (res *v1.RequestCommentListRes, err error) {
	out, err := service.CatalogCommentUsecase().ListForRequest(ctx, contexts.GetActor(ctx), req.RequestCommentListInp)
	if err != nil {
		return nil, err
	}
	return &v1.RequestCommentListRes{CommentListOut: *out}, nil
}

func (c *ControllerV1) RequestCommentToggleLike(ctx context.Context, req *v1.RequestCommentToggleLikeReq) (res *v1.RequestCommentToggleLikeRes, err error) {
	out, err := service.CatalogCommentUsecase().ToggleLikeForRequest(ctx, contexts.GetActor(ctx), req.RequestCommentActionInp)
	if err != nil {
		return nil, err
	}
	return &v1.RequestCommentToggleLikeRes{CommentToggleLikeOut: *out}, nil
}

func (c *ControllerV1) RequestCommentReward(ctx context.Context, req *v1.RequestCommentRewardReq) (res *v1.RequestCommentRewardRes, err error) {
	err = service.CatalogCommentUsecase().RewardForRequest(ctx, contexts.GetActor(ctx), req.RequestCommentRewardInp)
	return
}

func (c *ControllerV1) RequestCommentReport(ctx context.Context, req *v1.RequestCommentReportReq) (res *v1.RequestCommentReportRes, err error) {
	err = service.CatalogCommentUsecase().ReportForRequest(ctx, contexts.GetActor(ctx), req.RequestCommentReportInp)
	return
}
