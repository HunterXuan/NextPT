// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogCommentLike is the golang structure of table catalog_comment_like for DAO operations like Where/Data.
type CatalogCommentLike struct {
	g.Meta    `orm:"table:catalog_comment_like, do:true"`
	Id        any         //
	UserId    any         //
	CommentId any         //
	CreatedAt *gtime.Time //
}
