package v1

import (
	"server/internal/model/in/forumin"
	"server/internal/model/out/forumout"

	"github.com/gogf/gf/v2/frame/g"
)

type ReplyListReq struct {
	g.Meta `path:"/topics/{id}/replies" method:"get" tags:"Forum" summary:"获取主题回复列表" perm:"read:forum/reply:*"`
	forumin.ReplyListInp
}

type ReplyListRes struct {
	forumout.ReplyListOut
}

type ReplyCreateReq struct {
	g.Meta `path:"/topics/{id}/replies" method:"post" tags:"Forum" summary:"发表回复" perm:"create:forum/reply:*"`
	forumin.ReplyCreateInp
}

type ReplyCreateRes struct {
	forumout.ReplyCreateOut
}

type ReplyToggleLikeReq struct {
	g.Meta `path:"/replies/{id}:like" method:"post" tags:"Forum" summary:"点赞/取消点赞回复" perm:"read:forum/reply:*"`
	forumin.ReplyToggleLikeInp
}
type ReplyToggleLikeRes struct{}

type ReplyRewardReq struct {
	g.Meta `path:"/replies/{id}:reward" method:"post" tags:"Forum" summary:"打赏回复" perm:"read:forum/reply:*"`
	forumin.ReplyRewardInp
}
type ReplyRewardRes struct{}

type ReplyReportReq struct {
	g.Meta `path:"/replies/{id}:report" method:"post" tags:"Forum" summary:"举报回复" perm:"read:forum/reply:*"`
	forumin.ReplyReportInp
}
type ReplyReportRes struct{}
