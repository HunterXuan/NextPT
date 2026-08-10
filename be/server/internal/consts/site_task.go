package consts

const (
	SiteTaskHistoryRetentionDays = 60

	SiteTaskCycleOnce    = "once"
	SiteTaskCycleWeekly  = "weekly"
	SiteTaskCycleMonthly = "monthly"
)

const (
	SiteTaskStatusActive = iota
	SiteTaskStatusCompleted
	SiteTaskStatusRewarded
	SiteTaskStatusExpired
)

const (
	SiteTaskRuleTypeTorrentPublished = "catalog.torrent_published"
	SiteTaskRuleTypeSeedDuration     = "tracker.seed_duration"
	SiteTaskRuleTypeUploaded         = "tracker.uploaded"
	SiteTaskRuleTypeRoleLevelReached = "iam.role_level_reached"
)

const (
	SiteTaskRewardTypeBonus  = "bonus"
	SiteTaskRewardTypeVip    = "vip"
	SiteTaskRewardTypeInvite = "invite"

	SiteTaskVipRemark = "site_task_reward"
)
