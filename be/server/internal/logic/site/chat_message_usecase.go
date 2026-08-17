package site

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
)

const (
	siteChatMessageLimit = 200
	siteChatCacheTTL     = 10 * time.Minute
)

type sSiteChatMessageUsecase struct{}

func init() {
	service.RegisterSiteChatMessageUsecase(NewSiteChatMessageUsecase())
}

func NewSiteChatMessageUsecase() *sSiteChatMessageUsecase {
	return &sSiteChatMessageUsecase{}
}

func (s *sSiteChatMessageUsecase) List(ctx context.Context, actor *model.Actor, in sitein.ChatMessageListInp) (*siteout.ChatMessageListOut, error) {
	items, err := s.loadRecentItemsCache(ctx)
	if err != nil {
		return nil, err
	}
	if in.AfterId == 0 {
		return &siteout.ChatMessageListOut{List: items}, nil
	}
	return &siteout.ChatMessageListOut{List: s.itemsAfter(items, in.AfterId)}, nil
}

func (s *sSiteChatMessageUsecase) Create(ctx context.Context, actor *model.Actor, in sitein.ChatMessageCreateInp) (*siteout.ChatMessageItem, error) {
	if actor == nil || actor.Id == 0 {
		return nil, gerror.New("unauthorized")
	}
	content := strings.TrimSpace(in.Content)
	if utf8.RuneCountInString(content) == 0 || utf8.RuneCountInString(content) > 1000 {
		return nil, gerror.New("invalid chat message content")
	}

	message, err := service.SiteChatMessageDomain().Create(ctx, actor.Id, content)
	if err != nil {
		return nil, err
	}
	items := s.buildItems(ctx, []entity.SiteChatMessage{*message})
	if len(items) == 0 {
		return nil, gerror.New("create chat message failed")
	}

	cacheKey := service.SysCache().KeySiteChatRecentMessages(ctx)
	_, _ = gcache.Remove(ctx, cacheKey)
	_ = service.SysCache().PublishInvalidate(ctx, cacheKey)
	return items[0], nil
}

func (s *sSiteChatMessageUsecase) loadRecentItemsCache(ctx context.Context) ([]*siteout.ChatMessageItem, error) {
	cacheKey := service.SysCache().KeySiteChatRecentMessages(ctx)
	value, err := gcache.GetOrSetFuncLock(ctx, cacheKey, func(ctx context.Context) (any, error) {
		list, err := service.SiteChatMessageDomain().ListRecent(ctx, siteChatMessageLimit)
		if err != nil {
			return nil, err
		}
		return s.buildItems(ctx, list), nil
	}, siteChatCacheTTL)
	if err != nil {
		return nil, err
	}
	if items, ok := value.Val().([]*siteout.ChatMessageItem); ok {
		return items, nil
	}

	list, err := service.SiteChatMessageDomain().ListRecent(ctx, siteChatMessageLimit)
	if err != nil {
		return nil, err
	}
	return s.buildItems(ctx, list), nil
}

func (s *sSiteChatMessageUsecase) buildItems(ctx context.Context, list []entity.SiteChatMessage) []*siteout.ChatMessageItem {
	userMap := s.loadUserSummaryMap(ctx, s.userIds(list))
	items := make([]*siteout.ChatMessageItem, 0, len(list))
	for _, message := range list {
		items = append(items, &siteout.ChatMessageItem{
			Id:        message.Id,
			User:      userMap[message.UserId],
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		})
	}
	return items
}

func (s *sSiteChatMessageUsecase) userIds(list []entity.SiteChatMessage) []uint64 {
	seen := make(map[uint64]struct{}, len(list))
	ids := make([]uint64, 0, len(list))
	for _, message := range list {
		if message.UserId == 0 {
			continue
		}
		if _, ok := seen[message.UserId]; ok {
			continue
		}
		seen[message.UserId] = struct{}{}
		ids = append(ids, message.UserId)
	}
	return ids
}

func (s *sSiteChatMessageUsecase) loadUserSummaryMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
	userMap := make(map[uint64]model.IamUserSummary, len(userIds))
	for _, id := range userIds {
		userMap[id] = model.IamUserSummary{Id: id}
	}
	users, err := service.IamUserDomain().GetUsersByIds(ctx, userIds)
	if err != nil {
		return userMap
	}
	for _, user := range users {
		summary := userMap[user.Id]
		summary.Username = user.Username
		userMap[user.Id] = summary
	}
	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, userIds)
	if err != nil {
		return userMap
	}
	for _, profile := range profiles {
		summary := userMap[profile.UserId]
		summary.Avatar = profile.Avatar
		userMap[profile.UserId] = summary
	}
	return userMap
}

func (s *sSiteChatMessageUsecase) itemsAfter(items []*siteout.ChatMessageItem, afterId uint64) []*siteout.ChatMessageItem {
	start := len(items)
	for index, item := range items {
		if item.Id > afterId {
			start = index
			break
		}
	}
	return items[start:]
}
