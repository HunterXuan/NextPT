package admin

import (
	"context"

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

type sAdminCatalogCategoryUsecase struct{}

func init() {
	service.RegisterAdminCatalogCategoryUsecase(NewAdminCatalogCategoryUsecase())
}

func NewAdminCatalogCategoryUsecase() *sAdminCatalogCategoryUsecase {
	return &sAdminCatalogCategoryUsecase{}
}

func (s *sAdminCatalogCategoryUsecase) Create(ctx context.Context, actor *model.Actor, in adminin.CatalogCategoryCreateInp) error {
	enabled := true
	if in.Enabled != nil {
		enabled = *in.Enabled
	}

	nameI18N, _ := gjson.Encode(in.NameI18N)
	uploadConfig, err := s.encodeUploadConfig(in.UploadConfig)
	if err != nil {
		return err
	}
	id, err := service.CatalogCategoryDomain().AdminCreateCategory(ctx, nameI18N, in.Slug, in.SortOrder, enabled, uploadConfig)
	if err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionCreate,
		TargetType: consts.SiteAuditTargetTypeCatalogCategory,
		TargetId:   uint64(id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"slug":    in.Slug,
			"enabled": enabled,
		},
	})
	return nil
}

func (s *sAdminCatalogCategoryUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.CatalogCategoryUpdateInp) error {
	var nameI18N []byte
	if in.NameI18N != nil {
		nameI18N, _ = gjson.Encode(in.NameI18N)
	}
	var uploadConfig *[]byte
	if in.UploadConfig != nil {
		encodedConfig, err := s.encodeUploadConfig(in.UploadConfig)
		if err != nil {
			return err
		}
		uploadConfig = &encodedConfig
	}
	if err := service.CatalogCategoryDomain().AdminUpdateCategory(ctx, in.Id, nameI18N, in.Slug, in.SortOrder, in.Enabled, uploadConfig); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeCatalogCategory,
		TargetId:   uint64(in.Id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"slugChanged":         in.Slug != nil,
			"sortOrderChanged":    in.SortOrder != nil,
			"enabledChanged":      in.Enabled != nil,
			"nameChanged":         in.NameI18N != nil,
			"uploadConfigChanged": in.UploadConfig != nil,
		},
	})
	return nil
}

func (s *sAdminCatalogCategoryUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.CatalogCategoryDeleteInp) error {
	torrentCount, err := service.CatalogTorrentDomain().CountTorrentsByCategory(ctx, in.Id)
	if err != nil {
		return err
	}
	if torrentCount > 0 {
		return gerror.New(gi18n.T(ctx, "admin.catalog.category_has_torrents"))
	}

	if err := service.CatalogCategoryDomain().AdminDeleteCategory(ctx, in.Id); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionDelete,
		TargetType: consts.SiteAuditTargetTypeCatalogCategory,
		TargetId:   uint64(in.Id),
		Level:      consts.SiteAuditLevelCritical,
	})
	return nil
}

func (s *sAdminCatalogCategoryUsecase) List(ctx context.Context, actor *model.Actor, in adminin.CatalogCategoryListInp) (*adminout.CatalogCategoryListOut, error) {
	categories, err := service.CatalogCategoryDomain().AdminListCategories(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]adminout.CatalogCategoryItem, 0, len(categories))
	for _, category := range categories {
		uploadConfig := s.scanUploadConfig(category.UploadConfig)
		items = append(items, adminout.CatalogCategoryItem{
			Id:           category.Id,
			NameI18N:     category.NameI18N,
			Slug:         category.Slug,
			SortOrder:    category.SortOrder,
			Enabled:      category.Enabled,
			UploadConfig: uploadConfig,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
		})
	}
	return &adminout.CatalogCategoryListOut{Categories: items}, nil
}

func (s *sAdminCatalogCategoryUsecase) encodeUploadConfig(config any) ([]byte, error) {
	if config == nil {
		return nil, nil
	}
	return gjson.Encode(config)
}

func (s *sAdminCatalogCategoryUsecase) scanUploadConfig(value *gjson.Json) *model.CatalogUploadConfig {
	if value == nil {
		return nil
	}
	var config model.CatalogUploadConfig
	if err := value.Scan(&config); err != nil {
		return nil
	}
	return &config
}
