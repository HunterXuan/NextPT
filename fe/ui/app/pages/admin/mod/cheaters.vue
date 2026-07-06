<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_460px] xl:items-start">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.cheaters.list') }}</h2>
              <form class="grid gap-2 sm:w-[160px]" @submit.prevent="reloadFromFirstPage">
                <select v-model.number="query.status" class="h-9 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" @change="reloadFromFirstPage">
                  <option v-for="option in statusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
                </select>
              </form>
            </div>
          </div>

          <div v-if="pending" class="space-y-2 p-4">
            <div v-for="item in 6" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>

          <div v-else-if="cheaters.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.mod.cheaters.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[800px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[32%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.cheaters.table.subject') }}</th>
                  <th class="w-[22%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.cheaters.table.traffic') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.cheaters.table.swarm') }}</th>
                  <th class="w-[11%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.status') }}</th>
                  <th class="w-[15%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.mod.table.createdAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in cheaters"
                  :key="item.id"
                  class="cursor-pointer border-b border-slate-200 outline-none transition last:border-b-0 focus-visible:bg-slate-50 dark:border-slate-800 dark:focus-visible:bg-slate-950/70"
                  :class="selectedCheater?.id === item.id ? 'bg-amber-50/80 dark:bg-amber-950/25' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                  tabindex="0"
                  @click="selectCheater(item)"
                  @keydown.enter.prevent="selectCheater(item)"
                >
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <div class="flex min-w-0 items-center gap-3">
                      <IamUserAvatar :user="cheaterUser(item)" size="sm" />
                      <div class="min-w-0">
                        <p class="truncate font-semibold text-slate-950 dark:text-white">{{ userDisplayName(cheaterUser(item), item.user_id) }}</p>
                        <div class="mt-1 flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                          <UIcon name="i-lucide-database" class="size-3.5 shrink-0" />
                          <span class="truncate">{{ torrentDisplayName(item) }}</span>
                        </div>
                      </div>
                    </div>
                  </td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <div class="space-y-1">
                      <p class="flex items-center gap-2 font-medium text-emerald-700 dark:text-emerald-300">
                        <UIcon name="i-lucide-arrow-up" class="size-3.5 shrink-0" />
                        <span>{{ formatBytes(item.uploaded) }}</span>
                      </p>
                      <p class="flex items-center gap-2 text-xs font-medium text-sky-700 dark:text-sky-300">
                        <UIcon name="i-lucide-arrow-down" class="size-3.5 shrink-0" />
                        <span>{{ formatBytes(item.downloaded) }}</span>
                      </p>
                    </div>
                  </td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <p class="font-medium text-slate-800 dark:text-slate-100">{{ $t('admin.mod.cheaters.announceTime', { value: numberFormatter.format(item.announce_time) }) }}</p>
                    <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.swarm', { seeders: item.seeders, leechers: item.leechers, hits: item.hit_count }) }}</p>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <UBadge :color="item.is_dealt ? 'success' : 'warning'" variant="soft" class="whitespace-nowrap">
                      {{ item.is_dealt ? $t('admin.mod.status.resolved') : $t('admin.mod.status.pending') }}
                    </UBadge>
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    <UTooltip :text="formatDateTime(item.created_at, locale)" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                      <span>{{ relativeTime(item.created_at) }}</span>
                    </UTooltip>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <AppPager
            class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
            size="sm"
            :page="query.page"
            :total="total"
            :page-size="query.size"
            :disabled="pending"
            @page-change="changePage"
          />
        </section>

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900 xl:sticky xl:top-20">
          <div v-if="!selectedCheater" class="flex min-h-[360px] flex-col items-center justify-center px-6 py-12 text-center">
            <span class="flex size-12 items-center justify-center rounded-lg bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-500">
              <UIcon name="i-lucide-mouse-pointer-2" class="size-5" />
            </span>
            <h2 class="mt-4 text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.cheaters.select') }}</h2>
            <p class="mt-2 max-w-xs text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.selectHint') }}</p>
          </div>

          <template v-else>
            <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <div class="min-w-0">
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.current') }}</p>
                <h2 class="mt-1 text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.cheaters.logId', { id: selectedCheater.id }) }}</h2>
              </div>
              <UBadge :color="selectedCheater.is_dealt ? 'success' : 'warning'" variant="soft" class="whitespace-nowrap">
                {{ selectedCheater.is_dealt ? $t('admin.mod.status.resolved') : $t('admin.mod.status.pending') }}
              </UBadge>
            </div>

            <div class="divide-y divide-slate-200 dark:divide-slate-800">
              <div class="space-y-5 p-4">
                <section class="space-y-3">
                  <div class="flex items-start justify-between gap-3">
                    <div class="flex min-w-0 items-center gap-3">
                      <IamUserAvatar :user="cheaterUser(selectedCheater)" size="md" />
                      <div class="min-w-0">
                        <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ userDisplayName(cheaterUser(selectedCheater), selectedCheater.user_id) }}</p>
                        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.userId', { id: selectedCheater.user_id }) }}</p>
                      </div>
                    </div>
                  </div>

                  <div class="rounded-md border border-slate-200 bg-slate-50/70 px-3 py-2 dark:border-slate-800 dark:bg-slate-950/50">
                    <div class="flex items-start justify-between gap-3">
                      <div class="min-w-0">
                        <p class="line-clamp-2 text-sm font-semibold leading-5 text-slate-950 dark:text-white">{{ torrentDisplayName(selectedCheater) }}</p>
                        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.torrentId', { id: selectedCheater.torrent_id }) }}</p>
                      </div>
                      <UTooltip v-if="selectedCheater.torrent?.exist" :text="$t('admin.mod.cheaters.openTorrent')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                        <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-external-link" :to="localePath(`/catalog/torrents/${selectedCheater.torrent_id}`)" :aria-label="$t('admin.mod.cheaters.openTorrent')" />
                      </UTooltip>
                    </div>
                  </div>
                </section>

                <dl class="grid grid-cols-3 gap-3 text-sm">
                  <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.fields.uploaded') }}</dt>
                    <dd class="mt-1 font-semibold text-slate-950 dark:text-white">{{ formatBytes(selectedCheater.uploaded) }}</dd>
                  </div>
                  <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.fields.downloaded') }}</dt>
                    <dd class="mt-1 font-semibold text-slate-950 dark:text-white">{{ formatBytes(selectedCheater.downloaded) }}</dd>
                  </div>
                  <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.fields.announceTime') }}</dt>
                    <dd class="mt-1 font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(selectedCheater.announce_time) }}s</dd>
                  </div>
                  <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.fields.hitCount') }}</dt>
                    <dd class="mt-1 font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(selectedCheater.hit_count) }}</dd>
                  </div>
                  <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.fields.seeders') }}</dt>
                    <dd class="mt-1 font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(selectedCheater.seeders) }}</dd>
                  </div>
                  <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.fields.leechers') }}</dt>
                    <dd class="mt-1 font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(selectedCheater.leechers) }}</dd>
                  </div>
                </dl>

                <section class="space-y-2">
                  <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.mod.table.createdAt') }}</p>
                  <p class="text-sm font-medium text-slate-950 dark:text-white">{{ formatDateTime(selectedCheater.created_at, locale) }}</p>
                </section>

                <section v-if="selectedCheater.comment" class="space-y-2">
                  <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.recordComment') }}</p>
                  <p class="rounded-md border border-amber-200 bg-amber-50/60 px-3 py-2 text-sm leading-6 text-slate-800 dark:border-amber-900/60 dark:bg-amber-950/20 dark:text-slate-100">{{ selectedCheater.comment }}</p>
                </section>
              </div>

              <form v-if="!selectedCheater.is_dealt" class="space-y-4 p-4" @submit.prevent="resolveCheater">
                <label class="block">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.mod.form.comment') }}</span>
                  <textarea v-model="comment" rows="5" class="mt-2 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" :placeholder="$t('admin.mod.cheaters.commentPlaceholder')" :disabled="resolving" />
                </label>
                <p v-if="resolveError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ resolveError }}</p>
                <div class="flex justify-end">
                  <UButton type="submit" color="primary" icon="i-lucide-check" :loading="resolving">
                    {{ $t('admin.mod.actions.submitResolve') }}
                  </UButton>
                </div>
              </form>

              <div v-else class="space-y-3 p-4">
                <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.cheaters.resolvedDetail') }}</p>
                <dl class="grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.handledBy') }}</dt>
                    <dd class="mt-1 flex min-w-0 items-center gap-2 font-medium text-slate-950 dark:text-white">
                      <IamUserAvatar v-if="selectedCheater.dealt_by" :user="cheaterDealtUser(selectedCheater)" size="xs" />
                      <span class="truncate">{{ selectedCheater.dealt_by ? userDisplayName(cheaterDealtUser(selectedCheater), selectedCheater.dealt_by) : '-' }}</span>
                    </dd>
                  </div>
                  <div>
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.handledAt') }}</dt>
                    <dd class="mt-1 font-medium text-slate-950 dark:text-white">{{ selectedCheater.dealt_at ? formatDateTime(selectedCheater.dealt_at, locale) : '-' }}</dd>
                  </div>
                </dl>
                <div v-if="selectedCheater.dealt_comment" class="space-y-2">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.cheaters.resolvedComment') }}</p>
                  <p class="border-l-2 border-slate-200 pl-3 text-sm leading-6 text-slate-700 dark:border-slate-700 dark:text-slate-200">{{ selectedCheater.dealt_comment }}</p>
                </div>
              </div>
            </div>
          </template>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminModCheaterItem } from '~/composables/useAdmin'
import { formatBytes, formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const toast = useToast()
const adminApi = useAdmin()

const cheaters = ref<AdminModCheaterItem[]>([])
const total = ref(0)
const pending = ref(false)
const resolving = ref(false)
const errorMessage = ref('')
const resolveError = ref('')
const selectedCheater = ref<AdminModCheaterItem | null>(null)
const comment = ref('')

const query = reactive({
  page: 1,
  size: 30,
  status: -1
})

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const relativeTimeFormatter = computed(() => new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))
const statusOptions = computed(() => [
  { value: -1, label: t('admin.mod.status.all') },
  { value: 0, label: t('admin.mod.status.pending') },
  { value: 1, label: t('admin.mod.status.resolved') }
])

useHead({ title: t('admin.mod.cheaters.title') })
onMounted(loadCheaters)

async function loadCheaters() {
  pending.value = true
  errorMessage.value = ''
  const currentSelected = selectedCheater.value
  try {
    const data = await adminApi.listModCheaters(query)
    cheaters.value = data.list || []
    total.value = data.total || 0
    if (currentSelected) {
      selectedCheater.value = cheaters.value.find((item) => item.id === currentSelected.id) || currentSelected
    }
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function reloadFromFirstPage() {
  query.page = 1
  resetSelection()
  loadCheaters()
}

function changePage(page: number) {
  query.page = Math.min(Math.max(1, page), totalPages.value)
  resetSelection()
  loadCheaters()
}

function selectCheater(item: AdminModCheaterItem) {
  selectedCheater.value = item
  comment.value = item.is_dealt ? item.dealt_comment || '' : ''
  resolveError.value = ''
}

function resetSelection() {
  selectedCheater.value = null
  comment.value = ''
  resolveError.value = ''
}

async function resolveCheater() {
  if (!selectedCheater.value || selectedCheater.value.is_dealt) return
  const currentCheater = selectedCheater.value
  const nextComment = comment.value.trim()
  resolving.value = true
  resolveError.value = ''
  try {
    await adminApi.resolveModCheater(currentCheater.id, {
      status: 1,
      comment: nextComment
    })
    toast.add({ title: t('admin.mod.cheaters.resolved'), color: 'success', icon: 'i-lucide-check-circle' })
    if (query.status === 0) {
      resetSelection()
      await loadCheaters()
      return
    }
    selectedCheater.value = {
      ...currentCheater,
      is_dealt: true,
      dealt_comment: nextComment,
      dealt_at: new Date().toISOString()
    }
    await loadCheaters()
  } catch (error) {
    resolveError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    resolving.value = false
  }
}

function cheaterUser(item: AdminModCheaterItem) {
  return item.user || { id: item.user_id, username: '', avatar: '' }
}

function cheaterDealtUser(item: AdminModCheaterItem) {
  return item.dealt_user || { id: item.dealt_by, username: '', avatar: '' }
}

function userDisplayName(user: { id?: number, username?: string } | null | undefined, fallbackId?: number) {
  if (user?.username) return user.username
  const id = user?.id || fallbackId
  return id ? `#${id}` : '-'
}

function torrentDisplayName(item: AdminModCheaterItem) {
  return item.torrent?.name || t('admin.mod.cheaters.torrentId', { id: item.torrent_id })
}

function relativeTime(value?: string | null) {
  if (!value) return '-'
  const date = new Date(value)
  const timestamp = date.getTime()
  if (Number.isNaN(timestamp)) return '-'

  const diffSeconds = Math.round((timestamp - Date.now()) / 1000)
  const absSeconds = Math.abs(diffSeconds)
  if (absSeconds < 45) return relativeTimeFormatter.value.format(0, 'second')
  if (absSeconds < 45 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60), 'minute')
  if (absSeconds < 22 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60), 'hour')
  if (absSeconds < 30 * 24 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24), 'day')
  if (absSeconds < 12 * 30 * 24 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24 / 30), 'month')
  return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24 / 365), 'year')
}
</script>
