package adminout

import "server/internal/model/entity"

type ForumNodeListOut struct {
	Nodes []entity.ForumNode `json:"nodes"`
}
