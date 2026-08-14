// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ModStaffMessage is the golang structure for table mod_staff_message.
type ModStaffMessage struct {
	Id         uint64      `json:"id"         orm:"id"          description:""`
	SenderId   uint64      `json:"senderId"   orm:"sender_id"   description:""`
	Subject    string      `json:"subject"    orm:"subject"     description:""`
	Content    string      `json:"content"    orm:"content"     description:""`
	Status     int         `json:"status"     orm:"status"      description:"0=pending 1=processed"`
	AnsweredBy uint64      `json:"answeredBy" orm:"answered_by" description:""`
	Answer     string      `json:"answer"     orm:"answer"      description:""`
	AnsweredAt *gtime.Time `json:"answeredAt" orm:"answered_at" description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`
}
