package catalog

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sCatalogTagDomain struct{}

func init() {
	service.RegisterCatalogTagDomain(NewCatalogTagDomain())
}

func NewCatalogTagDomain() *sCatalogTagDomain {
	return &sCatalogTagDomain{}
}

func (s *sCatalogTagDomain) ListTagGroups(ctx context.Context) ([]entity.CatalogTagGroup, []entity.CatalogTag, error) {
	var groups []entity.CatalogTagGroup
	err := dao.CatalogTagGroup.Ctx(ctx).
		OrderAsc(dao.CatalogTagGroup.Columns().SortOrder).
		OrderAsc(dao.CatalogTagGroup.Columns().Id).
		Scan(&groups)
	if err != nil {
		return nil, nil, err
	}

	var tags []entity.CatalogTag
	err = dao.CatalogTag.Ctx(ctx).
		OrderAsc(dao.CatalogTag.Columns().GroupId).
		OrderAsc(dao.CatalogTag.Columns().SortOrder).
		OrderAsc(dao.CatalogTag.Columns().Id).
		Scan(&tags)
	return groups, tags, err
}

func (s *sCatalogTagDomain) GetTagsByIds(ctx context.Context, ids []uint) ([]entity.CatalogTag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var tags []entity.CatalogTag
	err := dao.CatalogTag.Ctx(ctx).
		WhereIn(dao.CatalogTag.Columns().Id, ids).
		Scan(&tags)
	return tags, err
}

func (s *sCatalogTagDomain) ReplaceTorrentTags(ctx context.Context, torrentId uint64, tagIds []uint) error {
	columns := dao.CatalogTorrentTag.Columns()
	if _, err := dao.CatalogTorrentTag.Ctx(ctx).Where(columns.TorrentId, torrentId).Delete(); err != nil {
		return err
	}
	if len(tagIds) == 0 {
		return nil
	}

	rows := make([]g.Map, 0, len(tagIds))
	for _, tagId := range tagIds {
		rows = append(rows, g.Map{
			columns.TorrentId: torrentId,
			columns.TagId:     tagId,
		})
	}
	_, err := dao.CatalogTorrentTag.Ctx(ctx).Data(rows).Insert()
	return err
}

func (s *sCatalogTagDomain) QueryTorrentTagsByTorrentIds(ctx context.Context, torrentIds []uint64) ([]entity.CatalogTorrentTag, []entity.CatalogTag, error) {
	if len(torrentIds) == 0 {
		return nil, nil, nil
	}

	var relations []entity.CatalogTorrentTag
	err := dao.CatalogTorrentTag.Ctx(ctx).
		WhereIn(dao.CatalogTorrentTag.Columns().TorrentId, torrentIds).
		OrderAsc(dao.CatalogTorrentTag.Columns().Id).
		Scan(&relations)
	if err != nil || len(relations) == 0 {
		return relations, nil, err
	}

	tagIds := make([]uint, 0, len(relations))
	for _, relation := range relations {
		tagIds = append(tagIds, relation.TagId)
	}
	tags, err := s.GetTagsByIds(ctx, tagIds)
	return relations, tags, err
}

func (s *sCatalogTagDomain) AdminGetTagGroupById(ctx context.Context, id uint) (*entity.CatalogTagGroup, error) {
	var group *entity.CatalogTagGroup
	err := dao.CatalogTagGroup.Ctx(ctx).Where(dao.CatalogTagGroup.Columns().Id, id).Scan(&group)
	return group, err
}

func (s *sCatalogTagDomain) AdminGetTagById(ctx context.Context, id uint) (*entity.CatalogTag, error) {
	var tag *entity.CatalogTag
	err := dao.CatalogTag.Ctx(ctx).Where(dao.CatalogTag.Columns().Id, id).Scan(&tag)
	return tag, err
}

func (s *sCatalogTagDomain) AdminCreateTagGroup(ctx context.Context, nameI18N []byte, slug string, categoryIds []byte, sortOrder int) (uint, error) {
	id, err := dao.CatalogTagGroup.Ctx(ctx).Data(g.Map{
		dao.CatalogTagGroup.Columns().NameI18N:    nameI18N,
		dao.CatalogTagGroup.Columns().Slug:        slug,
		dao.CatalogTagGroup.Columns().CategoryIds: categoryIds,
		dao.CatalogTagGroup.Columns().SortOrder:   sortOrder,
	}).InsertAndGetId()
	return uint(id), err
}

func (s *sCatalogTagDomain) AdminUpdateTagGroup(ctx context.Context, id uint, nameI18N []byte, categoryIds []byte, sortOrder int) error {
	_, err := dao.CatalogTagGroup.Ctx(ctx).
		Where(dao.CatalogTagGroup.Columns().Id, id).
		Data(g.Map{
			dao.CatalogTagGroup.Columns().NameI18N:    nameI18N,
			dao.CatalogTagGroup.Columns().CategoryIds: categoryIds,
			dao.CatalogTagGroup.Columns().SortOrder:   sortOrder,
		}).Update()
	return err
}

func (s *sCatalogTagDomain) AdminDeleteTagGroup(ctx context.Context, id uint) error {
	_, err := dao.CatalogTagGroup.Ctx(ctx).Where(dao.CatalogTagGroup.Columns().Id, id).Delete()
	return err
}

func (s *sCatalogTagDomain) AdminCreateTag(ctx context.Context, groupId uint, nameI18N []byte, value string, sortOrder int) (uint, error) {
	id, err := dao.CatalogTag.Ctx(ctx).Data(g.Map{
		dao.CatalogTag.Columns().GroupId:   groupId,
		dao.CatalogTag.Columns().NameI18N:  nameI18N,
		dao.CatalogTag.Columns().Value:     value,
		dao.CatalogTag.Columns().SortOrder: sortOrder,
	}).InsertAndGetId()
	return uint(id), err
}

func (s *sCatalogTagDomain) AdminUpdateTag(ctx context.Context, id uint, nameI18N []byte, sortOrder int) error {
	_, err := dao.CatalogTag.Ctx(ctx).
		Where(dao.CatalogTag.Columns().Id, id).
		Data(g.Map{
			dao.CatalogTag.Columns().NameI18N:  nameI18N,
			dao.CatalogTag.Columns().SortOrder: sortOrder,
		}).Update()
	return err
}

func (s *sCatalogTagDomain) AdminDeleteTag(ctx context.Context, id uint) error {
	_, err := dao.CatalogTag.Ctx(ctx).Where(dao.CatalogTag.Columns().Id, id).Delete()
	return err
}

func (s *sCatalogTagDomain) CountTagsByGroup(ctx context.Context, groupId uint) (int, error) {
	return dao.CatalogTag.Ctx(ctx).Where(dao.CatalogTag.Columns().GroupId, groupId).Count()
}

func (s *sCatalogTagDomain) CountTorrentTagsByTag(ctx context.Context, tagId uint) (int, error) {
	return dao.CatalogTorrentTag.Ctx(ctx).Where(dao.CatalogTorrentTag.Columns().TagId, tagId).Count()
}
