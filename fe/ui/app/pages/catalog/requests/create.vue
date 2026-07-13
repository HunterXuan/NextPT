<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <form class="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1fr)_320px] xl:items-start xl:gap-5" @submit.prevent="handleSubmit">
        <main class="min-w-0">
          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h1 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.create.formTitle') }}</h1>
            </div>

            <div class="space-y-5 p-4">
              <div class="grid gap-4" :class="form.requestType === CatalogRequestType.Torrent ? '2xl:grid-cols-[190px_230px_minmax(0,1fr)]' : '2xl:grid-cols-[190px_minmax(0,1fr)]'">
              <UFormField :label="$t('catalog.requests.fields.type')">
                <div class="grid h-10 grid-cols-2 rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950">
                  <button type="button" :class="typeButtonClass(form.requestType === CatalogRequestType.Torrent)" @click="form.requestType = CatalogRequestType.Torrent">
                    <UIcon name="i-lucide-package-search" class="size-4" />
                    {{ $t('catalog.requests.types.torrent') }}
                  </button>
                  <button type="button" :class="typeButtonClass(form.requestType === CatalogRequestType.Reseed)" @click="form.requestType = CatalogRequestType.Reseed">
                    <UIcon name="i-lucide-refresh-cw" class="size-4" />
                    {{ $t('catalog.requests.types.reseed') }}
                  </button>
                </div>
              </UFormField>

              <template v-if="form.requestType === CatalogRequestType.Torrent">
                <UFormField :label="$t('catalog.requests.fields.category')" required>
                  <select v-model.number="form.categoryId" class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="pending || categoriesPending">
                    <option :value="0">{{ $t('catalog.requests.create.categoryPlaceholder') }}</option>
                    <option v-for="category in categories" :key="category.id" :value="category.id">{{ categoryName(category) }}</option>
                  </select>
                </UFormField>

                <UFormField :label="$t('catalog.requests.fields.title')" required>
                  <UInput v-model="form.title" class="w-full" size="lg" :ui="{ base: 'h-10' }" :disabled="pending" :placeholder="$t('catalog.requests.create.titlePlaceholder')" />
                </UFormField>
              </template>

              <UFormField v-else :label="$t('catalog.requests.fields.targetTorrent')" required>
                <div v-if="selectedTarget" class="flex h-10 items-center gap-2 rounded-md border border-slate-200 bg-slate-50 px-3 dark:border-slate-800 dark:bg-slate-950">
                  <UIcon name="i-lucide-library" class="size-4 shrink-0 text-sky-500" />
                  <p class="min-w-0 flex-1 truncate text-sm font-medium text-slate-950 dark:text-white">{{ selectedTarget.name }}</p>
                  <span class="hidden shrink-0 text-xs text-slate-500 sm:inline dark:text-slate-400">#{{ selectedTarget.id }} / {{ $t('catalog.requests.create.seeders', { count: selectedTarget.seeders }) }}</span>
                  <UTooltip :text="$t('catalog.requests.create.changeTarget')" :delay-duration="300">
                    <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-pencil" :aria-label="$t('catalog.requests.create.changeTarget')" @click="clearTarget" />
                  </UTooltip>
                </div>

                <div v-else class="relative">
                  <UInput v-model="torrentKeyword" class="w-full" size="lg" :ui="{ base: 'h-10' }" icon="i-lucide-search" :loading="torrentSearchPending" :disabled="pending" :placeholder="$t('catalog.requests.create.targetPlaceholder')" />
                  <div v-if="torrentKeyword.trim() && !torrentSearchPending" class="mt-2 overflow-hidden rounded-md border border-slate-200 dark:border-slate-800">
                    <button
                      v-for="torrent in torrentResults"
                      :key="torrent.id"
                      type="button"
                      class="flex w-full items-start gap-3 border-t border-slate-100 px-3 py-2.5 text-left first:border-t-0 dark:border-slate-800"
                      :class="torrent.seeders > 0 ? 'cursor-not-allowed opacity-50' : 'hover:bg-slate-50 dark:hover:bg-slate-950'"
                      :disabled="torrent.seeders > 0"
                      @click="selectTarget(torrent)"
                    >
                      <span class="min-w-0 flex-1">
                        <span class="block truncate text-sm font-medium text-slate-950 dark:text-white">{{ torrent.name }}</span>
                        <span class="mt-0.5 block text-xs text-slate-500 dark:text-slate-400">#{{ torrent.id }} / {{ $t('catalog.requests.create.seeders', { count: torrent.seeders }) }}</span>
                      </span>
                      <UIcon v-if="torrent.seeders === 0" name="i-lucide-chevron-right" class="mt-0.5 size-4 shrink-0 text-slate-400" />
                    </button>
                    <p v-if="torrentResults.length === 0" class="px-3 py-4 text-center text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.create.noTorrentResults') }}</p>
                  </div>
                </div>
              </UFormField>
              </div>

              <div class="border-t border-slate-100 pt-5 dark:border-slate-800">
                <UFormField :label="$t('catalog.requests.fields.description')" required>
                  <RichTextComposer
                    v-model="form.description"
                    v-model:mode="editorMode"
                    :rows="7"
                    :disabled="pending"
                    :as-form="false"
                    :show-actions="false"
                    :submit-label="$t('catalog.requests.create.submit')"
                    :write-label="$t('common.editor.edit')"
                    :preview-label="$t('common.editor.preview')"
                    :preview-empty="$t('catalog.requests.create.previewEmpty')"
                    :placeholder="$t('catalog.requests.create.descriptionPlaceholder')"
                  />
                </UFormField>
              </div>

              <div class="border-t border-slate-100 pt-5 dark:border-slate-800">
                <div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
                  <div>
                    <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.fields.reward') }}</h2>
                    <p class="mt-1 text-xs leading-5 text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.create.rewardHelp') }}</p>
                  </div>
                  <div class="grid grid-cols-4 gap-2 sm:flex sm:shrink-0">
                    <button v-for="amount in rewardOptions" :key="amount" type="button" :class="rewardButtonClass(form.rewardAmount === amount)" @click="form.rewardAmount = amount">
                      {{ numberFormatter.format(amount) }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </main>

        <aside class="app-sticky-offset space-y-3 xl:sticky">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.create.publishTitle') }}</h2>
            <dl class="mt-3 divide-y divide-slate-100 text-sm dark:divide-slate-800">
              <div class="flex items-center justify-between gap-3 py-2">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.create.summary.currentBalance') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(user?.stat.bonus || 0) }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3 py-2">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.create.summary.balance') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(balanceAfterReward) }}</dd>
              </div>
            </dl>
            <div class="mt-3 border-t border-slate-100 pt-3 dark:border-slate-800">
              <AppPermissionButton :permission="Permission.CatalogRequestCreate" type="submit" color="primary" icon="i-lucide-send" block :loading="pending" :disabled="!canSubmit" :tooltip="$t('catalog.requests.create.submit')">
                {{ $t('catalog.requests.create.submit') }}
              </AppPermissionButton>
            </div>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.create.guide.title') }}</h2>
            <ul class="mt-3 space-y-2 text-sm leading-6 text-slate-600 dark:text-slate-300">
              <li class="flex gap-2"><UIcon name="i-lucide-search" class="mt-1 size-4 shrink-0" /><span>{{ $t('catalog.requests.create.guide.search') }}</span></li>
              <li class="flex gap-2"><UIcon name="i-lucide-circle-check" class="mt-1 size-4 shrink-0" /><span>{{ $t('catalog.requests.create.guide.confirm') }}</span></li>
              <li class="flex gap-2"><UIcon name="i-lucide-coins" class="mt-1 size-4 shrink-0" /><span>{{ $t('catalog.requests.create.guide.reward') }}</span></li>
            </ul>
          </section>

          <section v-if="errorMessage" class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">{{ errorMessage }}</section>
        </aside>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { CatalogRequestType } from '~/composables/useCatalogRequests'
import type { CatalogCategory, TorrentDetail, TorrentListItem } from '~/composables/useCatalogTorrents'
import { localizeI18nName } from '~/utils/format'

type EditorMode = 'write' | 'preview'

definePageMeta({ middleware: 'auth' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const requestApi = useCatalogRequests()
const catalogApi = useCatalogTorrents()
const { user, fetchUser, hasPermission } = useAuth()

const categories = ref<CatalogCategory[]>([])
const categoriesPending = ref(true)
const pending = ref(false)
const errorMessage = ref('')
const editorMode = ref<EditorMode>('write')
const selectedTarget = ref<TorrentDetail | TorrentListItem | null>(null)
const torrentKeyword = ref('')
const torrentResults = ref<TorrentListItem[]>([])
const torrentSearchPending = ref(false)
const rewardOptions = [10, 50, 100, 500]
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value, { maximumFractionDigits: 1 }))
let searchTimer: ReturnType<typeof setTimeout> | null = null

const initialTorrentId = readPositiveQuery('torrentId')
const form = reactive({
  requestType: readRequestTypeQuery(),
  categoryId: 0,
  title: '',
  description: '',
  rewardAmount: 50
})

const balanceAfterReward = computed(() => (user.value?.stat.bonus || 0) - form.rewardAmount)
const canSubmit = computed(() => {
  if (!hasPermission(Permission.CatalogRequestCreate) || pending.value) return false
  if (form.description.trim().length < 3 || form.rewardAmount <= 0 || (user.value?.stat.bonus || 0) < form.rewardAmount) return false
  if (form.requestType === CatalogRequestType.Reseed) return Boolean(selectedTarget.value && selectedTarget.value.seeders === 0)
  return form.categoryId > 0 && form.title.trim().length > 0
})

useHead(() => ({ title: t('catalog.requests.create.metaTitle') }))

watch(torrentKeyword, () => {
  if (searchTimer) clearTimeout(searchTimer)
  const keyword = torrentKeyword.value.trim()
  if (!keyword) {
    torrentResults.value = []
    return
  }
  searchTimer = setTimeout(searchTorrents, 350)
})

onMounted(async () => {
  await loadCategories()
  if (initialTorrentId > 0) await loadInitialTarget(initialTorrentId)
})

onBeforeUnmount(() => {
  if (searchTimer) clearTimeout(searchTimer)
})

async function loadCategories() {
  categoriesPending.value = true
  try {
    categories.value = (await catalogApi.listCategories()).list || []
    form.categoryId = categories.value[0]?.id || 0
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    categoriesPending.value = false
  }
}

async function loadInitialTarget(id: number) {
  try {
    const torrent = await catalogApi.getTorrent(id)
    selectedTarget.value = torrent
    form.requestType = CatalogRequestType.Reseed
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  }
}

async function searchTorrents() {
  const keyword = torrentKeyword.value.trim()
  if (!keyword) return
  torrentSearchPending.value = true
  try {
    torrentResults.value = (await catalogApi.listTorrents({ keyword, page: 1, size: 8 })).list || []
  } catch {
    torrentResults.value = []
  } finally {
    torrentSearchPending.value = false
  }
}

function selectTarget(torrent: TorrentListItem) {
  if (torrent.seeders > 0) return
  selectedTarget.value = torrent
  torrentKeyword.value = ''
  torrentResults.value = []
}

function clearTarget() {
  selectedTarget.value = null
  torrentKeyword.value = ''
}

async function handleSubmit() {
  if (!canSubmit.value) return
  pending.value = true
  errorMessage.value = ''
  try {
    const output = await requestApi.createRequest({
      requestType: form.requestType,
      categoryId: form.requestType === CatalogRequestType.Torrent ? form.categoryId : undefined,
      targetTorrentId: form.requestType === CatalogRequestType.Reseed ? selectedTarget.value?.id : undefined,
      title: form.requestType === CatalogRequestType.Torrent ? form.title.trim() : undefined,
      description: form.description.trim(),
      rewardAmount: form.rewardAmount
    })
    await fetchUser()
    toast.add({ title: t('catalog.requests.create.success'), color: 'success', icon: 'i-lucide-circle-check' })
    await navigateTo(localePath(`/catalog/requests/${output.id}`))
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function readPositiveQuery(key: string) {
  const raw = Array.isArray(route.query[key]) ? route.query[key]?.[0] : route.query[key]
  const value = Number(raw)
  return Number.isInteger(value) && value > 0 ? value : 0
}

function readRequestTypeQuery() {
  return route.query.type === 'reseed' || initialTorrentId > 0 ? CatalogRequestType.Reseed : CatalogRequestType.Torrent
}

function categoryName(category: CatalogCategory) {
  return localizeI18nName(category.name, locale.value)
}

function typeButtonClass(active: boolean) {
  return ['inline-flex h-9 items-center gap-2 rounded px-3 text-sm font-medium transition-colors', active ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white' : 'text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white']
}

function rewardButtonClass(active: boolean) {
  return ['inline-flex h-9 min-w-16 items-center justify-center rounded-md border px-3 text-sm font-semibold transition-colors', active ? 'border-amber-400 bg-amber-50 text-amber-700 dark:border-amber-600 dark:bg-amber-950/40 dark:text-amber-300' : 'border-slate-200 bg-white text-slate-600 hover:bg-slate-50 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-300 dark:hover:bg-slate-950']
}
</script>
