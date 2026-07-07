<template>
  <UTooltip
    v-if="tooltipText"
    :text="tooltipText"
    :content="{ side: tooltipSide, sideOffset: 8 }"
    :delay-duration="delayDuration"
  >
    <span :class="wrapperClass">
      <UButton v-bind="attrs" :to="effectiveTo" :disabled="effectiveDisabled">
        <slot />
      </UButton>
    </span>
  </UTooltip>
  <span v-else :class="wrapperClass">
    <UButton v-bind="attrs" :to="effectiveTo" :disabled="effectiveDisabled">
      <slot />
    </UButton>
  </span>
</template>

<script setup lang="ts">
defineOptions({
  inheritAttrs: false
})

const props = withDefaults(defineProps<{
  permission?: string
  allowed?: boolean
  disabled?: boolean
  to?: string
  tooltip?: string
  disabledTooltip?: string
  tooltipSide?: 'top' | 'right' | 'bottom' | 'left'
  delayDuration?: number
}>(), {
  permission: '',
  allowed: undefined,
  disabled: false,
  to: undefined,
  tooltip: '',
  disabledTooltip: '',
  tooltipSide: 'top',
  delayDuration: 120
})

const attrs = useAttrs()
const { t } = useI18n()
const { hasPermission } = useAuth()

const permissionAllowed = computed(() => props.permission ? hasPermission(props.permission) : true)
const canUse = computed(() => {
  const allowed = props.allowed ?? true
  return allowed && permissionAllowed.value
})
const effectiveDisabled = computed(() => props.disabled || !canUse.value)
const effectiveTo = computed(() => canUse.value ? props.to : undefined)
const tooltipText = computed(() => {
  if (!canUse.value) return props.disabledTooltip || t('common.noPermission')
  return props.tooltip
})
const wrapperClass = computed(() => attrs.block !== undefined ? 'block' : 'inline-flex')
</script>
