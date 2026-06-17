package v1

import (
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"

	"github.com/gogf/gf/v2/frame/g"
)

type CommentCreateReq struct {
	g.Meta `path:"/torrents/{id}/comments" method:"post" tags:"Catalog Comment" summary:"发表评论" perm:"create:catalog/comment:*"`
	catalogin.CommentCreateInp
}

type CommentCreateRes struct {
	Id uint64 `json:"id"`
}

type CommentListReq struct {
	g.Meta `path:"/torrents/{id}/comments" method:"get" tags:"Catalog Comment" summary:"获取评论列表"`
	catalogin.CommentListInp
}

type CommentListRes struct {
	catalogout.CommentListOut
}

type CommentReportReq struct {
	g.Meta `path:"/torrents/{id}/comments/{cid}:report" method:"post" tags:"Catalog Comment" summary:"举报评论" perm:"read:catalog/comment:*"`
	catalogin.CommentReportInp
}

type CommentReportRes struct{}

type CommentToggleLikeReq struct {
	g.Meta `path:"/torrents/{id}/comments/{cid}:like" method:"post" tags:"Catalog Comment" summary:"点赞/取消点赞评论" perm:"read:catalog/comment:*"`
	catalogin.CommentToggleLikeInp
}

type CommentToggleLikeRes struct {
	catalogout.CommentToggleLikeOut
}

type CommentRewardReq struct {
	g.Meta `path:"/torrents/{id}/comments/{cid}:reward" method:"post" tags:"Catalog Comment" summary:"打赏评论" perm:"read:catalog/comment:*"`
	catalogin.CommentRewardInp
}

type CommentRewardRes struct{}
