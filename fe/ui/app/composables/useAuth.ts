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

export interface PasskeyResetOut {
  passkey: string
}

export interface AuthPermissionListOut {
  permissions: string[]
}

export const Permission = {
  AdminIamUserManage: 'admin:iam/user:*',
  AdminIamRoleManage: 'admin:iam/role:*',
  AdminIamInviteManage: 'admin:iam/invite:*',
  AdminModReportManage: 'admin:mod/report:*',
  AdminModCheaterManage: 'admin:mod/cheater:*',
  AdminModUserManage: 'admin:mod/user:*',
  AdminSiteDashboard: 'admin:site/dashboard:*',
  AdminSiteConfig: 'admin:site/config:*',
  AdminSiteAudit: 'admin:site/audit:*',
  AdminSysCronManage: 'admin:sys/cron:*',
  AdminForumCategoryManage: 'admin:forum/category:*',
  AdminForumNodeManage: 'admin:forum/node:*',
  AdminForumTopicManage: 'admin:forum/topic:*',
  AdminForumReplyManage: 'admin:forum/reply:*',
  AdminCatalogCategoryManage: 'admin:catalog/category:*',
  AdminCatalogTorrentManage: 'admin:catalog/torrent:*',
  AdminCatalogSubtitleManage: 'admin:catalog/subtitle:*',
  AdminCatalogCommentManage: 'admin:catalog/comment:*',
  IamInviteRead: 'read:iam/invite:*',
  IamInviteCreate: 'create:iam/invite:*',
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
  CatalogCommentCreate: 'create:catalog/comment:*'
} as const

export const AdminPermissions = [
  Permission.AdminIamUserManage,
  Permission.AdminIamRoleManage,
  Permission.AdminIamInviteManage,
  Permission.AdminModReportManage,
  Permission.AdminModCheaterManage,
  Permission.AdminModUserManage,
  Permission.AdminSiteDashboard,
  Permission.AdminSiteConfig,
  Permission.AdminSiteAudit,
  Permission.AdminSysCronManage,
  Permission.AdminForumCategoryManage,
  Permission.AdminForumNodeManage,
  Permission.AdminForumTopicManage,
  Permission.AdminForumReplyManage,
  Permission.AdminCatalogCategoryManage,
  Permission.AdminCatalogTorrentManage,
  Permission.AdminCatalogSubtitleManage,
  Permission.AdminCatalogCommentManage
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
    const session = await fetchApi<{ token: string }>('/api/iam/sessions', {
      method: 'POST',
      body: { username, password }
    })

    setToken(session.token)
    await fetchUser()
    return user.value
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

  async function resetPasskey() {
    const data = await fetchApi<PasskeyResetOut>('/api/iam/users/me:resetPasskey', {
      method: 'POST'
    })
    await fetchUser()
    return data
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
    register,
    fetchUser,
    fetchPermissions,
    hasPermission,
    hasAnyPermission,
    updateProfile,
    changePassword,
    resetPasskey,
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
