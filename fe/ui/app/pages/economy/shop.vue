<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[260px_minmax(0,1fr)] xl:items-start">
        <aside class="app-sticky-offset xl:sticky">
          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-3 py-3 dark:border-slate-800">
              <div class="flex min-w-0 items-center gap-2.5">
                <span class="flex size-8 shrink-0 items-center justify-center rounded-md bg-amber-50 text-amber-600 dark:bg-amber-950 dark:text-amber-300">
                  <UIcon name="i-lucide-coins" class="size-4" />
                </span>
                <span class="truncate text-sm text-slate-600 dark:text-slate-300">{{ $t('user.shop.balance') }}</span>
              </div>
              <span class="shrink-0 text-base font-semibold tabular-nums text-slate-950 dark:text-white">{{ formatBonus(user?.stat.bonus) }}</span>
            </div>
            <nav class="flex gap-2 overflow-x-auto p-2 xl:block xl:space-y-1 xl:overflow-visible">
              <button
                v-for="section in shopSections"
                :key="section.value"
                type="button"
                class="flex min-w-max items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition xl:w-full xl:min-w-0"
                :class="activeSection === section.value
                  ? 'bg-slate-100 text-slate-950 dark:bg-slate-800 dark:text-white'
                  : 'text-slate-600 hover:bg-slate-50 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-950 dark:hover:text-white'"
                :aria-current="activeSection === section.value ? 'page' : undefined"
                @click="setActiveSection(section.value)"
              >
                <span
                  class="flex size-8 shrink-0 items-center justify-center rounded-md"
                  :class="activeSection === section.value
                    ? 'bg-white text-slate-950 dark:bg-slate-950 dark:text-white'
                    : 'bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-400'"
                >
                  <UIcon :name="section.icon" class="size-4" />
                </span>
                <span class="truncate font-medium">{{ section.label }}</span>
              </button>
            </nav>
          </section>
        </aside>

        <main class="min-w-0">
          <section v-show="activeSection === 'products'" class="min-w-0 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('user.shop.productsTitle') }}</h2>
          </div>

          <div v-if="productsPending" class="space-y-2 p-4">
            <div v-for="item in 3" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div v-else-if="productsError" class="flex flex-col items-center justify-center px-4 py-10 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ productsError }}</p>
          </div>
          <div v-else-if="products.length === 0" class="flex flex-col items-center justify-center px-4 py-10 text-center">
            <UIcon name="i-lucide-package-open" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.shop.emptyProducts') }}</p>
          </div>
          <div v-else class="divide-y divide-slate-100 dark:divide-slate-800">
            <div v-for="product in products" :key="product.key" class="grid gap-4 px-4 py-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center sm:px-5 sm:py-5">
              <div class="flex min-w-0 items-center gap-3.5">
                <span class="flex size-11 shrink-0 items-center justify-center rounded-md bg-amber-50 text-amber-700 dark:bg-amber-950 dark:text-amber-300">
                  <UIcon :name="productIcon(product.type)" class="size-5.5" />
                </span>
                <div class="min-w-0">
                  <h3 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ productTitle(product) }}</h3>
                  <p class="mt-1 max-w-2xl text-xs leading-5 text-slate-500 dark:text-slate-400">{{ productDescription(product) }}</p>
                </div>
              </div>
              <div class="flex items-center justify-between gap-3 sm:justify-end">
                <div class="min-w-24 text-left sm:text-right">
                  <p class="flex items-center gap-1 text-lg font-semibold tabular-nums text-amber-600 sm:justify-end dark:text-amber-300">
                    <UIcon name="i-lucide-coins" class="size-3.5" />
                    {{ formatBonus(product.price) }}
                  </p>
                  <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ productOptionLabel(product) }}</p>
                </div>
                <span class="hidden h-9 w-px bg-slate-200 sm:block dark:bg-slate-800" />
                <UPopover :content="{ side: 'left', align: 'center', sideOffset: 10 }" :ui="{ content: 'w-72 p-3' }">
                  <AppPermissionButton
                    :permission="Permission.EconomyShopOrderCreate"
                    color="primary"
                    variant="soft"
                    size="md"
                    icon="i-lucide-shopping-bag"
                    :loading="purchasingKey === product.key"
                    :disabled="purchasingKey !== '' || !canAfford(product)"
                    :tooltip="canAfford(product) ? $t('user.shop.redeem') : $t('user.shop.insufficient')"
                    :aria-label="$t('user.shop.redeem')"
                  />

                  <template #content="{ close }">
                    <div class="space-y-3">
                      <div>
                        <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('user.shop.confirmTitle') }}</p>
                        <p class="mt-1 text-xs leading-5 text-slate-500 dark:text-slate-400">
                          {{ $t('user.shop.confirmDescription', { product: productTitle(product), price: formatBonus(product.price) }) }}
                        </p>
                      </div>
                      <div class="flex justify-end gap-2">
                        <UButton color="neutral" variant="ghost" size="xs" @click="close()">{{ $t('common.cancel') }}</UButton>
                        <UButton color="primary" size="xs" icon="i-lucide-shopping-bag" :loading="purchasingKey === product.key" :disabled="!canCreateOrders || !canAfford(product)" @click="purchase(product, close)">
                          {{ $t('user.shop.redeem') }}
                        </UButton>
                      </div>
                    </div>
                  </template>
                </UPopover>
              </div>
            </div>
          </div>
          </section>

          <section v-if="canReadOrders" v-show="activeSection === 'orders'" class="min-w-0 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('user.shop.ordersTitle') }}</h2>
            </div>
            <div v-if="ordersPending" class="space-y-2 p-3">
              <div v-for="item in 3" :key="item" class="h-14 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>
            <div v-else-if="ordersError" class="px-4 py-10 text-center text-sm text-red-700 dark:text-red-200">{{ ordersError }}</div>
            <div v-else-if="orders.length === 0" class="flex flex-col items-center justify-center px-4 py-10 text-center">
              <UIcon name="i-lucide-receipt-text" class="size-8 text-slate-300 dark:text-slate-600" />
              <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.shop.emptyOrders') }}</p>
            </div>
            <div v-else>
              <div class="hidden grid-cols-[minmax(0,1fr)_140px_120px_180px] gap-4 border-b border-slate-200 bg-slate-50 px-4 py-2.5 text-xs font-medium text-slate-500 md:grid dark:border-slate-800 dark:bg-slate-950/60 dark:text-slate-400">
                <span>{{ $t('user.shop.table.product') }}</span>
                <span>{{ $t('user.shop.table.price') }}</span>
                <span>{{ $t('user.shop.table.status') }}</span>
                <span class="text-right">{{ $t('user.shop.table.time') }}</span>
              </div>
              <div v-for="order in orders" :key="order.id" class="grid gap-2 border-b border-slate-100 px-4 py-3.5 text-sm last:border-b-0 md:grid-cols-[minmax(0,1fr)_140px_120px_180px] md:items-center md:gap-4 dark:border-slate-800">
                <div class="flex min-w-0 items-center gap-3">
                  <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300">
                    <UIcon :name="productIcon(order.product.type)" class="size-4.5" />
                  </span>
                  <div class="min-w-0">
                    <p class="truncate font-medium text-slate-950 dark:text-white">{{ productTitle(order.product) }}</p>
                    <NuxtLink v-if="order.targetType === 'iam_invite'" :to="localePath('/iam/users/me/settings?section=invites')" class="mt-1 inline-flex items-center gap-1 text-xs text-sky-700 hover:text-sky-800 dark:text-sky-300 dark:hover:text-sky-200">
                      <UIcon name="i-lucide-ticket" class="size-3.5" />
                      {{ $t('user.shop.viewInvites') }}
                    </NuxtLink>
                  </div>
                </div>
                <span class="flex items-center gap-1 font-semibold tabular-nums text-amber-600 dark:text-amber-300">
                  <UIcon name="i-lucide-coins" class="size-3.5" />
                  -{{ formatBonus(order.price) }}
                </span>
                <UBadge color="success" variant="soft" class="w-fit">{{ $t('user.shop.status.completed') }}</UBadge>
                <span class="text-xs text-slate-500 md:text-right dark:text-slate-400">{{ formatDateTime(order.completedAt || order.createdAt, locale) }}</span>
              </div>
            </div>
            <AppPager
              v-if="orderTotal > 0"
              class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
              :page="orderPage"
              :total="orderTotal"
              :page-size="orderSize"
              :page-size-options="[20, 50, 100]"
              :disabled="ordersPending"
              @page-change="changeOrderPage"
              @page-size-change="changeOrderSize"
            />
          </section>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { ShopOrder, ShopProduct } from '~/composables/useEconomy'
import { formatDateTime } from '~/utils/format'

type ShopSection = 'products' | 'orders'

definePageMeta({ middleware: 'auth' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const toast = useToast()
const economyApi = useEconomy()
const { user, fetchUser, hasPermission } = useAuth()

const products = ref<ShopProduct[]>([])
const productsPending = ref(true)
const productsError = ref('')
const purchasingKey = ref('')
const orders = ref<ShopOrder[]>([])
const ordersPending = ref(false)
const ordersError = ref('')
const ordersLoaded = ref(false)
const orderPage = ref(1)
const orderSize = ref(20)
const orderTotal = ref(0)

const canReadOrders = computed(() => hasPermission(Permission.EconomyShopOrderRead))
const canCreateOrders = computed(() => hasPermission(Permission.EconomyShopOrderCreate))
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value, { maximumFractionDigits: 2 }))
const activeSection = ref<ShopSection>(normalizeSection(route.query.section))
const shopSections = computed(() => {
  const sections: Array<{ value: ShopSection; label: string; icon: string }> = [
    {
      value: 'products',
      label: t('user.shop.productsTitle'),
      icon: 'i-lucide-store'
    }
  ]

  if (canReadOrders.value) {
    sections.push({
      value: 'orders',
      label: t('user.shop.ordersTitle'),
      icon: 'i-lucide-receipt-text'
    })
  }

  return sections
})

useSeoMeta({ title: t('user.shop.title'), robots: 'noindex, nofollow' })

onMounted(async () => {
  await Promise.all([
    loadProducts(),
    activeSection.value === 'orders' && canReadOrders.value ? loadOrders() : Promise.resolve()
  ])
})

watch(
  () => route.query.section,
  section => setActiveSectionValue(normalizeSection(section))
)

watch(canReadOrders, (canRead) => {
  if (!canRead && activeSection.value === 'orders') {
    setActiveSection('products')
  }
})

async function loadProducts() {
  productsPending.value = true
  productsError.value = ''
  try {
    const data = await economyApi.listShopProducts()
    products.value = data.list || []
  } catch (error) {
    products.value = []
    productsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    productsPending.value = false
  }
}

async function loadOrders() {
  ordersPending.value = true
  ordersError.value = ''
  try {
    const data = await economyApi.listShopOrders(orderPage.value, orderSize.value)
    orders.value = data.list || []
    orderTotal.value = data.total || 0
    ordersLoaded.value = true
  } catch (error) {
    orders.value = []
    orderTotal.value = 0
    ordersError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    ordersPending.value = false
  }
}

async function purchase(product: ShopProduct, close: () => void) {
  if (purchasingKey.value || !canCreateOrders.value || !canAfford(product)) return
  purchasingKey.value = product.key
  try {
    await economyApi.createShopOrder(product.key)
    close()
    ordersLoaded.value = false
    await Promise.all([
      fetchUser(),
      canReadOrders.value && activeSection.value === 'orders' ? loadOrders() : Promise.resolve()
    ])
    toast.add({ title: t('user.shop.success'), color: 'success', icon: 'i-lucide-circle-check' })
  } catch (error) {
    toast.add({ title: error instanceof ApiError ? error.message : t('common.requestFailed'), color: 'error', icon: 'i-lucide-circle-alert' })
  } finally {
    purchasingKey.value = ''
  }
}

function canAfford(product: ShopProduct) {
  return Number(user.value?.stat.bonus || 0) >= Number(product.price || 0)
}

function formatBonus(value?: number | null) {
  return numberFormatter.value.format(Number(value || 0))
}

function productTitle(product: ShopProduct) {
  const key = `user.shop.products.${product.type}.title`
  const translated = t(key)
  return translated === key ? product.key : translated
}

function productDescription(product: ShopProduct) {
  const key = `user.shop.products.${product.type}.description`
  const translated = t(key)
  return translated === key ? product.type : translated
}

function productOptionLabel(product: ShopProduct) {
  if (product.type === 'invite') return t('user.shop.products.invite.amount', { count: Number(product.options?.amount || 1) })
  return product.type
}

function productIcon(type: string) {
  if (type === 'invite') return 'i-lucide-ticket-plus'
  if (type === 'vip') return 'i-lucide-crown'
  if (type === 'upload') return 'i-lucide-arrow-up-to-line'
  if (type === 'download') return 'i-lucide-arrow-down-to-line'
  return 'i-lucide-ticket-percent'
}

function normalizeSection(section: unknown): ShopSection {
  return section === 'orders' && canReadOrders.value ? 'orders' : 'products'
}

function setActiveSection(section: ShopSection) {
  if (activeSection.value === section) return

  const query = { ...route.query }
  if (section === 'products') {
    delete query.section
  } else {
    query.section = section
  }
  router.replace({ query })
}

function setActiveSectionValue(section: ShopSection) {
  activeSection.value = section
  if (section === 'orders' && canReadOrders.value && !ordersLoaded.value && !ordersPending.value) {
    loadOrders()
  }
}

function changeOrderPage(page: number) {
  orderPage.value = page
  loadOrders()
}

function changeOrderSize(size: number) {
  orderSize.value = size
  orderPage.value = 1
  loadOrders()
}
</script>
