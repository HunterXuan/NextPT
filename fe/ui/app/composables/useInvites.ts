export interface InviteItem {
  id: number
  inviterId: number
  inviteeEmail: string
  inviteeId: number
  inviteeName: string
  hash: string
  status: number
  isTemporary: boolean
  expireAt: string | null
  usedAt: string | null
  createdAt: string | null
}

export interface InviteListParams {
  page?: number
  size?: number
}

export interface InviteListOut {
  list: InviteItem[]
  total: number
}

export interface InviteCheckOut {
  hash: string
  inviterId: number
  inviterUsername: string
}

export function useInvites() {
  async function listInvites(params: InviteListParams = {}) {
    return await fetchApi<InviteListOut>('/api/iam/invites', {
      query: {
        page: params.page || 1,
        size: params.size || 20
      }
    })
  }

  async function sendInvite(hash: string, email: string) {
    await fetchApi('/api/iam/invites:send', {
      method: 'POST',
      body: { hash, email }
    })
  }

  async function checkInvite(hash: string) {
    return await fetchApi<InviteCheckOut>('/api/iam/invites:check', {
      query: { hash }
    })
  }

  return {
    listInvites,
    sendInvite,
    checkInvite
  }
}
