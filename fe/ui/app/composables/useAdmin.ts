import type { I18nName } from '~/types/i18n'

export interface AdminCatalogCategory {
  id: number
  nameI18N: I18nName
  slug: string
  sortOrder: number
  enabled: boolean
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminCatalogCategoryListOut {
  categories: AdminCatalogCategory[]
}

export interface AdminCatalogCategoryInput {
  nameI18N: I18nName
  slug: string
  sortOrder: number
  enabled: boolean
}

export interface AdminForumCategory {
  id: number
  nameI18N: I18nName
  descI18N: I18nName
  sortOrder: number
  minRoleView: number
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminForumCategoryListOut {
  categories: AdminForumCategory[]
}

export interface AdminForumCategoryInput {
  nameI18N: I18nName
  descI18N: I18nName
  sortOrder: number
  minRoleView: number
}

export interface AdminForumNode {
  id: number
  categoryId: number
  slug: string
  nameI18N: I18nName
  descI18N: I18nName
  sortOrder: number
  minRoleRead: number
  minRoleWrite: number
  minRoleCreate: number
  topicCount: number
  replyCount: number
  moderators?: number[] | null
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminForumNodeListOut {
  nodes: AdminForumNode[]
}

export interface AdminForumNodeInput {
  categoryId: number
  slug: string
  nameI18N: I18nName
  descI18N: I18nName
  sortOrder: number
  minRoleRead: number
  minRoleWrite: number
  minRoleCreate: number
  moderators: number[]
}

export interface AdminIamUser {
  id: number
  username: string
  email: string
  passkey: string
  status: number
  role: number
  vipUntil?: string | null
  vipRemark?: string
  invitedBy?: number
  lastLogin?: string | null
  lastIp?: string
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminIamUserListParams {
  search?: string
  order?: string
  page?: number
  size?: number
}

export interface AdminIamUserListOut {
  users: AdminIamUser[]
  total: number
}

export interface AdminIamUserUpdateInput {
  status?: number
  role?: number
  passkey?: string
}

export interface AdminIamUserBanInput {
  reason: string
  durationDays: number
}

export interface AdminIamUserPermissionInput {
  permKey: string
  isDeny: boolean
}

export interface AdminIamUserStat {
  id: number
  userId: number
  uploaded: number
  downloaded: number
  seedTime: number
  leechTime: number
  bonus: number
  bonusCharity: number
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminIamUserStatUpdateInput {
  uploadedDiff?: number
  downloadedDiff?: number
  bonusDiff?: number
}

export interface AdminIamInviteGrantInput {
  amount: number
  isTemp: boolean
  expireAt?: string | null
}

export interface AdminIamRole {
  id: number
  level: number
  nameI18N: I18nName | Record<string, unknown> | null
  rules: Record<string, unknown> | null
  permissions: string[] | Record<string, unknown> | null
  isStaff: boolean
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminIamRoleListOut {
  roles: AdminIamRole[]
}

export interface AdminIamRoleInput {
  level: number
  nameI18N: I18nName
  rules: Record<string, unknown>
  permissions: string[]
  isStaff: boolean
}

export interface AdminIamPermissionListOut {
  permissions: string[]
}

export interface AdminSiteConfig {
  id: number
  group: string
  key: string
  value: unknown
  valueType: 'string' | 'int' | 'float' | 'boolean' | 'json'
  description: string
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminSiteConfigListOut {
  configs: AdminSiteConfig[]
}

export interface AdminSiteAuditItem {
  id: number
  userId: number
  action: string
  targetType: string
  targetId: number
  level: number
  ip: string
  detail: string
  createdAt?: string | null
}

export interface AdminSiteAuditListOut {
  list: AdminSiteAuditItem[]
  total: number
  page: number
  size: number
}

export interface AdminSysCronItem {
  name: string
  status: number
  registerTime?: string | null
}

export interface AdminSysCronListOut {
  list: AdminSysCronItem[]
}

export interface AdminSysCronLogItem {
  id: number
  jobName: string
  nodeIp: string
  status: number
  durationMs: number
  errorMessage: string
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminSysCronLogListOut {
  list: AdminSysCronLogItem[]
  total: number
  page: number
  size: number
}

export interface AdminModReportItem {
  id: number
  reporter_id: number
  target_type: string
  target_id: number
  reason: string
  status: number
  dealt_by: number
  dealt_comment: string
  dealt_at?: string | null
  created_at?: string | null
}

export interface AdminModReportListOut {
  page: number
  size: number
  total: number
  list: AdminModReportItem[]
}

export interface AdminModReportListParams {
  page?: number
  size?: number
  status?: number
  targetType?: string
}

export interface AdminModResolveInput {
  status: number
  comment: string
}

export interface AdminModCheaterItem {
  id: number
  user_id: number
  torrent_id: number
  uploaded: number
  downloaded: number
  announce_time: number
  seeders: number
  leechers: number
  hit_count: number
  dealt_by: number
  is_dealt: boolean
  comment: string
  created_at?: string | null
}

export interface AdminModCheaterListOut {
  page: number
  size: number
  total: number
  list: AdminModCheaterItem[]
}

export interface AdminModCheaterListParams {
  page?: number
  size?: number
  status?: number
}

export interface AdminModUserItem {
  id: number
  user_id: number
  mod_type: number
  reason: string
  expire_at?: string | null
  mod_by: number
  mod_comment: string
  is_active: boolean
  created_at?: string | null
}

export interface AdminModUserListOut {
  list: AdminModUserItem[]
}

export interface AdminModUserApplyInput {
  type: number
  reason: string
  duration: number
}

function encodeI18n(value: I18nName) {
  return JSON.stringify(value)
}

export function useAdmin() {
  async function listCatalogCategories() {
    return await fetchApi<AdminCatalogCategoryListOut>('/api/admin/catalog/categories')
  }

  async function createCatalogCategory(input: AdminCatalogCategoryInput) {
    await fetchApi('/api/admin/catalog/categories', {
      method: 'POST',
      body: input
    })
  }

  async function updateCatalogCategory(id: number, input: AdminCatalogCategoryInput) {
    await fetchApi(`/api/admin/catalog/categories/${id}`, {
      method: 'PATCH',
      body: input
    })
  }

  async function deleteCatalogCategory(id: number) {
    await fetchApi(`/api/admin/catalog/categories/${id}`, {
      method: 'DELETE'
    })
  }

  async function listForumCategories() {
    return await fetchApi<AdminForumCategoryListOut>('/api/admin/forum/categories')
  }

  async function createForumCategory(input: AdminForumCategoryInput) {
    await fetchApi('/api/admin/forum/categories', {
      method: 'POST',
      body: {
        ...input,
        nameI18N: encodeI18n(input.nameI18N),
        descI18N: encodeI18n(input.descI18N)
      }
    })
  }

  async function updateForumCategory(id: number, input: AdminForumCategoryInput) {
    await fetchApi(`/api/admin/forum/categories/${id}`, {
      method: 'PATCH',
      body: {
        ...input,
        nameI18N: encodeI18n(input.nameI18N),
        descI18N: encodeI18n(input.descI18N)
      }
    })
  }

  async function deleteForumCategory(id: number) {
    await fetchApi(`/api/admin/forum/categories/${id}`, {
      method: 'DELETE'
    })
  }

  async function listForumNodes() {
    return await fetchApi<AdminForumNodeListOut>('/api/admin/forum/nodes')
  }

  async function createForumNode(input: AdminForumNodeInput) {
    await fetchApi('/api/admin/forum/nodes', {
      method: 'POST',
      body: {
        ...input,
        nameI18N: encodeI18n(input.nameI18N),
        descI18N: encodeI18n(input.descI18N)
      }
    })
  }

  async function updateForumNode(id: number, input: AdminForumNodeInput) {
    await fetchApi(`/api/admin/forum/nodes/${id}`, {
      method: 'PATCH',
      body: {
        ...input,
        nameI18N: encodeI18n(input.nameI18N),
        descI18N: encodeI18n(input.descI18N)
      }
    })
  }

  async function deleteForumNode(id: number) {
    await fetchApi(`/api/admin/forum/nodes/${id}`, {
      method: 'DELETE'
    })
  }

  async function listIamUsers(params: AdminIamUserListParams = {}) {
    return await fetchApi<AdminIamUserListOut>('/api/admin/iam/users', {
      query: params
    })
  }

  async function updateIamUser(id: number, input: AdminIamUserUpdateInput) {
    await fetchApi(`/api/admin/iam/users/${id}`, {
      method: 'PATCH',
      body: input
    })
  }

  async function banIamUser(id: number, input: AdminIamUserBanInput) {
    await fetchApi(`/api/admin/iam/users/${id}:ban`, {
      method: 'POST',
      body: input
    })
  }

  async function deleteIamUserSessions(id: number) {
    await fetchApi(`/api/admin/iam/users/${id}/sessions`, {
      method: 'DELETE'
    })
  }

  async function getIamUserStat(id: number) {
    return await fetchApi<AdminIamUserStat>(`/api/admin/iam/users/${id}/stat`)
  }

  async function incrementIamUserStat(id: number, input: AdminIamUserStatUpdateInput) {
    await fetchApi(`/api/admin/iam/users/${id}:incrementStat`, {
      method: 'POST',
      body: input
    })
  }

  async function grantIamUserPermission(id: number, input: AdminIamUserPermissionInput) {
    await fetchApi(`/api/admin/iam/users/${id}/permissions:grant`, {
      method: 'POST',
      body: input
    })
  }

  async function revokeIamUserPermission(id: number, input: AdminIamUserPermissionInput) {
    await fetchApi(`/api/admin/iam/users/${id}/permissions:revoke`, {
      method: 'POST',
      body: input
    })
  }

  async function grantIamInvites(input: AdminIamInviteGrantInput) {
    await fetchApi('/api/admin/iam/invites:grant', {
      method: 'POST',
      body: input
    })
  }

  async function listIamRoles() {
    return await fetchApi<AdminIamRoleListOut>('/api/admin/iam/roles')
  }

  async function createIamRole(input: AdminIamRoleInput) {
    return await fetchApi<{ id: number }>('/api/admin/iam/roles', {
      method: 'POST',
      body: input
    })
  }

  async function updateIamRole(id: number, input: AdminIamRoleInput) {
    await fetchApi(`/api/admin/iam/roles/${id}`, {
      method: 'PATCH',
      body: input
    })
  }

  async function deleteIamRole(id: number) {
    await fetchApi(`/api/admin/iam/roles/${id}`, {
      method: 'DELETE'
    })
  }

  async function listIamPermissions(roleId: number) {
    return await fetchApi<AdminIamPermissionListOut>('/api/admin/iam/permissions', {
      query: { roleId }
    })
  }

  async function listSiteConfigs(group: string) {
    return await fetchApi<AdminSiteConfigListOut>(`/api/admin/site/configs/${group}`)
  }

  async function updateSiteConfig(group: string, key: string, value: unknown) {
    await fetchApi(`/api/admin/site/configs/${group}/${key}`, {
      method: 'PUT',
      body: { value }
    })
  }

  async function listSiteAudits(params: { page?: number, size?: number } = {}) {
    return await fetchApi<AdminSiteAuditListOut>('/api/admin/site/audits', {
      query: params
    })
  }

  async function listSysCrons() {
    return await fetchApi<AdminSysCronListOut>('/api/admin/sys/crons')
  }

  async function listSysCronLogs(name: string, params: { page?: number, size?: number } = {}) {
    return await fetchApi<AdminSysCronLogListOut>(`/api/admin/sys/crons/${encodeURIComponent(name)}/logs`, {
      query: params
    })
  }

  async function listModReports(params: AdminModReportListParams = {}) {
    return await fetchApi<AdminModReportListOut>('/api/admin/mod/reports', {
      query: params
    })
  }

  async function resolveModReport(id: number, input: AdminModResolveInput) {
    await fetchApi(`/api/admin/mod/reports/${id}/resolve`, {
      method: 'POST',
      body: input
    })
  }

  async function listModCheaters(params: AdminModCheaterListParams = {}) {
    return await fetchApi<AdminModCheaterListOut>('/api/admin/mod/cheaters', {
      query: params
    })
  }

  async function resolveModCheater(id: number, input: AdminModResolveInput) {
    await fetchApi(`/api/admin/mod/cheaters/${id}/resolve`, {
      method: 'POST',
      body: input
    })
  }

  async function listUserMods(id: number) {
    return await fetchApi<AdminModUserListOut>(`/api/admin/mod/users/${id}/mods`)
  }

  async function applyUserMod(id: number, input: AdminModUserApplyInput) {
    await fetchApi(`/api/admin/mod/users/${id}/mods`, {
      method: 'POST',
      body: input
    })
  }

  async function removeUserMod(id: number, modId: number) {
    await fetchApi(`/api/admin/mod/users/${id}/mods/${modId}`, {
      method: 'DELETE'
    })
  }

  async function deleteCatalogTorrent(id: number) {
    await fetchApi(`/api/admin/catalog/torrents/${id}`, {
      method: 'DELETE'
    })
  }

  async function lockForumTopic(id: number) {
    await fetchApi(`/api/admin/forum/topics/${id}:lock`, { method: 'POST' })
  }

  async function unlockForumTopic(id: number) {
    await fetchApi(`/api/admin/forum/topics/${id}:unlock`, { method: 'POST' })
  }

  async function pinForumTopic(id: number) {
    await fetchApi(`/api/admin/forum/topics/${id}:pin`, { method: 'POST' })
  }

  async function unpinForumTopic(id: number) {
    await fetchApi(`/api/admin/forum/topics/${id}:unpin`, { method: 'POST' })
  }

  async function moveForumTopic(id: number, nodeId: number) {
    await fetchApi(`/api/admin/forum/topics/${id}:move`, {
      method: 'POST',
      body: { nodeId }
    })
  }

  return {
    listCatalogCategories,
    createCatalogCategory,
    updateCatalogCategory,
    deleteCatalogCategory,
    listForumCategories,
    createForumCategory,
    updateForumCategory,
    deleteForumCategory,
    listForumNodes,
    createForumNode,
    updateForumNode,
    deleteForumNode,
    listIamUsers,
    updateIamUser,
    banIamUser,
    deleteIamUserSessions,
    getIamUserStat,
    incrementIamUserStat,
    grantIamUserPermission,
    revokeIamUserPermission,
    grantIamInvites,
    listIamRoles,
    createIamRole,
    updateIamRole,
    deleteIamRole,
    listIamPermissions,
    listSiteConfigs,
    updateSiteConfig,
    listSiteAudits,
    listSysCrons,
    listSysCronLogs,
    listModReports,
    resolveModReport,
    listModCheaters,
    resolveModCheater,
    listUserMods,
    applyUserMod,
    removeUserMod,
    deleteCatalogTorrent,
    lockForumTopic,
    unlockForumTopic,
    pinForumTopic,
    unpinForumTopic,
    moveForumTopic
  }
}
