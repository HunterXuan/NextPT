package v1

import (
	"server/internal/model/out/forumout"

	"github.com/gogf/gf/v2/frame/g"
)

type NodeListReq struct {
	g.Meta `path:"/nodes" method:"get" tags:"Forum" summary:"获取社区节点与分类列表"`
}

type NodeListRes struct {
	forumout.NodeListOut
}
