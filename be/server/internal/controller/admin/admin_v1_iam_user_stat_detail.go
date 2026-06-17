package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserStatDetail(ctx context.Context, req *v1.IamUserStatDetailReq) (res *v1.IamUserStatDetailRes, err error) {
	out, err := service.AdminIamUserUsecase().StatDetail(ctx, contexts.GetActor(ctx), req.IamUserStatDetailInp)
	if err != nil {
		return nil, err
	}
	res = &v1.IamUserStatDetailRes{IamUserStatDetailOut: *out}
	return
}
