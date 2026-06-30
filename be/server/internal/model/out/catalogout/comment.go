package catalogout

import "server/internal/model"

type CommentListItem struct {
	Id          uint64            `json:"id"`
	Author      model.UserSummary `json:"author"`
	Content     string            `json:"content"`
	LikeCount   uint              `json:"likeCount"`
	RewardCount uint              `json:"rewardCount"`
	CreatedAt   string            `json:"createdAt"`
	IsLiked     bool              `json:"isLiked"`
}

type CommentListOut struct {
	List  []CommentListItem `json:"list"`
	Total int               `json:"total"`
}

type CommentToggleLikeOut struct {
	IsLiked bool `json:"isLiked"`
}

type CommentCreateOut struct {
	Id uint64 `json:"id"`
}
