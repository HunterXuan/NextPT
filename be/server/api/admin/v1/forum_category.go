package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type ForumCategoryCreateReq struct {
	g.Meta `path:"/forum/node-categories" method:"post" tags:"AdminForum" summary:"创建论坛节点分类" perm:"admin:forum/category:*"`
	adminin.ForumCategoryCreateInp
}

type ForumCategoryCreateRes struct{}

type ForumCategoryUpdateReq struct {
	g.Meta `path:"/forum/node-categories/{id}" method:"patch" tags:"AdminForum" summary:"更新论坛节点分类" perm:"admin:forum/category:*"`
	adminin.ForumCategoryUpdateInp
}

type ForumCategoryUpdateRes struct{}

type ForumCategoryDeleteReq struct {
	g.Meta `path:"/forum/node-categories/{id}" method:"delete" tags:"AdminForum" summary:"删除论坛节点分类" perm:"admin:forum/category:*"`
	adminin.ForumCategoryDeleteInp
}

type ForumCategoryDeleteRes struct{}

type ForumCategoryListReq struct {
	g.Meta `path:"/forum/node-categories" method:"get" tags:"AdminForum" summary:"获取论坛节点分类列表" perm:"admin:forum/category:*"`
	adminin.ForumCategoryListInp
}

type ForumCategoryListRes struct {
	adminout.ForumCategoryListOut
}
