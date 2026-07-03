<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="border-b border-slate-200 px-5 py-4 dark:border-slate-800">
      <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div class="min-w-0">
          <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.peers.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ peerSummary }}</p>
        </div>

        <div class="inline-flex w-fit rounded-md bg-slate-100 p-1 dark:bg-slate-800/80">
          <button
            v-for="filter in peerFilters"
            :key="filter.value"
            type="button"
            class="inline-flex h-8 items-center rounded px-3 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-60"
            :class="peerStatus === filter.value
              ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-950 dark:text-white'
              : 'text-slate-600 hover:text-slate-950 dark:text-slate-300 dark:hover:text-white'"
            :disabled="peersPending"
            @click="setPeerStatus(filter.value)"
          >
            {{ filter.label }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="peersError" class="m-5 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
      <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
      <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ peersError }}</p>
      <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadPeers">
        {{ $t('common.retry') }}
      </UButton>
    </div>
    <div v-else-if="peersPending && peers.length === 0" class="space-y-3 p-5">
      <div v-for="item in 3" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>
    <div v-else-if="peers.length === 0" class="m-5 rounded-md border border-dashed border-slate-200 px-4 py-12 text-center dark:border-slate-800">
      <UIcon name="i-lucide-radio-tower" class="mx-auto size-8 text-slate-400" />
      <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.peers.empty') }}</p>
    </div>
    <div v-else>
      <div class="hidden grid-cols-[minmax(0,1fr)_132px_110px_110px_160px_140px] gap-4 border-b border-slate-200 bg-slate-50 px-5 py-2.5 text-xs font-medium text-slate-500 xl:grid dark:border-slate-800 dark:bg-slate-950/60 dark:text-slate-400">
        <span>{{ $t('user.snatches.torrent') }}</span>
        <span class="text-center">{{ $t('user.peers.status') }}</span>
        <span class="text-right">{{ $t('user.snatches.uploaded') }}</span>
        <span class="text-right">{{ $t('user.snatches.downloaded') }}</span>
        <span class="text-right">{{ $t('user.peers.client') }}</span>
        <span class="text-right">{{ $t('user.snatches.lastAction') }}</span>
      </div>

      <div class="divide-y divide-slate-100 dark:divide-slate-800">
        <article
          v-for="peer in peers"
          :key="`${peer.torrentId}:${peer.agent}:${peer.lastActionAt || ''}`"
          class="grid gap-3 px-5 py-4 transition hover:bg-slate-50/70 xl:grid-cols-[minmax(0,1fr)_132px_110px_110px_160px_140px] xl:items-center xl:gap-4 dark:hover:bg-slate-950/50"
        >
          <div class="min-w-0">
            <NuxtLink
              :to="localePath(`/catalog/torrents/${peer.torrentId}`)"
              class="block truncate text-sm font-semibold text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300"
            >
              {{ peer.torrentName || $t('catalog.torrents.detail.titleFallback', { id: peer.torrentId }) }}
            </NuxtLink>
            <div class="mt-1.5 text-xs text-slate-500 dark:text-slate-400">
              <span class="tabular-nums">{{ formatBytes(peer.torrentSize) }}</span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-x-4 gap-y-3 text-xs sm:grid-cols-3 xl:contents">
            <div class="xl:text-center">
              <p class="text-slate-500 xl:hidden dark:text-slate-400">{{ $t('user.peers.status') }}</p>
              <div class="mt-1 flex items-center gap-2 xl:mt-0 xl:justify-center">
                <UBadge :color="peer.isSeeder ? 'success' : 'primary'" variant="soft" size="sm">
                  {{ peer.isSeeder ? $t('user.peers.seeding') : $t('user.peers.leeching') }}
                </UBadge>
                <span v-if="!peer.isSeeder" class="font-medium tabular-nums text-slate-950 dark:text-white">
                  {{ formatLeechProgress(peer) }}
                </span>
              </div>
            </div>
            <div class="xl:text-right">
              <p class="text-slate-500 xl:hidden dark:text-slate-400">{{ $t('user.snatches.uploaded') }}</p>
              <p class="mt-1 font-medium tabular-nums text-slate-950 xl:mt-0 dark:text-white">{{ formatBytes(peer.uploaded) }}</p>
            </div>
            <div class="xl:text-right">
              <p class="text-slate-500 xl:hidden dark:text-slate-400">{{ $t('user.snatches.downloaded') }}</p>
              <p class="mt-1 font-medium tabular-nums text-slate-950 xl:mt-0 dark:text-white">{{ formatBytes(peer.downloaded) }}</p>
            </div>
            <div class="min-w-0 xl:text-right">
              <p class="text-slate-500 xl:hidden dark:text-slate-400">{{ $t('user.peers.client') }}</p>
              <p class="mt-1 truncate font-medium text-slate-950 xl:mt-0 dark:text-white">{{ peer.agent || '-' }}</p>
            </div>
            <div class="xl:text-right">
              <p class="text-slate-500 xl:hidden dark:text-slate-400">{{ $t('user.snatches.lastAction') }}</p>
              <p class="mt-1 font-medium tabular-nums text-slate-950 xl:mt-0 dark:text-white">{{ formatDateTime(peer.lastActionAt, locale) }}</p>
            </div>
          </div>
        </article>
      </div>
    </div>

    <AppPager
      v-if="peerTotal > 0"
      class="border-t border-slate-200 px-5 py-4 dark:border-slate-800"
      size="sm"
      :page="peerPage"
      :total="peerTotal"
      :page-size="peerSize"
      :disabled="peersPending"
      @page-change="goToPeerPage"
    />
  </section>
</template>

<script setup lang="ts">
import { useAccounting, type PeerItem, type PeerStatus } from '~/composables/useAccounting'
import { ApiError } from '~/composables/useApi'
import { formatBytes, formatDateTime } from '~/utils/format'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const accounting = useAccounting()

const peers = ref<PeerItem[]>([])
const peerStatus = ref<PeerStatus>('all')
const peerPage = ref(1)
const peerSize = 10
const peerTotal = ref(0)
const seedingTotal = ref(0)
const leechingTotal = ref(0)
const peersPending = ref(true)
const peersError = ref('')

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const peerTotalPages = computed(() => Math.max(1, Math.ceil(peerTotal.value / peerSize)))

const peerSummary = computed(() => t('user.peers.summary', {
  seeding: numberFormatter.value.format(seedingTotal.value),
  leeching: numberFormatter.value.format(leechingTotal.value)
}))

const peerFilters = computed<Array<{ value: PeerStatus, label: string }>>(() => [
  {
    value: 'all',
    label: t('user.peers.filters.all')
  },
  {
    value: 'seeding',
    label: t('user.peers.filters.seeding')
  },
  {
    value: 'leeching',
    label: t('user.peers.filters.leeching')
  }
])

onMounted(loadPeers)

async function loadPeers() {
  if (peersPending.value && peers.value.length > 0) return

  peersPending.value = true
  peersError.value = ''

  try {
    let data = await accounting.listPeers({
      page: peerPage.value,
      size: peerSize,
      status: peerStatus.value
    })

    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / peerSize))
    if (peerPage.value > nextTotalPages) {
      peerPage.value = nextTotalPages
      data = await accounting.listPeers({
        page: peerPage.value,
        size: peerSize,
        status: peerStatus.value
      })
    }

    peers.value = data.list || []
    peerTotal.value = data.total || 0
    seedingTotal.value = data.seedingTotal || 0
    leechingTotal.value = data.leechingTotal || 0
  } catch (error) {
    peers.value = []
    peerTotal.value = 0
    seedingTotal.value = 0
    leechingTotal.value = 0
    peersError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    peersPending.value = false
  }
}

function setPeerStatus(status: PeerStatus) {
  if (peerStatus.value === status) return

  peerStatus.value = status
  peerPage.value = 1
  loadPeers()
}

function goToPeerPage(page: number) {
  peerPage.value = Math.min(Math.max(1, page), peerTotalPages.value)
  loadPeers()
}

function getPeerProgress(peer: PeerItem) {
  if (peer.isSeeder) return 100

  const size = Number(peer.torrentSize || 0)
  if (size <= 0) return 0

  const completed = Math.max(0, size - Number(peer.remaining || 0))
  return Math.min(100, Math.max(0, (completed / size) * 100))
}

function formatLeechProgress(peer: PeerItem) {
  const progress = getPeerProgress(peer)
  return `${numberFormatter.value.format(Number(progress.toFixed(progress >= 10 ? 0 : 1)))}%`
}
</script>
