<template>
  <select
    v-model.number="selectedLevel"
    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950 dark:disabled:bg-slate-900 dark:disabled:text-slate-500"
    :disabled="disabled || roles.length === 0"
  >
    <option v-if="!hasRoleLevel(selectedLevel)" :value="selectedLevel">{{ roleNameByLevel(selectedLevel) }}</option>
    <option v-for="role in roleOptions" :key="role.id" :value="role.level">{{ roleNameWithLevel(role) }}</option>
  </select>
</template>

<script setup lang="ts">
import type { AdminIamRole } from '~/composables/useAdmin'

const props = withDefaults(defineProps<{
  modelValue: number
  roles: AdminIamRole[]
  disabled?: boolean
}>(), {
  roles: () => [],
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const selectedLevel = computed({
  get: () => Number(props.modelValue || 0),
  set: (value) => emit('update:modelValue', Number(value || 0))
})
const rolesRef = toRef(props, 'roles')
const { roleOptions, roleNameWithLevel, roleNameByLevel, hasRoleLevel } = useAdminIamRoleLevels(rolesRef)
</script>
