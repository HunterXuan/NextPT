// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mod

import (
	"context"

	"server/api/mod/v1"
)

type IModV1 interface {
	StaffMessageList(ctx context.Context, req *v1.StaffMessageListReq) (res *v1.StaffMessageListRes, err error)
	StaffMessageCreate(ctx context.Context, req *v1.StaffMessageCreateReq) (res *v1.StaffMessageCreateRes, err error)
}
