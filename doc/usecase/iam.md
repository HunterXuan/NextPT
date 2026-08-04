# IAM (身份与访问管理) 域设计

## 核心实体 (Domain Entities)
- `User` (用户)
- `Role` (角色)
- `Permission` (权限)
- `Invite` (邀请码)
- `Session` / `LoginLog` (会话与登录记录)

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/v1/iam`

### 1. SessionUsecase (会话与认证应用服务)
* **登录 (Login)**
  * **Method/Path**: `POST /sessions`
  * **参数概述**: `username`, `password`
  * **核心逻辑**: 密码验证 -> 签发 JWT/创建 Session 记录。
* **登出 (Logout)**
  * **Method/Path**: `DELETE /sessions`

### 2. UserUsecase (用户应用服务)
* **注册 (Register)**
  * **Method/Path**: `POST /users`
  * **参数概述**: `username`, `password`, `email`, `invite_token`
  * **核心逻辑**: 校验邀请码 -> 创建待验证用户 -> 分配默认角色 -> 初始化用户资料和统计 -> 标记 invite 已使用 -> 异步发送邮箱验证邮件。
* **申请邮箱验证 (CreateEmailVerificationRequest)**
  * **Method/Path**: `POST /email-verification-requests`
  * **参数概述**: `email`
  * **核心逻辑**: 按 IP 和邮箱限流；未知邮箱、已验证邮箱和限流请求返回相同结果。待验证用户会获得新的 24 小时一次性令牌，Redis 只保存令牌摘要，新申请会使旧令牌失效。
* **确认邮箱验证 (CreateEmailVerification)**
  * **Method/Path**: `POST /email-verifications`
  * **参数概述**: `token`
  * **核心逻辑**: 原子校验并消费令牌 -> 仅将待验证用户更新为正常状态 -> 清理用户鉴权缓存。无效、已使用和过期令牌统一返回相同错误。
* **获取当前用户资料 (GetMyProfile)**
  * **Method/Path**: `GET /users/me`
* **获取用户公开资料 (GetPublicProfile)**
  * **Method/Path**: `GET /users/{id}`
  * **核心逻辑**: 聚合头像、用户名、签名、简介、角色、注册时间与公开分享统计，最终结果按用户和语言缓存 5 分钟；不返回邮箱、Passkey、真实流量、权限、登录记录或账号限制等敏感信息。
* **更新当前用户资料 (UpdateMyProfile)**
  * **Method/Path**: `PATCH /users/me`
  * **参数概述**: `avatar`, `signature` (支持局部更新)
* **修改密码 (ChangePassword)**
  * **Method/Path**: `POST /users/me:changePassword`
  * **参数概述**: `old_password`, `new_password`
* **申请重置密码 (CreatePasswordResetRequest)**
  * **Method/Path**: `POST /password-reset-requests`
  * **参数概述**: `email`
  * **核心逻辑**: 按 IP 和邮箱限流；无论邮箱是否存在都返回相同结果。命中用户时生成一次性令牌，仅在 Redis 保存令牌摘要和用户映射，有效期 30 分钟，并使用站点默认语言异步发送重置邮件。新申请会使该用户此前的令牌失效。
* **确认重置密码 (CreatePasswordReset)**
  * **Method/Path**: `POST /password-resets`
  * **参数概述**: `token`, `newPassword`
  * **核心逻辑**: 原子校验并消费一次性令牌 -> 更新密码哈希 -> 删除该用户全部登录会话 -> 清理用户鉴权缓存。无效、已使用和过期令牌统一返回相同错误。

> 待验证用户不能创建登录会话；登录接口会返回独立的邮箱未验证错误，不与封禁或停用状态混用。

### 3. InviteUsecase (邀请应用服务)
* **发送/生成邀请 (CreateInvite)**
  * **Method/Path**: `POST /invites`
  * **参数概述**: `email`
  * **核心逻辑**: 调用 Economy 扣除发送者魔力值 -> 生成 `invite` -> 发送邮件。
* **获取我的邀请记录 (ListInvites)**
  * **Method/Path**: `GET /invites`
* **核验邀请码有效性 (CheckInvite)**
  * **Method/Path**: `GET /invites/{token}:check`

### 4. PermissionDomain (内部权限域)
* **核心职责**：
  * **GetAllPermissions()**：暴露 `consts.IamPermissionAll` 的常量列表，供 Admin 域在呈现权限树时跨域调用。
  * **CheckPermission(ctx, userId, roleId, permKey)**：进行基于 Role 和 UserPermission ACL 结合通配符规则的鉴权。
