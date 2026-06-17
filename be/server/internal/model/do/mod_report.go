// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ModReport is the golang structure of table mod_report for DAO operations like Where/Data.
type ModReport struct {
	g.Meta       `orm:"table:mod_report, do:true"`
	Id           any         //
	ReporterId   any         //
	TargetType   any         // torrent/user/offer/request/forum_post/comment/subtitle
	TargetId     any         //
	Reason       any         //
	Status       any         // 0=pending 1=resolved 2=rejected
	DealtBy      any         //
	DealtComment any         //
	DealtAt      *gtime.Time //
	CreatedAt    *gtime.Time //
}
