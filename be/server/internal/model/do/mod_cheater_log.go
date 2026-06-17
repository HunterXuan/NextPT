// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ModCheaterLog is the golang structure of table mod_cheater_log for DAO operations like Where/Data.
type ModCheaterLog struct {
	g.Meta       `orm:"table:mod_cheater_log, do:true"`
	Id           any         //
	UserId       any         //
	TorrentId    any         //
	Uploaded     any         //
	Downloaded   any         //
	AnnounceTime any         //
	Seeders      any         //
	Leechers     any         //
	HitCount     any         //
	DealtBy      any         //
	IsDealt      any         //
	Comment      any         //
	CreatedAt    *gtime.Time //
}
