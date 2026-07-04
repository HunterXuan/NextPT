package iam

import (
	"context"

	"server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserPasskeyReset(ctx context.Context, req *v1.IamUserPasskeyResetReq) (res *v1.IamUserPasskeyResetRes, err error) {
	passkey, err := service.IamUserUsecase().ResetPasskey(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.IamUserPasskeyResetRes{Passkey: passkey}, nil
}
