package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type ForumNodeCreateReq struct {
	g.Meta `path:"/forum/nodes" method:"post" tags:"AdminForum" summary:"创建论坛节点" perm:"admin:forum/node:*"`
	adminin.ForumNodeCreateInp
}

type ForumNodeCreateRes struct{}

type ForumNodeUpdateReq struct {
	g.Meta `path:"/forum/nodes/{id}" method:"patch" tags:"AdminForum" summary:"更新论坛节点" perm:"admin:forum/node:*"`
	adminin.ForumNodeUpdateInp
}

type ForumNodeUpdateRes struct{}

type ForumNodeDeleteReq struct {
	g.Meta `path:"/forum/nodes/{id}" method:"delete" tags:"AdminForum" summary:"删除论坛节点" perm:"admin:forum/node:*"`
	adminin.ForumNodeDeleteInp
}

type ForumNodeDeleteRes struct{}

type ForumNodeListReq struct {
	g.Meta `path:"/forum/nodes" method:"get" tags:"AdminForum" summary:"获取论坛节点列表" perm:"admin:forum/node:*"`
	adminin.ForumNodeListInp
}

type ForumNodeListRes struct {
	adminout.ForumNodeListOut
}
