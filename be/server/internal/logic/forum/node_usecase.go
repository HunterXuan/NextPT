package forum

import (
	"context"

	"server/internal/model"
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

func (s *sForumNodeUsecase) List(ctx context.Context, actor *model.Actor) (*forumout.NodeListOut, error) {
	userLevel := 0
	if actor != nil {
		userLevel = actor.RoleLevel
	}

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
		if userLevel < int(node.MinRoleRead) {
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
