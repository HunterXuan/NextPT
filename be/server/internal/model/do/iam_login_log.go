// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamLoginLog is the golang structure of table iam_login_log for DAO operations like Where/Data.
type IamLoginLog struct {
	g.Meta     `orm:"table:iam_login_log, do:true"`
	Id         any         //
	UserId     any         //
	Ip         any         //
	UserAgent  any         //
	Result     any         // 1=success 0=fail
	FailReason any         //
	CreatedAt  *gtime.Time //
}
