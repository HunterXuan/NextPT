package iam

import (
	"context"

	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sIamLoginLogDomain struct{}

const iamLoginLogListMaxSize = 100

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

func (s *sIamLoginLogDomain) List(ctx context.Context, options model.IamLoginLogListOptions) ([]entity.IamLoginLog, int, error) {
	page, size := s.normalizeListPage(options.Page, options.Size)
	columns := dao.IamLoginLog.Columns()
	m := dao.IamLoginLog.Ctx(ctx)
	if options.UserId > 0 {
		m = m.Where(columns.UserId, options.UserId)
	}
	if options.Result != nil {
		m = m.Where(columns.Result, *options.Result)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var logs []entity.IamLoginLog
	if total > 0 {
		err = m.Page(page, size).OrderDesc(columns.Id).Scan(&logs)
	}
	return logs, total, err
}

func (s *sIamLoginLogDomain) normalizeListPage(page int, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > iamLoginLogListMaxSize {
		size = iamLoginLogListMaxSize
	}
	return page, size
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
