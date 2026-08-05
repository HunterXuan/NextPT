package iamout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type UserMeOut struct {
	User    UserMeAccountOut `json:"user"`
	Role    UserMeRoleOut    `json:"role"`
	Profile UserMeProfileOut `json:"profile"`
	Stat    UserMeStatOut    `json:"stat"`
}

type UserGetOut struct {
	User    UserGetAccountOut `json:"user"`
	Role    UserGetRoleOut    `json:"role"`
	Profile UserGetProfileOut `json:"profile"`
	Stat    UserGetStatOut    `json:"stat"`
}

type UserGetAccountOut struct {
	Id        uint64      `json:"id"`
	Username  string      `json:"username"`
	CreatedAt *gtime.Time `json:"createdAt"`
}

type UserGetRoleOut struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	IsStaff bool   `json:"isStaff"`
}

type UserGetProfileOut struct {
	Avatar    string `json:"avatar"`
	Info      string `json:"info"`
	Signature string `json:"signature"`
}

type UserGetStatOut struct {
	Uploaded   uint64  `json:"uploaded"`
	Downloaded uint64  `json:"downloaded"`
	ShareRatio float64 `json:"shareRatio"`
	SeedTime   uint64  `json:"seedTime"`
}

type UserMeAccountOut struct {
	Id             uint64      `json:"id"`
	Username       string      `json:"username"`
	Email          string      `json:"email"`
	Passkey        string      `json:"passkey"`
	Status         int         `json:"status"`
	TwoStepEnabled bool        `json:"twoStepEnabled"`
	VipUntil       *gtime.Time `json:"vipUntil"`
	CreatedAt      *gtime.Time `json:"createdAt"`
}

type UserTwoStepSetupOut struct {
	Challenge     string `json:"challenge"`
	QRCodeDataURL string `json:"qrCodeDataUrl"`
	Secret        string `json:"secret"`
}

type UserTwoStepRecoveryCodesOut struct {
	RecoveryCodes []string `json:"recoveryCodes"`
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
	Uploaded      uint64  `json:"uploaded"`
	Downloaded    uint64  `json:"downloaded"`
	RawUploaded   uint64  `json:"rawUploaded"`
	RawDownloaded uint64  `json:"rawDownloaded"`
	Bonus         float64 `json:"bonus"`
	ShareRatio    float64 `json:"shareRatio"`
}

type UserPermissionListOut struct {
	Permissions []string `json:"permissions"`
}

type UserLoginLogListOut struct {
	List  []model.IamLoginLogItem `json:"list"`
	Total int                     `json:"total"`
	Page  int                     `json:"page"`
	Size  int                     `json:"size"`
}
