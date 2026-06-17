package consts

const (
	IamInviteStatusUnused   = iota // 0=未分配/待发送
	IamInviteStatusSent            // 1=已发送/待注册
	IamInviteStatusUsed            // 2=已注册/已使用
	IamInviteStatusExpired         // 3=已过期
	IamInviteStatusRecycled        // 4=已回收
)
