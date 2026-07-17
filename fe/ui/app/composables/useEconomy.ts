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

export interface ShopProduct {
  key: string
  type: string
  enabled: boolean
  price: number
  sortOrder: number
  options: Record<string, unknown>
}

export interface ShopProductListOut {
  list: ShopProduct[]
}

export interface ShopOrderCreateOut {
  orderId: number
  productKey: string
  price: number
  targetType: string
  targetId: number
  balanceAfter: number
}

export interface ShopOrder {
  id: number
  product: ShopProduct
  price: number
  status: number
  targetType: string
  targetId: number
  createdAt?: string | null
  completedAt?: string | null
}

export interface ShopOrderListOut {
  page: number
  size: number
  total: number
  list: ShopOrder[]
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

  async function listShopProducts() {
    return await fetchApi<ShopProductListOut>('/api/economy/shop/products')
  }

  async function createShopOrder(productKey: string) {
    return await fetchApi<ShopOrderCreateOut>('/api/economy/shop/orders', {
      method: 'POST',
      body: { productKey }
    })
  }

  async function listShopOrders(page = 1, size = 20) {
    return await fetchApi<ShopOrderListOut>('/api/economy/shop/orders', {
      query: { page, size }
    })
  }

  return {
    listBonusLogs,
    getHourlyBonus,
    listShopProducts,
    createShopOrder,
    listShopOrders
  }
}
