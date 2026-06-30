import type { I18nName } from '~/types/i18n'

export interface CatalogCategory {
  id: number
  name: I18nName
  slug: string
  uploadConfig?: UploadConfig | null
}

export interface UploadConfig {
  title?: UploadTitleConfig
  fields?: UploadFieldConfig[]
}

export interface UploadTitleConfig {
  mode?: 'manual' | 'generated' | string
  allowManualOverride?: boolean
  parts?: UploadTitlePart[]
}

export interface UploadTitlePart {
  field: string
  prefix?: string
  suffix?: string
  separator?: string
}

export interface UploadFieldConfig {
  key: string
  type?: 'text' | 'textarea' | 'select' | 'multiSelect' | string
  label?: I18nName
  description?: I18nName
  placeholder?: I18nName
  required?: boolean
  options?: UploadFieldOptions
}

export interface UploadFieldOptions {
  source?: 'static' | 'tagGroup' | string
  slug?: string
  items?: UploadOptionItem[]
}

export interface UploadOptionItem {
  value: string
  label?: I18nName
}

export interface CatalogTagItem {
  id: number
  name: I18nName
  value: string
}

export interface CatalogTagGroup {
  id: number
  name: I18nName
  slug: string
  categories?: number[] | null
  tags: CatalogTagItem[]
}

export interface TorrentListItem {
  id: number
  name: string
  subTitle: string
  categoryId: number
  size: number
  fileCount: number
  spState: number
  spExpireAt: string
  isFeatured: boolean
  isPinned: boolean
  seeders: number
  leechers: number
  snatched: number
  likeCount: number
  ownerId: number
  ownerName: string
  anonymous: boolean
  createdAt: string
}

export interface TorrentDetail extends TorrentListItem {
  description: string
  releaseFields?: Record<string, unknown> | null
  isBookmarked: boolean
  isLiked: boolean
}

export interface TorrentFileItem {
  path: string
  size: number
}

export interface TorrentPeerItem {
  userId: number
  username: string
  isSeeder: boolean
  uploaded: number
  downloaded: number
  startedAt: string
}

export interface TorrentRewardItem {
  userId: number
  username: string
  amount: number
  rewardCount: number
  lastRewardAt: string
}

export interface UserSummary {
  id: number
  username: string
  avatar: string
}

export interface CommentItem {
  id: number
  author: UserSummary
  content: string
  likeCount: number
  rewardCount: number
  createdAt: string
  isLiked: boolean
}

export interface SubtitleItem {
  id: number
  torrentId: number
  userId: number
  username: string
  fileName: string
  language: string
  size: number
  createdAt: string
}

export interface TorrentListParams {
  page?: number
  size?: number
  keyword?: string
  categoryIds?: number[]
}

export interface TorrentUploadInput {
  file: File
  name?: string
  subTitle?: string
  categoryId: number
  description?: string
  releaseFields?: Record<string, unknown>
  anonymous?: boolean
}

export interface TorrentUpdateInput {
  name?: string
  subTitle?: string
  categoryId?: number
  description?: string
  releaseFields?: Record<string, unknown>
  anonymous?: boolean
}

export interface CatalogCategoryListOut {
  list: CatalogCategory[]
}

export interface CatalogTagGroupListOut {
  list: CatalogTagGroup[]
}

export interface TorrentListOut {
  list: TorrentListItem[]
  total: number
}

export interface TorrentBookmarkListOut {
  list: TorrentListItem[]
  total: number
}

export interface TorrentFileListOut {
  list: TorrentFileItem[]
}

export interface TorrentPeerListOut {
  list: TorrentPeerItem[]
}

export interface TorrentRewardListOut {
  list: TorrentRewardItem[]
  total: number
}

export interface CommentListOut {
  list: CommentItem[]
  total: number
}

export interface SubtitleListOut {
  list: SubtitleItem[]
  total: number
}

export interface TorrentUploadOut {
  torrentId: number
  infoHash: string
}

export interface TorrentToggleLikeOut {
  isLiked: boolean
}

export interface TorrentUpdateOut {
  success: boolean
}

export interface CommentCreateOut {
  id: number
}

export interface CommentToggleLikeOut {
  isLiked: boolean
}

export function useCatalogTorrents() {
  async function listCategories() {
    return await fetchApi<CatalogCategoryListOut>('/api/catalog/categories')
  }

  async function listTagGroups() {
    return await fetchApi<CatalogTagGroupListOut>('/api/catalog/tag-groups')
  }

  async function listTorrents(params: TorrentListParams) {
    const query: Record<string, unknown> = {}

    if (params.page) query.page = params.page
    if (params.size) query.size = params.size
    if (params.keyword?.trim()) query.keyword = params.keyword.trim()
    if (params.categoryIds?.length) query['categoryIds[]'] = params.categoryIds

    return await fetchApi<TorrentListOut>('/api/catalog/torrents', { query })
  }

  async function listBookmarks(page = 1, size = 20) {
    return await fetchApi<TorrentBookmarkListOut>('/api/catalog/bookmarks', {
      query: { page, size }
    })
  }

  async function getTorrent(id: number) {
    return await fetchApi<TorrentDetail>(`/api/catalog/torrents/${id}`)
  }

  async function listFiles(id: number) {
    return await fetchApi<TorrentFileListOut>(`/api/catalog/torrents/${id}/files`)
  }

  async function listPeers(id: number) {
    return await fetchApi<TorrentPeerListOut>(`/api/catalog/torrents/${id}/peers`)
  }

  async function listRewards(id: number, page = 1, size = 20) {
    return await fetchApi<TorrentRewardListOut>(`/api/catalog/torrents/${id}/rewards`, {
      query: { page, size }
    })
  }

  async function rewardTorrent(id: number, amount: number) {
    await fetchApi(`/api/catalog/torrents/${id}:reward`, {
      method: 'POST',
      body: { amount }
    })
  }

  async function listComments(id: number, page = 1, size = 20) {
    return await fetchApi<CommentListOut>(`/api/catalog/torrents/${id}/comments`, {
      query: { page, size }
    })
  }

  async function createComment(id: number, content: string) {
    return await fetchApi<CommentCreateOut>(`/api/catalog/torrents/${id}/comments`, {
      method: 'POST',
      body: { content }
    })
  }

  async function toggleCommentLike(id: number, commentId: number) {
    return await fetchApi<CommentToggleLikeOut>(`/api/catalog/torrents/${id}/comments/${commentId}:like`, {
      method: 'POST'
    })
  }

  async function rewardComment(id: number, commentId: number, amount: number) {
    await fetchApi(`/api/catalog/torrents/${id}/comments/${commentId}:reward`, {
      method: 'POST',
      body: { amount }
    })
  }

  async function reportComment(id: number, commentId: number, reason: string) {
    await fetchApi(`/api/catalog/torrents/${id}/comments/${commentId}:report`, {
      method: 'POST',
      body: { reason }
    })
  }

  async function listSubtitles(id: number, page = 1, size = 20) {
    return await fetchApi<SubtitleListOut>(`/api/catalog/torrents/${id}/subtitles`, {
      query: { page, size }
    })
  }

  async function listAllSubtitles(page = 1, size = 20) {
    return await fetchApi<SubtitleListOut>('/api/catalog/subtitles', {
      query: { page, size }
    })
  }

  async function uploadSubtitle(id: number, file: File, language: string) {
    const body = new FormData()
    body.append('file', file)
    body.append('language', language)

    return await fetchApi<{ id: number }>(`/api/catalog/torrents/${id}/subtitles`, {
      method: 'POST',
      body
    })
  }

  async function downloadSubtitle(id: number) {
    return await fetchApiBlob(`/api/catalog/subtitles/${id}:download`)
  }

  async function reportSubtitle(id: number, reason: string) {
    await fetchApi(`/api/catalog/subtitles/${id}:report`, {
      method: 'POST',
      body: { reason }
    })
  }

  async function uploadTorrent(input: TorrentUploadInput) {
    const body = new FormData()

    body.append('file', input.file)
    body.append('categoryId', String(input.categoryId))
    body.append('anonymous', input.anonymous ? 'true' : 'false')

    if (input.name?.trim()) body.append('name', input.name.trim())
    if (input.subTitle?.trim()) body.append('subTitle', input.subTitle.trim())
    if (input.description?.trim()) body.append('description', input.description.trim())
    if (input.releaseFields && Object.keys(input.releaseFields).length > 0) {
      body.append('releaseFields', JSON.stringify(input.releaseFields))
    }

    return await fetchApi<TorrentUploadOut>('/api/catalog/torrents', {
      method: 'POST',
      body
    })
  }

  async function downloadTorrent(id: number) {
    return await fetchApiBlob(`/api/catalog/torrents/${id}:download`)
  }

  async function updateTorrent(id: number, input: TorrentUpdateInput) {
    const body = {
      ...input,
      releaseFields: input.releaseFields === undefined ? undefined : JSON.stringify(input.releaseFields)
    }

    return await fetchApi<TorrentUpdateOut>(`/api/catalog/torrents/${id}`, {
      method: 'PATCH',
      body
    })
  }

  async function reportTorrent(id: number, reason: string) {
    await fetchApi(`/api/catalog/torrents/${id}:report`, {
      method: 'POST',
      body: { reason }
    })
  }

  async function toggleLike(id: number) {
    return await fetchApi<TorrentToggleLikeOut>(`/api/catalog/torrents/${id}:like`, {
      method: 'POST'
    })
  }

  async function bookmark(id: number) {
    await fetchApi(`/api/catalog/torrents/${id}:bookmark`, {
      method: 'POST'
    })
  }

  async function unbookmark(id: number) {
    await fetchApi(`/api/catalog/torrents/${id}:unbookmark`, {
      method: 'POST'
    })
  }

  return {
    listCategories,
    listTorrents,
    listBookmarks,
    getTorrent,
    listFiles,
    listPeers,
    listRewards,
    rewardTorrent,
    listComments,
    createComment,
    toggleCommentLike,
    rewardComment,
    reportComment,
    listSubtitles,
    listAllSubtitles,
    listTagGroups,
    uploadSubtitle,
    downloadSubtitle,
    reportSubtitle,
    uploadTorrent,
    downloadTorrent,
    updateTorrent,
    reportTorrent,
    toggleLike,
    bookmark,
    unbookmark
  }
}
