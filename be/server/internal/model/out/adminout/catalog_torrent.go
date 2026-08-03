package adminout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type CatalogTorrentReviewItem struct {
	Id          uint64               `json:"id"`
	Name        string               `json:"name"`
	SubTitle    string               `json:"subTitle"`
	CategoryId  uint                 `json:"categoryId"`
	Owner       model.IamUserSummary `json:"owner"`
	Size        uint64               `json:"size"`
	FileCount   uint                 `json:"fileCount"`
	Status      int                  `json:"status"`
	SubmittedAt *gtime.Time          `json:"submittedAt"`
	PublishedAt *gtime.Time          `json:"publishedAt"`
	ReviewedBy  uint64               `json:"reviewedBy"`
	ReviewedAt  *gtime.Time          `json:"reviewedAt"`
	Comment     string               `json:"comment"`
	CreatedAt   *gtime.Time          `json:"createdAt"`
}

type CatalogTorrentReviewListOut struct {
	List  []CatalogTorrentReviewItem `json:"list"`
	Total int                        `json:"total"`
	Page  int                        `json:"page"`
	Size  int                        `json:"size"`
}
