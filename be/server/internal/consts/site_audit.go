package consts

const (
	SiteAuditLevelNormal    = 0
	SiteAuditLevelImportant = 1
	SiteAuditLevelCritical  = 2
)

const (
	SiteAuditActionCreate           = "create"
	SiteAuditActionUpdate           = "update"
	SiteAuditActionDelete           = "delete"
	SiteAuditActionLock             = "lock"
	SiteAuditActionUnlock           = "unlock"
	SiteAuditActionPin              = "pin"
	SiteAuditActionUnpin            = "unpin"
	SiteAuditActionMove             = "move"
	SiteAuditActionGrant            = "grant"
	SiteAuditActionRecycle          = "recycle"
	SiteAuditActionGrantPermission  = "grant_permission"
	SiteAuditActionRevokePermission = "revoke_permission"
	SiteAuditActionRemoveSession    = "remove_session"
	SiteAuditActionUpdateStat       = "update_stat"
	SiteAuditActionResolve          = "resolve"
	SiteAuditActionApply            = "apply"
	SiteAuditActionRemove           = "remove"
)

const (
	SiteAuditTargetTypeSiteConfig        = "site_config"
	SiteAuditTargetTypeCatalogCategory   = "catalog_category"
	SiteAuditTargetTypeCatalogTorrent    = "catalog_torrent"
	SiteAuditTargetTypeForumCategory     = "forum_category"
	SiteAuditTargetTypeForumNode         = "forum_node"
	SiteAuditTargetTypeForumTopic        = "forum_topic"
	SiteAuditTargetTypeIamRole           = "iam_role"
	SiteAuditTargetTypeIamUser           = "iam_user"
	SiteAuditTargetTypeIamInvite         = "iam_invite"
	SiteAuditTargetTypeIamUserPermission = "iam_user_permission"
	SiteAuditTargetTypeIamSession        = "iam_session"
	SiteAuditTargetTypeModReport         = "mod_report"
	SiteAuditTargetTypeModCheaterLog     = "mod_cheater_log"
	SiteAuditTargetTypeModUserLog        = "mod_user_log"
)
