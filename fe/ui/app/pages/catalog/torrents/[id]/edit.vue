<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div v-if="pending" class="h-96 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />

      <div v-else-if="errorMessage" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadPage">
          {{ $t('common.retry') }}
        </UButton>
      </div>

      <form v-else class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_340px]" @submit.prevent="handleSubmit">
        <div class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.edit.sections.info') }}</h2>
            </template>

            <div class="grid grid-cols-1 gap-4">
              <UFormField :label="$t('catalog.torrents.upload.fields.category')" required>
                <select
                  v-model="form.categoryId"
                  class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="savePending || categories.length === 0"
                >
                  <option value="0">{{ $t('catalog.torrents.upload.fields.categoryPlaceholder') }}</option>
                  <option v-for="category in categories" :key="category.id" :value="String(category.id)">
                    {{ categoryDisplayName(category) }}
                  </option>
                </select>
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.name')">
                <UInput v-model="form.name" class="w-full" :disabled="savePending" />
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.subTitle')">
                <UInput v-model="form.subTitle" class="w-full" :disabled="savePending" />
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.description')">
                <UTextarea v-model="form.description" class="w-full" :rows="10" :disabled="savePending" />
              </UFormField>
            </div>
          </UCard>
        </div>

        <aside class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.edit.sections.publish') }}</h2>
            </template>

            <div class="space-y-4">
              <label class="flex items-center justify-between gap-4 rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('catalog.torrents.upload.fields.anonymous') }}</span>
                <input
                  v-model="form.anonymous"
                  type="checkbox"
                  class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                  :disabled="savePending"
                >
              </label>

              <dl class="space-y-3 text-sm">
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.edit.summary.torrentId') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">#{{ torrentId }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.edit.summary.category') }}</dt>
                  <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ selectedCategoryName }}</dd>
                </div>
              </dl>

              <UButton type="submit" color="primary" icon="i-lucide-save" block :loading="savePending" :disabled="!canSubmit">
                {{ $t('catalog.torrents.edit.submit') }}
              </UButton>
            </div>
          </UCard>
        </aside>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { CatalogCategory, TorrentDetail } from '~/composables/useCatalogTorrents'
import { localizeI18nName } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()
const { user, isStaff, fetchUser } = useAuth()

const torrent = ref<TorrentDetail | null>(null)
const categories = ref<CatalogCategory[]>([])
const pending = ref(true)
const savePending = ref(false)
const errorMessage = ref('')

const form = reactive({
  categoryId: '0',
  name: '',
  subTitle: '',
  description: '',
  anonymous: false
})

const torrentId = computed(() => readRouteId())
const canEditTorrent = computed(() => Boolean(torrent.value && (isStaff.value || user.value?.id === torrent.value.ownerId)))
const canSubmit = computed(() => Number(form.categoryId) > 0 && Boolean(torrent.value) && canEditTorrent.value && !savePending.value)
const selectedCategoryName = computed(() => {
  const category = categories.value.find((item) => item.id === Number(form.categoryId))
  return category ? categoryDisplayName(category) : '-'
})

useHead(() => ({
  title: `${t('catalog.torrents.edit.metaTitle')} - NextPT`
}))

onMounted(loadPage)

function readRouteId() {
  const raw = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id
  const id = Number(raw)
  return Number.isInteger(id) && id > 0 ? id : 0
}

async function loadPage() {
  if (torrentId.value <= 0) {
    errorMessage.value = t('catalog.torrents.detail.invalidId')
    pending.value = false
    return
  }

  pending.value = true
  errorMessage.value = ''
  try {
    if (!user.value) {
      await fetchUser()
    }

    const [detail, categoryList] = await Promise.all([
      catalogTorrents.getTorrent(torrentId.value),
      catalogTorrents.listCategories().catch(() => ({ list: [] }))
    ])

    torrent.value = detail
    categories.value = categoryList.list || []
    form.categoryId = String(detail.categoryId || 0)
    form.name = detail.name || ''
    form.subTitle = detail.subTitle || ''
    form.description = detail.description || ''
    form.anonymous = Boolean(detail.anonymous)

    if (!canEditTorrent.value) {
      errorMessage.value = t('catalog.torrents.edit.forbidden')
    }
  } catch (error) {
    torrent.value = null
    categories.value = []
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

async function handleSubmit() {
  if (!canSubmit.value) return

  savePending.value = true
  try {
    await catalogTorrents.updateTorrent(torrentId.value, {
      categoryId: Number(form.categoryId),
      name: form.name.trim(),
      subTitle: form.subTitle.trim(),
      description: form.description.trim(),
      anonymous: form.anonymous
    })
    toast.add({
      title: t('catalog.torrents.edit.saved'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await navigateTo(localePath(`/catalog/torrents/${torrentId.value}`))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    savePending.value = false
  }
}

function categoryDisplayName(category: CatalogCategory) {
  return localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`)
}
</script>
