<template>
  <section class="overflow-hidden rounded-lg border border-amber-200 bg-white dark:border-amber-900/70 dark:bg-slate-900">
    <div class="flex items-center justify-between gap-3 border-b border-amber-100 bg-amber-50/70 px-4 py-3 dark:border-amber-900/60 dark:bg-amber-950/20">
      <div class="flex min-w-0 items-center gap-2">
        <UIcon name="i-lucide-clipboard-check" class="size-4 shrink-0 text-amber-600 dark:text-amber-400" />
        <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.reviews.title') }}</h2>
      </div>
      <div class="flex shrink-0 items-center gap-1.5">
        <UBadge color="warning" variant="soft" size="sm">{{ $t('admin.catalog.reviews.status.pending') }}</UBadge>
        <UTooltip :text="$t('admin.catalog.reviews.backToList')" :delay-duration="600">
          <UButton
            color="neutral"
            variant="ghost"
            size="xs"
            icon="i-lucide-arrow-left"
            :to="localePath('/admin/catalog/reviews')"
            :aria-label="$t('admin.catalog.reviews.backToList')"
          />
        </UTooltip>
      </div>
    </div>

    <form class="space-y-3 p-4" @submit.prevent="resolveReview('approve')">
      <label class="block">
        <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.reviews.fields.comment') }}</span>
        <UTextarea
          v-model="comment"
          class="mt-2 w-full"
          :rows="4"
          :placeholder="$t('admin.catalog.reviews.commentPlaceholder')"
          :disabled="Boolean(reviewing)"
        />
      </label>

      <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-300">{{ errorMessage }}</p>

      <div class="flex items-center justify-between gap-3">
        <UButton type="button" color="error" variant="soft" icon="i-lucide-x" :loading="reviewing === 'reject'" :disabled="Boolean(reviewing)" @click="resolveReview('reject')">
          {{ $t('admin.catalog.reviews.reject') }}
        </UButton>
        <UButton type="submit" color="primary" icon="i-lucide-check" :loading="reviewing === 'approve'" :disabled="Boolean(reviewing)">
          {{ $t('admin.catalog.reviews.approve') }}
        </UButton>
      </div>
    </form>
  </section>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'

const props = defineProps<{
  torrentId: number
}>()

const emit = defineEmits<{
  reviewed: [action: 'approve' | 'reject']
}>()

const { t } = useI18n()
const localePath = useLocalePath()
const toast = useToast()
const adminApi = useAdmin()
const comment = ref('')
const reviewing = ref<'' | 'approve' | 'reject'>('')
const errorMessage = ref('')

async function resolveReview(action: 'approve' | 'reject') {
  if (reviewing.value) return
  const normalizedComment = comment.value.trim()
  if (action === 'reject' && !normalizedComment) {
    errorMessage.value = t('admin.catalog.reviews.commentRequired')
    return
  }

  reviewing.value = action
  errorMessage.value = ''
  try {
    if (action === 'approve') await adminApi.approveCatalogTorrent(props.torrentId, normalizedComment)
    else await adminApi.rejectCatalogTorrent(props.torrentId, normalizedComment)
    toast.add({ title: t(`admin.catalog.reviews.${action}Success`), color: 'success', icon: 'i-lucide-check-circle' })
    emit('reviewed', action)
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    reviewing.value = ''
  }
}
</script>
