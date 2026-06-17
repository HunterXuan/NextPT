// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerSnatch is the golang structure for table tracker_snatch.
type TrackerSnatch struct {
	Id          uint64      `json:"id"          orm:"id"           description:""`
	TorrentId   uint64      `json:"torrentId"   orm:"torrent_id"   description:""`
	UserId      uint64      `json:"userId"      orm:"user_id"      description:""`
	Ipv4        string      `json:"ipv4"        orm:"ipv4"         description:""`
	Ipv6        string      `json:"ipv6"        orm:"ipv6"         description:""`
	Port        uint        `json:"port"        orm:"port"         description:""`
	Uploaded    uint64      `json:"uploaded"    orm:"uploaded"     description:""`
	Downloaded  uint64      `json:"downloaded"  orm:"downloaded"   description:""`
	Remaining   uint64      `json:"remaining"   orm:"remaining"    description:""`
	SeedTime    uint        `json:"seedTime"    orm:"seed_time"    description:"做种时间 (秒)"`
	LeechTime   uint        `json:"leechTime"   orm:"leech_time"   description:"下载时间 (秒)"`
	IsFinished  bool        `json:"isFinished"  orm:"is_finished"  description:""`
	StartedAt   *gtime.Time `json:"startedAt"   orm:"started_at"   description:""`
	CompletedAt *gtime.Time `json:"completedAt" orm:"completed_at" description:""`
	LastAction  *gtime.Time `json:"lastAction"  orm:"last_action"  description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
}
