package model

type ResponseType string

const (
	ResponseTypeJSON    ResponseType = "json"
	ResponseTypeTracker ResponseType = "tracker"
)

// Context 请求上下文结构
type Context struct {
	Actor        *Actor       // 上下文操作者信息
	ResponseType ResponseType // 响应格式
}

// Actor 通用操作者模型
type Actor struct {
	Id        uint64 `json:"id"              description:"用户ID"`
	Email     string `json:"email"           description:"用户邮箱"`
	RoleId    uint   `json:"roleId"          description:"用户角色ID"`
	RoleLevel int    `json:"roleLevel"       description:"用户等级数字"`
	IsStaff   bool   `json:"isStaff"         description:"是否为管理员"`

	RoleVersion int64 `json:"-"`
}
