package catalogout

type SubtitleListItem struct {
	Id        uint64 `json:"id"`
	TorrentId uint64 `json:"torrentId"`
	UserId    uint64 `json:"userId"`
	Username  string `json:"username"`
	FileName  string `json:"fileName"`
	Language  string `json:"language"`
	Size      uint64 `json:"size"`
	CreatedAt string `json:"createdAt"`
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
