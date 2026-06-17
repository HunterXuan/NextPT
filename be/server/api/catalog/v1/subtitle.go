package v1

import (
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"

	"github.com/gogf/gf/v2/frame/g"
)

type SubtitleListReq struct {
	g.Meta `path:"/subtitles" method:"get" tags:"Catalog Subtitle" summary:"获取全局字幕列表"`
	catalogin.SubtitleListInp
}

type SubtitleListRes struct {
	catalogout.SubtitleListOut
}

type TorrentSubtitleListReq struct {
	g.Meta `path:"/torrents/{id}/subtitles" method:"get" tags:"Catalog Subtitle" summary:"获取种子的字幕列表"`
	catalogin.TorrentSubtitleListInp
}

type TorrentSubtitleListRes struct {
	catalogout.SubtitleListOut
}

type SubtitleUploadReq struct {
	g.Meta `path:"/torrents/{id}/subtitles" method:"post" tags:"Catalog Subtitle" summary:"上传字幕" mime:"multipart/form-data" perm:"create:catalog/subtitle:*"`
	catalogin.SubtitleUploadInp
}

type SubtitleUploadRes struct {
	Id uint64 `json:"id"`
}

type SubtitleDownloadReq struct {
	g.Meta `path:"/subtitles/{id}:download" method:"get" tags:"Catalog Subtitle" summary:"下载字幕" perm:"download:catalog/subtitle:*"`
	catalogin.SubtitleDownloadInp
}

type SubtitleDownloadRes struct {
	catalogout.SubtitleDownloadOut
}

type SubtitleUpdateReq struct {
	g.Meta `path:"/subtitles/{id}" method:"patch" tags:"Catalog Subtitle" summary:"更新字幕"`
	catalogin.SubtitleUpdateInp
}

type SubtitleUpdateRes struct{}

type SubtitleReportReq struct {
	g.Meta `path:"/subtitles/{id}:report" method:"post" tags:"Catalog Subtitle" summary:"举报字幕"`
	catalogin.SubtitleReportInp
}

type SubtitleReportRes struct{}
