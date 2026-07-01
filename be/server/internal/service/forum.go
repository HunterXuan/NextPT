// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/forumin"
	"server/internal/model/out/forumout"
)

type (
	IForumCategoryDomain interface {
		AdminCreateCategory(ctx context.Context, nameI18N string, descI18N string, sortOrder int, minRoleView int) error
		AdminUpdateCategory(ctx context.Context, id uint, nameI18N *string, descI18N *string, sortOrder *int, minRoleView *int) error
		AdminDeleteCategory(ctx context.Context, id uint) (int, error)
		AdminListCategories(ctx context.Context) ([]entity.ForumCategory, error)
	}
	IForumNodeDomain interface {
		GetCategories(ctx context.Context) ([]entity.ForumCategory, error)
		GetNodes(ctx context.Context) ([]entity.ForumNode, error)
		CheckNodeReadPolicy(ctx context.Context, actor *model.Actor, node *entity.ForumNode) error
		CheckNodeWritePolicy(ctx context.Context, actor *model.Actor, node *entity.ForumNode) error
		CheckNodeCreatePolicy(ctx context.Context, actor *model.Actor, node *entity.ForumNode) error
		UpdateStats(ctx context.Context, nodeId uint, topicDelta int, replyDelta int) error
		AdminCreateNode(ctx context.Context, categoryId uint, slug string, nameI18N string, descI18N string, sortOrder int, minRoleRead int, minRoleWrite int, minRoleCreate int, moderators []byte) error
		AdminUpdateNode(ctx context.Context, id uint, categoryId *uint, slug *string, nameI18N *string, descI18N *string, sortOrder *int, minRoleRead *int, minRoleWrite *int, minRoleCreate *int, moderators []byte) error
		AdminDeleteNode(ctx context.Context, id uint) (int, error)
		AdminListNodes(ctx context.Context) ([]entity.ForumNode, error)
		GetNodeBySlug(ctx context.Context, slug string) (*entity.ForumNode, error)
		GetNodeById(ctx context.Context, id uint) (*entity.ForumNode, error)
		IncrementNodeStats(ctx context.Context, nodeId uint, topicId uint64) error
	}
	IForumNodeUsecase interface {
		List(ctx context.Context, actor *model.Actor, in forumin.NodeListInp) (*forumout.NodeListOut, error)
	}
	IForumReplyDomain interface {
		InsertReply(ctx context.Context, actor *model.Actor, in forumin.ReplyCreateInp) (uint64, error)
		QueryRepliesByTopic(ctx context.Context, topicId uint64, page int, size int) ([]entity.ForumReply, int, error)
		GetReplyById(ctx context.Context, replyId uint64) (*entity.ForumReply, error)
		ToggleLike(ctx context.Context, actor *model.Actor, replyId uint64) (bool, error)
		GetReplyLikesByUser(ctx context.Context, userId uint64, replyIds []uint64) ([]entity.ForumReplyLike, error)
	}
	IForumReplyUsecase interface {
		List(ctx context.Context, actor *model.Actor, in forumin.ReplyListInp) (*forumout.ReplyListOut, error)
		Create(ctx context.Context, actor *model.Actor, in forumin.ReplyCreateInp) (uint64, error)
		ToggleReplyLike(ctx context.Context, actor *model.Actor, in forumin.ReplyToggleLikeInp) error
		RewardReply(ctx context.Context, actor *model.Actor, in forumin.ReplyRewardInp) error
		ReportReply(ctx context.Context, actor *model.Actor, in forumin.ReplyReportInp) error
	}
	IForumTopicDomain interface {
		GetTopicById(ctx context.Context, topicId uint64) (*entity.ForumTopic, error)
		CheckTopicWritePolicy(ctx context.Context, actor *model.Actor, topic *entity.ForumTopic) error
		CheckTopicEditPolicy(ctx context.Context, actor *model.Actor, topic *entity.ForumTopic) error
		CheckTopicAppendPolicy(ctx context.Context, actor *model.Actor, topic *entity.ForumTopic) error
		InsertTopic(ctx context.Context, actor *model.Actor, in forumin.TopicCreateInp) (uint64, error)
		AppendContent(ctx context.Context, topic *entity.ForumTopic, content string) error
		UpdateTopic(ctx context.Context, id uint64, nodeId uint, subject string, content string) error
		IncrementTopicViews(ctx context.Context, topicId uint64) error
		UpdateTopicReplyStats(ctx context.Context, topicId uint64, replyId uint64, lastReplyBy uint64) error
		ToggleLike(ctx context.Context, actor *model.Actor, topicId uint64) (bool, error)
		Bookmark(ctx context.Context, actor *model.Actor, topicId uint64) error
		Unbookmark(ctx context.Context, actor *model.Actor, topicId uint64) error
		QueryBookmarkedTopics(ctx context.Context, actor *model.Actor, page int, size int) ([]entity.ForumTopic, int, error)
		AdminSetTopicLock(ctx context.Context, id uint64, isLocked bool) error
		AdminSetTopicSticky(ctx context.Context, id uint64, isSticky bool) error
		AdminMoveTopic(ctx context.Context, id uint64, newNodeId uint) error
		QueryTopicsByNode(ctx context.Context, nodeId uint, page int, size int) ([]entity.ForumTopic, int, error)
		CheckTopicLiked(ctx context.Context, topicId uint64, userId uint64) (bool, error)
		CheckTopicBookmarked(ctx context.Context, topicId uint64, userId uint64) (bool, error)
	}
	IForumTopicUsecase interface {
		List(ctx context.Context, actor *model.Actor, in forumin.TopicListInp) (*forumout.TopicListOut, error)
		Detail(ctx context.Context, actor *model.Actor, in forumin.TopicDetailInp) (*forumout.TopicDetailOut, error)
		Create(ctx context.Context, actor *model.Actor, in forumin.TopicCreateInp) (uint64, error)
		Update(ctx context.Context, actor *model.Actor, in forumin.TopicUpdateInp) error
		Append(ctx context.Context, actor *model.Actor, in forumin.TopicAppendInp) error
		ToggleTopicLike(ctx context.Context, actor *model.Actor, in forumin.TopicToggleLikeInp) error
		RewardTopic(ctx context.Context, actor *model.Actor, in forumin.TopicRewardInp) error
		ReportTopic(ctx context.Context, actor *model.Actor, in forumin.TopicReportInp) error
		BookmarkTopic(ctx context.Context, actor *model.Actor, in forumin.TopicBookmarkInp) error
		UnbookmarkTopic(ctx context.Context, actor *model.Actor, in forumin.TopicUnbookmarkInp) error
		ListBookmarkedTopics(ctx context.Context, actor *model.Actor, in forumin.TopicBookmarkListInp) (*forumout.TopicBookmarkListOut, error)
	}
)

var (
	localForumCategoryDomain IForumCategoryDomain
	localForumNodeDomain     IForumNodeDomain
	localForumNodeUsecase    IForumNodeUsecase
	localForumReplyDomain    IForumReplyDomain
	localForumReplyUsecase   IForumReplyUsecase
	localForumTopicDomain    IForumTopicDomain
	localForumTopicUsecase   IForumTopicUsecase
)

func ForumCategoryDomain() IForumCategoryDomain {
	if localForumCategoryDomain == nil {
		panic("implement not found for interface IForumCategoryDomain, forgot register?")
	}
	return localForumCategoryDomain
}

func RegisterForumCategoryDomain(i IForumCategoryDomain) {
	localForumCategoryDomain = i
}

func ForumNodeDomain() IForumNodeDomain {
	if localForumNodeDomain == nil {
		panic("implement not found for interface IForumNodeDomain, forgot register?")
	}
	return localForumNodeDomain
}

func RegisterForumNodeDomain(i IForumNodeDomain) {
	localForumNodeDomain = i
}

func ForumNodeUsecase() IForumNodeUsecase {
	if localForumNodeUsecase == nil {
		panic("implement not found for interface IForumNodeUsecase, forgot register?")
	}
	return localForumNodeUsecase
}

func RegisterForumNodeUsecase(i IForumNodeUsecase) {
	localForumNodeUsecase = i
}

func ForumReplyDomain() IForumReplyDomain {
	if localForumReplyDomain == nil {
		panic("implement not found for interface IForumReplyDomain, forgot register?")
	}
	return localForumReplyDomain
}

func RegisterForumReplyDomain(i IForumReplyDomain) {
	localForumReplyDomain = i
}

func ForumReplyUsecase() IForumReplyUsecase {
	if localForumReplyUsecase == nil {
		panic("implement not found for interface IForumReplyUsecase, forgot register?")
	}
	return localForumReplyUsecase
}

func RegisterForumReplyUsecase(i IForumReplyUsecase) {
	localForumReplyUsecase = i
}

func ForumTopicDomain() IForumTopicDomain {
	if localForumTopicDomain == nil {
		panic("implement not found for interface IForumTopicDomain, forgot register?")
	}
	return localForumTopicDomain
}

func RegisterForumTopicDomain(i IForumTopicDomain) {
	localForumTopicDomain = i
}

func ForumTopicUsecase() IForumTopicUsecase {
	if localForumTopicUsecase == nil {
		panic("implement not found for interface IForumTopicUsecase, forgot register?")
	}
	return localForumTopicUsecase
}

func RegisterForumTopicUsecase(i IForumTopicUsecase) {
	localForumTopicUsecase = i
}
