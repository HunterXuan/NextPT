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

type AuditListInp struct {
	Page       int
	Size       int
	Level      *int
	Action     string
	TargetType string
	UserId     uint64
	StartAt    *gtime.Time
	EndAt      *gtime.Time
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
