package trackerin

import "github.com/gogf/gf/v2/os/gtime"

type AnnounceEvent struct {
	Event      string      `json:"event"`
	TorrentId  uint64      `json:"torrentId"`
	Uploaded   int64       `json:"uploaded"`
	Downloaded int64       `json:"downloaded"`
	Left       int64       `json:"left"`
	UserId     uint64      `json:"userId"`
	PeerId     string      `json:"peerId"`
	Ipv4       string      `json:"ipv4"`
	Ipv6       string      `json:"ipv6"`
	Port       int         `json:"port"`
	IsSeeder   bool        `json:"isSeeder"`
	Now        *gtime.Time `json:"now"`
	UserAgent  string      `json:"userAgent"`
}
