export interface AuthUser {
  user: AuthUserAccount
  role: AuthUserRole
  profile: AuthUserProfile
  stat: AuthUserStat
}

export interface AuthUserAccount {
  id: number
  username: string
  email: string
  passkey: string
  status: number
  twoStepEnabled: boolean
  vipUntil?: string | null
  createdAt?: string | null
}

export interface AuthUserRole {
  id: number
  name: string
  level: number
  isStaff: boolean
}

export interface AuthUserProfile {
  avatar?: string
  info?: string
  signature?: string
}

export interface AuthUserStat {
  uploaded: number
  downloaded: number
  rawUploaded: number
  rawDownloaded: number
  bonus: number
  shareRatio: number
}

export interface RegisterInput {
  username: string
  email: string
  password: string
  inviteHash?: string
}

export interface ProfileInput {
  avatar?: string
  info?: string
  signature?: string
}

export interface PasswordInput {
  oldPassword: string
  newPassword: string
}

export interface PasswordResetInput {
  token: string
  newPassword: string
}

export interface EmailVerificationInput {
  token: string
}

export interface PasskeyResetOut {
  passkey: string
}

export interface AuthPermissionListOut {
  permissions: string[]
}

export interface AuthLoginOut {
  token: string
  twoStepRequired: boolean
  twoStepChallenge: string
}

export interface AuthTwoStepSetupOut {
  challenge: string
  qrCodeDataUrl: string
  secret: string
}

export interface AuthTwoStepRecoveryCodesOut {
  recoveryCodes: string[]
}

export interface AuthLoginLog {
  id: number
  userId: number
  ip: string
  userAgent: string
  result: number
  failReason: string
  createdAt?: string | null
}

export interface AuthLoginLogListParams {
  result?: number
  page?: number
  size?: number
}

export interface AuthLoginLogListOut {
  list: AuthLoginLog[]
  total: number
  page: number
  size: number
}

export interface AuthSessionItem {
  id: string
  ip: string
  userAgent: string
  createdAt?: string | null
  lastSeenAt?: string | null
  current: boolean
}

export interface AuthSessionListOut {
  list: AuthSessionItem[]
}

export interface AuthRoleRuleCondition {
  accountAgeDaysGte?: number
  downloadedGiBGte?: number
  downloadedGiBGt?: number
  downloadedGiBLte?: number
  ratioGt?: number
  ratioGte?: number
  ratioLt?: number
  [key: string]: unknown
}

export interface AuthRoleRules {
  promotion?: AuthRoleRuleCondition[]
  demotion?: AuthRoleRuleCondition[]
  [key: string]: unknown
}

export interface AuthRoleItem {
  id: number
  level: number
  name: string
  nameI18N: Record<string, string>
  rules: AuthRoleRules
  isStaff: boolean
}

export interface AuthRoleListOut {
  roles: AuthRoleItem[]
}

export const Permission = {
  AdminIamUserManage: 'admin:iam/user:*',
  AdminIamRoleManage: 'admin:iam/role:*',
  AdminIamInviteManage: 'admin:iam/invite:*',
  AdminModReportManage: 'admin:mod/report:*',
  AdminModCheaterManage: 'admin:mod/cheater:*',
  AdminModUserManage: 'admin:mod/user:*',
  AdminModStaffMessage: 'admin:mod/staff-message:*',
  AdminSiteConfig: 'admin:site/config:*',
  AdminSiteAudit: 'admin:site/audit:*',
  AdminSiteAnnouncement: 'admin:site/announcement:*',
  AdminSiteMessage: 'admin:site/message:*',
  AdminSysCronManage: 'admin:sys/cron:*',
  AdminForumCategoryManage: 'admin:forum/category:*',
  AdminForumNodeManage: 'admin:forum/node:*',
  AdminForumTopicManage: 'admin:forum/topic:*',
  AdminForumReplyManage: 'admin:forum/reply:*',
  AdminCatalogCategoryManage: 'admin:catalog/category:*',
  AdminCatalogTagManage: 'admin:catalog/tag:*',
  AdminCatalogTorrentManage: 'admin:catalog/torrent:*',
  AdminCatalogSubtitleManage: 'admin:catalog/subtitle:*',
  AdminCatalogCommentManage: 'admin:catalog/comment:*',
  AdminCatalogRequestManage: 'admin:catalog/request:*',
  IamInviteRead: 'read:iam/invite:*',
  IamInviteCreate: 'create:iam/invite:*',
  EconomyShopProductRead: 'read:economy/shop-product:*',
  EconomyShopOrderRead: 'read:economy/shop-order:*',
  EconomyShopOrderCreate: 'create:economy/shop-order:*',
  SiteAnnouncementRead: 'read:site/announcement:*',
  SiteMessageRead: 'read:site/message:*',
  ModStaffMessageRead: 'read:mod/staff-message:*',
  ModStaffMessageCreate: 'create:mod/staff-message:*',
  ForumTopicRead: 'read:forum/topic:*',
  ForumTopicCreate: 'create:forum/topic:*',
  ForumTopicUpdate: 'update:forum/topic:*',
  ForumReplyRead: 'read:forum/reply:*',
  ForumReplyCreate: 'create:forum/reply:*',
  ForumReplyUpdate: 'update:forum/reply:*',
  CatalogTorrentRead: 'read:catalog/torrent:*',
  CatalogTorrentCreate: 'create:catalog/torrent:*',
  CatalogTorrentDownload: 'download:catalog/torrent:*',
  CatalogSubtitleRead: 'read:catalog/subtitle:*',
  CatalogSubtitleCreate: 'create:catalog/subtitle:*',
  CatalogSubtitleDownload: 'download:catalog/subtitle:*',
  CatalogCommentRead: 'read:catalog/comment:*',
  CatalogCommentCreate: 'create:catalog/comment:*',
  CatalogRequestRead: 'read:catalog/request:*',
  CatalogRequestCreate: 'create:catalog/request:*',
  CatalogRequestUpdate: 'update:catalog/request:*'
} as const

export const AdminPermissions = [
  Permission.AdminIamUserManage,
  Permission.AdminIamRoleManage,
  Permission.AdminIamInviteManage,
  Permission.AdminModReportManage,
  Permission.AdminModCheaterManage,
  Permission.AdminModUserManage,
  Permission.AdminModStaffMessage,
  Permission.AdminSiteConfig,
  Permission.AdminSiteAudit,
  Permission.AdminSiteAnnouncement,
  Permission.AdminSiteMessage,
  Permission.AdminSysCronManage,
  Permission.AdminForumCategoryManage,
  Permission.AdminForumNodeManage,
  Permission.AdminForumTopicManage,
  Permission.AdminForumReplyManage,
  Permission.AdminCatalogCategoryManage,
  Permission.AdminCatalogTagManage,
  Permission.AdminCatalogTorrentManage,
  Permission.AdminCatalogSubtitleManage,
  Permission.AdminCatalogCommentManage,
  Permission.AdminCatalogRequestManage
]

export function useAuth() {
  const token = useCookie<string | null>('nextpt_token', {
    sameSite: 'lax',
    maxAge: 60 * 60 * 24 * 10,
    watch: true
  })
  const user = useState<AuthUser | null>('auth:user', () => null)
  const permissions = useState<string[]>('auth:permissions', () => [])
  const permissionsLoaded = useState<boolean>('auth:permissionsLoaded', () => false)
  const isLoggedIn = computed(() => Boolean(token.value))
  const isStaff = computed(() => Boolean(user.value?.role.isStaff))

  setApiAuthToken(token.value)

  function setToken(value: string | null) {
    token.value = value
    setApiAuthToken(value)
  }

  function clearLocalSession() {
    user.value = null
    permissions.value = []
    permissionsLoaded.value = false
    setToken(null)
  }

  async function login(username: string, password: string) {
    const session = await fetchApi<AuthLoginOut>('/api/iam/sessions', {
      method: 'POST',
      body: { username, password }
    })

    if (session.twoStepRequired) return session

    setToken(session.token)
    await fetchUser()
    return session
  }

  async function verifyTwoStepLogin(challenge: string, code: string) {
    const session = await fetchApi<AuthLoginOut>('/api/iam/sessions:verifyTwoStep', {
      method: 'POST',
      body: { challenge, code }
    })
    setToken(session.token)
    await fetchUser()
    return session
  }

  async function register(input: RegisterInput) {
    return await fetchApi<{ id: number }>('/api/iam/users', {
      method: 'POST',
      body: input
    })
  }

  async function fetchUser() {
    if (!token.value) {
      user.value = null
      permissions.value = []
      permissionsLoaded.value = false
      return null
    }

    const data = await fetchApi<AuthUser>('/api/iam/users/me')
    user.value = data
    permissions.value = []
    permissionsLoaded.value = false
    try {
      await fetchPermissions(true)
    } catch {
      permissions.value = []
      permissionsLoaded.value = true
    }
    return data
  }

  async function fetchPermissions(force = false) {
    if (!token.value) {
      permissions.value = []
      permissionsLoaded.value = false
      return []
    }
    if (!force && permissionsLoaded.value) {
      return permissions.value
    }

    const data = await fetchApi<AuthPermissionListOut>('/api/iam/users/me/permissions')
    permissions.value = data.permissions || []
    permissionsLoaded.value = true
    return permissions.value
  }

  async function listRoles() {
    return await fetchApi<AuthRoleListOut>('/api/iam/roles')
  }

  async function listLoginLogs(params: AuthLoginLogListParams = {}) {
    return await fetchApi<AuthLoginLogListOut>('/api/iam/users/me/login-logs', {
      query: params
    })
  }

  function hasPermission(permission: string) {
    return matchesPermissionList(permissions.value, permission)
  }

  function hasAnyPermission(requiredPermissions: string[]) {
    return requiredPermissions.some((permission) => hasPermission(permission))
  }

  async function updateProfile(input: ProfileInput) {
    await fetchApi('/api/iam/users/me', {
      method: 'PATCH',
      body: input
    })
    await fetchUser()
  }

  async function changePassword(input: PasswordInput) {
    await fetchApi('/api/iam/users/me:changePassword', {
      method: 'POST',
      body: input
    })
    clearLocalSession()
  }

  async function requestEmailVerification(email: string) {
    await fetchApi('/api/iam/email-verification-requests', {
      method: 'POST',
      body: { email }
    })
  }

  async function verifyEmail(input: EmailVerificationInput) {
    await fetchApi('/api/iam/email-verifications', {
      method: 'POST',
      body: input
    })
  }

  async function requestPasswordReset(email: string) {
    await fetchApi('/api/iam/password-reset-requests', {
      method: 'POST',
      body: { email }
    })
  }

  async function resetPassword(input: PasswordResetInput) {
    await fetchApi('/api/iam/password-resets', {
      method: 'POST',
      body: input
    })
  }

  async function resetPasskey() {
    const data = await fetchApi<PasskeyResetOut>('/api/iam/users/me:resetPasskey', {
      method: 'POST'
    })
    await fetchUser()
    return data
  }

  async function setupTwoStep(password: string) {
    return await fetchApi<AuthTwoStepSetupOut>('/api/iam/users/me/two-step:setup', {
      method: 'POST',
      body: { password }
    })
  }

  async function confirmTwoStep(challenge: string, code: string) {
    const data = await fetchApi<AuthTwoStepRecoveryCodesOut>('/api/iam/users/me/two-step:confirm', {
      method: 'POST',
      body: { challenge, code }
    })
    await fetchUser()
    return data
  }

  async function createTwoStepRecoveryCodes(code: string) {
    return await fetchApi<AuthTwoStepRecoveryCodesOut>('/api/iam/users/me/two-step:recoveryCodes', {
      method: 'POST',
      body: { code }
    })
  }

  async function disableTwoStep(password: string, code: string) {
    await fetchApi('/api/iam/users/me/two-step', {
      method: 'DELETE',
      body: { password, code }
    })
    await fetchUser()
  }

  async function listSessions() {
    return await fetchApi<AuthSessionListOut>('/api/iam/sessions')
  }

  async function deleteSession(id: string) {
    await fetchApi(`/api/iam/sessions/${encodeURIComponent(id)}`, { method: 'DELETE' })
  }

  async function logout(remote = true) {
    if (remote && token.value) {
      try {
        await fetchApi('/api/iam/sessions', { method: 'DELETE' })
      } catch {
        // Local session cleanup still matters if the server token is already gone.
      }
    }

    clearLocalSession()
  }

  return {
    token,
    user,
    permissions,
    permissionsLoaded,
    isLoggedIn,
    isStaff,
    login,
    verifyTwoStepLogin,
    register,
    fetchUser,
    fetchPermissions,
    listRoles,
    listLoginLogs,
    hasPermission,
    hasAnyPermission,
    updateProfile,
    changePassword,
    requestEmailVerification,
    verifyEmail,
    requestPasswordReset,
    resetPassword,
    resetPasskey,
    setupTwoStep,
    confirmTwoStep,
    createTwoStepRecoveryCodes,
    disableTwoStep,
    listSessions,
    deleteSession,
    logout,
    clearLocalSession
  }
}

function matchesPermissionList(permissions: string[], permission: string) {
  if (!permission) return false

  return permissions.some((pattern) => {
    if (!pattern) return false
    if (pattern === permission) return true
    if (pattern.endsWith('*')) {
      return permission.startsWith(pattern.slice(0, -1))
    }
    return false
  })
}
