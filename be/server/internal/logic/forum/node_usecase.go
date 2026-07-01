package forum

import (
	"context"
	"strings"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/forumin"
	"server/internal/model/out/forumout"
	"server/internal/service"
)

type sForumNodeUsecase struct{}

func init() {
	service.RegisterForumNodeUsecase(NewForumNodeUsecase())
}

func NewForumNodeUsecase() *sForumNodeUsecase {
	return &sForumNodeUsecase{}
}

func (s *sForumNodeUsecase) List(ctx context.Context, actor *model.Actor, in forumin.NodeListInp) (*forumout.NodeListOut, error) {
	userLevel := 0
	if actor != nil {
		userLevel = actor.RoleLevel
	}
	scope := s.normalizeListScope(in.Scope)

	categories, err := service.ForumNodeDomain().GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := service.ForumNodeDomain().GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[uint][]forumout.NodeItem)
	for _, node := range nodes {
		if !s.canUseNodeForScope(userLevel, &node, scope) {
			continue
		}
		item := forumout.NodeItem{
			Id:         node.Id,
			Slug:       node.Slug,
			NameI18N:   node.NameI18N,
			DescI18N:   node.DescI18N,
			TopicCount: node.TopicCount,
			ReplyCount: node.ReplyCount,
		}
		nodeMap[node.CategoryId] = append(nodeMap[node.CategoryId], item)
	}

	var list []forumout.NodeCategoryItem
	for _, cat := range categories {
		if userLevel < int(cat.MinRoleView) {
			continue
		}

		catNodes := nodeMap[cat.Id]
		if catNodes == nil {
			catNodes = []forumout.NodeItem{}
		}

		list = append(list, forumout.NodeCategoryItem{
			Id:       cat.Id,
			NameI18N: cat.NameI18N,
			DescI18N: cat.DescI18N,
			Nodes:    catNodes,
		})
	}

	return &forumout.NodeListOut{
		List: list,
	}, nil
}

func (s *sForumNodeUsecase) normalizeListScope(scope string) string {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case consts.ForumNodeListScopeCreate:
		return consts.ForumNodeListScopeCreate
	default:
		return consts.ForumNodeListScopeRead
	}
}

func (s *sForumNodeUsecase) canUseNodeForScope(userLevel int, node *entity.ForumNode, scope string) bool {
	switch scope {
	case consts.ForumNodeListScopeCreate:
		return userLevel >= int(node.MinRoleCreate)
	default:
		return userLevel >= int(node.MinRoleRead)
	}
}
