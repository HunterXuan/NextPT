// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerSnatch is the golang structure of table tracker_snatch for DAO operations like Where/Data.
type TrackerSnatch struct {
	g.Meta      `orm:"table:tracker_snatch, do:true"`
	Id          any         //
	TorrentId   any         //
	UserId      any         //
	Ipv4        any         //
	Ipv6        any         //
	Port        any         //
	Uploaded    any         //
	Downloaded  any         //
	Remaining   any         //
	SeedTime    any         // 做种时间 (秒)
	LeechTime   any         // 下载时间 (秒)
	IsFinished  any         //
	StartedAt   *gtime.Time //
	CompletedAt *gtime.Time //
	LastAction  *gtime.Time //
	CreatedAt   *gtime.Time //
}
