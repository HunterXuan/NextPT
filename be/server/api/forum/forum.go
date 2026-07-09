// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package forum

import (
	"context"

	"server/api/forum/v1"
)

type IForumV1 interface {
	BookmarkList(ctx context.Context, req *v1.BookmarkListReq) (res *v1.BookmarkListRes, err error)
	NodeList(ctx context.Context, req *v1.NodeListReq) (res *v1.NodeListRes, err error)
	ReplyList(ctx context.Context, req *v1.ReplyListReq) (res *v1.ReplyListRes, err error)
	ReplyCreate(ctx context.Context, req *v1.ReplyCreateReq) (res *v1.ReplyCreateRes, err error)
	ReplyToggleLike(ctx context.Context, req *v1.ReplyToggleLikeReq) (res *v1.ReplyToggleLikeRes, err error)
	ReplyReward(ctx context.Context, req *v1.ReplyRewardReq) (res *v1.ReplyRewardRes, err error)
	ReplyReport(ctx context.Context, req *v1.ReplyReportReq) (res *v1.ReplyReportRes, err error)
	TopicList(ctx context.Context, req *v1.TopicListReq) (res *v1.TopicListRes, err error)
	TopicGetHot(ctx context.Context, req *v1.TopicGetHotReq) (res *v1.TopicGetHotRes, err error)
	TopicDetail(ctx context.Context, req *v1.TopicDetailReq) (res *v1.TopicDetailRes, err error)
	TopicCreate(ctx context.Context, req *v1.TopicCreateReq) (res *v1.TopicCreateRes, err error)
	TopicUpdate(ctx context.Context, req *v1.TopicUpdateReq) (res *v1.TopicUpdateRes, err error)
	TopicAppend(ctx context.Context, req *v1.TopicAppendReq) (res *v1.TopicAppendRes, err error)
	TopicToggleLike(ctx context.Context, req *v1.TopicToggleLikeReq) (res *v1.TopicToggleLikeRes, err error)
	TopicReward(ctx context.Context, req *v1.TopicRewardReq) (res *v1.TopicRewardRes, err error)
	TopicRewardList(ctx context.Context, req *v1.TopicRewardListReq) (res *v1.TopicRewardListRes, err error)
	TopicReport(ctx context.Context, req *v1.TopicReportReq) (res *v1.TopicReportRes, err error)
	TopicBookmark(ctx context.Context, req *v1.TopicBookmarkReq) (res *v1.TopicBookmarkRes, err error)
	TopicUnbookmark(ctx context.Context, req *v1.TopicUnbookmarkReq) (res *v1.TopicUnbookmarkRes, err error)
}
