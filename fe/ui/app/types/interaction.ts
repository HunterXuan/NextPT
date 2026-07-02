import type { UserSummary } from './iam'

export interface InteractionCommentItem {
  id: number
  author: UserSummary
  content: string
  createdAt: string
  likeCount: number
  rewardCount: number
  isLiked: boolean
}
