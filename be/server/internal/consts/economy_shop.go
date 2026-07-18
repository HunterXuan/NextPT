package consts

const (
	EconomyShopProductTypeInvite   = "invite"
	EconomyShopProductTypeVip      = "vip"
	EconomyShopProductTypeUpload   = "upload"
	EconomyShopProductTypeDownload = "download"
)

var EconomyShopProductTypes = []string{
	EconomyShopProductTypeInvite,
	EconomyShopProductTypeVip,
	EconomyShopProductTypeUpload,
	EconomyShopProductTypeDownload,
}

const (
	EconomyShopOrderStatusPending   = 0
	EconomyShopOrderStatusCompleted = 1
)

const (
	EconomyShopOrderTargetTypeIamInvite   = "iam_invite"
	EconomyShopOrderTargetTypeIamUser     = "iam_user"
	EconomyShopOrderTargetTypeIamUserStat = "iam_user_stat"
)

const EconomyShopVipRemark = "economy_shop"
