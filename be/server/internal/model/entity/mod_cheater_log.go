// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ModCheaterLog is the golang structure for table mod_cheater_log.
type ModCheaterLog struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`
	UserId       uint64      `json:"userId"       orm:"user_id"       description:""`
	TorrentId    uint64      `json:"torrentId"    orm:"torrent_id"    description:""`
	Uploaded     uint64      `json:"uploaded"     orm:"uploaded"      description:""`
	Downloaded   uint64      `json:"downloaded"   orm:"downloaded"    description:""`
	AnnounceTime uint        `json:"announceTime" orm:"announce_time" description:""`
	Seeders      uint        `json:"seeders"      orm:"seeders"       description:""`
	Leechers     uint        `json:"leechers"     orm:"leechers"      description:""`
	HitCount     uint        `json:"hitCount"     orm:"hit_count"     description:""`
	DealtBy      uint64      `json:"dealtBy"      orm:"dealt_by"      description:""`
	IsDealt      bool        `json:"isDealt"      orm:"is_dealt"      description:""`
	Comment      string      `json:"comment"      orm:"comment"       description:"检测备注"`
	DealtComment string      `json:"dealtComment" orm:"dealt_comment" description:"处理说明"`
	DealtAt      *gtime.Time `json:"dealtAt"      orm:"dealt_at"      description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
}
