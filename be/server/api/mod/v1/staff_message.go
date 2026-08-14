package v1

import (
	"server/internal/model/in/modin"
	"server/internal/model/out/modout"

	"github.com/gogf/gf/v2/frame/g"
)

type StaffMessageListReq struct {
	g.Meta `path:"/staff-messages" method:"get" tags:"Mod" summary:"查看发给管理组的消息" perm:"read:mod/staff-message:*"`
	modin.StaffMessageListInp
}

type StaffMessageListRes struct {
	modout.StaffMessageListOut
}

type StaffMessageCreateReq struct {
	g.Meta `path:"/staff-messages" method:"post" tags:"Mod" summary:"给管理组发消息" perm:"create:mod/staff-message:*"`
	modin.StaffMessageCreateInp
}

type StaffMessageCreateRes struct {
	modout.StaffMessageCreateOut
}
