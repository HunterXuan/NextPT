// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerEventIdempotency is the golang structure of table tracker_event_idempotency for DAO operations like Where/Data.
type TrackerEventIdempotency struct {
	g.Meta    `orm:"table:tracker_event_idempotency, do:true"`
	MsgId     any         //
	CreatedAt *gtime.Time //
}
