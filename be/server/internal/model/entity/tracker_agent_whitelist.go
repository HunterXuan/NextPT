// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerAgentWhitelist is the golang structure for table tracker_agent_whitelist.
type TrackerAgentWhitelist struct {
	Id            uint        `json:"id"            orm:"id"              description:""`
	Family        string      `json:"family"        orm:"family"          description:"客户端家族名"`
	PeerIdPattern string      `json:"peerIdPattern" orm:"peer_id_pattern" description:""`
	AgentPattern  string      `json:"agentPattern"  orm:"agent_pattern"   description:""`
	Exceptions    *gjson.Json `json:"exceptions"    orm:"exceptions"      description:""`
	Enabled       bool        `json:"enabled"       orm:"enabled"         description:""`
	Comment       string      `json:"comment"       orm:"comment"         description:""`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:""`
}
