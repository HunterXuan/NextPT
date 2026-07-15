<template>
  <USelect
    :model-value="selectedLevel"
    class="mt-1 w-full"
    size="lg"
    :ui="{ base: 'h-10 w-full' }"
    :items="selectOptions"
    value-key="value"
    :disabled="disabled || roles.length === 0"
    @update:model-value="updateSelectedLevel"
  />
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
const selectOptions = computed(() => {
  const options = roleOptions.value.map(role => ({ value: role.level, label: roleNameWithLevel(role) }))
  if (!hasRoleLevel(selectedLevel.value)) {
    options.unshift({ value: selectedLevel.value, label: roleNameByLevel(selectedLevel.value) })
  }
  return options
})

function updateSelectedLevel(value: unknown) {
  selectedLevel.value = Number(value || 0)
}
</script>
