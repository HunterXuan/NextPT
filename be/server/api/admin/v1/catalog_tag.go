package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type CatalogTagGroupListReq struct {
	g.Meta `path:"/catalog/tag-groups" method:"get" tags:"AdminCatalog" summary:"获取标签组列表" perm:"admin:catalog/tag:*"`
	adminin.CatalogTagGroupListInp
}

type CatalogTagGroupListRes struct {
	adminout.CatalogTagGroupListOut
}

type CatalogTagGroupCreateReq struct {
	g.Meta `path:"/catalog/tag-groups" method:"post" tags:"AdminCatalog" summary:"创建标签组" perm:"admin:catalog/tag:*"`
	adminin.CatalogTagGroupCreateInp
}

type CatalogTagGroupCreateRes struct{}

type CatalogTagGroupUpdateReq struct {
	g.Meta `path:"/catalog/tag-groups/{id}" method:"patch" tags:"AdminCatalog" summary:"更新标签组" perm:"admin:catalog/tag:*"`
	adminin.CatalogTagGroupUpdateInp
}

type CatalogTagGroupUpdateRes struct{}

type CatalogTagGroupDeleteReq struct {
	g.Meta `path:"/catalog/tag-groups/{id}" method:"delete" tags:"AdminCatalog" summary:"删除标签组" perm:"admin:catalog/tag:*"`
	adminin.CatalogTagGroupDeleteInp
}

type CatalogTagGroupDeleteRes struct{}

type CatalogTagCreateReq struct {
	g.Meta `path:"/catalog/tag-groups/{groupId}/tags" method:"post" tags:"AdminCatalog" summary:"创建标签" perm:"admin:catalog/tag:*"`
	adminin.CatalogTagCreateInp
}

type CatalogTagCreateRes struct{}

type CatalogTagUpdateReq struct {
	g.Meta `path:"/catalog/tags/{id}" method:"patch" tags:"AdminCatalog" summary:"更新标签" perm:"admin:catalog/tag:*"`
	adminin.CatalogTagUpdateInp
}

type CatalogTagUpdateRes struct{}

type CatalogTagDeleteReq struct {
	g.Meta `path:"/catalog/tags/{id}" method:"delete" tags:"AdminCatalog" summary:"删除标签" perm:"admin:catalog/tag:*"`
	adminin.CatalogTagDeleteInp
}

type CatalogTagDeleteRes struct{}
