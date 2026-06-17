package v1

import (
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"

	"github.com/gogf/gf/v2/frame/g"
)

type TorrentUploadReq struct {
	g.Meta `path:"/torrents" method:"post" tags:"Catalog" summary:"上传种子文件" mime:"multipart/form-data" perm:"create:catalog/torrent:*"`
	catalogin.TorrentUploadInp
}

type TorrentUploadRes struct {
	catalogout.TorrentUploadOut
}

type TorrentListReq struct {
	g.Meta `path:"/torrents" method:"get" tags:"Catalog" summary:"获取种子列表"`
	catalogin.TorrentListInp
}

type TorrentListRes struct {
	catalogout.TorrentListOut
}

type TorrentGetReq struct {
	g.Meta `path:"/torrents/{id}" method:"get" tags:"Catalog" summary:"获取种子详情"`
	catalogin.TorrentGetInp
}

type TorrentGetRes struct {
	catalogout.TorrentDetailOut
}

type TorrentDownloadReq struct {
	g.Meta `path:"/torrents/{id}:download" method:"get" tags:"Catalog" summary:"下载种子文件" perm:"download:catalog/torrent:*"`
	catalogin.TorrentDownloadInp
}

type TorrentDownloadRes struct {
	// 空结构体，因为是通过 r.Response.ServeContent 返回二进制文件
}

type TorrentRewardReq struct {
	g.Meta `path:"/torrents/{id}:reward" method:"post" tags:"Catalog" summary:"赞赏种子" perm:"read:catalog/torrent:*"`
	catalogin.TorrentRewardInp
}

type TorrentRewardRes struct {
	catalogout.TorrentRewardOut
}

type TorrentRewardListReq struct {
	g.Meta `path:"/torrents/{id}/rewards" method:"get" tags:"Catalog" summary:"获取种子赞赏列表"`
	catalogin.TorrentRewardListInp
}

type TorrentRewardListRes struct {
	catalogout.TorrentRewardListOut
}

type TorrentBookmarkReq struct {
	g.Meta `path:"/torrents/{id}:bookmark" method:"post" tags:"Catalog" summary:"收藏种子" perm:"read:catalog/torrent:*"`
	catalogin.TorrentBookmarkInp
}

type TorrentBookmarkRes struct {
	// 空
}

type TorrentUnbookmarkReq struct {
	g.Meta `path:"/torrents/{id}:unbookmark" method:"post" tags:"Catalog" summary:"取消收藏种子" perm:"read:catalog/torrent:*"`
	catalogin.TorrentUnbookmarkInp
}

type TorrentUnbookmarkRes struct {
	// 空
}

type TorrentToggleLikeReq struct {
	g.Meta `path:"/torrents/{id}:like" method:"post" tags:"Catalog" summary:"点赞/取消点赞种子" perm:"read:catalog/torrent:*"`
	catalogin.TorrentToggleLikeInp
}

type TorrentToggleLikeRes struct {
	catalogout.TorrentToggleLikeOut
}

type TorrentLikeListReq struct {
	g.Meta `path:"/torrents/{id}/likes" method:"get" tags:"Catalog" summary:"获取种子点赞列表"`
	catalogin.TorrentLikeListInp
}

type TorrentLikeListRes struct {
	catalogout.TorrentLikeListOut
}

type TorrentUpdateReq struct {
	g.Meta `path:"/torrents/{id}" method:"patch" tags:"Catalog" summary:"局部更新种子信息"`
	catalogin.TorrentUpdateInp
}

type TorrentUpdateRes struct {
	catalogout.TorrentUpdateOut
}

type TorrentFileListReq struct {
	g.Meta `path:"/torrents/{id}/files" method:"get" tags:"Catalog" summary:"获取种子内部文件列表"`
	catalogin.TorrentFileListInp
}

type TorrentFileListRes struct {
	catalogout.TorrentFileListOut
}

type TorrentPeerListReq struct {
	g.Meta `path:"/torrents/{id}/peers" method:"get" tags:"Catalog" summary:"获取种子当前在线同伴列表"`
	catalogin.TorrentPeerListInp
}

type TorrentPeerListRes struct {
	catalogout.TorrentPeerListOut
}

type TorrentReportReq struct {
	g.Meta `path:"/torrents/{id}:report" method:"post" tags:"Catalog" summary:"举报种子" perm:"read:catalog/torrent:*"`
	catalogin.TorrentReportInp
}

type TorrentReportRes struct {
	catalogout.TorrentReportOut
}
