export interface SiteAnnouncement {
  id: number
  title: string
  content: string
  status: number
  isRead: boolean
  createdBy: number
  updatedBy: number
  publishedAt?: string | null
  createdAt?: string | null
  updatedAt?: string | null
}

export type SiteAdvertisementPlacement = 'home' | 'catalog_list' | 'forum_list'

export interface SiteAdvertisement {
  enabled: boolean
  title: string
  image: string
  url: string
}

export interface SiteAdvertisementListOut {
  placements: Partial<Record<SiteAdvertisementPlacement, SiteAdvertisement>>
}

export interface SiteAnnouncementListOut {
  list: SiteAnnouncement[]
  total: number
  page: number
  size: number
}

export interface SiteAnnouncementListParams {
  page?: number
  size?: number
  isRead?: boolean
}

export interface SiteMessage {
  id: number
  senderId: number
  sender: {
    id: number
    username: string
    avatar: string
  }
  receiverId: number
  receiver: {
    id: number
    username: string
    avatar: string
  }
  title: string
  content: string
  targetType: string
  targetId: number
  isRead: boolean
  readAt?: string | null
  createdAt?: string | null
}

export interface SiteMessageListParams {
  page?: number
  size?: number
  isRead?: boolean
}

export interface SiteMessageListOut {
  list: SiteMessage[]
  total: number
  page: number
  size: number
}

export function useSite() {
  async function listAdvertisements() {
    return await fetchApi<SiteAdvertisementListOut>('/api/site/advertisements')
  }

  async function listAnnouncements(params: SiteAnnouncementListParams = {}) {
    return await fetchApi<SiteAnnouncementListOut>('/api/site/announcements', {
      query: params
    })
  }

  async function markAnnouncementRead(id: number) {
    await fetchApi(`/api/site/announcements/${id}:read`, {
      method: 'POST'
    })
  }

  async function listMessages(params: SiteMessageListParams = {}) {
    return await fetchApi<SiteMessageListOut>('/api/site/messages', {
      query: params
    })
  }

  async function unreadMessageCount() {
    const data = await listMessages({ page: 1, size: 1, isRead: false })
    return data.total || 0
  }

  async function markMessageRead(id: number) {
    await fetchApi(`/api/site/messages/${id}:read`, {
      method: 'POST'
    })
  }

  async function markAllMessagesRead() {
    await fetchApi('/api/site/messages:readAll', {
      method: 'POST'
    })
  }

  return {
    listAdvertisements,
    listAnnouncements,
    markAnnouncementRead,
    listMessages,
    unreadMessageCount,
    markMessageRead,
    markAllMessagesRead
  }
}
