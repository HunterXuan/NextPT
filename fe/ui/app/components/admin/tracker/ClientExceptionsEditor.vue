<template>
  <div class="border-t border-neutral-200 pt-4 dark:border-neutral-800">
    <div class="flex items-center justify-between gap-3">
      <h3 class="text-sm font-medium">{{ t('trackerAdmin.exceptions') }}</h3>
      <UTooltip :text="t('trackerAdmin.addException')"><UButton type="button" icon="i-lucide-plus" color="neutral" variant="soft" :aria-label="t('trackerAdmin.addException')" :disabled="disabled || exceptions.length >= 50" @click="add" /></UTooltip>
    </div>
    <div class="divide-y divide-neutral-200 dark:divide-neutral-800">
      <details v-for="(exception, index) in exceptions" :key="index" open class="mt-3 pt-3">
        <summary class="flex cursor-pointer list-none items-center gap-2 text-sm"><UIcon name="i-lucide-ban" class="size-4 shrink-0 text-muted" /><span class="min-w-0 flex-1 truncate">{{ exception.name || t('trackerAdmin.exceptionName') }}</span><UBadge :color="exception.enabled ? 'error' : 'neutral'" variant="soft">{{ t(exception.enabled ? 'trackerAdmin.enabled' : 'trackerAdmin.disabled') }}</UBadge><UIcon name="i-lucide-chevron-down" class="size-4 shrink-0 text-muted" /></summary>
        <div class="mt-3 space-y-3">
          <UFormField :label="t('trackerAdmin.exceptionName')" required><UInput :model-value="exception.name" class="w-full" :maxlength="60" :disabled="disabled" @update:model-value="value => update(index, 'name', String(value))" /></UFormField>
          <UFormField :label="t('trackerAdmin.peerPattern')"><UInput :model-value="exception.peerIdPattern" class="w-full" :maxlength="200" :disabled="disabled" @update:model-value="value => update(index, 'peerIdPattern', String(value))" /></UFormField>
          <UFormField :label="t('trackerAdmin.pattern')"><UInput :model-value="exception.agentPattern" class="w-full" :maxlength="200" :disabled="disabled" @update:model-value="value => update(index, 'agentPattern', String(value))" /></UFormField>
          <UFormField :label="t('trackerAdmin.comment')"><UTextarea :model-value="exception.comment" class="w-full" :rows="2" :maxlength="200" :disabled="disabled" @update:model-value="value => update(index, 'comment', String(value))" /></UFormField>
          <div class="flex items-center justify-between gap-3"><USwitch :model-value="exception.enabled" :label="t('trackerAdmin.enabled')" :disabled="disabled" @update:model-value="value => update(index, 'enabled', value)" /><UTooltip :text="t('trackerAdmin.removeException')"><UButton type="button" icon="i-lucide-trash-2" color="error" variant="soft" :aria-label="t('trackerAdmin.removeException')" :disabled="disabled" @click="remove(index)" /></UTooltip></div>
        </div>
      </details>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AdminTrackerClientException } from '~/composables/useAdmin'
defineProps<{ disabled?: boolean }>()
const exceptions = defineModel<AdminTrackerClientException[]>({ required: true })
const { t } = useI18n()
function add() { exceptions.value = [...exceptions.value, { name: '', peerIdPattern: '', agentPattern: '', enabled: true, comment: '' }] }
function remove(index: number) { exceptions.value = exceptions.value.filter((_, i) => i !== index) }
function update(index: number, key: keyof AdminTrackerClientException, value: string | boolean) {
  exceptions.value = exceptions.value.map((exception, i) => i === index ? { ...exception, [key]: value } : exception)
}
</script>
