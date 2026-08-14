package consts

const (
	SiteAuditLevelNormal    = 0
	SiteAuditLevelImportant = 1
	SiteAuditLevelCritical  = 2
)

const (
	SiteAuditActionCreate = "create"
	SiteAuditActionUpdate = "update"
	SiteAuditActionDelete = "delete"
)

const (
	SiteAuditOperationLock              = "lock"
	SiteAuditOperationUnlock            = "unlock"
	SiteAuditOperationPin               = "pin"
	SiteAuditOperationUnpin             = "unpin"
	SiteAuditOperationFeature           = "feature"
	SiteAuditOperationUnfeature         = "unfeature"
	SiteAuditOperationSetPromotion      = "set_promotion"
	SiteAuditOperationClearPromotion    = "clear_promotion"
	SiteAuditOperationMove              = "move"
	SiteAuditOperationGrant             = "grant"
	SiteAuditOperationRecycle           = "recycle"
	SiteAuditOperationGrantPermission   = "grant_permission"
	SiteAuditOperationRevokePermission  = "revoke_permission"
	SiteAuditOperationRemoveSession     = "remove_session"
	SiteAuditOperationUpdateStat        = "update_stat"
	SiteAuditOperationComplete          = "complete"
	SiteAuditOperationCancel            = "cancel"
	SiteAuditOperationResolve           = "resolve"
	SiteAuditOperationApply             = "apply"
	SiteAuditOperationRemove            = "remove"
	SiteAuditOperationApprove           = "approve"
	SiteAuditOperationReject            = "reject"
	SiteAuditOperationEdit              = "edit"
	SiteAuditOperationEnableTwoStep     = "enable_two_step"
	SiteAuditOperationDisableTwoStep    = "disable_two_step"
	SiteAuditOperationRenewTwoStepCodes = "renew_two_step_codes"
)

const (
	SiteAuditTargetTypeSiteConfig       = "site_config"
	SiteAuditTargetTypeSiteAnnouncement = "site_announcement"
	SiteAuditTargetTypeSiteMessage      = "site_message"
	SiteAuditTargetTypeCatalogCategory  = "catalog_category"
	SiteAuditTargetTypeCatalogTagGroup  = "catalog_tag_group"
	SiteAuditTargetTypeCatalogTag       = "catalog_tag"
	SiteAuditTargetTypeCatalogRequest   = "catalog_request"
	SiteAuditTargetTypeCatalogTorrent   = "catalog_torrent"
	SiteAuditTargetTypeForumCategory    = "forum_category"
	SiteAuditTargetTypeForumNode        = "forum_node"
	SiteAuditTargetTypeForumTopic       = "forum_topic"
	SiteAuditTargetTypeIamRole          = "iam_role"
	SiteAuditTargetTypeIamUser          = "iam_user"
	SiteAuditTargetTypeIamInvite        = "iam_invite"
	SiteAuditTargetTypeIamSession       = "iam_session"
	SiteAuditTargetTypeModReport        = "mod_report"
	SiteAuditTargetTypeModCheaterLog    = "mod_cheater_log"
	SiteAuditTargetTypeModUserLog       = "mod_user_log"
	SiteAuditTargetTypeModStaffMessage  = "mod_staff_message"
)
