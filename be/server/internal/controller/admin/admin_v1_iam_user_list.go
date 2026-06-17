package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserList(ctx context.Context, req *v1.IamUserListReq) (res *v1.IamUserListRes, err error) {
	out, err := service.AdminIamUserUsecase().List(ctx, contexts.GetActor(ctx), req.IamUserListInp)
	if err != nil {
		return nil, err
	}
	res = &v1.IamUserListRes{IamUserListOut: *out}
	return res, nil
}
