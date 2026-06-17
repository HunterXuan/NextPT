// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package accounting

import (
	"context"

	"server/api/accounting/v1"
)

type IAccountingV1 interface {
	SnatchList(ctx context.Context, req *v1.SnatchListReq) (res *v1.SnatchListRes, err error)
	SnatchGet(ctx context.Context, req *v1.SnatchGetReq) (res *v1.SnatchGetRes, err error)
	TrafficGetMe(ctx context.Context, req *v1.TrafficGetMeReq) (res *v1.TrafficGetMeRes, err error)
	TrafficHistoryList(ctx context.Context, req *v1.TrafficHistoryListReq) (res *v1.TrafficHistoryListRes, err error)
}
