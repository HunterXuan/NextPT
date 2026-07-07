package modout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type CheaterLogItem struct {
	Id           uint64                      `json:"id"`
	UserId       uint64                      `json:"userId"`
	User         model.IamUserSummary        `json:"user"`
	TorrentId    uint64                      `json:"torrentId"`
	Torrent      model.CatalogTorrentSummary `json:"torrent"`
	Uploaded     uint64                      `json:"uploaded"`
	Downloaded   uint64                      `json:"downloaded"`
	AnnounceTime uint                        `json:"announceTime"`
	Seeders      uint                        `json:"seeders"`
	Leechers     uint                        `json:"leechers"`
	HitCount     uint                        `json:"hitCount"`
	DealtBy      uint64                      `json:"dealtBy"`
	DealtUser    model.IamUserSummary        `json:"dealtUser"`
	IsDealt      bool                        `json:"isDealt"`
	Comment      string                      `json:"comment"`
	DealtComment string                      `json:"dealtComment"`
	DealtAt      *gtime.Time                 `json:"dealtAt"`
	CreatedAt    *gtime.Time                 `json:"createdAt"`
}

type ListCheaterLogsOut struct {
	Page  int              `json:"page"`
	Size  int              `json:"size"`
	Total int              `json:"total"`
	List  []CheaterLogItem `json:"list"`
}
