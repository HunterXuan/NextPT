// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TrackerPeer is the golang structure of table tracker_peer for DAO operations like Where/Data.
type TrackerPeer struct {
	g.Meta         `orm:"table:tracker_peer, do:true"`
	Id             any         //
	TorrentId      any         //
	UserId         any         //
	PeerId         []byte      //
	Ipv4           any         //
	Ipv6           any         //
	Port           any         //
	Uploaded       any         //
	Downloaded     any         //
	Remaining      any         // 剩余大小
	IsSeeder       any         //
	IsConnectable  any         //
	Agent          any         // BT 客户端
	Passkey        any         //
	UploadOffset   any         //
	DownloadOffset any         //
	StartedAt      *gtime.Time //
	LastAction     *gtime.Time //
	FinishedAt     *gtime.Time //
}
