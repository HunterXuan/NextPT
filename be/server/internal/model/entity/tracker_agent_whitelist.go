// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerAgentWhitelist is the golang structure for table tracker_agent_whitelist.
type TrackerAgentWhitelist struct {
	Id           uint        `json:"id"           orm:"id"             description:""`
	Family       string      `json:"family"       orm:"family"         description:"客户端家族名"`
	PeerIdPrefix string      `json:"peerIdPrefix" orm:"peer_id_prefix" description:""`
	AgentPattern string      `json:"agentPattern" orm:"agent_pattern"  description:""`
	MinVersion   string      `json:"minVersion"   orm:"min_version"    description:""`
	MaxVersion   string      `json:"maxVersion"   orm:"max_version"    description:""`
	AllowHttps   bool        `json:"allowHttps"   orm:"allow_https"    description:""`
	Enabled      bool        `json:"enabled"      orm:"enabled"        description:""`
	Comment      string      `json:"comment"      orm:"comment"        description:""`
	Hits         uint        `json:"hits"         orm:"hits"           description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"     description:""`
}
