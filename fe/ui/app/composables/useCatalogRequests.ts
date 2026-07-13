import type { CommentCreateOut, CommentListOut, CommentToggleLikeOut } from '~/composables/useCatalogTorrents'
import type { UserSummary } from '~/types/iam'

export const CatalogRequestType = {
  Torrent: 1,
  Reseed: 2
} as const

export const CatalogRequestStatus = {
  Open: 0,
  Claimed: 1,
  Submitted: 2,
  Completed: 3,
  Cancelled: 4
} as const

export type CatalogRequestView = 'all' | 'created' | 'claimed'

export interface CatalogRequestActions {
  canClaim: boolean
  canAbandon: boolean
  canSubmit: boolean
  canComplete: boolean
  canCancel: boolean
}

export interface CatalogTorrentSummary {
  id: number
  name: string
  size: number
  exist: boolean
}

export interface CatalogRequestListItem {
  id: number
  requestType: number
  categoryId: number
  targetTorrentId: number
  resultTorrentId: number
  title: string
  rewardAmount: number
  status: number
  requester: UserSummary
  claimer?: UserSummary | null
  claimExpiresAt: string
  submittedAt: string
  completedAt: string
  createdAt: string
  updatedAt: string
}

export interface CatalogRequestDetail extends CatalogRequestListItem {
  description: string
  targetTorrent?: CatalogTorrentSummary | null
  resultTorrent?: CatalogTorrentSummary | null
  cancelReason: string
  actions: CatalogRequestActions
}

export interface CatalogRequestListParams {
  page?: number
  size?: number
  keyword?: string
  requestType?: number
  status?: number
  categoryId?: number
  view?: CatalogRequestView
}

export interface CatalogRequestCreateInput {
  requestType: number
  categoryId?: number
  targetTorrentId?: number
  title?: string
  description: string
  rewardAmount: number
}

export function useCatalogRequests() {
  async function listRequests(params: CatalogRequestListParams) {
    const query: Record<string, unknown> = {}
    if (params.page) query.page = params.page
    if (params.size) query.size = params.size
    if (params.keyword?.trim()) query.keyword = params.keyword.trim()
    if (params.requestType) query.requestType = params.requestType
    if (params.status !== undefined) query.status = params.status
    if (params.categoryId) query.categoryId = params.categoryId
    if (params.view && params.view !== 'all') query.view = params.view
    return await fetchApi<{ list: CatalogRequestListItem[], total: number }>('/api/catalog/requests', { query })
  }

  async function getRequest(id: number) {
    return await fetchApi<CatalogRequestDetail>(`/api/catalog/requests/${id}`)
  }

  async function createRequest(input: CatalogRequestCreateInput) {
    return await fetchApi<{ id: number }>('/api/catalog/requests', {
      method: 'POST',
      body: input
    })
  }

  async function claimRequest(id: number) {
    await fetchApi(`/api/catalog/requests/${id}:claim`, { method: 'POST' })
  }

  async function abandonRequest(id: number) {
    await fetchApi(`/api/catalog/requests/${id}:abandon`, { method: 'POST' })
  }

  async function submitRequest(id: number, resultTorrentId = 0) {
    await fetchApi(`/api/catalog/requests/${id}:submit`, {
      method: 'POST',
      body: resultTorrentId > 0 ? { resultTorrentId } : {}
    })
  }

  async function completeRequest(id: number) {
    await fetchApi(`/api/catalog/requests/${id}:complete`, { method: 'POST' })
  }

  async function completeRequestByAdmin(id: number) {
    await fetchApi(`/api/admin/catalog/requests/${id}:complete`, { method: 'POST' })
  }

  async function cancelRequest(id: number, reason = '') {
    await fetchApi(`/api/catalog/requests/${id}:cancel`, {
      method: 'POST',
      body: reason.trim() ? { reason: reason.trim() } : {}
    })
  }

  async function cancelRequestByAdmin(id: number, reason = '') {
    await fetchApi(`/api/admin/catalog/requests/${id}:cancel`, {
      method: 'POST',
      body: reason.trim() ? { reason: reason.trim() } : {}
    })
  }

  async function listComments(id: number, page = 1, size = 20) {
    return await fetchApi<CommentListOut>(`/api/catalog/requests/${id}/comments`, {
      query: { page, size }
    })
  }

  async function createComment(id: number, content: string) {
    return await fetchApi<CommentCreateOut>(`/api/catalog/requests/${id}/comments`, {
      method: 'POST',
      body: { content }
    })
  }

  async function toggleCommentLike(id: number, commentId: number) {
    return await fetchApi<CommentToggleLikeOut>(`/api/catalog/requests/${id}/comments/${commentId}:like`, {
      method: 'POST'
    })
  }

  async function rewardComment(id: number, commentId: number, amount: number) {
    await fetchApi(`/api/catalog/requests/${id}/comments/${commentId}:reward`, {
      method: 'POST',
      body: { amount }
    })
  }

  async function reportComment(id: number, commentId: number, reason: string) {
    await fetchApi(`/api/catalog/requests/${id}/comments/${commentId}:report`, {
      method: 'POST',
      body: { reason }
    })
  }

  return {
    listRequests,
    getRequest,
    createRequest,
    claimRequest,
    abandonRequest,
    submitRequest,
    completeRequest,
    completeRequestByAdmin,
    cancelRequest,
    cancelRequestByAdmin,
    listComments,
    createComment,
    toggleCommentLike,
    rewardComment,
    reportComment
  }
}
