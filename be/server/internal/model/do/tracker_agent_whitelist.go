// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerAgentWhitelist is the golang structure of table tracker_agent_whitelist for DAO operations like Where/Data.
type TrackerAgentWhitelist struct {
	g.Meta       `orm:"table:tracker_agent_whitelist, do:true"`
	Id           any         //
	Family       any         // 客户端家族名
	PeerIdPrefix any         //
	AgentPattern any         //
	MinVersion   any         //
	MaxVersion   any         //
	AllowHttps   any         //
	Enabled      any         //
	Comment      any         //
	Hits         any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
