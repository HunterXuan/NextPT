package sitein

import "github.com/gogf/gf/v2/os/gtime"

type AuditRecordInp struct {
	UserId     uint64
	Action     string
	TargetType string
	TargetId   uint64
	Detail     any
	Level      int
}

type AuditCreateInp struct {
	UserId     uint64
	Action     string
	TargetType string
	TargetId   uint64
	Detail     string
	Ip         string
	Level      int
	CreatedAt  *gtime.Time
}
