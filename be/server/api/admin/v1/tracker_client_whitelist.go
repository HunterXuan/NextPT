package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type TrackerClientWhitelistListReq struct {
	g.Meta `path:"/tracker/client-whitelists" method:"get" tags:"AdminTracker" summary:"客户端白名单列表" perm:"admin:tracker/client-whitelist:*"`
}
type TrackerClientWhitelistListRes struct {
	adminout.TrackerClientWhitelistListOut
}

type TrackerClientWhitelistCreateReq struct {
	g.Meta `path:"/tracker/client-whitelists" method:"post" tags:"AdminTracker" summary:"创建客户端白名单" perm:"admin:tracker/client-whitelist:*"`
	adminin.TrackerClientWhitelistCreateInp
}
type TrackerClientWhitelistCreateRes struct {
	adminout.TrackerClientWhitelistCreateOut
}

type TrackerClientWhitelistUpdateReq struct {
	g.Meta `path:"/tracker/client-whitelists/{id}" method:"patch" tags:"AdminTracker" summary:"更新客户端白名单" perm:"admin:tracker/client-whitelist:*"`
	adminin.TrackerClientWhitelistUpdateInp
}
type TrackerClientWhitelistUpdateRes struct{}

type TrackerClientWhitelistDeleteReq struct {
	g.Meta `path:"/tracker/client-whitelists/{id}" method:"delete" tags:"AdminTracker" summary:"删除客户端白名单" perm:"admin:tracker/client-whitelist:*"`
	adminin.TrackerClientWhitelistDeleteInp
}
type TrackerClientWhitelistDeleteRes struct{}

type TrackerClientWhitelistCheckReq struct {
	g.Meta `path:"/tracker/client-whitelists:check" method:"post" tags:"AdminTracker" summary:"测试客户端匹配规则" perm:"admin:tracker/client-whitelist:*"`
	adminin.TrackerClientWhitelistCheckInp
}
type TrackerClientWhitelistCheckRes struct {
	adminout.TrackerClientWhitelistCheckOut
}
