// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerPeer is the golang structure for table tracker_peer.
type TrackerPeer struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`
	TorrentId      uint64      `json:"torrentId"      orm:"torrent_id"      description:""`
	UserId         uint64      `json:"userId"         orm:"user_id"         description:""`
	PeerId         []byte      `json:"peerId"         orm:"peer_id"         description:""`
	Ipv4           string      `json:"ipv4"           orm:"ipv4"            description:""`
	Ipv6           string      `json:"ipv6"           orm:"ipv6"            description:""`
	Port           uint        `json:"port"           orm:"port"            description:""`
	Uploaded       uint64      `json:"uploaded"       orm:"uploaded"        description:""`
	Downloaded     uint64      `json:"downloaded"     orm:"downloaded"      description:""`
	Remaining      uint64      `json:"remaining"      orm:"remaining"       description:"剩余大小"`
	IsSeeder       bool        `json:"isSeeder"       orm:"is_seeder"       description:""`
	IsConnectable  bool        `json:"isConnectable"  orm:"is_connectable"  description:""`
	Agent          string      `json:"agent"          orm:"agent"           description:"BT 客户端"`
	Passkey        string      `json:"passkey"        orm:"passkey"         description:""`
	UploadOffset   uint64      `json:"uploadOffset"   orm:"upload_offset"   description:""`
	DownloadOffset uint64      `json:"downloadOffset" orm:"download_offset" description:""`
	StartedAt      *gtime.Time `json:"startedAt"      orm:"started_at"      description:""`
	LastAction     *gtime.Time `json:"lastAction"     orm:"last_action"     description:""`
	FinishedAt     *gtime.Time `json:"finishedAt"     orm:"finished_at"     description:""`
}
