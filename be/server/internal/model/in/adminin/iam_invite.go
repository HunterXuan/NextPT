package adminin

import "github.com/gogf/gf/v2/os/gtime"

type IamInviteGrantInp struct {
	Amount   int
	IsTemp   bool
	ExpireAt *gtime.Time
}
