import { fetchApi } from '~/composables/useApi'
import type { UserSummary } from '~/types/iam'

export const ModStaffMessageStatus = {
  Pending: 0,
  Processed: 1
} as const

export type ModStaffMessageStatusValue = typeof ModStaffMessageStatus[keyof typeof ModStaffMessageStatus]

export interface StaffMessage {
  id: number
  senderId: number
  sender: UserSummary
  subject: string
  content: string
  status: ModStaffMessageStatusValue
  answeredBy: number
  answeredByUser: UserSummary
  answer: string
  answeredAt?: string | null
  createdAt?: string | null
  updatedAt?: string | null
}

export interface StaffMessageListParams {
  page?: number
  size?: number
  status?: number
}

export interface StaffMessageListOut {
  list: StaffMessage[]
  total: number
  page: number
  size: number
}

export interface StaffMessageCreateInput {
  subject: string
  content: string
}

export function useMod() {
  async function listStaffMessages(params: StaffMessageListParams = {}) {
    return await fetchApi<StaffMessageListOut>('/api/mod/staff-messages', { query: params })
  }

  async function createStaffMessage(input: StaffMessageCreateInput) {
    return await fetchApi<{ id: number }>('/api/mod/staff-messages', {
      method: 'POST',
      body: input
    })
  }

  return {
    listStaffMessages,
    createStaffMessage
  }
}
