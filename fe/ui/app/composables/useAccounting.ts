export interface TrafficSummary {
  uploaded: number
  downloaded: number
  shareRatio: number
  seedTime: number
  leechTime: number
}

export type TrafficHistoryPeriod = 'daily' | 'monthly'

export interface TrafficHistoryItem {
  date: string
  uploaded: number
  downloaded: number
  seedTime: number
  leechTime: number
  bonus: string
}

export interface TrafficHistoryParams {
  period: TrafficHistoryPeriod
  startDate?: string
  endDate?: string
}

export interface TrafficHistoryListOut {
  list: TrafficHistoryItem[]
}

export interface SnatchItem {
  id: number
  torrentId: number
  torrentName: string
  torrentSize: number
  uploaded: number
  downloaded: number
  seedTime: number
  leechTime: number
  isFinished: boolean
  startedAt: string | null
  finishedAt: string | null
  lastActionAt: string | null
}

export interface SnatchListParams {
  page?: number
  size?: number
  isFinished?: boolean
}

export interface SnatchListOut {
  page: number
  size: number
  total: number
  list: SnatchItem[]
}

export type PeerStatus = 'all' | 'seeding' | 'leeching'

export interface PeerItem {
  torrentId: number
  torrentName: string
  torrentSize: number
  uploaded: number
  downloaded: number
  remaining: number
  isSeeder: boolean
  agent: string
  startedAt: string | null
  finishedAt: string | null
  lastActionAt: string | null
}

export interface PeerListParams {
  page?: number
  size?: number
  status?: PeerStatus
}

export interface PeerListOut {
  page: number
  size: number
  total: number
  seedingTotal: number
  leechingTotal: number
  list: PeerItem[]
}

export function useAccounting() {
  async function getTraffic() {
    return await fetchApi<TrafficSummary>('/api/accounting/users/me/traffic')
  }

  async function listTrafficHistory(params: TrafficHistoryParams) {
    const query: Record<string, unknown> = {
      period: params.period
    }

    if (params.startDate) query.startDate = params.startDate
    if (params.endDate) query.endDate = params.endDate

    return await fetchApi<TrafficHistoryListOut>('/api/accounting/users/me/traffic-history', { query })
  }

  async function listSnatches(params: SnatchListParams = {}) {
    const query: Record<string, unknown> = {}

    if (params.page) query.page = params.page
    if (params.size) query.size = params.size
    if (typeof params.isFinished === 'boolean') query.isFinished = params.isFinished

    return await fetchApi<SnatchListOut>('/api/accounting/users/me/snatches', { query })
  }

  async function listPeers(params: PeerListParams = {}) {
    const query: Record<string, unknown> = {}

    if (params.page) query.page = params.page
    if (params.size) query.size = params.size
    if (params.status && params.status !== 'all') query.status = params.status

    return await fetchApi<PeerListOut>('/api/accounting/users/me/peers', { query })
  }

  return {
    getTraffic,
    listTrafficHistory,
    listSnatches,
    listPeers
  }
}
