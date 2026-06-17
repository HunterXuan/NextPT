package adminout

import "server/internal/model/entity"

type IamRoleListOut struct {
	Roles []entity.IamRole `json:"roles"`
}

type IamRoleCreateOut struct {
	Id uint `json:"id"`
}
