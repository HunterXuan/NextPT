# IAM (身份与访问管理) 域设计

## 核心实体 (Domain Entities)
- `User` (用户)
- `Role` (角色)
- `Permission` (权限)
- `Invite` (邀请码)
- `Session` / `LoginLog` (有效会话与登录记录)

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/iam`

### 1. SessionUsecase (会话与认证应用服务)
* **登录 (Login)**
  * **Method/Path**: `POST /sessions`
  * **参数概述**: `username`, `password`
  * **核心逻辑**: 密码验证；未启用两步验证时直接签发 Session，启用 TOTP 时只返回一次性 `twoStepChallenge`。
  * 前端初始化长期 `deviceId` 并通过 `X-Device-Id` 请求头随浏览器请求发送；服务端仅保存摘要，将其绑定到新建 Session，并在后续鉴权时要求摘要一致。首次出现时发送新设备登录通知。
* **完成两步验证登录 (VerifyTwoStep)**
  * **Method/Path**: `POST /sessions:verifyTwoStep`
  * **参数概述**: `challenge`, `code`
  * **核心逻辑**: 校验短时、单次 challenge 及 TOTP 或恢复码；成功后才签发 Session、更新登录记录和最后登录信息。恢复码只可使用一次。
* **登出 (Logout)**
  * **Method/Path**: `DELETE /sessions`
  * **核心逻辑**: 只销毁当前请求对应的 Session，不影响同一用户的其它设备。
* **获取有效会话 (ListSessions)**
  * **Method/Path**: `GET /sessions`
  * **核心逻辑**: 返回当前用户尚未过期的设备会话，包含 IP、User-Agent、创建时间、最近活跃时间和当前设备标识。
* **销毁指定会话 (DeleteSession)**
  * **Method/Path**: `DELETE /sessions/{id}`
  * **核心逻辑**: 仅允许用户销毁属于自己的 Session；修改密码、重置 Passkey、密码找回和管理员踢出仍销毁该用户全部 Session。

> Web Token 使用随机 `sessionId` 作为 `gtoken` user key，不再直接使用用户 ID。Redis 保存 Session 元数据、`userId -> sessionId` 集合和一年期已知设备摘要集合；会话列表会清理已过期的索引，最近活跃时间按 5 分钟节流更新。设备 ID 不单独代表用户身份，但作为绑定条件参与 Session 鉴权。该结构不需要新增数据库表。

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
* **获取当前权限与登录记录**
  * **Method/Path**: `GET /users/me/permissions`, `GET /users/me/login-logs`
  * 权限列表供前端做交互限制，实际授权仍由后端 RBAC 判定；登录记录只返回当前用户自身数据。
* **重置 Passkey**
  * **Method/Path**: `POST /users/me:resetPasskey`
  * 重置后会销毁该用户现有会话，Tracker 客户端需要使用新的 Passkey 重新配置。
* **申请重置密码 (CreatePasswordResetRequest)**
  * **Method/Path**: `POST /password-reset-requests`
  * **参数概述**: `email`
  * **核心逻辑**: 按 IP 和邮箱限流；无论邮箱是否存在都返回相同结果。命中用户时生成一次性令牌，仅在 Redis 保存令牌摘要和用户映射，有效期 30 分钟，并使用站点默认语言异步发送重置邮件。新申请会使该用户此前的令牌失效。
* **确认重置密码 (CreatePasswordReset)**
  * **Method/Path**: `POST /password-resets`
  * **参数概述**: `token`, `newPassword`
  * **核心逻辑**: 原子校验并消费一次性令牌 -> 更新密码哈希 -> 删除该用户全部登录会话 -> 清理用户鉴权缓存。无效、已使用和过期令牌统一返回相同错误。
* **开始绑定两步验证 (SetupTwoStep)**
  * **Method/Path**: `POST /users/me/two-step:setup`
  * **参数概述**: `password`
  * **核心逻辑**: 验证当前密码，生成 RFC 6238 TOTP 密钥和仅短时保存在 Redis 的绑定 challenge，返回二维码与手动密钥。
* **确认两步验证绑定 (ConfirmTwoStep)**
  * **Method/Path**: `POST /users/me/two-step:confirm`
  * **参数概述**: `challenge`, `code`
  * **核心逻辑**: 验证认证器验证码后加密保存 TOTP 密钥，并原子生成恢复码哈希；恢复码明文仅在本次响应中返回。当前已验证会话保持有效，后续新登录必须完成两步验证。
* **重新生成恢复码 (CreateTwoStepRecoveryCodes)**
  * **Method/Path**: `POST /users/me/two-step:recoveryCodes`
  * **参数概述**: `code`
  * **核心逻辑**: 使用当前认证器验证码或一个有效恢复码重新生成整组恢复码，旧恢复码立即失效。
* **关闭两步验证 (DeleteTwoStep)**
  * **Method/Path**: `DELETE /users/me/two-step`
  * **参数概述**: `password`, `code`
  * **核心逻辑**: 同时验证当前密码和认证器验证码/恢复码，清除加密密钥与恢复码；当前会话保持有效。

> 待验证用户不能创建登录会话；登录接口会返回独立的邮箱未验证错误，不与封禁或停用状态混用。

### 3. InviteUsecase (邀请应用服务)
* **发送/生成邀请 (CreateInvite)**
  * **Method/Path**: `POST /invites:send`
  * **参数概述**: `email`
  * **核心逻辑**: 调用 Economy 扣除发送者魔力值 -> 生成 `invite` -> 发送邮件。
* **获取我的邀请记录 (ListInvites)**
  * **Method/Path**: `GET /invites`
* **核验邀请码有效性 (CheckInvite)**
  * **Method/Path**: `GET /invites:check`

### 4. PermissionDomain (内部权限域)
* **核心职责**：
  * **GetAllPermissions()**：暴露 `consts.IamPermissionAll` 的常量列表，供 Admin 域在呈现权限树时跨域调用。
  * **CheckPermission(ctx, userId, roleId, permKey)**：进行基于 Role 和 UserPermission ACL 结合通配符规则的鉴权。
