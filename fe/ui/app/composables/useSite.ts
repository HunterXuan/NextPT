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

export type SiteTaskCycle = 'once' | 'weekly' | 'monthly'
export type SiteTaskRuleType = 'catalog.torrent_published' | 'tracker.seed_duration' | 'tracker.uploaded' | 'iam.role_level_reached'
export type SiteTaskRewardType = 'bonus' | 'vip' | 'invite'

export interface SiteTaskReward {
  type: SiteTaskRewardType
  amount: number
}

export interface SiteTaskRule {
  type: SiteTaskRuleType
  target: number
}

export interface SiteUserTask {
  id: number
  taskKey: string
  cycleKey: string
  status: number
  progress: number
  target: number
  cycleStartedAt?: string | null
  cycleEndedAt?: string | null
  claimedAt?: string | null
  completedAt?: string | null
  rewardedAt?: string | null
}

export interface SiteTask {
  key: string
  enabled: boolean
  cycle: SiteTaskCycle
  nameI18n: Record<string, string>
  descriptionI18n: Record<string, string>
  rule: SiteTaskRule
  rewards: SiteTaskReward[]
  userTask?: SiteUserTask | null
}

export interface SiteTaskListOut {
  list: SiteTask[]
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

  async function listTasks() {
    return await fetchApi<SiteTaskListOut>('/api/site/tasks')
  }

  async function claimTask(key: string) {
    return await fetchApi<{ userTask: SiteUserTask }>(`/api/site/tasks/${key}:claim`, { method: 'POST' })
  }

  async function claimTaskReward(id: number) {
    await fetchApi(`/api/site/user-tasks/${id}:claimReward`, { method: 'POST' })
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
    listTasks,
    claimTask,
    claimTaskReward,
    listAnnouncements,
    markAnnouncementRead,
    listMessages,
    unreadMessageCount,
    markMessageRead,
    markAllMessagesRead
  }
}
