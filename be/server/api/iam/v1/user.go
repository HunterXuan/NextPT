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

type IamUserGetReq struct {
	g.Meta `path:"/users/{id}" method:"get" tags:"IamUser" summary:"获取用户公开资料"`
	iamin.UserGetInp
}

type IamUserGetRes struct {
	iamout.UserGetOut
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

type IamUserTwoStepSetupReq struct {
	g.Meta `path:"/users/me/two-step:setup" method:"post" tags:"IamUser" summary:"开始绑定认证器"`
	iamin.UserTwoStepSetupInp
}

type IamUserTwoStepSetupRes struct {
	iamout.UserTwoStepSetupOut
}

type IamUserTwoStepConfirmReq struct {
	g.Meta `path:"/users/me/two-step:confirm" method:"post" tags:"IamUser" summary:"确认绑定认证器"`
	iamin.UserTwoStepConfirmInp
}

type IamUserTwoStepConfirmRes struct {
	iamout.UserTwoStepRecoveryCodesOut
}

type IamUserTwoStepRecoveryCodesCreateReq struct {
	g.Meta `path:"/users/me/two-step:recoveryCodes" method:"post" tags:"IamUser" summary:"重新生成两步验证恢复码"`
	iamin.UserTwoStepRecoveryCodesCreateInp
}

type IamUserTwoStepRecoveryCodesCreateRes struct {
	iamout.UserTwoStepRecoveryCodesOut
}

type IamUserTwoStepDeleteReq struct {
	g.Meta `path:"/users/me/two-step" method:"delete" tags:"IamUser" summary:"关闭两步验证"`
	iamin.UserTwoStepDisableInp
}

type IamUserTwoStepDeleteRes struct{}

type IamUserPasskeyResetReq struct {
	g.Meta `path:"/users/me:resetPasskey" method:"post" tags:"IamUser" summary:"重置当前用户 Passkey"`
}

type IamUserPasskeyResetRes struct {
	Passkey string `json:"passkey"`
}
