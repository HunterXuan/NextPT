// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package site

import (
	"context"

	"server/api/site/v1"
)

type ISiteV1 interface {
	AnnouncementList(ctx context.Context, req *v1.AnnouncementListReq) (res *v1.AnnouncementListRes, err error)
	AnnouncementRead(ctx context.Context, req *v1.AnnouncementReadReq) (res *v1.AnnouncementReadRes, err error)
	MessageList(ctx context.Context, req *v1.MessageListReq) (res *v1.MessageListRes, err error)
	MessageRead(ctx context.Context, req *v1.MessageReadReq) (res *v1.MessageReadRes, err error)
	MessageReadAll(ctx context.Context, req *v1.MessageReadAllReq) (res *v1.MessageReadAllRes, err error)
}
