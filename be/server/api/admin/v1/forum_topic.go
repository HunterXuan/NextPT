package v1

import (
	"server/internal/model/in/adminin"

	"github.com/gogf/gf/v2/frame/g"
)

type ForumTopicLockReq struct {
	g.Meta `path:"/forum/topics/{id}:lock" method:"post" tags:"AdminForum" summary:"锁定主题" perm:"admin:forum/topic:*"`
	adminin.ForumTopicLockInp
}
type ForumTopicLockRes struct{}

type ForumTopicUnlockReq struct {
	g.Meta `path:"/forum/topics/{id}:unlock" method:"post" tags:"AdminForum" summary:"解锁主题" perm:"admin:forum/topic:*"`
	adminin.ForumTopicUnlockInp
}
type ForumTopicUnlockRes struct{}

type ForumTopicPinReq struct {
	g.Meta `path:"/forum/topics/{id}:pin" method:"post" tags:"AdminForum" summary:"置顶主题" perm:"admin:forum/topic:*"`
	adminin.ForumTopicPinInp
}
type ForumTopicPinRes struct{}

type ForumTopicUnpinReq struct {
	g.Meta `path:"/forum/topics/{id}:unpin" method:"post" tags:"AdminForum" summary:"取消置顶主题" perm:"admin:forum/topic:*"`
	adminin.ForumTopicUnpinInp
}
type ForumTopicUnpinRes struct{}

type ForumTopicMoveReq struct {
	g.Meta `path:"/forum/topics/{id}:move" method:"post" tags:"AdminForum" summary:"移动主题" perm:"admin:forum/topic:*"`
	adminin.ForumTopicMoveInp
}
type ForumTopicMoveRes struct{}
