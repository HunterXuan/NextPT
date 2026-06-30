<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <form class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_340px]" @submit.prevent="handleSubmit">
        <div class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.upload.sections.file') }}</h2>
            </template>

            <label
              class="flex min-h-40 cursor-pointer flex-col items-center justify-center rounded-lg border border-dashed px-4 py-8 text-center transition"
              :class="selectedFile
                ? 'border-sky-300 bg-sky-50 dark:border-sky-700 dark:bg-sky-950'
                : 'border-slate-300 bg-slate-50 hover:border-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:hover:border-slate-600'"
              @dragover.prevent
              @drop.prevent="handleDrop"
            >
              <input
                class="sr-only"
                type="file"
                accept=".torrent,application/x-bittorrent"
                :disabled="pending"
                @change="handleFileChange"
              >
              <UIcon name="i-lucide-file-up" class="size-9 text-slate-400" />
              <span class="mt-3 text-sm font-medium text-slate-950 dark:text-white">
                {{ selectedFile?.name || $t('catalog.torrents.upload.file.choose') }}
              </span>
              <span v-if="selectedFile" class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                {{ formatBytes(selectedFile.size) }}
              </span>
            </label>
          </UCard>

          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.upload.sections.info') }}</h2>
            </template>

            <div class="grid grid-cols-1 gap-4">
              <UFormField :label="$t('catalog.torrents.upload.fields.category')" required>
                <select
                  v-model="form.categoryId"
                  class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="pending || categories.length === 0"
                >
                  <option value="0">{{ $t('catalog.torrents.upload.fields.categoryPlaceholder') }}</option>
                  <option v-for="category in categories" :key="category.id" :value="String(category.id)">
                    {{ categoryDisplayName(category) }}
                  </option>
                </select>
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.name')">
                <UInput v-model="form.name" class="w-full" :disabled="pending" />
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.subTitle')">
                <UInput v-model="form.subTitle" class="w-full" :disabled="pending" />
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.description')">
                <UTextarea v-model="form.description" class="w-full" :rows="8" :disabled="pending" />
              </UFormField>
            </div>
          </UCard>
        </div>

        <aside class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.upload.sections.publish') }}</h2>
            </template>

            <div class="space-y-4">
              <label class="flex items-center justify-between gap-4 rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('catalog.torrents.upload.fields.anonymous') }}</span>
                <input
                  v-model="form.anonymous"
                  type="checkbox"
                  class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                  :disabled="pending"
                >
              </label>

              <dl class="space-y-3 text-sm">
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.upload.summary.file') }}</dt>
                  <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ selectedFile?.name || '-' }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.upload.summary.category') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">{{ selectedCategoryName }}</dd>
                </div>
              </dl>

              <UButton type="submit" color="primary" icon="i-lucide-upload" block :loading="pending" :disabled="!canSubmit">
                {{ $t('catalog.torrents.upload.submit') }}
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
import type { CatalogCategory } from '~/composables/useCatalogTorrents'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()

const categories = ref<CatalogCategory[]>([])
const selectedFile = ref<File | null>(null)
const pending = ref(false)

const form = reactive({
  categoryId: '0',
  name: '',
  subTitle: '',
  description: '',
  anonymous: false
})

const selectedCategoryName = computed(() => {
  const categoryId = Number(form.categoryId)
  const category = categories.value.find((item) => item.id === categoryId)
  return category ? categoryDisplayName(category) : '-'
})

const canSubmit = computed(() => Boolean(selectedFile.value && Number(form.categoryId) > 0 && !pending.value))

onMounted(loadCategories)

async function loadCategories() {
  try {
    const data = await catalogTorrents.listCategories()
    categories.value = data.list || []
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  }
}

function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  setFile(input.files?.[0] || null)
}

function handleDrop(event: DragEvent) {
  setFile(event.dataTransfer?.files?.[0] || null)
}

function setFile(file: File | null) {
  if (!file) {
    selectedFile.value = null
    return
  }

  if (!file.name.toLowerCase().endsWith('.torrent')) {
    toast.add({
      title: t('catalog.torrents.upload.file.invalid'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
    return
  }

  selectedFile.value = file
  if (!form.name.trim()) {
    form.name = file.name.replace(/\.torrent$/i, '')
  }
}

async function handleSubmit() {
  if (!selectedFile.value || Number(form.categoryId) <= 0) return

  pending.value = true
  try {
    const out = await catalogTorrents.uploadTorrent({
      file: selectedFile.value,
      categoryId: Number(form.categoryId),
      name: form.name,
      subTitle: form.subTitle,
      description: form.description,
      anonymous: form.anonymous
    })

    toast.add({
      title: t('catalog.torrents.upload.success', { id: out.torrentId }),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await navigateTo(localePath(`/catalog/torrents/${out.torrentId}`))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    pending.value = false
  }
}

function categoryDisplayName(category: CatalogCategory) {
  return localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`)
}

useSeoMeta({
  title: t('catalog.torrents.upload.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
