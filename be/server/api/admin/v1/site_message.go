package v1

import (
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"

	"github.com/gogf/gf/v2/frame/g"
)

type SiteMessageListReq struct {
	g.Meta `path:"/site/messages" method:"get" tags:"AdminSite" summary:"查询站内消息" perm:"admin:site/message:*"`
	sitein.AdminMessageListInp
}

type SiteMessageListRes struct {
	siteout.MessageListOut
}

type SiteMessageCreateReq struct {
	g.Meta `path:"/site/messages" method:"post" tags:"AdminSite" summary:"发送站内消息" perm:"admin:site/message:*"`
	sitein.AdminMessageCreateInp
}

type SiteMessageCreateRes struct {
	siteout.MessageCreateOut
}
