// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package iam

import (
	"context"

	"server/api/iam/v1"
)

type IIamV1 interface {
	InviteList(ctx context.Context, req *v1.InviteListReq) (res *v1.InviteListRes, err error)
	InviteSend(ctx context.Context, req *v1.InviteSendReq) (res *v1.InviteSendRes, err error)
	InviteCheck(ctx context.Context, req *v1.InviteCheckReq) (res *v1.InviteCheckRes, err error)
	IamSessionCreate(ctx context.Context, req *v1.IamSessionCreateReq) (res *v1.IamSessionCreateRes, err error)
	IamSessionDelete(ctx context.Context, req *v1.IamSessionDeleteReq) (res *v1.IamSessionDeleteRes, err error)
	IamUserCreate(ctx context.Context, req *v1.IamUserCreateReq) (res *v1.IamUserCreateRes, err error)
	IamUserMe(ctx context.Context, req *v1.IamUserMeReq) (res *v1.IamUserMeRes, err error)
	IamUserPermissionList(ctx context.Context, req *v1.IamUserPermissionListReq) (res *v1.IamUserPermissionListRes, err error)
	IamUserProfileUpdate(ctx context.Context, req *v1.IamUserProfileUpdateReq) (res *v1.IamUserProfileUpdateRes, err error)
	IamUserPasswordChange(ctx context.Context, req *v1.IamUserPasswordChangeReq) (res *v1.IamUserPasswordChangeRes, err error)
	IamUserPasskeyReset(ctx context.Context, req *v1.IamUserPasskeyResetReq) (res *v1.IamUserPasskeyResetRes, err error)
}
