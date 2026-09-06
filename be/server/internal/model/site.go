package model

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"server/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type SiteMaintenance struct {
	Enabled bool
	Message string
}

func (m SiteMaintenance) DisplayMessage(defaultMessage string) string {
	if message := strings.TrimSpace(m.Message); message != "" {
		return message
	}
	return defaultMessage
}

type SiteMaintenanceError struct {
	Message string
}

func (e *SiteMaintenanceError) Error() string {
	return e.Message
}

type SiteAdvertisementConfig struct {
	Enabled     bool   `json:"enabled"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	AspectRatio string `json:"aspectRatio"`
}

type SiteAdvertisements map[string]SiteAdvertisementConfig

func (ads SiteAdvertisements) Normalized() SiteAdvertisements {
	result := make(SiteAdvertisements, len(consts.SiteAdvertisementPlacements))
	for _, placement := range consts.SiteAdvertisementPlacements {
		advertisement := ads[placement]
		advertisement.Title = strings.TrimSpace(advertisement.Title)
		advertisement.Image = strings.TrimSpace(advertisement.Image)
		advertisement.URL = strings.TrimSpace(advertisement.URL)
		advertisement.AspectRatio = normalizeSiteAdvertisementAspectRatio(advertisement.AspectRatio)
		if advertisement.AspectRatio == "" {
			advertisement.AspectRatio = consts.SiteAdvertisementDefaultAspectRatio(placement)
		}
		result[placement] = advertisement
	}
	return result
}

func (ads SiteAdvertisements) Enabled() SiteAdvertisements {
	result := make(SiteAdvertisements)
	for _, placement := range consts.SiteAdvertisementPlacements {
		if advertisement := ads[placement]; advertisement.Enabled {
			result[placement] = advertisement
		}
	}
	return result
}

func (ads SiteAdvertisements) Validate() error {
	for placement := range ads {
		if !slices.Contains(consts.SiteAdvertisementPlacements, placement) {
			return gerror.New("unsupported advertisement placement")
		}
	}

	for _, placement := range consts.SiteAdvertisementPlacements {
		advertisement := ads[placement]
		if advertisement.AspectRatio != "" && !isSiteAdvertisementAspectRatio(normalizeSiteAdvertisementAspectRatio(advertisement.AspectRatio)) {
			return gerror.New("advertisement aspect ratio must use width:height format")
		}
		if !advertisement.Enabled {
			continue
		}
		if advertisement.Title == "" || len(advertisement.Title) > 120 {
			return gerror.New("advertisement title is required and must not exceed 120 characters")
		}
		if !isSiteAdvertisementURL(advertisement.Image) {
			return gerror.New("advertisement image must be an internal path or http(s) URL")
		}
		if !isSiteAdvertisementURL(advertisement.URL) {
			return gerror.New("advertisement URL must be an internal path or http(s) URL")
		}
	}
	return nil
}

func normalizeSiteAdvertisementAspectRatio(value string) string {
	parts := strings.Split(strings.TrimSpace(value), ":")
	if len(parts) != 2 {
		return strings.TrimSpace(value)
	}
	return strings.TrimSpace(parts[0]) + ":" + strings.TrimSpace(parts[1])
}

func isSiteAdvertisementAspectRatio(value string) bool {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return false
	}
	width, widthErr := strconv.Atoi(parts[0])
	height, heightErr := strconv.Atoi(parts[1])
	return widthErr == nil && heightErr == nil && width > 0 && height > 0 && width <= 100 && height <= 100
}

func isSiteAdvertisementURL(value string) bool {
	if strings.HasPrefix(value, "/") {
		return !strings.HasPrefix(value, "//")
	}
	parsed, err := url.ParseRequestURI(value)
	return err == nil && parsed.Host != "" && (parsed.Scheme == "http" || parsed.Scheme == "https")
}

type SiteTaskReward struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

type SiteTaskRule struct {
	Type   string `json:"type"`
	Target uint64 `json:"target"`
}

type SiteTaskDefinition struct {
	Key             string            `json:"key"`
	Enabled         bool              `json:"enabled"`
	Cycle           string            `json:"cycle"`
	NameI18N        map[string]string `json:"nameI18n"`
	DescriptionI18N map[string]string `json:"descriptionI18n"`
	Rule            SiteTaskRule      `json:"rule"`
	Rewards         []SiteTaskReward  `json:"rewards"`
}

func (task SiteTaskDefinition) CyclePeriod(now time.Time) (string, *gtime.Time, *gtime.Time) {
	switch task.Cycle {
	case consts.SiteTaskCycleWeekly:
		started := task.beginningOfWeek(now)
		year, week := started.ISOWeek()
		return fmt.Sprintf("%04d-W%02d", year, week), gtime.NewFromTime(started), gtime.NewFromTime(started.AddDate(0, 0, 7))
	case consts.SiteTaskCycleMonthly:
		started := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return started.Format("2006-01"), gtime.NewFromTime(started), gtime.NewFromTime(started.AddDate(0, 1, 0))
	default:
		return consts.SiteTaskCycleOnce, gtime.NewFromTime(now), nil
	}
}

func (task SiteTaskDefinition) beginningOfWeek(now time.Time) time.Time {
	day := int(now.Weekday())
	if day == 0 {
		day = 7
	}
	return time.Date(now.Year(), now.Month(), now.Day()-day+1, 0, 0, 0, 0, now.Location())
}

type SiteTasks []SiteTaskDefinition

var siteTaskKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{2,63}$`)

func (tasks SiteTasks) Normalized() SiteTasks {
	result := make(SiteTasks, len(tasks))
	for index, task := range tasks {
		task.Key = strings.TrimSpace(task.Key)
		task.Cycle = strings.TrimSpace(task.Cycle)
		task.Rule.Type = strings.TrimSpace(task.Rule.Type)
		task.NameI18N = normalizeSiteTaskI18N(task.NameI18N)
		task.DescriptionI18N = normalizeSiteTaskI18N(task.DescriptionI18N)
		for rewardIndex := range task.Rewards {
			task.Rewards[rewardIndex].Type = strings.TrimSpace(task.Rewards[rewardIndex].Type)
			task.Rewards[rewardIndex].Amount = math.Round(task.Rewards[rewardIndex].Amount*10) / 10
		}
		result[index] = task
	}
	return result
}

func (tasks SiteTasks) Supported() SiteTasks {
	result := make(SiteTasks, 0, len(tasks))
	for _, task := range tasks {
		if task.SupportsRule() {
			result = append(result, task)
		}
	}
	return result
}

func (task SiteTaskDefinition) SupportsRule() bool {
	return slices.Contains([]string{
		consts.SiteTaskRuleTypeTorrentPublished,
		consts.SiteTaskRuleTypeSeedDuration,
		consts.SiteTaskRuleTypeUploaded,
		consts.SiteTaskRuleTypeRoleLevelReached,
	}, task.Rule.Type)
}

func (tasks SiteTasks) Find(key string) *SiteTaskDefinition {
	key = strings.TrimSpace(key)
	for index := range tasks {
		if tasks[index].Key == key {
			return &tasks[index]
		}
	}
	return nil
}

func (tasks SiteTasks) Validate() error {
	seen := make(map[string]struct{}, len(tasks))
	for _, task := range tasks {
		if !siteTaskKeyPattern.MatchString(task.Key) {
			return gerror.New("task key must use lowercase letters, numbers, and underscores")
		}
		if _, exists := seen[task.Key]; exists {
			return gerror.New("task key must be unique")
		}
		seen[task.Key] = struct{}{}
		if !slices.Contains([]string{consts.SiteTaskCycleOnce, consts.SiteTaskCycleWeekly, consts.SiteTaskCycleMonthly}, task.Cycle) {
			return gerror.New("unsupported task cycle")
		}
		if !hasSiteTaskI18NValue(task.NameI18N) {
			return gerror.New("task name is required")
		}
		if err := task.validateRule(); err != nil {
			return err
		}
		if len(task.Rewards) == 0 {
			return gerror.New("task needs at least one reward")
		}
		if err := task.validateRewards(); err != nil {
			return err
		}
	}
	return nil
}

func (task SiteTaskDefinition) validateRule() error {
	if !task.SupportsRule() {
		return gerror.New("unsupported task rule type")
	}
	switch task.Rule.Type {
	case consts.SiteTaskRuleTypeTorrentPublished:
		if task.Rule.Target == 0 {
			return gerror.New("torrent published task target must be greater than zero")
		}
	case consts.SiteTaskRuleTypeSeedDuration, consts.SiteTaskRuleTypeUploaded:
		if task.Rule.Target == 0 {
			return gerror.New("tracker task target must be greater than zero")
		}
	case consts.SiteTaskRuleTypeRoleLevelReached:
		if task.Rule.Target == 0 {
			return gerror.New("role level task target must be greater than zero")
		}
	}
	return nil
}

func (task SiteTaskDefinition) validateRewards() error {
	seen := make(map[string]struct{}, len(task.Rewards))
	for _, reward := range task.Rewards {
		if _, exists := seen[reward.Type]; exists {
			return gerror.New("task reward type must be unique")
		}
		seen[reward.Type] = struct{}{}
		if reward.Amount <= 0 {
			return gerror.New("task reward amount must be greater than zero")
		}
		switch reward.Type {
		case consts.SiteTaskRewardTypeBonus:
		case consts.SiteTaskRewardTypeVip, consts.SiteTaskRewardTypeInvite:
			if math.Trunc(reward.Amount) != reward.Amount || reward.Amount > float64(math.MaxInt) {
				return gerror.New("VIP and invite reward amounts must be positive integers")
			}
		default:
			return gerror.New("unsupported task reward type")
		}
	}
	return nil
}

func normalizeSiteTaskI18N(values map[string]string) map[string]string {
	result := make(map[string]string, len(values))
	for locale, value := range values {
		if locale = strings.TrimSpace(locale); locale != "" {
			result[locale] = strings.TrimSpace(value)
		}
	}
	return result
}

func hasSiteTaskI18NValue(values map[string]string) bool {
	for _, value := range values {
		if value != "" {
			return true
		}
	}
	return false
}
