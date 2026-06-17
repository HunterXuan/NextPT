package catalogin

import "github.com/gogf/gf/v2/net/ghttp"

type SubtitleListInp struct {
	Page int `json:"page" d:"1" v:"min:1" description:"页码"`
	Size int `json:"size" d:"20" v:"max:100" description:"每页数量"`
}

type TorrentSubtitleListInp struct {
	Id   uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Page int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size int    `json:"size" d:"20" v:"max:100" description:"每页数量"`
}

type SubtitleUploadInp struct {
	Id       uint64            `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	File     *ghttp.UploadFile `json:"file" type:"file" v:"required#{#catalog.subtitle.file_req}" description:"字幕文件"`
	Language string            `json:"language" v:"required#{#catalog.subtitle.lang_req}" description:"语言代码"`
}

type SubtitleDownloadInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.subtitle.id_req}" description:"字幕ID"`
}

type SubtitleUpdateInp struct {
	Id       uint64 `json:"id" in:"path" v:"required#{#catalog.subtitle.id_req}" description:"字幕ID"`
	Language string `json:"language" description:"语言代码"`
}

type SubtitleReportInp struct {
	Id     uint64 `json:"id" in:"path" v:"required#{#catalog.subtitle.id_req}" description:"字幕ID"`
	Reason string `json:"reason" v:"required|length:5,500#{#catalog.comment.report_reason_req}|{#catalog.comment.report_reason_len}" description:"举报原因"`
}
