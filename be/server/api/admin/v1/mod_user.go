package v1

import (
	"server/internal/model/in/adminin"

	"github.com/gogf/gf/v2/frame/g"
)

type ModUserApplyReq struct {
	g.Meta `path:"/users/{id}/applyMod" method:"post" tags:"Admin Mod" summary:"Apply moderation to a user" perm:"admin:mod/user:*"`
	adminin.ModUserApplyInp
}

type ModUserApplyRes struct{}

type ModUserRemoveReq struct {
	g.Meta `path:"/users/{id}/removeMod" method:"post" tags:"Admin Mod" summary:"Remove moderation from a user" perm:"admin:mod/user:*"`
	adminin.ModUserRemoveInp
}

type ModUserRemoveRes struct{}
