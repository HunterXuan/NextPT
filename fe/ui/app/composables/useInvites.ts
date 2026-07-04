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
  status?: number
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
    const query: Record<string, number> = {
      page: params.page || 1,
      size: params.size || 20
    }
    if (params.status !== undefined) {
      query.status = params.status
    }

    return await fetchApi<InviteListOut>('/api/iam/invites', {
      query
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
