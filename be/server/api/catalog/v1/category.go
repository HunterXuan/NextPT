package v1

import (
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"

	"github.com/gogf/gf/v2/frame/g"
)

type CategoryListReq struct {
	g.Meta `path:"/categories" method:"get" tags:"Catalog" summary:"获取分类列表"`
	catalogin.CategoryListInp
}

type CategoryListRes struct {
	catalogout.CategoryListOut
}

type TagGroupListReq struct {
	g.Meta `path:"/tag-groups" method:"get" tags:"Catalog" summary:"获取标签分组列表"`
	catalogin.TagGroupListInp
}

type TagGroupListRes struct {
	catalogout.TagGroupListOut
}
