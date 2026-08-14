// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ModStaffMessage is the golang structure of table mod_staff_message for DAO operations like Where/Data.
type ModStaffMessage struct {
	g.Meta     `orm:"table:mod_staff_message, do:true"`
	Id         any         //
	SenderId   any         //
	Subject    any         //
	Content    any         //
	Status     any         // 0=pending 1=processed
	AnsweredBy any         //
	Answer     any         //
	AnsweredAt *gtime.Time //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
