// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogCommentLike is the golang structure for table catalog_comment_like.
type CatalogCommentLike struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	CommentId uint64      `json:"commentId" orm:"comment_id" description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
