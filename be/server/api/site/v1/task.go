package v1

import (
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"

	"github.com/gogf/gf/v2/frame/g"
)

type TaskListReq struct {
	g.Meta `path:"/tasks" method:"get" tags:"Site" summary:"获取用户任务" noPerm:"true"`
}

type TaskListRes struct{ siteout.TaskListOut }

type TaskClaimReq struct {
	g.Meta `path:"/tasks/{key}:claim" method:"post" tags:"Site" summary:"认领用户任务" noPerm:"true"`
	sitein.TaskClaimInp
}

type TaskClaimRes struct{ siteout.TaskClaimOut }

type UserTaskRewardClaimReq struct {
	g.Meta `path:"/user-tasks/{id}:claimReward" method:"post" tags:"Site" summary:"领取任务奖励" noPerm:"true"`
	sitein.UserTaskRewardClaimInp
}

type UserTaskRewardClaimRes struct{}
