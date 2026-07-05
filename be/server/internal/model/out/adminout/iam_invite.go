package adminout

import "github.com/gogf/gf/v2/os/gtime"

type IamInviteListOut struct {
	Invites []IamInviteItem `json:"invites"`
	Total   int             `json:"total"`
}

type IamInviteItem struct {
	Id           uint64      `json:"id"`
	InviterId    uint64      `json:"inviterId"`
	InviteeEmail string      `json:"inviteeEmail"`
	InviteeId    uint64      `json:"inviteeId"`
	InviteeName  string      `json:"inviteeName"`
	Hash         string      `json:"hash"`
	Status       int         `json:"status"`
	IsTemporary  bool        `json:"isTemporary"`
	ExpireAt     *gtime.Time `json:"expireAt"`
	UsedAt       *gtime.Time `json:"usedAt"`
	CreatedAt    *gtime.Time `json:"createdAt"`
}
