// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteConfig is the golang structure for table site_config.
type SiteConfig struct {
	Id        uint        `json:"id"        orm:"id"         description:""`
	Group     string      `json:"group"     orm:"group"      description:""`
	Key       string      `json:"key"       orm:"key"        description:""`
	Value     *gjson.Json `json:"value"     orm:"value"      description:"支持存 boolean/number/array"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
