export interface BonusLogItem {
  id: number
  userId: number
  amount: number
  balanceAfter: number
  action: string
  targetType: string
  targetId: number
  period: string
  remark: string
  createdAt: string | null
}

export interface BonusLogsParams {
  page?: number
  size?: number
  action?: string
}

export interface BonusLogsOut {
  page: number
  size: number
  total: number
  list: BonusLogItem[]
}

export interface HourlyBonusOut {
  hourlyBonus: number
}

export function useEconomy() {
  async function listBonusLogs(params: BonusLogsParams = {}) {
    const query: Record<string, unknown> = {}

    if (params.page) query.page = params.page
    if (params.size) query.size = params.size
    if (params.action?.trim()) query.action = params.action.trim()

    return await fetchApi<BonusLogsOut>('/api/economy/users/me/bonus-logs', { query })
  }

  async function getHourlyBonus() {
    return await fetchApi<HourlyBonusOut>('/api/economy/users/me/hourly-bonus')
  }

  return {
    listBonusLogs,
    getHourlyBonus
  }
}
