package adminout

import "server/internal/model/entity"

type IamUserListOut struct {
	Users []*entity.IamUser `json:"users"`
	Total int               `json:"total"`
}

type IamUserStatDetailOut struct {
	entity.IamUserStat
}
