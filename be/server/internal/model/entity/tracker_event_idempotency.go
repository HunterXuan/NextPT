// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerEventIdempotency is the golang structure for table tracker_event_idempotency.
type TrackerEventIdempotency struct {
	MsgId     string      `json:"msgId"     orm:"msg_id"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
