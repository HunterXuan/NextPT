package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/catalogin"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

const (
	uploadTitleModeGenerated = "generated"

	uploadFieldTypeText        = "text"
	uploadFieldTypeTextarea    = "textarea"
	uploadFieldTypeSelect      = "select"
	uploadFieldTypeMultiSelect = "multiSelect"

	uploadOptionSourceStatic   = "static"
	uploadOptionSourceTagGroup = "tagGroup"
)

type uploadReleaseData struct {
	Name   string
	Fields map[string]any
}

func (s *sCatalogTorrentUsecase) prepareUploadReleaseData(ctx context.Context, category *entity.CatalogCategory, infoName string, in catalogin.TorrentUploadInp) (*uploadReleaseData, error) {
	manualName := strings.TrimSpace(in.Name)
	if category == nil || category.UploadConfig == nil {
		return &uploadReleaseData{Name: s.resolveUploadName(manualName, infoName), Fields: nil}, nil
	}

	config, err := s.scanUploadConfig(ctx, category)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return &uploadReleaseData{Name: s.resolveUploadName(manualName, infoName), Fields: nil}, nil
	}

	rawFields, err := s.parseReleaseFields(ctx, in.ReleaseFields)
	if err != nil {
		return nil, err
	}

	fields, err := s.validateReleaseFields(ctx, category.Id, config, rawFields)
	if err != nil {
		return nil, err
	}

	title := s.resolveUploadName(manualName, infoName)
	if config.Title.Mode == uploadTitleModeGenerated {
		generatedTitle := s.buildGeneratedTitle(config.Title.Parts, fields)
		title = s.resolveGeneratedUploadName(manualName, generatedTitle, infoName, config.Title.AllowManualOverride)
	}

	return &uploadReleaseData{Name: title, Fields: fields}, nil
}

func (s *sCatalogTorrentUsecase) prepareUpdateReleaseData(ctx context.Context, category *entity.CatalogCategory, torrent *entity.CatalogTorrent, in catalogin.TorrentUpdateInp) (*uploadReleaseData, error) {
	manualName := strings.TrimSpace(in.Name)
	currentName := ""
	if torrent != nil {
		currentName = torrent.Name
	}
	if category == nil || category.UploadConfig == nil {
		return &uploadReleaseData{Name: s.resolveUploadName(manualName, currentName), Fields: nil}, nil
	}

	config, err := s.scanUploadConfig(ctx, category)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return &uploadReleaseData{Name: s.resolveUploadName(manualName, currentName), Fields: nil}, nil
	}

	rawFields, err := s.parseUpdateReleaseFields(ctx, torrent, in.ReleaseFields)
	if err != nil {
		return nil, err
	}

	fields, err := s.validateReleaseFields(ctx, category.Id, config, rawFields)
	if err != nil {
		return nil, err
	}

	title := s.resolveUploadName(manualName, currentName)
	if config.Title.Mode == uploadTitleModeGenerated {
		generatedTitle := s.buildGeneratedTitle(config.Title.Parts, fields)
		title = s.resolveGeneratedUploadName(manualName, generatedTitle, currentName, config.Title.AllowManualOverride)
	}

	return &uploadReleaseData{Name: title, Fields: fields}, nil
}

func (s *sCatalogTorrentUsecase) parseUpdateReleaseFields(ctx context.Context, torrent *entity.CatalogTorrent, raw string) (map[string]any, error) {
	if strings.TrimSpace(raw) != "" {
		return s.parseReleaseFields(ctx, raw)
	}
	if torrent == nil || torrent.ReleaseFields == nil {
		return map[string]any{}, nil
	}

	var fields map[string]any
	if err := torrent.ReleaseFields.Scan(&fields); err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.release_fields_invalid"))
	}
	if fields == nil {
		fields = map[string]any{}
	}
	return fields, nil
}

func (s *sCatalogTorrentUsecase) scanUploadConfig(ctx context.Context, category *entity.CatalogCategory) (*model.CatalogUploadConfig, error) {
	if category == nil || category.UploadConfig == nil {
		return nil, nil
	}
	var config model.CatalogUploadConfig
	if err := category.UploadConfig.Scan(&config); err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.upload_config_invalid"))
	}
	return &config, nil
}

func (s *sCatalogTorrentUsecase) parseReleaseFields(ctx context.Context, raw string) (map[string]any, error) {
	source := strings.TrimSpace(raw)
	if source == "" {
		return map[string]any{}, nil
	}

	var fields map[string]any
	if err := json.Unmarshal([]byte(source), &fields); err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.release_fields_invalid"))
	}
	if fields == nil {
		fields = map[string]any{}
	}
	return fields, nil
}

func (s *sCatalogTorrentUsecase) validateReleaseFields(ctx context.Context, categoryId uint, config *model.CatalogUploadConfig, rawFields map[string]any) (map[string]any, error) {
	if config == nil || len(config.Fields) == 0 {
		return nil, nil
	}

	tagGroupOptions, err := s.loadUploadTagGroupOptions(ctx, categoryId, config.Fields)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]any, len(config.Fields))
	for _, field := range config.Fields {
		if strings.TrimSpace(field.Key) == "" {
			continue
		}

		value, exists := rawFields[field.Key]
		switch field.Type {
		case uploadFieldTypeSelect:
			normalized := strings.TrimSpace(s.stringifyReleaseValue(value))
			if field.Required && normalized == "" {
				return nil, s.releaseFieldRequiredError(ctx, field.Key)
			}
			if normalized != "" && !s.isAllowedReleaseOption(field, normalized, tagGroupOptions) {
				return nil, s.releaseFieldOptionError(ctx, field.Key)
			}
			if normalized != "" {
				fields[field.Key] = normalized
			}
		case uploadFieldTypeMultiSelect:
			normalized := s.stringifyReleaseValueList(value)
			if field.Required && len(normalized) == 0 {
				return nil, s.releaseFieldRequiredError(ctx, field.Key)
			}
			for _, item := range normalized {
				if !s.isAllowedReleaseOption(field, item, tagGroupOptions) {
					return nil, s.releaseFieldOptionError(ctx, field.Key)
				}
			}
			if len(normalized) > 0 {
				fields[field.Key] = normalized
			}
		default:
			normalized := strings.TrimSpace(s.stringifyReleaseValue(value))
			if field.Required && (!exists || normalized == "") {
				return nil, s.releaseFieldRequiredError(ctx, field.Key)
			}
			if normalized != "" {
				fields[field.Key] = normalized
			}
		}
	}

	if len(fields) == 0 {
		return nil, nil
	}
	return fields, nil
}

func (s *sCatalogTorrentUsecase) loadUploadTagGroupOptions(ctx context.Context, categoryId uint, fields []model.CatalogUploadFieldConfig) (map[string]map[string]struct{}, error) {
	needsTagGroups := false
	for _, field := range fields {
		if field.Options != nil && field.Options.Source == uploadOptionSourceTagGroup {
			needsTagGroups = true
			break
		}
	}
	if !needsTagGroups {
		return nil, nil
	}

	groups, tags, err := service.CatalogCategoryDomain().ListTagGroups(ctx)
	if err != nil {
		return nil, err
	}

	groupById := make(map[uint]entity.CatalogTagGroup, len(groups))
	for _, group := range groups {
		if !s.uploadTagGroupAppliesToCategory(group, categoryId) {
			continue
		}
		groupById[group.Id] = group
	}

	options := make(map[string]map[string]struct{})
	for _, tag := range tags {
		group, ok := groupById[tag.GroupId]
		if !ok || tag.Value == "" {
			continue
		}
		if options[group.Slug] == nil {
			options[group.Slug] = make(map[string]struct{})
		}
		options[group.Slug][tag.Value] = struct{}{}
	}
	return options, nil
}

func (s *sCatalogTorrentUsecase) uploadTagGroupAppliesToCategory(group entity.CatalogTagGroup, categoryId uint) bool {
	if group.CategoryIds == nil {
		return true
	}
	var categoryIds []uint
	if err := group.CategoryIds.Scan(&categoryIds); err != nil || len(categoryIds) == 0 {
		return true
	}
	for _, id := range categoryIds {
		if id == categoryId {
			return true
		}
	}
	return false
}

func (s *sCatalogTorrentUsecase) isAllowedReleaseOption(field model.CatalogUploadFieldConfig, value string, tagGroupOptions map[string]map[string]struct{}) bool {
	if field.Options == nil {
		return true
	}

	switch field.Options.Source {
	case uploadOptionSourceStatic:
		for _, item := range field.Options.Items {
			if item.Value == value {
				return true
			}
		}
		return false
	case uploadOptionSourceTagGroup:
		values := tagGroupOptions[field.Options.Slug]
		_, ok := values[value]
		return ok
	default:
		return true
	}
}

func (s *sCatalogTorrentUsecase) buildGeneratedTitle(parts []model.CatalogUploadTitlePart, fields map[string]any) string {
	var builder strings.Builder
	for _, part := range parts {
		value := strings.TrimSpace(s.releaseTitlePartValue(fields[part.Field], part.Separator))
		if value == "" {
			continue
		}
		builder.WriteString(part.Prefix)
		builder.WriteString(value)
		builder.WriteString(part.Suffix)
	}
	return strings.TrimSpace(builder.String())
}

func (s *sCatalogTorrentUsecase) releaseTitlePartValue(value any, separator string) string {
	if separator == "" {
		separator = "/"
	}
	switch typed := value.(type) {
	case []string:
		return strings.Join(typed, separator)
	default:
		return s.stringifyReleaseValue(value)
	}
}

func (s *sCatalogTorrentUsecase) stringifyReleaseValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return fmt.Sprintf("%g", typed)
	case float32:
		return fmt.Sprintf("%g", typed)
	case int:
		return fmt.Sprintf("%d", typed)
	case int64:
		return fmt.Sprintf("%d", typed)
	case uint:
		return fmt.Sprintf("%d", typed)
	case uint64:
		return fmt.Sprintf("%d", typed)
	case bool:
		if typed {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func (s *sCatalogTorrentUsecase) stringifyReleaseValueList(value any) []string {
	switch typed := value.(type) {
	case []string:
		return s.cleanReleaseValueList(typed)
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			items = append(items, s.stringifyReleaseValue(item))
		}
		return s.cleanReleaseValueList(items)
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil
		}
		return []string{strings.TrimSpace(typed)}
	default:
		return nil
	}
}

func (s *sCatalogTorrentUsecase) cleanReleaseValueList(values []string) []string {
	items := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := strings.TrimSpace(value)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		items = append(items, normalized)
	}
	return items
}

func (s *sCatalogTorrentUsecase) resolveUploadName(manualName string, infoName string) string {
	if strings.TrimSpace(manualName) != "" {
		return strings.TrimSpace(manualName)
	}
	return strings.TrimSpace(infoName)
}

func (s *sCatalogTorrentUsecase) resolveGeneratedUploadName(manualName string, generatedName string, infoName string, allowManualOverride bool) string {
	if allowManualOverride && strings.TrimSpace(manualName) != "" {
		return strings.TrimSpace(manualName)
	}
	if strings.TrimSpace(generatedName) != "" {
		return strings.TrimSpace(generatedName)
	}
	return s.resolveUploadName(manualName, infoName)
}

func (s *sCatalogTorrentUsecase) releaseFieldRequiredError(ctx context.Context, key string) error {
	return gerror.Newf("%s: %s", gi18n.T(ctx, "catalog.torrent.release_field_required"), key)
}

func (s *sCatalogTorrentUsecase) releaseFieldOptionError(ctx context.Context, key string) error {
	return gerror.Newf("%s: %s", gi18n.T(ctx, "catalog.torrent.release_field_option_invalid"), key)
}
