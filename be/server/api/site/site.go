// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package site

import (
	"context"

	"server/api/site/v1"
)

type ISiteV1 interface {
	AdvertisementList(ctx context.Context, req *v1.AdvertisementListReq) (res *v1.AdvertisementListRes, err error)
	AnnouncementList(ctx context.Context, req *v1.AnnouncementListReq) (res *v1.AnnouncementListRes, err error)
	AnnouncementRead(ctx context.Context, req *v1.AnnouncementReadReq) (res *v1.AnnouncementReadRes, err error)
	MessageList(ctx context.Context, req *v1.MessageListReq) (res *v1.MessageListRes, err error)
	MessageRead(ctx context.Context, req *v1.MessageReadReq) (res *v1.MessageReadRes, err error)
	MessageReadAll(ctx context.Context, req *v1.MessageReadAllReq) (res *v1.MessageReadAllRes, err error)
	ChatMessageList(ctx context.Context, req *v1.ChatMessageListReq) (res *v1.ChatMessageListRes, err error)
	ChatMessageCreate(ctx context.Context, req *v1.ChatMessageCreateReq) (res *v1.ChatMessageCreateRes, err error)
	TaskList(ctx context.Context, req *v1.TaskListReq) (res *v1.TaskListRes, err error)
	TaskClaim(ctx context.Context, req *v1.TaskClaimReq) (res *v1.TaskClaimRes, err error)
	UserTaskRewardClaim(ctx context.Context, req *v1.UserTaskRewardClaimReq) (res *v1.UserTaskRewardClaimRes, err error)
}
