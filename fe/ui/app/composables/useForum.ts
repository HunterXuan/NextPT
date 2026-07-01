import type { I18nName } from '~/types/i18n'
import type { UserSummary } from '~/types/iam'

export interface ForumNodeCategory {
  id: number
  nameI18n: I18nName
  descI18n: I18nName
  nodes: ForumNode[]
}

export interface ForumNode {
  id: number
  slug: string
  nameI18n: I18nName
  descI18n: I18nName
  topicCount: number
  replyCount: number
}

export interface ForumTopicListItem {
  id: number
  subject: string
  author: UserSummary
  isLocked: boolean
  isSticky: boolean
  views: number
  replyCount: number
  lastReplyAt: string
  lastReplyUser: UserSummary
  createdAt: string
}

export interface ForumTopicAppend {
  content?: string
  created_at?: string
  createdAt?: string
}

export interface ForumTopicDetail extends ForumTopicListItem {
  content: string
  appends: ForumTopicAppend[]
  nodeId: number
  isLiked: boolean
  isBookmarked: boolean
}

export interface ForumReplyItem {
  id: number
  author: UserSummary
  content: string
  createdAt: string
  isLiked: boolean
}

export interface ForumNodeListOut {
  list: ForumNodeCategory[]
}

export interface ForumTopicListOut {
  list: ForumTopicListItem[]
  total: number
  node: ForumNode
}

export interface ForumTopicBookmarkListOut {
  list: ForumTopicListItem[]
  total: number
}

export interface ForumReplyListOut {
  list: ForumReplyItem[]
  total: number
}

export interface ForumTopicCreateOut {
  id: number
}

export interface ForumReplyCreateOut {
  id: number
}

export interface ForumTopicListParams {
  page?: number
  size?: number
}

export interface ForumTopicCreateInput {
  nodeId: number
  subject: string
  content: string
}

export function useForum() {
  async function listNodes() {
    return await fetchApi<ForumNodeListOut>('/api/forum/nodes')
  }

  async function listTopics(slug: string, params: ForumTopicListParams = {}) {
    return await fetchApi<ForumTopicListOut>(`/api/forum/nodes/${slug}/topics`, {
      query: {
        page: params.page || 1,
        size: params.size || 20
      }
    })
  }

  async function getTopic(id: number) {
    return await fetchApi<ForumTopicDetail>(`/api/forum/topics/${id}`)
  }

  async function createTopic(input: ForumTopicCreateInput) {
    return await fetchApi<ForumTopicCreateOut>('/api/forum/topics', {
      method: 'POST',
      body: input
    })
  }

  async function appendTopic(id: number, content: string) {
    await fetchApi(`/api/forum/topics/${id}:append`, {
      method: 'POST',
      body: { content }
    })
  }

  async function toggleTopicLike(id: number) {
    await fetchApi(`/api/forum/topics/${id}:like`, {
      method: 'POST'
    })
  }

  async function rewardTopic(id: number, amount: number) {
    await fetchApi(`/api/forum/topics/${id}:reward`, {
      method: 'POST',
      body: { amount }
    })
  }

  async function reportTopic(id: number, reason: string) {
    await fetchApi(`/api/forum/topics/${id}:report`, {
      method: 'POST',
      body: { reason }
    })
  }

  async function bookmarkTopic(id: number) {
    await fetchApi(`/api/forum/topics/${id}:bookmark`, {
      method: 'POST'
    })
  }

  async function unbookmarkTopic(id: number) {
    await fetchApi(`/api/forum/topics/${id}:unbookmark`, {
      method: 'POST'
    })
  }

  async function listReplies(id: number, page = 1, size = 50) {
    return await fetchApi<ForumReplyListOut>(`/api/forum/topics/${id}/replies`, {
      query: { page, size }
    })
  }

  async function createReply(id: number, content: string, replyTo = 0) {
    return await fetchApi<ForumReplyCreateOut>(`/api/forum/topics/${id}/replies`, {
      method: 'POST',
      body: { content, replyTo }
    })
  }

  async function toggleReplyLike(id: number) {
    await fetchApi(`/api/forum/replies/${id}:like`, {
      method: 'POST'
    })
  }

  async function rewardReply(id: number, amount: number) {
    await fetchApi(`/api/forum/replies/${id}:reward`, {
      method: 'POST',
      body: { amount }
    })
  }

  async function reportReply(id: number, reason: string) {
    await fetchApi(`/api/forum/replies/${id}:report`, {
      method: 'POST',
      body: { reason }
    })
  }

  async function listBookmarks(page = 1, size = 20) {
    return await fetchApi<ForumTopicBookmarkListOut>('/api/forum/bookmarks', {
      query: { page, size }
    })
  }

  return {
    listNodes,
    listTopics,
    getTopic,
    createTopic,
    appendTopic,
    toggleTopicLike,
    rewardTopic,
    reportTopic,
    bookmarkTopic,
    unbookmarkTopic,
    listReplies,
    createReply,
    toggleReplyLike,
    rewardReply,
    reportReply,
    listBookmarks
  }
}
