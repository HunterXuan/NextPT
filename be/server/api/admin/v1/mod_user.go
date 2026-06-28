package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/modout"

	"github.com/gogf/gf/v2/frame/g"
)

type ModUserListReq struct {
	g.Meta `path:"/mod/users/{id}/mods" method:"get" tags:"Admin Mod" summary:"List user moderation records" perm:"admin:mod/user:*"`
	adminin.ModUserListInp
}

type ModUserListRes struct {
	modout.ListUserOut
}

type ModUserApplyReq struct {
	g.Meta `path:"/mod/users/{id}/mods" method:"post" tags:"Admin Mod" summary:"Apply moderation to a user" perm:"admin:mod/user:*"`
	adminin.ModUserApplyInp
}

type ModUserApplyRes struct{}

type ModUserRemoveReq struct {
	g.Meta `path:"/mod/users/{id}/mods/{modId}" method:"delete" tags:"Admin Mod" summary:"Remove moderation from a user" perm:"admin:mod/user:*"`
	adminin.ModUserRemoveInp
}

type ModUserRemoveRes struct{}
