package trackerin

type AnnounceInp struct {
	InfoHash   string `json:"info_hash" in:"query" v:"required"`
	PeerId     string `json:"peer_id" in:"query" v:"required"`
	Port       int    `json:"port" in:"query" v:"required|between:1,65535"`
	Uploaded   int64  `json:"uploaded" in:"query"`
	Downloaded int64  `json:"downloaded" in:"query"`
	Left       int64  `json:"left" in:"query"`
	Event      string `json:"event" in:"query"`
	Compact    int    `json:"compact" in:"query"`
	NumWant    int    `json:"numwant" in:"query"`
}

type ScrapeInp struct {
	InfoHash []string `json:"info_hash" in:"query"`
}

type DownloadInp struct {
	Id uint64 `json:"id" in:"query" v:"required#种子ID不能为空"`
}
