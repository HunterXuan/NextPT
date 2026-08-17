package consts

const (
	IamUserPermissionSourceManual  = 1
	IamUserPermissionSourceUserMod = 2
)

const (
	// Format: {act}:{domain}/{resource}:{id}
	// Act: read/create/update/download/admin
	// Id: could be number or *

	// All permissions wildcard. It can be granted as "*" or denied as "-*".
	IamPermissionAll = "*"

	// IAM Management
	IamPermissionAdminIamUserManage   = "admin:iam/user:*"
	IamPermissionAdminIamRoleManage   = "admin:iam/role:*"
	IamPermissionAdminIamInviteManage = "admin:iam/invite:*"

	// Mod Management
	IamPermissionAdminModReportManage  = "admin:mod/report:*"
	IamPermissionAdminModCheaterManage = "admin:mod/cheater:*"
	IamPermissionAdminModUserManage    = "admin:mod/user:*"
	IamPermissionAdminModStaffMessage  = "admin:mod/staff-message:*"

	// Site Management
	IamPermissionAdminSiteConfig       = "admin:site/config:*"
	IamPermissionAdminSiteAudit        = "admin:site/audit:*"
	IamPermissionAdminSiteAnnouncement = "admin:site/announcement:*"
	IamPermissionAdminSiteMessage      = "admin:site/message:*"

	// Sys Management
	IamPermissionAdminSysCronManage = "admin:sys/cron:*"

	// Forum Management
	IamPermissionAdminForumCategoryManage = "admin:forum/category:*"
	IamPermissionAdminForumNodeManage     = "admin:forum/node:*"
	IamPermissionAdminForumTopicManage    = "admin:forum/topic:*"
	IamPermissionAdminForumReplyManage    = "admin:forum/reply:*"

	// Catalog Management
	IamPermissionAdminCatalogCategoryManage = "admin:catalog/category:*"
	IamPermissionAdminCatalogTagManage      = "admin:catalog/tag:*"
	IamPermissionAdminCatalogTorrentManage  = "admin:catalog/torrent:*"
	IamPermissionAdminCatalogSubtitleManage = "admin:catalog/subtitle:*"
	IamPermissionAdminCatalogCommentManage  = "admin:catalog/comment:*"
	IamPermissionAdminCatalogRequestManage  = "admin:catalog/request:*"

	// IAM
	IamPermissionIamInviteRead   = "read:iam/invite:*"
	IamPermissionIamInviteCreate = "create:iam/invite:*"

	// Economy
	IamPermissionEconomyShopProductRead = "read:economy/shop-product:*"
	IamPermissionEconomyShopOrderRead   = "read:economy/shop-order:*"
	IamPermissionEconomyShopOrderCreate = "create:economy/shop-order:*"

	// Site
	IamPermissionSiteAnnouncementRead  = "read:site/announcement:*"
	IamPermissionSiteMessageRead       = "read:site/message:*"
	IamPermissionSiteChatMessageRead   = "read:site/chat-message:*"
	IamPermissionSiteChatMessageCreate = "create:site/chat-message:*"
	IamPermissionModStaffMessageRead   = "read:mod/staff-message:*"
	IamPermissionModStaffMessageCreate = "create:mod/staff-message:*"

	// Forum
	IamPermissionForumTopicRead   = "read:forum/topic:*"
	IamPermissionForumTopicCreate = "create:forum/topic:*"
	IamPermissionForumTopicUpdate = "update:forum/topic:*"
	IamPermissionForumReplyRead   = "read:forum/reply:*"
	IamPermissionForumReplyCreate = "create:forum/reply:*"
	IamPermissionForumReplyUpdate = "update:forum/reply:*"

	// Catalog
	IamPermissionCatalogTorrentRead      = "read:catalog/torrent:*"
	IamPermissionCatalogTorrentCreate    = "create:catalog/torrent:*"
	IamPermissionCatalogTorrentDownload  = "download:catalog/torrent:*"
	IamPermissionCatalogSubtitleRead     = "read:catalog/subtitle:*"
	IamPermissionCatalogSubtitleCreate   = "create:catalog/subtitle:*"
	IamPermissionCatalogSubtitleDownload = "download:catalog/subtitle:*"
	IamPermissionCatalogCommentRead      = "read:catalog/comment:*"
	IamPermissionCatalogCommentCreate    = "create:catalog/comment:*"
	IamPermissionCatalogRequestRead      = "read:catalog/request:*"
	IamPermissionCatalogRequestCreate    = "create:catalog/request:*"
	IamPermissionCatalogRequestUpdate    = "update:catalog/request:*"
)

var (
	IamPermissionList = []string{
		IamPermissionAll,

		IamPermissionAdminIamUserManage,
		IamPermissionAdminIamRoleManage,
		IamPermissionAdminIamInviteManage,

		IamPermissionAdminModReportManage,
		IamPermissionAdminModCheaterManage,
		IamPermissionAdminModUserManage,
		IamPermissionAdminModStaffMessage,

		IamPermissionAdminSiteConfig,
		IamPermissionAdminSiteAudit,
		IamPermissionAdminSiteAnnouncement,
		IamPermissionAdminSiteMessage,

		IamPermissionAdminSysCronManage,

		IamPermissionAdminForumCategoryManage,
		IamPermissionAdminForumNodeManage,
		IamPermissionAdminForumTopicManage,
		IamPermissionAdminForumReplyManage,

		IamPermissionAdminCatalogCategoryManage,
		IamPermissionAdminCatalogTagManage,
		IamPermissionAdminCatalogTorrentManage,
		IamPermissionAdminCatalogSubtitleManage,
		IamPermissionAdminCatalogCommentManage,
		IamPermissionAdminCatalogRequestManage,

		IamPermissionIamInviteRead,
		IamPermissionIamInviteCreate,

		IamPermissionEconomyShopProductRead,
		IamPermissionEconomyShopOrderRead,
		IamPermissionEconomyShopOrderCreate,

		IamPermissionSiteAnnouncementRead,
		IamPermissionSiteMessageRead,
		IamPermissionSiteChatMessageRead,
		IamPermissionSiteChatMessageCreate,
		IamPermissionModStaffMessageRead,
		IamPermissionModStaffMessageCreate,

		IamPermissionForumTopicRead,
		IamPermissionForumTopicCreate,
		IamPermissionForumTopicUpdate,
		IamPermissionForumReplyRead,
		IamPermissionForumReplyCreate,
		IamPermissionForumReplyUpdate,

		IamPermissionCatalogTorrentRead,
		IamPermissionCatalogTorrentCreate,
		IamPermissionCatalogTorrentDownload,
		IamPermissionCatalogSubtitleRead,
		IamPermissionCatalogSubtitleCreate,
		IamPermissionCatalogSubtitleDownload,
		IamPermissionCatalogCommentRead,
		IamPermissionCatalogCommentCreate,
		IamPermissionCatalogRequestRead,
		IamPermissionCatalogRequestCreate,
		IamPermissionCatalogRequestUpdate,
	}
)
