<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto]">
        <form class="flex gap-2" @submit.prevent="reloadFromFirstPage">
          <select v-model.number="query.status" class="h-10 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-900" @change="reloadFromFirstPage">
            <option v-for="option in statusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
          <UButton type="submit" color="primary" icon="i-lucide-search" :loading="pending">
            {{ $t('admin.mod.actions.filter') }}
          </UButton>
        </form>

        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.stats.total') }}</p>
          <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(total) }}</p>
        </div>
      </div>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.cheaters.list') }}</h2>
            <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadCheaters">
              {{ $t('common.refresh') }}
            </UButton>
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
            <table class="min-w-[980px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[9%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">ID</th>
                  <th class="w-[17%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.cheaters.table.user') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.cheaters.table.traffic') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.cheaters.table.swarm') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.status') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.createdAt') }}</th>
                  <th class="w-[8%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.mod.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in cheaters"
                  :key="item.id"
                  class="border-b border-slate-200 last:border-b-0 dark:border-slate-800"
                  :class="selectedCheater?.id === item.id ? 'bg-amber-50/70 dark:bg-amber-950/20' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-4 py-3 text-sm font-medium text-slate-950 dark:text-white">#{{ item.id }}</td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <p>{{ $t('admin.mod.cheaters.userId', { id: item.user_id }) }}</p>
                    <p class="mt-1 text-xs text-slate-500">{{ $t('admin.mod.cheaters.torrentId', { id: item.torrent_id }) }}</p>
                  </td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <p>{{ $t('admin.mod.cheaters.uploaded', { value: formatBytes(item.uploaded) }) }}</p>
                    <p class="mt-1">{{ $t('admin.mod.cheaters.downloaded', { value: formatBytes(item.downloaded) }) }}</p>
                  </td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <p>{{ $t('admin.mod.cheaters.announceTime', { value: numberFormatter.format(item.announce_time) }) }}</p>
                    <p class="mt-1 text-xs text-slate-500">{{ $t('admin.mod.cheaters.swarm', { seeders: item.seeders, leechers: item.leechers, hits: item.hit_count }) }}</p>
                  </td>
                  <td class="px-4 py-3">
                    <UBadge :color="item.is_dealt ? 'success' : 'warning'" variant="soft">
                      {{ item.is_dealt ? $t('admin.mod.status.resolved') : $t('admin.mod.status.pending') }}
                    </UBadge>
                  </td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">{{ formatDateTime(item.created_at, locale) }}</td>
                  <td class="px-4 py-3 text-right">
                    <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-check-check" :disabled="item.is_dealt" @click="selectCheater(item)">
                      {{ $t('admin.mod.actions.resolve') }}
                    </UButton>
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

        <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
              {{ selectedCheater ? $t('admin.mod.cheaters.resolveTitle', { id: selectedCheater.id }) : $t('admin.mod.cheaters.select') }}
            </h2>
          </div>
          <form class="space-y-4 p-4" @submit.prevent="resolveCheater">
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.mod.form.comment') }}</span>
              <textarea v-model="comment" rows="6" class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="!selectedCheater || resolving" />
            </label>
            <p v-if="resolveError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ resolveError }}</p>
            <UButton type="submit" color="primary" icon="i-lucide-check" :loading="resolving" :disabled="!selectedCheater">
              {{ $t('admin.mod.actions.submitResolve') }}
            </UButton>
          </form>
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
  try {
    const data = await adminApi.listModCheaters(query)
    cheaters.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function reloadFromFirstPage() {
  query.page = 1
  loadCheaters()
}

function changePage(page: number) {
  query.page = Math.min(Math.max(1, page), totalPages.value)
  loadCheaters()
}

function selectCheater(item: AdminModCheaterItem) {
  selectedCheater.value = item
  comment.value = item.comment || ''
  resolveError.value = ''
}

async function resolveCheater() {
  if (!selectedCheater.value) return
  resolving.value = true
  resolveError.value = ''
  try {
    await adminApi.resolveModCheater(selectedCheater.value.id, {
      status: 1,
      comment: comment.value.trim()
    })
    toast.add({ title: t('admin.mod.cheaters.resolved'), color: 'success', icon: 'i-lucide-check-circle' })
    selectedCheater.value = null
    comment.value = ''
    await loadCheaters()
  } catch (error) {
    resolveError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    resolving.value = false
  }
}
</script>
