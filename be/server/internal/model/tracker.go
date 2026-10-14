package model

import (
	libtracker "server/internal/library/tracker"

	"github.com/gogf/gf/v2/os/gtime"
)

type TrackerClientException = libtracker.ClientException

type TrackerClientWhitelist struct {
	Id            uint                     `json:"id" orm:"id"`
	Family        string                   `json:"family" orm:"family"`
	PeerIdPattern string                   `json:"peerIdPattern" orm:"peer_id_pattern"`
	AgentPattern  string                   `json:"agentPattern" orm:"agent_pattern"`
	Exceptions    []TrackerClientException `json:"exceptions" orm:"exceptions"`
	Enabled       bool                     `json:"enabled" orm:"enabled"`
	Comment       string                   `json:"comment" orm:"comment"`
	CreatedAt     *gtime.Time              `json:"createdAt" orm:"created_at"`
	UpdatedAt     *gtime.Time              `json:"updatedAt" orm:"updated_at"`
}
