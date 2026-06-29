// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrent is the golang structure for table catalog_torrent.
type CatalogTorrent struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`
	InfoHash       []byte      `json:"infoHash"       orm:"info_hash"       description:"BitTorrent info_hash"`
	Name           string      `json:"name"           orm:"name"            description:"种子标题"`
	SubTitle       string      `json:"subTitle"       orm:"sub_title"       description:"副标题"`
	CategoryId     uint        `json:"categoryId"     orm:"category_id"     description:""`
	Description    string      `json:"description"    orm:"description"     description:"详情描述 (BBCode/Markdown)"`
	Nfo            []byte      `json:"nfo"            orm:"nfo"             description:"NFO 文件内容"`
	FileName       string      `json:"fileName"       orm:"file_name"       description:"种子文件名"`
	Size           uint64      `json:"size"           orm:"size"            description:"总大小 (bytes)"`
	FileCount      uint        `json:"fileCount"      orm:"file_count"      description:"文件数量"`
	OwnerId        uint64      `json:"ownerId"        orm:"owner_id"        description:"上传者"`
	Anonymous      bool        `json:"anonymous"      orm:"anonymous"       description:"匿名上传"`
	SpState        int         `json:"spState"        orm:"sp_state"        description:"0=normal 1=free 2=2x 3=2xfree 4=50%off 5=2x50% 6=30%off 7=custom"`
	SpExpireAt     *gtime.Time `json:"spExpireAt"     orm:"sp_expire_at"    description:"促销到期时间"`
	IsFeatured     bool        `json:"isFeatured"     orm:"is_featured"     description:"是否推荐"`
	IsPinned       bool        `json:"isPinned"       orm:"is_pinned"       description:"是否置顶"`
	PinWeight      int         `json:"pinWeight"      orm:"pin_weight"      description:"置顶权重 (越大越靠前)"`
	Seeders        uint        `json:"seeders"        orm:"seeders"         description:""`
	Leechers       uint        `json:"leechers"       orm:"leechers"        description:""`
	TimesCompleted uint        `json:"timesCompleted" orm:"times_completed" description:""`
	CommentsCount  uint        `json:"commentsCount"  orm:"comments_count"  description:""`
	Views          uint        `json:"views"          orm:"views"           description:""`
	LikeCount      uint        `json:"likeCount"      orm:"like_count"      description:""`
	RewardsCount   uint        `json:"rewardsCount"   orm:"rewards_count"   description:"收到的赞赏次数"`
	RewardsAmount  float64     `json:"rewardsAmount"  orm:"rewards_amount"  description:"收到的赞赏总金额(Bonus)"`
	Visible        bool        `json:"visible"        orm:"visible"         description:""`
	Banned         bool        `json:"banned"         orm:"banned"          description:""`
	LastAction     *gtime.Time `json:"lastAction"     orm:"last_action"     description:"Tracker 最后活动时间"`
	LastReseed     *gtime.Time `json:"lastReseed"     orm:"last_reseed"     description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""`
}
