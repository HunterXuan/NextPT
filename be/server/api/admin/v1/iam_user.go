package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type IamUserListReq struct {
	g.Meta `path:"/iam/users" method:"get" tags:"AdminIAM" summary:"获取用户列表" perm:"admin:iam/user:*"`
	adminin.IamUserListInp
}

type IamUserListRes struct {
	adminout.IamUserListOut
}

type IamUserUpdateReq struct {
	g.Meta `path:"/iam/users/{id}" method:"patch" tags:"AdminIAM" summary:"部分更新用户信息" perm:"admin:iam/user:*"`
	adminin.IamUserUpdateInp
}

type IamUserUpdateRes struct{}

type IamUserStatDetailReq struct {
	g.Meta `path:"/iam/users/{id}/stat" method:"get" tags:"AdminIAM" summary:"获取用户统计详情" perm:"admin:iam/user:*"`
	adminin.IamUserStatDetailInp
}

type IamUserStatDetailRes struct {
	adminout.IamUserStatDetailOut
}

type IamUserStatUpdateReq struct {
	g.Meta `path:"/iam/users/{id}:incrementStat" method:"post" tags:"AdminIAM" summary:"增量更新用户统计数据" perm:"admin:iam/user:*"`
	adminin.IamUserStatUpdateInp
}

type IamUserStatUpdateRes struct{}
