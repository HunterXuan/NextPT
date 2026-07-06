package modout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type CheaterLogItem struct {
	Id           uint64                      `json:"id"`
	UserId       uint64                      `json:"user_id"`
	User         model.IamUserSummary        `json:"user"`
	TorrentId    uint64                      `json:"torrent_id"`
	Torrent      model.CatalogTorrentSummary `json:"torrent"`
	Uploaded     uint64                      `json:"uploaded"`
	Downloaded   uint64                      `json:"downloaded"`
	AnnounceTime uint                        `json:"announce_time"`
	Seeders      uint                        `json:"seeders"`
	Leechers     uint                        `json:"leechers"`
	HitCount     uint                        `json:"hit_count"`
	DealtBy      uint64                      `json:"dealt_by"`
	DealtUser    model.IamUserSummary        `json:"dealt_user"`
	IsDealt      bool                        `json:"is_dealt"`
	Comment      string                      `json:"comment"`
	DealtComment string                      `json:"dealt_comment"`
	DealtAt      *gtime.Time                 `json:"dealt_at"`
	CreatedAt    *gtime.Time                 `json:"created_at"`
}

type ListCheaterLogsOut struct {
	Page  int              `json:"page"`
	Size  int              `json:"size"`
	Total int              `json:"total"`
	List  []CheaterLogItem `json:"list"`
}
