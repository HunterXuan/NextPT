<template>
  <div class="space-y-3">
    <p class="text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.actions.confirmDeleteTitle') }}</p>
    <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.actions.deleteIrreversible') }}</p>
    <div class="flex justify-end gap-2">
      <UButton color="neutral" variant="ghost" size="xs" type="button" :disabled="pending" @click="$emit('cancel')">{{ $t('common.cancel') }}</UButton>
      <UButton color="error" size="xs" type="button" icon="i-lucide-trash-2" :loading="pending" @click="$emit('confirm')">{{ $t('admin.actions.delete') }}</UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ pending?: boolean }>()
defineEmits<{ cancel: [], confirm: [] }>()
</script>
