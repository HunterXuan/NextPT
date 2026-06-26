package consts

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
	// SiteConfigIamDefaultRole 新注册用户的默认角色标识(字符串)
	SiteConfigIamDefaultRole = "iam.default_role"

	// SiteConfigIamDefaultRegisterRole 新注册用户的默认角色ID(整型)
	SiteConfigIamDefaultRegisterRole = "iam.default_register_role"

	// SiteConfigIamRegisterEnabled 是否开放全站开放注册
	SiteConfigIamRegisterEnabled = "iam.register_enabled"
)

// ==============================================================================
// 3. Catalog 资源系统相关配置 (Group: catalog)
// ==============================================================================
const (
	// SiteConfigCatalogTorrentSource 私有种子专属标识（注入到 Bencode 的 source 字段）
	SiteConfigCatalogTorrentSource = "catalog.torrent_source"
)

// ==============================================================================
// 后台业务配置的统一默认值字典
// 供 AdminConfig.GetByPath 取不到数据库值时兜底使用
// ==============================================================================
var SiteConfigDefaults = map[string]any{
	SiteConfigTrackerAnnounceInterval:    1800,
	SiteConfigTrackerAnnounceMinInterval: 900,
	SiteConfigTrackerUrl:                 "http://127.0.0.1:8000/api/tracker/announce",
	SiteConfigTrackerBonusT0:             8.0,
	SiteConfigTrackerBonusN0:             7.0,
	SiteConfigTrackerBonusB0:             100.0,
	SiteConfigTrackerBonusL:              300.0,
	SiteConfigTrackerBonusBase:           0.4,
	SiteConfigIamDefaultRole:             "user",
	SiteConfigIamDefaultRegisterRole:     1,
	SiteConfigIamRegisterEnabled:         true,
	SiteConfigCatalogTorrentSource:       "NextPT",
}
