package v1

import (
	"server/internal/model/in/adminin"

	"github.com/gogf/gf/v2/frame/g"
)

type CatalogTorrentPinReq struct {
	g.Meta `path:"/catalog/torrents/{id}:pin" method:"post" tags:"AdminCatalog" summary:"置顶种子" perm:"admin:catalog/torrent:*"`
	adminin.CatalogTorrentPinInp
}
type CatalogTorrentPinRes struct{}

type CatalogTorrentUnpinReq struct {
	g.Meta `path:"/catalog/torrents/{id}:unpin" method:"post" tags:"AdminCatalog" summary:"取消置顶种子" perm:"admin:catalog/torrent:*"`
	adminin.CatalogTorrentUnpinInp
}
type CatalogTorrentUnpinRes struct{}

type CatalogTorrentFeatureReq struct {
	g.Meta `path:"/catalog/torrents/{id}:feature" method:"post" tags:"AdminCatalog" summary:"推荐种子" perm:"admin:catalog/torrent:*"`
	adminin.CatalogTorrentFeatureInp
}
type CatalogTorrentFeatureRes struct{}

type CatalogTorrentUnfeatureReq struct {
	g.Meta `path:"/catalog/torrents/{id}:unfeature" method:"post" tags:"AdminCatalog" summary:"取消推荐种子" perm:"admin:catalog/torrent:*"`
	adminin.CatalogTorrentUnfeatureInp
}
type CatalogTorrentUnfeatureRes struct{}

type CatalogTorrentPromotionReq struct {
	g.Meta `path:"/catalog/torrents/{id}:promotion" method:"post" tags:"AdminCatalog" summary:"设置种子优惠" perm:"admin:catalog/torrent:*"`
	adminin.CatalogTorrentPromotionInp
}
type CatalogTorrentPromotionRes struct{}

type CatalogTorrentClearPromotionReq struct {
	g.Meta `path:"/catalog/torrents/{id}:clearPromotion" method:"post" tags:"AdminCatalog" summary:"取消种子优惠" perm:"admin:catalog/torrent:*"`
	adminin.CatalogTorrentClearPromotionInp
}
type CatalogTorrentClearPromotionRes struct{}

type CatalogTorrentDeleteReq struct {
	g.Meta `path:"/catalog/torrents/{id}" method:"delete" tags:"AdminCatalog" summary:"强制删除种子" perm:"admin:catalog/torrent:*"`
	adminin.CatalogTorrentDeleteInp
}

type CatalogTorrentDeleteRes struct{}
