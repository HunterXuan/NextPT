import type { UploadConfig } from '~/composables/useCatalogTorrents'
import type { I18nName } from '~/types/i18n'

export interface AdminCatalogCategory {
  id: number
  nameI18N: I18nName
  slug: string
  sortOrder: number
  enabled: boolean
  uploadConfig?: UploadConfig | null
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
  uploadConfig?: UploadConfig | null
}

export interface AdminCatalogTorrentPinInput {
  pinWeight: number
}

export interface AdminCatalogTorrentPromotionInput {
  spState: number
  spExpireAt?: string | null
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
  avatar?: string | null
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

export interface AdminIamUserPermissionInput {
  permKeys: string[]
  isDeny: boolean
}

export interface AdminIamUserPermissionRevokeInput {
  ids: number[]
}

export interface AdminIamUserPermissionQuery {
  sourceType?: number
  wildcardOnly?: boolean
  page?: number
  size?: number
}

export interface AdminIamUserAcl {
  id: number
  userId: number
  permKey: string
  rawPermKey: string
  isDeny: boolean
  sourceType: number
  sourceId: number
  expireAt?: string | null
  isActive: boolean
  createdAt?: string | null
}

export interface AdminIamUserPermissionDetailOut {
  userAcls: AdminIamUserAcl[]
  total: number
  page: number
  size: number
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

export interface AdminIamLoginLog {
  id: number
  userId: number
  ip: string
  userAgent: string
  result: number
  failReason: string
  createdAt?: string | null
}

export interface AdminIamLoginLogListParams {
  userId?: number
  result?: number
  page?: number
  size?: number
}

export interface AdminIamLoginLogListOut {
  list: AdminIamLoginLog[]
  total: number
  page: number
  size: number
}

export interface AdminIamInviteGrantInput {
  amount: number
  targetMode: 'site' | 'roles'
  isTemp: boolean
  expireAt?: string | null
  roleIds?: number[]
}

export interface AdminIamInvite {
  id: number
  inviterId: number
  inviteeEmail: string
  inviteeId: number
  inviteeName: string
  hash: string
  status: number
  isTemporary: boolean
  expireAt?: string | null
  usedAt?: string | null
  createdAt?: string | null
}

export interface AdminIamInviteListParams {
  page?: number
  size?: number
  status?: number
}

export interface AdminIamInviteListOut {
  invites: AdminIamInvite[]
  total: number
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
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminSiteConfigListOut {
  configs: AdminSiteConfig[]
}

export interface AdminSiteAuditItem {
  id: number
  userId: number
  actor: AdminUserSummary
  action: string
  targetType: string
  targetId: number
  level: number
  ip: string
  detail: string
  createdAt?: string | null
}

export interface AdminSiteAuditListParams {
  page?: number
  size?: number
  level?: number
  action?: string
  targetType?: string
  userId?: number
  startAt?: string
  endAt?: string
}

export interface AdminSiteAuditListOut {
  list: AdminSiteAuditItem[]
  total: number
  page: number
  size: number
}

export interface AdminSiteAnnouncement {
  id: number
  title: string
  content: string
  status: number
  isRead?: boolean
  createdBy: number
  updatedBy: number
  publishedAt?: string | null
  createdAt?: string | null
  updatedAt?: string | null
}

export interface AdminSiteAnnouncementListParams {
  page?: number
  size?: number
  status?: number
}

export interface AdminSiteAnnouncementListOut {
  list: AdminSiteAnnouncement[]
  total: number
  page: number
  size: number
}

export interface AdminSiteAnnouncementInput {
  title: string
  content: string
  status: number
  publishedAt?: string | null
}

export interface AdminSiteMessage {
  id: number
  senderId: number
  sender: AdminUserSummary
  receiverId: number
  receiver: AdminUserSummary
  title: string
  content: string
  targetType: string
  targetId: number
  isRead: boolean
  readAt?: string | null
  createdAt?: string | null
}

export interface AdminSiteMessageListParams {
  page?: number
  size?: number
  receiverId?: number
  isRead?: boolean
}

export interface AdminSiteMessageListOut {
  list: AdminSiteMessage[]
  total: number
  page: number
  size: number
}

export interface AdminSiteMessageCreateInput {
  receiverIds: number[]
  title: string
  content: string
  targetType?: string
  targetId?: number
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

export interface AdminSysCronLogListParams {
  page?: number
  size?: number
  status?: number
}

export interface AdminSysCronLogListOut {
  list: AdminSysCronLogItem[]
  total: number
  page: number
  size: number
}

export interface AdminModReportItem {
  id: number
  reporterId: number
  reporter: AdminUserSummary
  targetType: string
  targetId: number
  target: AdminModReportTarget
  reason: string
  status: number
  dealtBy: number
  dealtUser: AdminUserSummary
  dealtComment: string
  dealtAt?: string | null
  createdAt?: string | null
}

export interface AdminUserSummary {
  id: number
  username: string
  avatar: string
}

export interface AdminModReportTarget {
  type: string
  id: number
  title: string
  parentType: string
  parentId: number
  status: 'normal' | 'deleted' | 'unavailable' | string
  author: AdminUserSummary
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
  userId: number
  user: AdminUserSummary
  torrentId: number
  torrent: AdminCatalogTorrentSummary
  uploaded: number
  downloaded: number
  announceTime: number
  seeders: number
  leechers: number
  hitCount: number
  dealtBy: number
  dealtUser: AdminUserSummary
  isDealt: boolean
  comment: string
  dealtComment: string
  dealtAt?: string | null
  createdAt?: string | null
}

export interface AdminCatalogTorrentSummary {
  id: number
  name: string
  size: number
  exist: boolean
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
  userId: number
  modType: number
  reason: string
  expireAt?: string | null
  modBy: number
  modComment: string
  isActive: boolean
  createdAt?: string | null
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

  async function getIamUserPermissions(id: number, query?: AdminIamUserPermissionQuery) {
    return await fetchApi<AdminIamUserPermissionDetailOut>(`/api/admin/iam/users/${id}/permissions`, { query })
  }

  async function grantIamUserPermission(id: number, input: AdminIamUserPermissionInput) {
    await fetchApi(`/api/admin/iam/users/${id}/permissions:grant`, {
      method: 'POST',
      body: input
    })
  }

  async function revokeIamUserPermission(id: number, input: AdminIamUserPermissionRevokeInput) {
    await fetchApi(`/api/admin/iam/users/${id}/permissions:revoke`, {
      method: 'POST',
      body: input
    })
  }

  async function listIamLoginLogs(params: AdminIamLoginLogListParams = {}) {
    return await fetchApi<AdminIamLoginLogListOut>('/api/admin/iam/login-logs', {
      query: params
    })
  }

  async function grantIamInvites(input: AdminIamInviteGrantInput) {
    await fetchApi('/api/admin/iam/invites:grant', {
      method: 'POST',
      body: input
    })
  }

  async function listIamInvites(params: AdminIamInviteListParams = {}) {
    return await fetchApi<AdminIamInviteListOut>('/api/admin/iam/invites', {
      query: params
    })
  }

  async function recycleIamInvite(id: number) {
    await fetchApi(`/api/admin/iam/invites/${id}:recycle`, {
      method: 'POST'
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

  async function listSiteAudits(params: AdminSiteAuditListParams = {}) {
    return await fetchApi<AdminSiteAuditListOut>('/api/admin/site/audits', {
      query: params
    })
  }

  async function listSiteAnnouncements(params: AdminSiteAnnouncementListParams = {}) {
    return await fetchApi<AdminSiteAnnouncementListOut>('/api/admin/site/announcements', {
      query: params
    })
  }

  async function createSiteAnnouncement(input: AdminSiteAnnouncementInput) {
    return await fetchApi<{ id: number }>('/api/admin/site/announcements', {
      method: 'POST',
      body: input
    })
  }

  async function updateSiteAnnouncement(id: number, input: AdminSiteAnnouncementInput) {
    await fetchApi(`/api/admin/site/announcements/${id}`, {
      method: 'PATCH',
      body: input
    })
  }

  async function deleteSiteAnnouncement(id: number) {
    await fetchApi(`/api/admin/site/announcements/${id}`, {
      method: 'DELETE'
    })
  }

  async function listSiteMessages(params: AdminSiteMessageListParams = {}) {
    return await fetchApi<AdminSiteMessageListOut>('/api/admin/site/messages', {
      query: params
    })
  }

  async function createSiteMessage(input: AdminSiteMessageCreateInput) {
    return await fetchApi<{ count: number }>('/api/admin/site/messages', {
      method: 'POST',
      body: input
    })
  }

  async function listSysCrons() {
    return await fetchApi<AdminSysCronListOut>('/api/admin/sys/crons')
  }

  async function listSysCronLogs(name: string, params: AdminSysCronLogListParams = {}) {
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

  async function pinCatalogTorrent(id: number, input: AdminCatalogTorrentPinInput) {
    await fetchApi(`/api/admin/catalog/torrents/${id}:pin`, {
      method: 'POST',
      body: input
    })
  }

  async function unpinCatalogTorrent(id: number) {
    await fetchApi(`/api/admin/catalog/torrents/${id}:unpin`, {
      method: 'POST'
    })
  }

  async function featureCatalogTorrent(id: number) {
    await fetchApi(`/api/admin/catalog/torrents/${id}:feature`, {
      method: 'POST'
    })
  }

  async function unfeatureCatalogTorrent(id: number) {
    await fetchApi(`/api/admin/catalog/torrents/${id}:unfeature`, {
      method: 'POST'
    })
  }

  async function setCatalogTorrentPromotion(id: number, input: AdminCatalogTorrentPromotionInput) {
    await fetchApi(`/api/admin/catalog/torrents/${id}:promotion`, {
      method: 'POST',
      body: input
    })
  }

  async function clearCatalogTorrentPromotion(id: number) {
    await fetchApi(`/api/admin/catalog/torrents/${id}:clearPromotion`, {
      method: 'POST'
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

  async function deleteForumTopic(id: number) {
    await fetchApi(`/api/admin/forum/topics/${id}`, {
      method: 'DELETE'
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
    deleteIamUserSessions,
    getIamUserStat,
    incrementIamUserStat,
    getIamUserPermissions,
    grantIamUserPermission,
    revokeIamUserPermission,
    listIamLoginLogs,
    grantIamInvites,
    listIamInvites,
    recycleIamInvite,
    listIamRoles,
    createIamRole,
    updateIamRole,
    deleteIamRole,
    listIamPermissions,
    listSiteConfigs,
    updateSiteConfig,
    listSiteAudits,
    listSiteAnnouncements,
    createSiteAnnouncement,
    updateSiteAnnouncement,
    deleteSiteAnnouncement,
    listSiteMessages,
    createSiteMessage,
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
    pinCatalogTorrent,
    unpinCatalogTorrent,
    featureCatalogTorrent,
    unfeatureCatalogTorrent,
    setCatalogTorrentPromotion,
    clearCatalogTorrentPromotion,
    lockForumTopic,
    unlockForumTopic,
    pinForumTopic,
    unpinForumTopic,
    moveForumTopic,
    deleteForumTopic
  }
}
