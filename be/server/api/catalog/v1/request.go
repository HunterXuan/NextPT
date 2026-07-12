package v1

import (
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"

	"github.com/gogf/gf/v2/frame/g"
)

type RequestListReq struct {
	g.Meta `path:"/requests" method:"get" tags:"Catalog Request" summary:"获取求种与续种请求" perm:"read:catalog/request:*"`
	catalogin.RequestListInp
}

type RequestListRes struct {
	catalogout.RequestListOut
}

type RequestGetReq struct {
	g.Meta `path:"/requests/{id}" method:"get" tags:"Catalog Request" summary:"获取请求详情" perm:"read:catalog/request:*"`
	catalogin.RequestGetInp
}

type RequestGetRes struct {
	catalogout.RequestDetailOut
}

type RequestCreateReq struct {
	g.Meta `path:"/requests" method:"post" tags:"Catalog Request" summary:"创建求种或续种请求" perm:"create:catalog/request:*"`
	catalogin.RequestCreateInp
}

type RequestCreateRes struct {
	catalogout.RequestCreateOut
}

type RequestClaimReq struct {
	g.Meta `path:"/requests/{id}:claim" method:"post" tags:"Catalog Request" summary:"认领请求" perm:"read:catalog/request:*"`
	catalogin.RequestClaimInp
}

type RequestClaimRes struct{}

type RequestAbandonReq struct {
	g.Meta `path:"/requests/{id}:abandon" method:"post" tags:"Catalog Request" summary:"放弃认领" perm:"read:catalog/request:*"`
	catalogin.RequestAbandonInp
}

type RequestAbandonRes struct{}

type RequestSubmitReq struct {
	g.Meta `path:"/requests/{id}:submit" method:"post" tags:"Catalog Request" summary:"提交请求结果" perm:"update:catalog/request:{id}"`
	catalogin.RequestSubmitInp
}

type RequestSubmitRes struct{}

type RequestCompleteReq struct {
	g.Meta `path:"/requests/{id}:complete" method:"post" tags:"Catalog Request" summary:"确认请求完成" perm:"update:catalog/request:{id}"`
	catalogin.RequestCompleteInp
}

type RequestCompleteRes struct{}

type RequestCancelReq struct {
	g.Meta `path:"/requests/{id}:cancel" method:"post" tags:"Catalog Request" summary:"取消请求" perm:"update:catalog/request:{id}"`
	catalogin.RequestCancelInp
}

type RequestCancelRes struct{}

type RequestCommentCreateReq struct {
	g.Meta `path:"/requests/{id}/comments" method:"post" tags:"Catalog Request Comment" summary:"发表请求评论" perm:"create:catalog/comment:*"`
	catalogin.RequestCommentCreateInp
}

type RequestCommentCreateRes struct {
	catalogout.CommentCreateOut
}

type RequestCommentListReq struct {
	g.Meta `path:"/requests/{id}/comments" method:"get" tags:"Catalog Request Comment" summary:"获取请求评论" perm:"read:catalog/comment:*"`
	catalogin.RequestCommentListInp
}

type RequestCommentListRes struct {
	catalogout.CommentListOut
}

type RequestCommentToggleLikeReq struct {
	g.Meta `path:"/requests/{id}/comments/{cid}:like" method:"post" tags:"Catalog Request Comment" summary:"点赞请求评论" perm:"read:catalog/comment:*"`
	catalogin.RequestCommentActionInp
}

type RequestCommentToggleLikeRes struct {
	catalogout.CommentToggleLikeOut
}

type RequestCommentRewardReq struct {
	g.Meta `path:"/requests/{id}/comments/{cid}:reward" method:"post" tags:"Catalog Request Comment" summary:"赞赏请求评论" perm:"read:catalog/comment:*"`
	catalogin.RequestCommentRewardInp
}

type RequestCommentRewardRes struct{}

type RequestCommentReportReq struct {
	g.Meta `path:"/requests/{id}/comments/{cid}:report" method:"post" tags:"Catalog Request Comment" summary:"举报请求评论" perm:"read:catalog/comment:*"`
	catalogin.RequestCommentReportInp
}

type RequestCommentReportRes struct{}
