package adminout

import (
	"server/internal/model"
	"server/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

type IamUserListOut struct {
	Users []IamUserItem `json:"users"`
	Total int           `json:"total"`
}

type IamUserItem struct {
	Id        uint64      `json:"id"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	Passkey   string      `json:"passkey"`
	Status    int         `json:"status"`
	Role      uint        `json:"role"`
	VipUntil  *gtime.Time `json:"vipUntil"`
	VipRemark string      `json:"vipRemark"`
	InvitedBy uint64      `json:"invitedBy"`
	LastLogin *gtime.Time `json:"lastLogin"`
	LastIp    string      `json:"lastIp"`
	Avatar    string      `json:"avatar"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

type IamUserStatDetailOut struct {
	entity.IamUserStat
}

type IamLoginLogListOut struct {
	List  []model.IamLoginLogItem `json:"list"`
	Total int                     `json:"total"`
	Page  int                     `json:"page"`
	Size  int                     `json:"size"`
}
