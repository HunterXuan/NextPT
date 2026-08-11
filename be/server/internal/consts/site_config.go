package consts

type SiteConfigValueType string

const (
	SiteConfigValueTypeString  SiteConfigValueType = "string"
	SiteConfigValueTypeInt     SiteConfigValueType = "int"
	SiteConfigValueTypeFloat   SiteConfigValueType = "float"
	SiteConfigValueTypeBoolean SiteConfigValueType = "boolean"
	SiteConfigValueTypeJSON    SiteConfigValueType = "json"
)

// ==============================================================================
// 1. Tracker 相关系统配置 (Group: tracker)
// 格式: group.key
// ==============================================================================
const (
	// SiteConfigTrackerAnnounceInterval 客户端心跳汇报间隔（秒）
	SiteConfigTrackerAnnounceInterval = "tracker.announce_interval"

	// SiteConfigTrackerAnnounceMinInterval 客户端心跳最小汇报间隔（秒）
	SiteConfigTrackerAnnounceMinInterval = "tracker.announce_min_interval"

	// SiteConfigTrackerUrl 默认的 Tracker 宣告地址 (Passkey 注入前缀)
	SiteConfigTrackerUrl = "tracker.announce_url"

	// --------------------------------------------------------------------------
	// 魔力值计算公式参数 (NexusPHP 还原)
	// --------------------------------------------------------------------------

	// SiteConfigTrackerBonusT0 衰减参数 T0（单位：周）
	SiteConfigTrackerBonusT0 = "tracker.bonus_T0"

	// SiteConfigTrackerBonusN0 做种人数拥挤惩罚基数 N0
	SiteConfigTrackerBonusN0 = "tracker.bonus_N0"

	// SiteConfigTrackerBonusB0 用户每小时通过体积做种获得魔力的硬上限 B0
	SiteConfigTrackerBonusB0 = "tracker.bonus_B0"

	// SiteConfigTrackerBonusL 收益收敛平滑参数 L
	SiteConfigTrackerBonusL = "tracker.bonus_L"

	// SiteConfigTrackerBonusBase 每个达标种子（>=1GB）的基础每小时奖励
	SiteConfigTrackerBonusBase = "tracker.bonus_base"
)

// ==============================================================================
// 2. IAM 用户系统相关配置 (Group: iam)
// ==============================================================================
const (
	// SiteConfigIamDefaultRegisterRole 新注册用户的默认角色ID(整型)
	SiteConfigIamDefaultRegisterRole = "iam.default_register_role"

	// SiteConfigIamRegisterEnabled 是否开放全站开放注册
	SiteConfigIamRegisterEnabled = "iam.register_enabled"

	// SiteConfigIamInviteBypassEmailPattern 邀请制下可免邀请码注册的邮箱正则
	SiteConfigIamInviteBypassEmailPattern = "iam.invite_bypass_email_pattern"
)

// ==============================================================================
// 3. Catalog 资源系统相关配置 (Group: catalog)
// ==============================================================================
const (
	// SiteConfigCatalogTorrentSource 私有种子专属标识（注入到 Bencode 的 source 字段）
	SiteConfigCatalogTorrentSource = "catalog.torrent_source"

	// SiteConfigCatalogTorrentDirectPublishLevel 允许免审核直接发布种子的最低角色等级
	SiteConfigCatalogTorrentDirectPublishLevel = "catalog.torrent_direct_publish_level"

	// SiteConfigCatalogGlobalPromotion 全站种子优惠配置
	SiteConfigCatalogGlobalPromotion = "catalog.global_promotion"

	// SiteConfigCatalogNewTorrentPromotion 新发布种子自动优惠配置
	SiteConfigCatalogNewTorrentPromotion = "catalog.new_torrent_promotion"
)

const SiteConfigEconomyShopProducts = "economy.shop_products"

const (
	SiteConfigSiteMaintenanceEnabled = "site.maintenance_enabled"
	SiteConfigSiteMaintenanceMessage = "site.maintenance_message"
	SiteConfigSiteAdvertisements     = "site.advertisements"
	SiteConfigSiteTasks              = "site.tasks"
)

const (
	SiteAdvertisementPlacementHome        = "home"
	SiteAdvertisementPlacementCatalogList = "catalog_list"
	SiteAdvertisementPlacementForumList   = "forum_list"
)

var SiteAdvertisementPlacements = []string{
	SiteAdvertisementPlacementHome,
	SiteAdvertisementPlacementCatalogList,
	SiteAdvertisementPlacementForumList,
}

// ==============================================================================
// 后台业务配置的统一默认值字典
// 供 AdminConfig.GetByPath 取不到数据库值时兜底使用
// ==============================================================================
var SiteConfigDefaults = map[string]any{
	SiteConfigTrackerAnnounceInterval:     1800,
	SiteConfigTrackerAnnounceMinInterval:  900,
	SiteConfigTrackerUrl:                  "http://127.0.0.1:8000/api/tracker/announce",
	SiteConfigTrackerBonusT0:              8.0,
	SiteConfigTrackerBonusN0:              7.0,
	SiteConfigTrackerBonusB0:              100.0,
	SiteConfigTrackerBonusL:               300.0,
	SiteConfigTrackerBonusBase:            0.4,
	SiteConfigIamDefaultRegisterRole:      2,
	SiteConfigIamRegisterEnabled:          true,
	SiteConfigIamInviteBypassEmailPattern: "",
	SiteConfigSiteMaintenanceEnabled:      false,
	SiteConfigSiteMaintenanceMessage:      "",
	SiteConfigSiteTasks:                   []map[string]any{},
	SiteConfigSiteAdvertisements: map[string]any{
		SiteAdvertisementPlacementHome: map[string]any{
			"enabled": false,
			"title":   "",
			"image":   "",
			"url":     "",
		},
		SiteAdvertisementPlacementCatalogList: map[string]any{
			"enabled": false,
			"title":   "",
			"image":   "",
			"url":     "",
		},
		SiteAdvertisementPlacementForumList: map[string]any{
			"enabled": false,
			"title":   "",
			"image":   "",
			"url":     "",
		},
	},
	SiteConfigCatalogTorrentSource:             "NextPT",
	SiteConfigCatalogTorrentDirectPublishLevel: 20,
	SiteConfigCatalogGlobalPromotion: map[string]any{
		"enabled":  false,
		"state":    ResourceTorrentPromotionStateFree,
		"expireAt": "",
	},
	SiteConfigCatalogNewTorrentPromotion: map[string]any{
		"enabled": true,
		"rules": []map[string]any{
			{
				"minGiB":        0,
				"durationHours": 72,
				"options": []map[string]any{
					{"state": ResourceTorrentPromotionStateNormal, "weight": 70},
					{"state": ResourceTorrentPromotionStateFree, "weight": 10},
					{"state": ResourceTorrentPromotionState2x, "weight": 10},
					{"state": ResourceTorrentPromotionState50Percent, "weight": 10},
				},
			},
			{
				"minGiB":        10,
				"durationHours": 96,
				"options": []map[string]any{
					{"state": ResourceTorrentPromotionStateNormal, "weight": 50},
					{"state": ResourceTorrentPromotionStateFree, "weight": 20},
					{"state": ResourceTorrentPromotionState2x, "weight": 20},
					{"state": ResourceTorrentPromotionState2xFree, "weight": 10},
				},
			},
			{
				"minGiB":        50,
				"durationHours": 168,
				"options": []map[string]any{
					{"state": ResourceTorrentPromotionStateFree, "weight": 40},
					{"state": ResourceTorrentPromotionState2xFree, "weight": 30},
					{"state": ResourceTorrentPromotionState2x50Percent, "weight": 30},
				},
			},
		},
	},
	SiteConfigEconomyShopProducts: []map[string]any{
		{
			"key":       EconomyShopProductTypeInvite,
			"type":      EconomyShopProductTypeInvite,
			"enabled":   true,
			"price":     1000.0,
			"sortOrder": 10,
			"options":   map[string]any{"amount": 1},
		},
		{
			"key":       "upload_100_gib",
			"type":      EconomyShopProductTypeUpload,
			"enabled":   true,
			"price":     500.0,
			"sortOrder": 20,
			"options":   map[string]any{"amountGiB": 100},
		},
		{
			"key":       "download_50_gib",
			"type":      EconomyShopProductTypeDownload,
			"enabled":   true,
			"price":     800.0,
			"sortOrder": 30,
			"options":   map[string]any{"amountGiB": 50},
		},
		{
			"key":       "vip_30d",
			"type":      EconomyShopProductTypeVip,
			"enabled":   true,
			"price":     3000.0,
			"sortOrder": 40,
			"options":   map[string]any{"durationDays": 30},
		},
	},
}
