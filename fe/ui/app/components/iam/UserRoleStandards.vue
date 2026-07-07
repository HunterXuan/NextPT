<template>
  <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
      <div>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.roleStandards.title') }}</h2>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.roleStandards.subtitle') }}</p>
      </div>
    </div>

    <div v-if="pending" class="mt-4 h-64 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />

    <div v-else-if="error" class="mt-4 rounded-md border border-red-200 bg-red-50 px-4 py-6 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
      {{ error }}
    </div>

    <div v-else-if="orderedRoles.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-8 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
      {{ $t('user.roleStandards.empty') }}
    </div>

    <div v-else-if="activeRole" class="mt-4">
      <div class="rounded-md bg-slate-50 p-4 dark:bg-slate-950/50">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-2">
              <span
                class="flex size-9 shrink-0 items-center justify-center rounded-md"
                :class="activeRole.isStaff ? 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-200' : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-200'"
              >
                <UIcon :name="activeRole.isStaff ? 'i-lucide-shield-check' : 'i-lucide-medal'" class="size-5" />
              </span>
              <h3 class="truncate text-lg font-semibold text-slate-950 dark:text-white">{{ activeRole.name }}</h3>
              <UBadge v-if="activeRole.id === user?.role.id" color="primary" variant="soft">
                {{ $t('user.roleStandards.currentShort') }}
              </UBadge>
            </div>
            <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">
              {{ $t('user.roleStandards.position', { current: numberFormatter.format(activeIndex + 1), total: numberFormatter.format(orderedRoles.length) }) }}
            </p>
          </div>

          <div class="flex shrink-0 items-center gap-1.5">
            <UTooltip :text="$t('user.roleStandards.previous')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
              <UButton
                color="neutral"
                variant="ghost"
                size="sm"
                icon="i-lucide-chevron-left"
                class="rounded-md"
                :disabled="activeIndex <= 0"
                :aria-label="$t('user.roleStandards.previous')"
                @click="goPrevious"
              />
            </UTooltip>
            <span class="min-w-14 text-center text-sm tabular-nums text-slate-500 dark:text-slate-400">
              {{ numberFormatter.format(activeIndex + 1) }} / {{ numberFormatter.format(orderedRoles.length) }}
            </span>
            <UTooltip :text="$t('user.roleStandards.next')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
              <UButton
                color="neutral"
                variant="ghost"
                size="sm"
                icon="i-lucide-chevron-right"
                class="rounded-md"
                :disabled="activeIndex >= orderedRoles.length - 1"
                :aria-label="$t('user.roleStandards.next')"
                @click="goNext"
              />
            </UTooltip>
          </div>
        </div>

        <div class="mt-4">
          <div v-if="activeRole.isStaff" class="rounded-md bg-white px-4 py-4 dark:bg-slate-900">
            <div class="flex min-w-0 items-start gap-3">
              <span class="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-md bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-200">
                <UIcon name="i-lucide-shield-check" class="size-4" />
              </span>
              <div>
                <p class="font-medium text-slate-950 dark:text-white">{{ $t('user.roleStandards.staffTitle') }}</p>
                <p class="mt-1 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('user.roleStandards.notAuto') }}</p>
              </div>
            </div>
          </div>

          <div v-else class="grid gap-3 lg:grid-cols-2">
            <RulePanel :title="$t('user.roleStandards.promotion')" :groups="ruleGroups(activeRole, 'promotion')" />
            <RulePanel :title="$t('user.roleStandards.demotion')" :groups="ruleGroups(activeRole, 'demotion')" />
          </div>
        </div>
      </div>

      <div class="mt-3 flex flex-wrap items-center justify-center gap-1.5">
        <UTooltip
          v-for="(role, index) in orderedRoles"
          :key="role.id"
          :text="role.name"
          :content="{ side: 'top', sideOffset: 8 }"
          :delay-duration="600"
        >
          <button
            type="button"
            class="h-2.5 rounded-full transition"
            :class="roleDotClass(role, index)"
            :aria-label="$t('user.roleStandards.viewRole', { role: role.name })"
            @click="goTo(index)"
          />
        </UTooltip>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { defineComponent, h, type PropType } from 'vue'
import type { AuthRoleItem, AuthRoleRuleCondition, AuthUser } from '~/composables/useAuth'

const RulePanel = defineComponent({
  props: {
    title: {
      type: String,
      required: true
    },
    groups: {
      type: Array as PropType<string[][]>,
      required: true
    }
  },
  setup(props) {
    const { t } = useI18n()
    return () => h('div', { class: 'rounded-md bg-white/80 p-3 ring-1 ring-slate-200/70 dark:bg-slate-900/80 dark:ring-slate-800' }, [
      h('p', { class: 'text-sm font-semibold text-slate-950 dark:text-white' }, props.title),
      props.groups.length > 0
        ? h('div', { class: 'mt-2 max-h-28 space-y-1.5 overflow-y-auto pr-1' }, props.groups.map((group, index) => h('p', { class: 'text-sm leading-6 text-slate-600 dark:text-slate-300' }, [
            index > 0 ? h('span', { class: 'mr-1 text-slate-400 dark:text-slate-500' }, t('user.roleStandards.or')) : null,
            group.join(` ${t('user.roleStandards.and')} `)
          ])))
        : h('p', { class: 'mt-2 text-sm leading-6 text-slate-400 dark:text-slate-500' }, t('user.roleStandards.noRules'))
    ])
  }
})

const props = withDefaults(defineProps<{
  user?: AuthUser | null
  roles?: AuthRoleItem[]
  pending?: boolean
  error?: string
}>(), {
  user: null,
  roles: () => [],
  pending: false,
  error: ''
})

const { t, locale } = useI18n()
const activeIndex = ref(0)
const initializedRoleId = ref<number | null>(null)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const orderedRoles = computed(() => [...props.roles].sort(compareRole))
const activeRole = computed(() => orderedRoles.value[activeIndex.value] || null)

watch(
  () => [orderedRoles.value.length, props.user?.role.id] as const,
  () => {
    const roles = orderedRoles.value
    if (roles.length === 0) {
      activeIndex.value = 0
      initializedRoleId.value = null
      return
    }

    const currentRoleId = props.user?.role.id || null
    if (currentRoleId && initializedRoleId.value !== currentRoleId) {
      const index = roles.findIndex((role) => role.id === currentRoleId)
      activeIndex.value = index >= 0 ? index : 0
      initializedRoleId.value = currentRoleId
      return
    }

    if (activeIndex.value >= roles.length) {
      activeIndex.value = roles.length - 1
    }
  },
  { immediate: true }
)

function compareRole(a: AuthRoleItem, b: AuthRoleItem) {
  return a.level - b.level || a.id - b.id
}

function goPrevious() {
  activeIndex.value = Math.max(0, activeIndex.value - 1)
}

function goNext() {
  activeIndex.value = Math.min(orderedRoles.value.length - 1, activeIndex.value + 1)
}

function goTo(index: number) {
  activeIndex.value = Math.min(Math.max(index, 0), orderedRoles.value.length - 1)
}

function roleDotClass(role: AuthRoleItem, index: number) {
  const classes = []
  if (index === activeIndex.value) {
    classes.push('w-8 bg-sky-500 dark:bg-sky-400')
  } else {
    classes.push('w-2.5 bg-slate-300 hover:bg-slate-400 dark:bg-slate-700 dark:hover:bg-slate-600')
  }
  if (role.id === props.user?.role.id && index !== activeIndex.value) {
    classes.push('bg-sky-300 dark:bg-sky-700')
  }
  if (role.isStaff) {
    classes.push('opacity-70')
  }
  return classes
}

function ruleGroups(role: AuthRoleItem, type: 'promotion' | 'demotion') {
  const groups = role.rules?.[type]
  if (!Array.isArray(groups)) return []
  return groups
    .map((group) => conditionLabels(group as AuthRoleRuleCondition))
    .filter((labels) => labels.length > 0)
}

function conditionLabels(condition: AuthRoleRuleCondition) {
  const labels: string[] = []

  pushNumber(labels, condition.accountAgeDaysGte, 'accountAgeDaysGte')
  pushGiB(labels, condition.downloadedGiBGte, 'downloadedGiBGte')
  pushGiB(labels, condition.downloadedGiBGt, 'downloadedGiBGt')
  pushGiB(labels, condition.downloadedGiBLte, 'downloadedGiBLte')
  pushRatio(labels, condition.ratioGt, 'ratioGt')
  pushRatio(labels, condition.ratioGte, 'ratioGte')
  pushRatio(labels, condition.ratioLt, 'ratioLt')

  for (const [key, value] of Object.entries(condition)) {
    if (knownConditionKeys.has(key) || value === null || value === undefined || value === '') continue
    labels.push(`${key}: ${String(value)}`)
  }

  return labels
}

const knownConditionKeys = new Set([
  'accountAgeDaysGte',
  'downloadedGiBGte',
  'downloadedGiBGt',
  'downloadedGiBLte',
  'ratioGt',
  'ratioGte',
  'ratioLt'
])

function pushNumber(labels: string[], value: unknown, key: string) {
  if (!isFiniteNumber(value)) return
  labels.push(t(`user.roleStandards.conditions.${key}`, { count: numberFormatter.value.format(value) }))
}

function pushGiB(labels: string[], value: unknown, key: string) {
  if (!isFiniteNumber(value)) return
  labels.push(t(`user.roleStandards.conditions.${key}`, { size: `${numberFormatter.value.format(value)} GiB` }))
}

function pushRatio(labels: string[], value: unknown, key: string) {
  if (!isFiniteNumber(value)) return
  labels.push(t(`user.roleStandards.conditions.${key}`, { ratio: numberFormatter.value.format(value) }))
}

function isFiniteNumber(value: unknown): value is number {
  return typeof value === 'number' && Number.isFinite(value)
}
</script>
