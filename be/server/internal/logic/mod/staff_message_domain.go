package mod

import (
	"context"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/modin"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

const modStaffMessageListMaxSize = 100

type sModStaffMessageDomain struct{}

func init() {
	service.RegisterModStaffMessageDomain(NewModStaffMessageDomain())
}

func NewModStaffMessageDomain() *sModStaffMessageDomain {
	return &sModStaffMessageDomain{}
}

func (s *sModStaffMessageDomain) Create(ctx context.Context, message entity.ModStaffMessage) (uint64, error) {
	now := gtime.Now()
	id, err := dao.ModStaffMessage.Ctx(ctx).Data(do.ModStaffMessage{
		SenderId:  message.SenderId,
		Subject:   message.Subject,
		Content:   message.Content,
		Status:    consts.ModStaffMessageStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}).InsertAndGetId()
	return uint64(id), err
}

func (s *sModStaffMessageDomain) GetById(ctx context.Context, id uint64) (*entity.ModStaffMessage, error) {
	var message *entity.ModStaffMessage
	err := dao.ModStaffMessage.Ctx(ctx).WherePri(id).Scan(&message)
	return message, err
}

func (s *sModStaffMessageDomain) ListBySender(ctx context.Context, senderId uint64, in modin.StaffMessageListInp) ([]entity.ModStaffMessage, int, error) {
	page, size := s.normalizePage(in.Page, in.Size)
	columns := dao.ModStaffMessage.Columns()
	m := dao.ModStaffMessage.Ctx(ctx).Where(columns.SenderId, senderId)
	if in.Status != nil {
		m = m.Where(columns.Status, *in.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.ModStaffMessage
	err = m.Page(page, size).OrderDesc(columns.Id).Scan(&list)
	return list, total, err
}

func (s *sModStaffMessageDomain) AdminList(ctx context.Context, in modin.AdminStaffMessageListInp) ([]entity.ModStaffMessage, int, error) {
	page, size := s.normalizePage(in.Page, in.Size)
	columns := dao.ModStaffMessage.Columns()
	m := dao.ModStaffMessage.Ctx(ctx)
	if in.Status != nil {
		m = m.Where(columns.Status, *in.Status)
	}
	if in.SenderId > 0 {
		m = m.Where(columns.SenderId, in.SenderId)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.ModStaffMessage
	err = m.Page(page, size).OrderAsc(columns.Status).OrderDesc(columns.Id).Scan(&list)
	return list, total, err
}

func (s *sModStaffMessageDomain) Update(ctx context.Context, id uint64, data interface{}) error {
	_, err := dao.ModStaffMessage.Ctx(ctx).WherePri(id).Data(data).Update()
	return err
}

func (s *sModStaffMessageDomain) normalizePage(page int, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > modStaffMessageListMaxSize {
		size = modStaffMessageListMaxSize
	}
	return page, size
}
