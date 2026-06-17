package adminout

import "server/internal/model/entity"

type ForumCategoryListOut struct {
	Categories []entity.ForumCategory `json:"categories"`
}
