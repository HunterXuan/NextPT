package accountingout

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type SnatchItem struct {
	Id           uint64      `json:"id"`
	TorrentId    uint64      `json:"torrentId"`
	TorrentName  string      `json:"torrentName"`
	TorrentSize  uint64      `json:"torrentSize"`
	Uploaded     uint64      `json:"uploaded"`
	Downloaded   uint64      `json:"downloaded"`
	SeedTime     uint64      `json:"seedTime"`
	LeechTime    uint64      `json:"leechTime"`
	IsFinished   bool        `json:"isFinished"`
	StartedAt    *gtime.Time `json:"startedAt"`
	FinishedAt   *gtime.Time `json:"finishedAt"`
	LastActionAt *gtime.Time `json:"lastActionAt"`
}

type SnatchListOut struct {
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Total int          `json:"total"`
	List  []SnatchItem `json:"list"`
}

type SnatchGetOut struct {
	SnatchItem
}
