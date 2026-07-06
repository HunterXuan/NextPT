package forum

import (
	"context"
	"fmt"

	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sForumNodeDomain struct{}

func init() {
	service.RegisterForumNodeDomain(NewForumNodeDomain())
}

func NewForumNodeDomain() *sForumNodeDomain {
	return &sForumNodeDomain{}
}

func (s *sForumNodeDomain) GetCategories(ctx context.Context) ([]entity.ForumCategory, error) {
	var categories []entity.ForumCategory
	err := dao.ForumCategory.Ctx(ctx).OrderAsc(dao.ForumCategory.Columns().SortOrder).Scan(&categories)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "forum.node.query_categories_failed"))
	}
	return categories, nil
}

func (s *sForumNodeDomain) GetNodes(ctx context.Context) ([]entity.ForumNode, error) {
	var nodes []entity.ForumNode
	err := dao.ForumNode.Ctx(ctx).OrderAsc(dao.ForumNode.Columns().SortOrder).Scan(&nodes)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "forum.node.query_nodes_failed"))
	}
	return nodes, nil
}

func (s *sForumNodeDomain) CheckNodeReadPolicy(ctx context.Context, actor *model.Actor, node *entity.ForumNode) error {
	userLevel := 0
	if actor != nil {
		userLevel = actor.RoleLevel
	}
	if userLevel < int(node.MinRoleRead) {
		return gerror.New(gi18n.T(ctx, "forum.node.read_permission_denied"))
	}
	return nil
}

func (s *sForumNodeDomain) CheckNodeWritePolicy(ctx context.Context, actor *model.Actor, node *entity.ForumNode) error {
	userLevel := 0
	if actor != nil {
		userLevel = actor.RoleLevel
	}
	if userLevel < int(node.MinRoleWrite) {
		return gerror.New(gi18n.T(ctx, "forum.node.write_permission_denied"))
	}
	return nil
}

func (s *sForumNodeDomain) CheckNodeCreatePolicy(ctx context.Context, actor *model.Actor, node *entity.ForumNode) error {
	userLevel := 0
	if actor != nil {
		userLevel = actor.RoleLevel
	}
	if userLevel < int(node.MinRoleCreate) {
		return gerror.New(gi18n.T(ctx, "forum.node.create_permission_denied"))
	}
	return nil
}

func (s *sForumNodeDomain) UpdateStats(ctx context.Context, nodeId uint, topicDelta int, replyDelta int) error {
	_, err := dao.ForumNode.Ctx(ctx).Where(dao.ForumNode.Columns().Id, nodeId).Update(g.Map{
		dao.ForumNode.Columns().TopicCount: gdb.Raw(fmt.Sprintf("topic_count + %d", topicDelta)),
		dao.ForumNode.Columns().ReplyCount: gdb.Raw(fmt.Sprintf("reply_count + %d", replyDelta)),
	})
	return err
}

func (s *sForumNodeDomain) AdminCreateNode(ctx context.Context, categoryId uint, slug, nameI18N, descI18N string, sortOrder, minRoleRead, minRoleWrite, minRoleCreate int, moderators []byte) (uint, error) {
	id, err := dao.ForumNode.Ctx(ctx).Data(g.Map{
		dao.ForumNode.Columns().CategoryId:    categoryId,
		dao.ForumNode.Columns().Slug:          slug,
		dao.ForumNode.Columns().NameI18N:      nameI18N,
		dao.ForumNode.Columns().DescI18N:      descI18N,
		dao.ForumNode.Columns().SortOrder:     sortOrder,
		dao.ForumNode.Columns().MinRoleRead:   minRoleRead,
		dao.ForumNode.Columns().MinRoleWrite:  minRoleWrite,
		dao.ForumNode.Columns().MinRoleCreate: minRoleCreate,
		dao.ForumNode.Columns().Moderators:    moderators,
	}).InsertAndGetId()
	return uint(id), err
}

func (s *sForumNodeDomain) AdminUpdateNode(ctx context.Context, id uint, categoryId *uint, slug, nameI18N, descI18N *string, sortOrder, minRoleRead, minRoleWrite, minRoleCreate *int, moderators []byte) error {
	data := g.Map{}
	if categoryId != nil {
		data[dao.ForumNode.Columns().CategoryId] = *categoryId
	}
	if slug != nil {
		data[dao.ForumNode.Columns().Slug] = *slug
	}
	if nameI18N != nil {
		data[dao.ForumNode.Columns().NameI18N] = *nameI18N
	}
	if descI18N != nil {
		data[dao.ForumNode.Columns().DescI18N] = *descI18N
	}
	if sortOrder != nil {
		data[dao.ForumNode.Columns().SortOrder] = *sortOrder
	}
	if minRoleRead != nil {
		data[dao.ForumNode.Columns().MinRoleRead] = *minRoleRead
	}
	if minRoleWrite != nil {
		data[dao.ForumNode.Columns().MinRoleWrite] = *minRoleWrite
	}
	if minRoleCreate != nil {
		data[dao.ForumNode.Columns().MinRoleCreate] = *minRoleCreate
	}
	if moderators != nil {
		data[dao.ForumNode.Columns().Moderators] = moderators
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.ForumNode.Ctx(ctx).Where(dao.ForumNode.Columns().Id, id).Data(data).Update()
	return err
}

func (s *sForumNodeDomain) AdminDeleteNode(ctx context.Context, id uint) (int, error) {
	count, err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().NodeId, id).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return count, nil
	}
	_, err = dao.ForumNode.Ctx(ctx).Where(dao.ForumNode.Columns().Id, id).Delete()
	return 0, err
}

func (s *sForumNodeDomain) AdminListNodes(ctx context.Context) ([]entity.ForumNode, error) {
	var list []entity.ForumNode
	err := dao.ForumNode.Ctx(ctx).OrderDesc(dao.ForumNode.Columns().SortOrder).Scan(&list)
	return list, err
}

func (s *sForumNodeDomain) GetNodeBySlug(ctx context.Context, slug string) (*entity.ForumNode, error) {
	var node *entity.ForumNode
	err := dao.ForumNode.Ctx(ctx).Where(dao.ForumNode.Columns().Slug, slug).Scan(&node)
	return node, err
}

func (s *sForumNodeDomain) GetNodeById(ctx context.Context, id uint) (*entity.ForumNode, error) {
	var node *entity.ForumNode
	err := dao.ForumNode.Ctx(ctx).Where(dao.ForumNode.Columns().Id, id).Scan(&node)
	return node, err
}

func (s *sForumNodeDomain) IncrementNodeStats(ctx context.Context, nodeId uint, topicId uint64) error {
	_, err := dao.ForumNode.Ctx(ctx).Where(dao.ForumNode.Columns().Id, nodeId).Update(g.Map{
		dao.ForumNode.Columns().TopicCount:  gdb.Raw("topic_count + 1"),
		dao.ForumNode.Columns().LastTopicId: topicId,
		dao.ForumNode.Columns().LastReplyAt: gtime.Now(),
	})
	return err
}
