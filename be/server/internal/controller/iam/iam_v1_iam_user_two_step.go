package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserTwoStepSetup(ctx context.Context, req *v1.IamUserTwoStepSetupReq) (res *v1.IamUserTwoStepSetupRes, err error) {
	out, err := service.IamTwoStepUsecase().Setup(ctx, contexts.GetActor(ctx), req.UserTwoStepSetupInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserTwoStepSetupRes{UserTwoStepSetupOut: *out}, nil
}

func (c *ControllerV1) IamUserTwoStepConfirm(ctx context.Context, req *v1.IamUserTwoStepConfirmReq) (res *v1.IamUserTwoStepConfirmRes, err error) {
	out, err := service.IamTwoStepUsecase().Confirm(ctx, contexts.GetActor(ctx), req.UserTwoStepConfirmInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserTwoStepConfirmRes{UserTwoStepRecoveryCodesOut: *out}, nil
}

func (c *ControllerV1) IamUserTwoStepRecoveryCodesCreate(ctx context.Context, req *v1.IamUserTwoStepRecoveryCodesCreateReq) (res *v1.IamUserTwoStepRecoveryCodesCreateRes, err error) {
	out, err := service.IamTwoStepUsecase().CreateRecoveryCodes(ctx, contexts.GetActor(ctx), req.UserTwoStepRecoveryCodesCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserTwoStepRecoveryCodesCreateRes{UserTwoStepRecoveryCodesOut: *out}, nil
}

func (c *ControllerV1) IamUserTwoStepDelete(ctx context.Context, req *v1.IamUserTwoStepDeleteReq) (res *v1.IamUserTwoStepDeleteRes, err error) {
	if err := service.IamTwoStepUsecase().Disable(ctx, contexts.GetActor(ctx), req.UserTwoStepDisableInp); err != nil {
		return nil, err
	}
	return &v1.IamUserTwoStepDeleteRes{}, nil
}
