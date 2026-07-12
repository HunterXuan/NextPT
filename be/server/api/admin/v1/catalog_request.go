package v1

import (
	"server/internal/model/in/catalogin"

	"github.com/gogf/gf/v2/frame/g"
)

type CatalogRequestCompleteReq struct {
	g.Meta `path:"/catalog/requests/{id}:complete" method:"post" tags:"AdminCatalog" summary:"确认请求完成" perm:"admin:catalog/request:*"`
	catalogin.RequestCompleteInp
}

type CatalogRequestCompleteRes struct{}

type CatalogRequestCancelReq struct {
	g.Meta `path:"/catalog/requests/{id}:cancel" method:"post" tags:"AdminCatalog" summary:"取消请求" perm:"admin:catalog/request:*"`
	catalogin.RequestCancelInp
}

type CatalogRequestCancelRes struct{}
