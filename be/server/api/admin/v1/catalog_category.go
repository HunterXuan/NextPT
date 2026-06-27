package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type CatalogCategoryCreateReq struct {
	g.Meta `path:"/catalog/categories" method:"post" tags:"AdminCatalog" summary:"创建种子分类" perm:"admin:catalog/category:*"`
	adminin.CatalogCategoryCreateInp
}

type CatalogCategoryCreateRes struct{}

type CatalogCategoryUpdateReq struct {
	g.Meta `path:"/catalog/categories/{id}" method:"patch" tags:"AdminCatalog" summary:"更新种子分类" perm:"admin:catalog/category:*"`
	adminin.CatalogCategoryUpdateInp
}

type CatalogCategoryUpdateRes struct{}

type CatalogCategoryDeleteReq struct {
	g.Meta `path:"/catalog/categories/{id}" method:"delete" tags:"AdminCatalog" summary:"删除种子分类" perm:"admin:catalog/category:*"`
	adminin.CatalogCategoryDeleteInp
}

type CatalogCategoryDeleteRes struct{}

type CatalogCategoryListReq struct {
	g.Meta `path:"/catalog/categories" method:"get" tags:"AdminCatalog" summary:"获取种子分类列表" perm:"admin:catalog/category:*"`
	adminin.CatalogCategoryListInp
}

type CatalogCategoryListRes struct {
	adminout.CatalogCategoryListOut
}
