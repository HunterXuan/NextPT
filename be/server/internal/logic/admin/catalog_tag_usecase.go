package admin

import (
	"context"
	"strings"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sAdminCatalogTagUsecase struct{}

func init() {
	service.RegisterAdminCatalogTagUsecase(NewAdminCatalogTagUsecase())
}

func NewAdminCatalogTagUsecase() *sAdminCatalogTagUsecase {
	return &sAdminCatalogTagUsecase{}
}

func (s *sAdminCatalogTagUsecase) List(ctx context.Context, actor *model.Actor, in adminin.CatalogTagGroupListInp) (*adminout.CatalogTagGroupListOut, error) {
	groups, tags, err := service.CatalogTagDomain().ListTagGroups(ctx)
	if err != nil {
		return nil, err
	}
	tagMap := make(map[uint][]adminout.CatalogTagItem)
	for _, tag := range tags {
		tagMap[tag.GroupId] = append(tagMap[tag.GroupId], adminout.CatalogTagItem{
			Id: tag.Id, GroupId: tag.GroupId, NameI18N: tag.NameI18N, Value: tag.Value,
			SortOrder: tag.SortOrder, CreatedAt: tag.CreatedAt, UpdatedAt: tag.UpdatedAt,
		})
	}
	items := make([]adminout.CatalogTagGroupItem, 0, len(groups))
	for _, group := range groups {
		var categoryIds []uint
		_ = group.CategoryIds.Scan(&categoryIds)
		items = append(items, adminout.CatalogTagGroupItem{
			Id: group.Id, NameI18N: group.NameI18N, Slug: group.Slug, CategoryIds: categoryIds,
			SortOrder: group.SortOrder, Tags: tagMap[group.Id], CreatedAt: group.CreatedAt, UpdatedAt: group.UpdatedAt,
		})
	}
	return &adminout.CatalogTagGroupListOut{Groups: items}, nil
}

func (s *sAdminCatalogTagUsecase) CreateGroup(ctx context.Context, actor *model.Actor, in adminin.CatalogTagGroupCreateInp) error {
	nameI18N, _ := gjson.Encode(in.NameI18N)
	categoryIds, _ := gjson.Encode(in.CategoryIds)
	id, err := service.CatalogTagDomain().AdminCreateTagGroup(ctx, nameI18N, strings.TrimSpace(in.Slug), categoryIds, in.SortOrder)
	if err != nil {
		return err
	}
	s.recordTagAudit(ctx, actor, consts.SiteAuditActionCreate, consts.SiteAuditTargetTypeCatalogTagGroup, uint64(id), map[string]any{"slug": in.Slug})
	return nil
}

func (s *sAdminCatalogTagUsecase) UpdateGroup(ctx context.Context, actor *model.Actor, in adminin.CatalogTagGroupUpdateInp) error {
	group, err := service.CatalogTagDomain().AdminGetTagGroupById(ctx, in.Id)
	if err != nil {
		return err
	}
	if group == nil {
		return gerror.New(gi18n.T(ctx, "admin.catalog.tag_group_not_found"))
	}
	nameI18N, _ := gjson.Encode(in.NameI18N)
	categoryIds, _ := gjson.Encode(in.CategoryIds)
	if err := service.CatalogTagDomain().AdminUpdateTagGroup(ctx, in.Id, nameI18N, categoryIds, in.SortOrder); err != nil {
		return err
	}
	s.recordTagAudit(ctx, actor, consts.SiteAuditActionUpdate, consts.SiteAuditTargetTypeCatalogTagGroup, uint64(in.Id), nil)
	return nil
}

func (s *sAdminCatalogTagUsecase) DeleteGroup(ctx context.Context, actor *model.Actor, in adminin.CatalogTagGroupDeleteInp) error {
	count, err := service.CatalogTagDomain().CountTagsByGroup(ctx, in.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "admin.catalog.tag_group_has_tags"))
	}
	if err := service.CatalogTagDomain().AdminDeleteTagGroup(ctx, in.Id); err != nil {
		return err
	}
	s.recordTagAudit(ctx, actor, consts.SiteAuditActionDelete, consts.SiteAuditTargetTypeCatalogTagGroup, uint64(in.Id), nil)
	return nil
}

func (s *sAdminCatalogTagUsecase) CreateTag(ctx context.Context, actor *model.Actor, in adminin.CatalogTagCreateInp) error {
	group, err := service.CatalogTagDomain().AdminGetTagGroupById(ctx, in.GroupId)
	if err != nil {
		return err
	}
	if group == nil {
		return gerror.New(gi18n.T(ctx, "admin.catalog.tag_group_not_found"))
	}
	nameI18N, _ := gjson.Encode(in.NameI18N)
	id, err := service.CatalogTagDomain().AdminCreateTag(ctx, in.GroupId, nameI18N, strings.TrimSpace(in.Value), in.SortOrder)
	if err != nil {
		return err
	}
	s.recordTagAudit(ctx, actor, consts.SiteAuditActionCreate, consts.SiteAuditTargetTypeCatalogTag, uint64(id), map[string]any{"groupId": in.GroupId, "value": in.Value})
	return nil
}

func (s *sAdminCatalogTagUsecase) UpdateTag(ctx context.Context, actor *model.Actor, in adminin.CatalogTagUpdateInp) error {
	tag, err := service.CatalogTagDomain().AdminGetTagById(ctx, in.Id)
	if err != nil {
		return err
	}
	if tag == nil {
		return gerror.New(gi18n.T(ctx, "admin.catalog.tag_not_found"))
	}
	nameI18N, _ := gjson.Encode(in.NameI18N)
	if err := service.CatalogTagDomain().AdminUpdateTag(ctx, in.Id, nameI18N, in.SortOrder); err != nil {
		return err
	}
	s.recordTagAudit(ctx, actor, consts.SiteAuditActionUpdate, consts.SiteAuditTargetTypeCatalogTag, uint64(in.Id), nil)
	return nil
}

func (s *sAdminCatalogTagUsecase) DeleteTag(ctx context.Context, actor *model.Actor, in adminin.CatalogTagDeleteInp) error {
	count, err := service.CatalogTagDomain().CountTorrentTagsByTag(ctx, in.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "admin.catalog.tag_in_use"))
	}
	if err := service.CatalogTagDomain().AdminDeleteTag(ctx, in.Id); err != nil {
		return err
	}
	s.recordTagAudit(ctx, actor, consts.SiteAuditActionDelete, consts.SiteAuditTargetTypeCatalogTag, uint64(in.Id), nil)
	return nil
}

func (s *sAdminCatalogTagUsecase) recordTagAudit(ctx context.Context, actor *model.Actor, action, targetType string, targetId uint64, detail any) {
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action: action, TargetType: targetType, TargetId: targetId, Level: consts.SiteAuditLevelCritical, Detail: detail,
	})
}
