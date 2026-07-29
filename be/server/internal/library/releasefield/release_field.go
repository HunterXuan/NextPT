package releasefield

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"server/internal/model"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

const (
	titleModeGenerated = "generated"

	fieldTypeSelect      = "select"
	fieldTypeMultiSelect = "multiSelect"

	optionSourceStatic   = "static"
	optionSourceTagGroup = "tagGroup"
)

type TagGroup struct {
	Id          uint
	Slug        string
	CategoryIds *gjson.Json
}

type Tag struct {
	Id      uint
	GroupId uint
	Value   string
}

type Input struct {
	ManualName    string
	FallbackName  string
	ReleaseFields string
	StoredFields  *gjson.Json
}

type Result struct {
	Name   string
	Fields map[string]any

	tagIds      []uint
	tagGroupIds []uint
}

type tagOption struct {
	id uint
}

type Resolver struct {
	config      *model.CatalogUploadConfig
	categoryId  uint
	groupsById  map[uint]TagGroup
	tagsById    map[uint]Tag
	tagOptions  map[string]map[string]tagOption
	tagGroupIds map[string]uint
}

func NewResolver(ctx context.Context, value *gjson.Json) (*Resolver, error) {
	resolver := &Resolver{}
	if value == nil {
		return resolver, nil
	}

	var config model.CatalogUploadConfig
	if err := value.Scan(&config); err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.upload_config_invalid"))
	}
	resolver.config = &config
	return resolver, nil
}

func (r *Resolver) UsesTagGroups() bool {
	if r == nil || r.config == nil {
		return false
	}
	for _, field := range r.config.Fields {
		if field.Options != nil && field.Options.Source == optionSourceTagGroup {
			return true
		}
	}
	return false
}

func (r *Resolver) SetTagCatalog(categoryId uint, groups []TagGroup, tags []Tag) {
	r.categoryId = categoryId
	r.groupsById = make(map[uint]TagGroup, len(groups))
	r.tagsById = make(map[uint]Tag, len(tags))
	r.tagOptions = make(map[string]map[string]tagOption)
	r.tagGroupIds = make(map[string]uint)

	for _, group := range groups {
		r.groupsById[group.Id] = group
		if !tagGroupAppliesToCategory(group, categoryId) {
			continue
		}
		r.tagGroupIds[group.Slug] = group.Id
	}

	for _, tag := range tags {
		r.tagsById[tag.Id] = tag
		group, ok := r.groupsById[tag.GroupId]
		if !ok || !tagGroupAppliesToCategory(group, categoryId) || tag.Value == "" {
			continue
		}
		if r.tagOptions[group.Slug] == nil {
			r.tagOptions[group.Slug] = make(map[string]tagOption)
		}
		r.tagOptions[group.Slug][tag.Value] = tagOption{id: tag.Id}
	}
}

func (r *Resolver) Resolve(ctx context.Context, in Input) (Result, error) {
	name := resolveName(in.ManualName, in.FallbackName)
	if r == nil || r.config == nil {
		return Result{Name: name}, nil
	}

	rawFields, err := parseFields(ctx, in)
	if err != nil {
		return Result{}, err
	}
	fields, tagIds, tagGroupIds, err := r.validateFields(ctx, rawFields)
	if err != nil {
		return Result{}, err
	}

	if r.config.Title.Mode == titleModeGenerated {
		generatedName := buildGeneratedTitle(r.config.Title.Parts, fields)
		name = resolveGeneratedName(in.ManualName, generatedName, in.FallbackName, r.config.Title.AllowManualOverride)
	}

	return Result{
		Name:        name,
		Fields:      fields,
		tagIds:      tagIds,
		tagGroupIds: tagGroupIds,
	}, nil
}

func (r *Resolver) ResolveTagIds(ctx context.Context, selectedIds []uint, result Result) ([]uint, error) {
	selectedIds = normalizeIds(selectedIds)
	controlledGroups := make(map[uint]struct{}, len(result.tagGroupIds))
	for _, groupId := range result.tagGroupIds {
		controlledGroups[groupId] = struct{}{}
	}

	resolved := make([]uint, 0, len(selectedIds)+len(result.tagIds))
	for _, id := range selectedIds {
		tag, ok := r.tagsById[id]
		if !ok {
			return nil, gerror.New(gi18n.T(ctx, "catalog.tag.invalid"))
		}
		group, ok := r.groupsById[tag.GroupId]
		if !ok || !tagGroupAppliesToCategory(group, r.categoryId) {
			return nil, gerror.New(gi18n.T(ctx, "catalog.tag.category_invalid"))
		}
		if _, controlled := controlledGroups[tag.GroupId]; controlled {
			continue
		}
		resolved = append(resolved, id)
	}
	resolved = append(resolved, result.tagIds...)
	return normalizeIds(resolved), nil
}

func (r *Resolver) validateFields(ctx context.Context, rawFields map[string]any) (map[string]any, []uint, []uint, error) {
	if len(r.config.Fields) == 0 {
		return nil, nil, nil, nil
	}

	fields := make(map[string]any, len(r.config.Fields))
	tagIds := make([]uint, 0)
	tagGroupIds := make([]uint, 0)
	seenGroups := make(map[uint]struct{})
	for _, field := range r.config.Fields {
		if strings.TrimSpace(field.Key) == "" {
			continue
		}

		value, exists := rawFields[field.Key]
		if field.Options != nil && field.Options.Source == optionSourceTagGroup {
			if groupId := r.tagGroupIds[field.Options.Slug]; groupId > 0 {
				if _, ok := seenGroups[groupId]; !ok {
					seenGroups[groupId] = struct{}{}
					tagGroupIds = append(tagGroupIds, groupId)
				}
			}
		}

		switch field.Type {
		case fieldTypeSelect:
			normalized := strings.TrimSpace(stringifyValue(value))
			if field.Required && normalized == "" {
				return nil, nil, nil, requiredFieldError(ctx, field.Key)
			}
			if normalized != "" && !r.isAllowedOption(field, normalized) {
				return nil, nil, nil, invalidOptionError(ctx, field.Key)
			}
			if normalized != "" {
				fields[field.Key] = normalized
				if field.Options != nil {
					if option := r.tagOptions[field.Options.Slug][normalized]; option.id > 0 {
						tagIds = append(tagIds, option.id)
					}
				}
			}
		case fieldTypeMultiSelect:
			normalized := stringifyValueList(value)
			if field.Required && len(normalized) == 0 {
				return nil, nil, nil, requiredFieldError(ctx, field.Key)
			}
			for _, item := range normalized {
				if !r.isAllowedOption(field, item) {
					return nil, nil, nil, invalidOptionError(ctx, field.Key)
				}
				if field.Options != nil {
					if option := r.tagOptions[field.Options.Slug][item]; option.id > 0 {
						tagIds = append(tagIds, option.id)
					}
				}
			}
			if len(normalized) > 0 {
				fields[field.Key] = normalized
			}
		default:
			normalized := strings.TrimSpace(stringifyValue(value))
			if field.Required && (!exists || normalized == "") {
				return nil, nil, nil, requiredFieldError(ctx, field.Key)
			}
			if normalized != "" {
				fields[field.Key] = normalized
			}
		}
	}

	if len(fields) == 0 {
		fields = nil
	}
	return fields, tagIds, tagGroupIds, nil
}

func (r *Resolver) isAllowedOption(field model.CatalogUploadFieldConfig, value string) bool {
	if field.Options == nil {
		return true
	}

	switch field.Options.Source {
	case optionSourceStatic:
		for _, item := range field.Options.Items {
			if item.Value == value {
				return true
			}
		}
		return false
	case optionSourceTagGroup:
		_, ok := r.tagOptions[field.Options.Slug][value]
		return ok
	default:
		return true
	}
}

func parseFields(ctx context.Context, in Input) (map[string]any, error) {
	if source := strings.TrimSpace(in.ReleaseFields); source != "" {
		var fields map[string]any
		if err := json.Unmarshal([]byte(source), &fields); err != nil {
			return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.release_fields_invalid"))
		}
		if fields != nil {
			return fields, nil
		}
		return map[string]any{}, nil
	}

	if in.StoredFields != nil {
		var fields map[string]any
		if err := in.StoredFields.Scan(&fields); err != nil {
			return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.release_fields_invalid"))
		}
		if fields != nil {
			return fields, nil
		}
	}
	return map[string]any{}, nil
}

func tagGroupAppliesToCategory(group TagGroup, categoryId uint) bool {
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

func buildGeneratedTitle(parts []model.CatalogUploadTitlePart, fields map[string]any) string {
	var builder strings.Builder
	for _, part := range parts {
		value := strings.TrimSpace(titlePartValue(fields[part.Field], part.Separator))
		if value == "" {
			continue
		}
		builder.WriteString(part.Prefix)
		builder.WriteString(value)
		builder.WriteString(part.Suffix)
	}
	return strings.TrimSpace(builder.String())
}

func titlePartValue(value any, separator string) string {
	if separator == "" {
		separator = "/"
	}
	if values, ok := value.([]string); ok {
		return strings.Join(values, separator)
	}
	return stringifyValue(value)
}

func stringifyValue(value any) string {
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

func stringifyValueList(value any) []string {
	switch typed := value.(type) {
	case []string:
		return cleanValueList(typed)
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			items = append(items, stringifyValue(item))
		}
		return cleanValueList(items)
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil
		}
		return []string{strings.TrimSpace(typed)}
	default:
		return nil
	}
}

func cleanValueList(values []string) []string {
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

func normalizeIds(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func resolveName(manualName string, fallbackName string) string {
	if strings.TrimSpace(manualName) != "" {
		return strings.TrimSpace(manualName)
	}
	return strings.TrimSpace(fallbackName)
}

func resolveGeneratedName(manualName string, generatedName string, fallbackName string, allowManualOverride bool) string {
	if allowManualOverride && strings.TrimSpace(manualName) != "" {
		return strings.TrimSpace(manualName)
	}
	if strings.TrimSpace(generatedName) != "" {
		return strings.TrimSpace(generatedName)
	}
	return resolveName(manualName, fallbackName)
}

func requiredFieldError(ctx context.Context, key string) error {
	return gerror.Newf("%s: %s", gi18n.T(ctx, "catalog.torrent.release_field_required"), key)
}

func invalidOptionError(ctx context.Context, key string) error {
	return gerror.Newf("%s: %s", gi18n.T(ctx, "catalog.torrent.release_field_option_invalid"), key)
}
