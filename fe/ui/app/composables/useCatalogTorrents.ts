import type { I18nName } from '~/types/i18n'
import type { UserSummary } from '~/types/iam'

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

export type ReleaseFieldValue = string | string[]

export interface ReleaseFieldsState {
  generatedTitle: string
  valid: boolean
  firstError: string
  missingLabels: string[]
  output: Record<string, unknown>
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
  pinWeight: number
  seeders: number
  leechers: number
  snatched: number
  likeCount: number
  owner: UserSummary
  anonymous: boolean
  createdAt: string
}

export type TorrentHotItem = Omit<TorrentListItem, 'owner' | 'anonymous'> & {
  category?: CatalogCategory | null
}

export interface TorrentDetail extends TorrentListItem {
  description: string
  releaseFields?: Record<string, unknown> | null
  metadata?: TorrentMetadataOut | null
  isBookmarked: boolean
  isLiked: boolean
}

export interface TorrentMetadataBinding {
  imdbId: string
  imdbRating?: number
  doubanId: string
  doubanRating?: number
  bangumiId: string
  bangumiRating?: number
  tmdbId: string
  tmdbType: 'movie' | 'tv' | ''
  tmdbRating?: number
}

export interface TorrentMetadataItem {
  provider: string
  providerId: string
  tmdbType: 'movie' | 'tv' | ''
  title: string
  originalTitle: string
  year: string
  releaseDate: string
  overview: string
  posterUrl: string
  backdropUrl: string
  rating: number
  genres: string[]
  imdbId: string
}

export interface TorrentMetadataOut {
  binding: TorrentMetadataBinding
  data?: TorrentMetadataItem | null
  sources: TorrentMetadataItem[]
}

export interface TorrentMetadataSearchOut {
  list: TorrentMetadataItem[]
  page: number
  totalPages: number
  totalResults: number
}

export interface TorrentFileItem {
  path: string
  size: number
}

export interface TorrentPeerItem {
  user: UserSummary
  isSeeder: boolean
  uploaded: number
  downloaded: number
  startedAt: string
}

export interface TorrentRewardItem {
  user: UserSummary
  amount: number
  rewardCount: number
  lastRewardAt: string
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
  uploader: UserSummary
  anonymous: boolean
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
  promotion?: string
  seedStatus?: string
  featuredOnly?: boolean
  minSize?: number
  maxSize?: number
  publishedWithin?: number
  sort?: string
  imdbId?: string
  doubanId?: string
  bangumiId?: string
  tmdbId?: string
  tmdbType?: string
}

export interface TorrentAdvancedFilters {
  promotion: string
  seedStatus: string
  featuredOnly: boolean
  minSize: number
  maxSize: number
  publishedWithin: number
  sort: string
  imdbId: string
  doubanId: string
  bangumiId: string
  tmdbId: string
  tmdbType: string
}

export interface TorrentUploadInput {
  file: File
  name?: string
  subTitle?: string
  categoryId: number
  description?: string
  releaseFields?: Record<string, unknown>
  metadata?: TorrentMetadataBinding
  anonymous?: boolean
}

export interface TorrentUpdateInput {
  name?: string
  subTitle?: string
  categoryId?: number
  description?: string
  releaseFields?: Record<string, unknown>
  metadata?: TorrentMetadataBinding
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

export interface TorrentHotListOut {
  list: TorrentHotItem[]
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
    if (params.promotion && params.promotion !== 'all') query.promotion = params.promotion
    if (params.seedStatus && params.seedStatus !== 'all') query.seedStatus = params.seedStatus
    if (params.featuredOnly) query.featuredOnly = true
    if (params.minSize && params.minSize > 0) query.minSize = params.minSize
    if (params.maxSize && params.maxSize > 0) query.maxSize = params.maxSize
    if (params.publishedWithin && params.publishedWithin > 0) query.publishedWithin = params.publishedWithin
    if (params.sort && params.sort !== 'newest') query.sort = params.sort
    if (params.imdbId?.trim()) query.imdbId = params.imdbId.trim().toLowerCase()
    if (params.doubanId?.trim()) query.doubanId = params.doubanId.trim()
    if (params.bangumiId?.trim()) query.bangumiId = params.bangumiId.trim()
    if (params.tmdbId?.trim()) {
      query.tmdbId = params.tmdbId.trim()
      if (params.tmdbType && params.tmdbType !== 'all') query.tmdbType = params.tmdbType
    }

    return await fetchApi<TorrentListOut>('/api/catalog/torrents', { query })
  }

  async function listHotTorrents(size = 5) {
    return await fetchApi<TorrentHotListOut>('/api/catalog/torrents:getHot', {
      query: { size }
    })
  }

  async function searchMetadata(query: string, tmdbType: 'movie' | 'tv', page = 1) {
    return await fetchApi<TorrentMetadataSearchOut>('/api/catalog/torrent-metadata:search', {
      query: { query, tmdbType, page }
    })
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

  async function uploadSubtitle(id: number, file: File, language: string, anonymous = false) {
    const body = new FormData()
    body.append('file', file)
    body.append('language', language)
    body.append('anonymous', anonymous ? 'true' : 'false')

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
    if (input.metadata) body.append('metadata', JSON.stringify(input.metadata))

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
      releaseFields: input.releaseFields === undefined ? undefined : JSON.stringify(input.releaseFields),
      metadata: input.metadata === undefined ? undefined : JSON.stringify(input.metadata)
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
    listHotTorrents,
    searchMetadata,
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
