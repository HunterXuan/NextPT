package modin

type RecordCheaterLogInp struct {
	UserId       uint64 `json:"userId"`
	TorrentId    uint64 `json:"torrentId"`
	Uploaded     uint64 `json:"uploaded"`
	Downloaded   uint64 `json:"downloaded"`
	AnnounceTime uint   `json:"announceTime"`
	Seeders      uint   `json:"seeders"`
	Leechers     uint   `json:"leechers"`
	HitCount     uint   `json:"hitCount"`
	Comment      string `json:"comment"`
}

type ListCheaterLogsInp struct {
	Page    int  `json:"page" d:"1" v:"min:1"`
	Size    int  `json:"size" d:"20" v:"min:1|max:100"`
	IsDealt *int `json:"isDealt" d:"-1" description:"-1=All, 0=Unresolved, 1=Resolved"`
}

type ResolveCheaterLogInp struct {
	Id      uint64 `json:"id" v:"required"`
	Comment string `json:"comment"`
}
