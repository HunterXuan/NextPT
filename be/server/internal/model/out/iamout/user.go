package iamout

import "github.com/gogf/gf/v2/os/gtime"

type UserMeOut struct {
	Id         uint64      `json:"id"`
	Username   string      `json:"username"`
	Email      string      `json:"email"`
	Passkey    string      `json:"passkey"`
	Status     int         `json:"status"`
	Role       uint        `json:"role"`
	RoleName   string      `json:"roleName"`
	RoleLevel  int         `json:"roleLevel"`
	IsStaff    bool        `json:"isStaff"`
	VipUntil   *gtime.Time `json:"vipUntil"`
	Avatar     string      `json:"avatar"`
	Info       string      `json:"info"`
	Signature  string      `json:"signature"`
	Uploaded   uint64      `json:"uploaded"`
	Downloaded uint64      `json:"downloaded"`
	Bonus      float64     `json:"bonus"`
	Invites    int         `json:"invites"`
	ShareRatio float64     `json:"shareRatio"`
	CreatedAt  *gtime.Time `json:"createdAt"`
}
