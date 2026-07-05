package adminin

import "github.com/gogf/gf/v2/os/gtime"

type IamInviteGrantInp struct {
	Amount     int         `json:"amount" v:"required|min:1#{#iam.invite.count_req}|{#iam.invite.count_min}"`
	TargetMode string      `json:"targetMode" v:"required|in:site,roles#{#admin.invite.target_required}|{#admin.invite.invalid_target}"`
	IsTemp     bool        `json:"isTemp"`
	ExpireAt   *gtime.Time `json:"expireAt"`
	RoleIds    []uint      `json:"roleIds"`
}

type IamInviteListInp struct {
	Page   int   `json:"page" d:"1" v:"min:1#{#iam.invite.page_min}"`
	Size   int   `json:"size" d:"10" v:"max:100#{#iam.invite.size_max}"`
	Status *uint `json:"status" v:"max:4#{#iam.invite.status_max}"`
}

type IamInviteRecycleInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}
