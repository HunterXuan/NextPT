package v1

import (
	"server/internal/model/in/modin"
	"server/internal/model/out/modout"

	"github.com/gogf/gf/v2/frame/g"
)

type ModStaffMessageListReq struct {
	g.Meta `path:"/mod/staff-messages" method:"get" tags:"Admin Mod" summary:"查看管理组信箱" perm:"admin:mod/staff-message:*"`
	modin.AdminStaffMessageListInp
}

type ModStaffMessageListRes struct {
	modout.StaffMessageListOut
}

type ModStaffMessageUpdateReq struct {
	g.Meta `path:"/mod/staff-messages/{id}" method:"patch" tags:"Admin Mod" summary:"处理管理组消息" perm:"admin:mod/staff-message:*"`
	modin.StaffMessageProcessInp
}

type ModStaffMessageUpdateRes struct{}
