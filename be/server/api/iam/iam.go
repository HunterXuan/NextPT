// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package iam

import (
	"context"

	"server/api/iam/v1"
)

type IIamV1 interface {
	EmailVerificationRequestCreate(ctx context.Context, req *v1.EmailVerificationRequestCreateReq) (res *v1.EmailVerificationRequestCreateRes, err error)
	EmailVerificationCreate(ctx context.Context, req *v1.EmailVerificationCreateReq) (res *v1.EmailVerificationCreateRes, err error)
	InviteList(ctx context.Context, req *v1.InviteListReq) (res *v1.InviteListRes, err error)
	InviteSend(ctx context.Context, req *v1.InviteSendReq) (res *v1.InviteSendRes, err error)
	InviteCheck(ctx context.Context, req *v1.InviteCheckReq) (res *v1.InviteCheckRes, err error)
	PasswordResetRequestCreate(ctx context.Context, req *v1.PasswordResetRequestCreateReq) (res *v1.PasswordResetRequestCreateRes, err error)
	PasswordResetCreate(ctx context.Context, req *v1.PasswordResetCreateReq) (res *v1.PasswordResetCreateRes, err error)
	RoleList(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error)
	IamSessionCreate(ctx context.Context, req *v1.IamSessionCreateReq) (res *v1.IamSessionCreateRes, err error)
	IamSessionList(ctx context.Context, req *v1.IamSessionListReq) (res *v1.IamSessionListRes, err error)
	IamSessionTwoStepVerify(ctx context.Context, req *v1.IamSessionTwoStepVerifyReq) (res *v1.IamSessionTwoStepVerifyRes, err error)
	IamSessionDelete(ctx context.Context, req *v1.IamSessionDeleteReq) (res *v1.IamSessionDeleteRes, err error)
	IamSessionDeleteById(ctx context.Context, req *v1.IamSessionDeleteByIdReq) (res *v1.IamSessionDeleteByIdRes, err error)
	IamUserCreate(ctx context.Context, req *v1.IamUserCreateReq) (res *v1.IamUserCreateRes, err error)
	IamUserGet(ctx context.Context, req *v1.IamUserGetReq) (res *v1.IamUserGetRes, err error)
	IamUserMe(ctx context.Context, req *v1.IamUserMeReq) (res *v1.IamUserMeRes, err error)
	IamUserPermissionList(ctx context.Context, req *v1.IamUserPermissionListReq) (res *v1.IamUserPermissionListRes, err error)
	IamUserLoginLogList(ctx context.Context, req *v1.IamUserLoginLogListReq) (res *v1.IamUserLoginLogListRes, err error)
	IamUserProfileUpdate(ctx context.Context, req *v1.IamUserProfileUpdateReq) (res *v1.IamUserProfileUpdateRes, err error)
	IamUserPasswordChange(ctx context.Context, req *v1.IamUserPasswordChangeReq) (res *v1.IamUserPasswordChangeRes, err error)
	IamUserTwoStepSetup(ctx context.Context, req *v1.IamUserTwoStepSetupReq) (res *v1.IamUserTwoStepSetupRes, err error)
	IamUserTwoStepConfirm(ctx context.Context, req *v1.IamUserTwoStepConfirmReq) (res *v1.IamUserTwoStepConfirmRes, err error)
	IamUserTwoStepRecoveryCodesCreate(ctx context.Context, req *v1.IamUserTwoStepRecoveryCodesCreateReq) (res *v1.IamUserTwoStepRecoveryCodesCreateRes, err error)
	IamUserTwoStepDelete(ctx context.Context, req *v1.IamUserTwoStepDeleteReq) (res *v1.IamUserTwoStepDeleteRes, err error)
	IamUserPasskeyReset(ctx context.Context, req *v1.IamUserPasskeyResetReq) (res *v1.IamUserPasskeyResetRes, err error)
}
