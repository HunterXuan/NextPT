package consts

const (
	EconomyShopProductTypeInvite   = "invite"
	EconomyShopProductTypeVip      = "vip"
	EconomyShopProductTypeUpload   = "upload"
	EconomyShopProductTypeDownload = "download"
	EconomyShopProductTypeCoupon   = "download_coupon"
)

var EconomyShopProductTypes = []string{
	EconomyShopProductTypeInvite,
	EconomyShopProductTypeVip,
	EconomyShopProductTypeUpload,
	EconomyShopProductTypeDownload,
	EconomyShopProductTypeCoupon,
}

var EconomyShopImplementedProductTypes = []string{
	EconomyShopProductTypeInvite,
}

const (
	EconomyShopOrderStatusPending   = 0
	EconomyShopOrderStatusCompleted = 1
)

const (
	EconomyShopOrderTargetTypeIamInvite = "iam_invite"
)
