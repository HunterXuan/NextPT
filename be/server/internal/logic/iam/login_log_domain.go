package iam

import (
	"context"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sIamLoginLogDomain struct{}

func init() {
	service.RegisterIamLoginLogDomain(NewIamLoginLogDomain())
}

func NewIamLoginLogDomain() *sIamLoginLogDomain {
	return &sIamLoginLogDomain{}
}

func (s *sIamLoginLogDomain) Create(ctx context.Context, data do.IamLoginLog) error {
	columns := dao.IamLoginLog.Columns()
	if data.CreatedAt == nil {
		data.CreatedAt = gtime.Now()
	}
	data.Ip = s.limitLoginLogString(data.Ip, 64)
	data.UserAgent = s.limitLoginLogString(data.UserAgent, 500)
	data.FailReason = s.limitLoginLogString(data.FailReason, 100)
	_, err := dao.IamLoginLog.Ctx(ctx).Data(g.Map{
		columns.UserId:     data.UserId,
		columns.Ip:         data.Ip,
		columns.UserAgent:  data.UserAgent,
		columns.Result:     data.Result,
		columns.FailReason: data.FailReason,
		columns.CreatedAt:  data.CreatedAt,
	}).Insert()
	return err
}

func (s *sIamLoginLogDomain) limitLoginLogString(value any, max int) any {
	text, ok := value.(string)
	if !ok || max <= 0 {
		return value
	}
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return string(runes[:max])
}
