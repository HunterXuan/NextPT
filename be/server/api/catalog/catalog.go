// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package catalog

import (
	"context"

	"server/api/catalog/v1"
)

type ICatalogV1 interface {
	BookmarkList(ctx context.Context, req *v1.BookmarkListReq) (res *v1.BookmarkListRes, err error)
	CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error)
	TagGroupList(ctx context.Context, req *v1.TagGroupListReq) (res *v1.TagGroupListRes, err error)
	CommentCreate(ctx context.Context, req *v1.CommentCreateReq) (res *v1.CommentCreateRes, err error)
	CommentList(ctx context.Context, req *v1.CommentListReq) (res *v1.CommentListRes, err error)
	CommentReport(ctx context.Context, req *v1.CommentReportReq) (res *v1.CommentReportRes, err error)
	CommentToggleLike(ctx context.Context, req *v1.CommentToggleLikeReq) (res *v1.CommentToggleLikeRes, err error)
	CommentReward(ctx context.Context, req *v1.CommentRewardReq) (res *v1.CommentRewardRes, err error)
	RequestList(ctx context.Context, req *v1.RequestListReq) (res *v1.RequestListRes, err error)
	RequestGet(ctx context.Context, req *v1.RequestGetReq) (res *v1.RequestGetRes, err error)
	RequestCreate(ctx context.Context, req *v1.RequestCreateReq) (res *v1.RequestCreateRes, err error)
	RequestClaim(ctx context.Context, req *v1.RequestClaimReq) (res *v1.RequestClaimRes, err error)
	RequestAbandon(ctx context.Context, req *v1.RequestAbandonReq) (res *v1.RequestAbandonRes, err error)
	RequestSubmit(ctx context.Context, req *v1.RequestSubmitReq) (res *v1.RequestSubmitRes, err error)
	RequestComplete(ctx context.Context, req *v1.RequestCompleteReq) (res *v1.RequestCompleteRes, err error)
	RequestCancel(ctx context.Context, req *v1.RequestCancelReq) (res *v1.RequestCancelRes, err error)
	RequestCommentCreate(ctx context.Context, req *v1.RequestCommentCreateReq) (res *v1.RequestCommentCreateRes, err error)
	RequestCommentList(ctx context.Context, req *v1.RequestCommentListReq) (res *v1.RequestCommentListRes, err error)
	RequestCommentToggleLike(ctx context.Context, req *v1.RequestCommentToggleLikeReq) (res *v1.RequestCommentToggleLikeRes, err error)
	RequestCommentReward(ctx context.Context, req *v1.RequestCommentRewardReq) (res *v1.RequestCommentRewardRes, err error)
	RequestCommentReport(ctx context.Context, req *v1.RequestCommentReportReq) (res *v1.RequestCommentReportRes, err error)
	SubtitleList(ctx context.Context, req *v1.SubtitleListReq) (res *v1.SubtitleListRes, err error)
	TorrentSubtitleList(ctx context.Context, req *v1.TorrentSubtitleListReq) (res *v1.TorrentSubtitleListRes, err error)
	SubtitleUpload(ctx context.Context, req *v1.SubtitleUploadReq) (res *v1.SubtitleUploadRes, err error)
	SubtitleDownload(ctx context.Context, req *v1.SubtitleDownloadReq) (res *v1.SubtitleDownloadRes, err error)
	SubtitleUpdate(ctx context.Context, req *v1.SubtitleUpdateReq) (res *v1.SubtitleUpdateRes, err error)
	SubtitleReport(ctx context.Context, req *v1.SubtitleReportReq) (res *v1.SubtitleReportRes, err error)
	TorrentUpload(ctx context.Context, req *v1.TorrentUploadReq) (res *v1.TorrentUploadRes, err error)
	TorrentList(ctx context.Context, req *v1.TorrentListReq) (res *v1.TorrentListRes, err error)
	TorrentGetHot(ctx context.Context, req *v1.TorrentGetHotReq) (res *v1.TorrentGetHotRes, err error)
	TorrentMetadataSearch(ctx context.Context, req *v1.TorrentMetadataSearchReq) (res *v1.TorrentMetadataSearchRes, err error)
	TorrentGet(ctx context.Context, req *v1.TorrentGetReq) (res *v1.TorrentGetRes, err error)
	TorrentDownload(ctx context.Context, req *v1.TorrentDownloadReq) (res *v1.TorrentDownloadRes, err error)
	TorrentReward(ctx context.Context, req *v1.TorrentRewardReq) (res *v1.TorrentRewardRes, err error)
	TorrentRewardList(ctx context.Context, req *v1.TorrentRewardListReq) (res *v1.TorrentRewardListRes, err error)
	TorrentBookmark(ctx context.Context, req *v1.TorrentBookmarkReq) (res *v1.TorrentBookmarkRes, err error)
	TorrentUnbookmark(ctx context.Context, req *v1.TorrentUnbookmarkReq) (res *v1.TorrentUnbookmarkRes, err error)
	TorrentToggleLike(ctx context.Context, req *v1.TorrentToggleLikeReq) (res *v1.TorrentToggleLikeRes, err error)
	TorrentLikeList(ctx context.Context, req *v1.TorrentLikeListReq) (res *v1.TorrentLikeListRes, err error)
	TorrentUpdate(ctx context.Context, req *v1.TorrentUpdateReq) (res *v1.TorrentUpdateRes, err error)
	TorrentFileList(ctx context.Context, req *v1.TorrentFileListReq) (res *v1.TorrentFileListRes, err error)
	TorrentPeerList(ctx context.Context, req *v1.TorrentPeerListReq) (res *v1.TorrentPeerListRes, err error)
	TorrentReport(ctx context.Context, req *v1.TorrentReportReq) (res *v1.TorrentReportRes, err error)
}
