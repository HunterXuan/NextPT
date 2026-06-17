// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserIpHistory is the golang structure of table iam_user_ip_history for DAO operations like Where/Data.
type IamUserIpHistory struct {
	g.Meta    `orm:"table:iam_user_ip_history, do:true"`
	Id        any         //
	UserId    any         //
	Ip        any         //
	Type      any         // login/tracker/register
	CreatedAt *gtime.Time //
}
