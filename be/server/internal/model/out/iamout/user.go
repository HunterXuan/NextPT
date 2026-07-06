package iamout

import "github.com/gogf/gf/v2/os/gtime"

type UserMeOut struct {
	User    UserMeAccountOut `json:"user"`
	Role    UserMeRoleOut    `json:"role"`
	Profile UserMeProfileOut `json:"profile"`
	Stat    UserMeStatOut    `json:"stat"`
}

type UserMeAccountOut struct {
	Id        uint64      `json:"id"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	Passkey   string      `json:"passkey"`
	Status    int         `json:"status"`
	VipUntil  *gtime.Time `json:"vipUntil"`
	CreatedAt *gtime.Time `json:"createdAt"`
}

type UserMeRoleOut struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	IsStaff bool   `json:"isStaff"`
}

type UserMeProfileOut struct {
	Avatar    string `json:"avatar"`
	Info      string `json:"info"`
	Signature string `json:"signature"`
}

type UserMeStatOut struct {
	Uploaded   uint64  `json:"uploaded"`
	Downloaded uint64  `json:"downloaded"`
	Bonus      float64 `json:"bonus"`
	ShareRatio float64 `json:"shareRatio"`
}

type UserPermissionListOut struct {
	Permissions []string `json:"permissions"`
}
