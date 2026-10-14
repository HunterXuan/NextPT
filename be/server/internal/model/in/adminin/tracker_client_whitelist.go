package adminin

import "server/internal/model"

type TrackerClientWhitelistFields struct {
	Family        string                         `json:"family" v:"required|length:1,60#{#admin.tracker.client_family_required}|{#admin.tracker.client_family_length}"`
	PeerIdPattern string                         `json:"peerIdPattern" v:"length:0,200#{#admin.tracker.client_peer_pattern_length}"`
	Exceptions    []model.TrackerClientException `json:"exceptions"`
	AgentPattern  string                         `json:"agentPattern" v:"length:0,200#{#admin.tracker.client_pattern_length}"`
	Enabled       bool                           `json:"enabled"`
	Comment       string                         `json:"comment" v:"length:0,200#{#admin.tracker.client_comment_length}"`
}

type TrackerClientWhitelistCreateInp struct{ TrackerClientWhitelistFields }
type TrackerClientWhitelistUpdateInp struct {
	Id            uint                            `json:"id" in:"path" v:"required|min:1"`
	Family        *string                         `json:"family" v:"length:1,60#{#admin.tracker.client_family_length}"`
	PeerIdPattern *string                         `json:"peerIdPattern" v:"length:0,200#{#admin.tracker.client_peer_pattern_length}"`
	Exceptions    *[]model.TrackerClientException `json:"exceptions"`
	AgentPattern  *string                         `json:"agentPattern" v:"length:0,200#{#admin.tracker.client_pattern_length}"`
	Enabled       *bool                           `json:"enabled"`
	Comment       *string                         `json:"comment" v:"length:0,200#{#admin.tracker.client_comment_length}"`
}
type TrackerClientWhitelistDeleteInp struct {
	Id uint `json:"id" in:"path" v:"required|min:1"`
}
type TrackerClientWhitelistCheckInp struct {
	PeerIdPattern string                         `json:"peerIdPattern" v:"length:0,200#{#admin.tracker.client_peer_pattern_length}"`
	Exceptions    []model.TrackerClientException `json:"exceptions"`
	AgentPattern  string                         `json:"agentPattern" v:"length:0,200#{#admin.tracker.client_pattern_length}"`
	PeerId        string                         `json:"peerId" v:"length:0,200"`
	UserAgent     string                         `json:"userAgent" v:"length:0,500"`
}
