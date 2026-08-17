package site

import (
	"context"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

type sSiteChatMessageDomain struct{}

func init() {
	service.RegisterSiteChatMessageDomain(NewSiteChatMessageDomain())
}

func NewSiteChatMessageDomain() *sSiteChatMessageDomain {
	return &sSiteChatMessageDomain{}
}

func (s *sSiteChatMessageDomain) ListRecent(ctx context.Context, size int) ([]entity.SiteChatMessage, error) {
	columns := dao.SiteChatMessage.Columns()
	var list []entity.SiteChatMessage
	if err := dao.SiteChatMessage.Ctx(ctx).OrderDesc(columns.Id).Limit(size).Scan(&list); err != nil {
		return nil, err
	}
	s.reverse(list)
	return list, nil
}

func (s *sSiteChatMessageDomain) Create(ctx context.Context, userId uint64, content string) (*entity.SiteChatMessage, error) {
	createdAt := gtime.Now()
	id, err := dao.SiteChatMessage.Ctx(ctx).Data(do.SiteChatMessage{
		UserId:    userId,
		Content:   content,
		CreatedAt: createdAt,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &entity.SiteChatMessage{
		Id:        uint64(id),
		UserId:    userId,
		Content:   content,
		CreatedAt: createdAt,
	}, nil
}

func (s *sSiteChatMessageDomain) reverse(list []entity.SiteChatMessage) {
	for left, right := 0, len(list)-1; left < right; left, right = left+1, right-1 {
		list[left], list[right] = list[right], list[left]
	}
}
