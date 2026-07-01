package catalogout

import "server/internal/model"

type SubtitleListItem struct {
	Id        uint64               `json:"id"`
	TorrentId uint64               `json:"torrentId"`
	Uploader  model.IamUserSummary `json:"uploader"`
	Anonymous bool                 `json:"anonymous"`
	FileName  string               `json:"fileName"`
	Language  string               `json:"language"`
	Size      uint64               `json:"size"`
	CreatedAt string               `json:"createdAt"`
}

type SubtitleListOut struct {
	List  []SubtitleListItem `json:"list"`
	Total int                `json:"total"`
}

type SubtitleDownloadOut struct {
	Bytes    []byte `json:"-"`
	FileName string `json:"-"`
	MimeType string `json:"-"`
}
