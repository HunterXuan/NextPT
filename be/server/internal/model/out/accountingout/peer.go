package accountingout

import "github.com/gogf/gf/v2/os/gtime"

type PeerItem struct {
	TorrentId    uint64      `json:"torrentId"`
	TorrentName  string      `json:"torrentName"`
	TorrentSize  uint64      `json:"torrentSize"`
	Uploaded     uint64      `json:"uploaded"`
	Downloaded   uint64      `json:"downloaded"`
	Remaining    uint64      `json:"remaining"`
	IsSeeder     bool        `json:"isSeeder"`
	Agent        string      `json:"agent"`
	StartedAt    *gtime.Time `json:"startedAt"`
	FinishedAt   *gtime.Time `json:"finishedAt"`
	LastActionAt *gtime.Time `json:"lastActionAt"`
}

type PeerListOut struct {
	Page          int        `json:"page"`
	Size          int        `json:"size"`
	Total         int        `json:"total"`
	SeedingTotal  int        `json:"seedingTotal"`
	LeechingTotal int        `json:"leechingTotal"`
	List          []PeerItem `json:"list"`
}
