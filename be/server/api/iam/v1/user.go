package v1

import (
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"

	"github.com/gogf/gf/v2/frame/g"
)

type IamUserCreateReq struct {
	g.Meta `path:"/users" method:"post" tags:"IamUser" summary:"创建用户 (注册)" noAuth:"true"`
	iamin.UserCreateInp
}

type IamUserCreateRes struct {
	Id uint64 `json:"id"`
}

type IamUserMeReq struct {
	g.Meta `path:"/users/me" method:"get" tags:"IamUser" summary:"获取当前登录用户基础信息"`
}

type IamUserMeRes struct {
	iamout.UserMeOut
}

type IamUserPermissionListReq struct {
	g.Meta `path:"/users/me/permissions" method:"get" tags:"IamUser" summary:"获取当前用户通配权限列表"`
}

type IamUserPermissionListRes struct {
	iamout.UserPermissionListOut
}

type IamUserLoginLogListReq struct {
	g.Meta `path:"/users/me/login-logs" method:"get" tags:"IamUser" summary:"获取当前用户登录记录"`
	iamin.UserLoginLogListInp
}

type IamUserLoginLogListRes struct {
	iamout.UserLoginLogListOut
}

type IamUserProfileUpdateReq struct {
	g.Meta `path:"/users/me" method:"patch" tags:"IamUser" summary:"更新当前用户资料"`
	iamin.UserProfileUpdateInp
}

type IamUserProfileUpdateRes struct{}

type IamUserPasswordChangeReq struct {
	g.Meta `path:"/users/me:changePassword" method:"post" tags:"IamUser" summary:"修改当前用户密码"`
	iamin.UserPasswordChangeInp
}

type IamUserPasswordChangeRes struct{}

type IamUserPasskeyResetReq struct {
	g.Meta `path:"/users/me:resetPasskey" method:"post" tags:"IamUser" summary:"重置当前用户 Passkey"`
}

type IamUserPasskeyResetRes struct {
	Passkey string `json:"passkey"`
}
