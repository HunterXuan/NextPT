<template>
  <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
    <div class="flex min-w-0 flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
      <div class="flex min-w-0 items-start gap-4">
        <IamUserAvatar
          :id="user?.user.id"
          :username="user?.user.username"
          :avatar="user?.profile.avatar"
          :alt="displayName"
          size="xl"
        />
        <div class="min-w-0 pt-0.5">
          <div class="flex flex-wrap items-center gap-2">
            <p class="truncate text-xl font-semibold text-slate-950 dark:text-white">{{ displayName }}</p>
            <UBadge :color="user?.role.isStaff ? 'primary' : 'neutral'" variant="soft">
              {{ roleDisplayName }}
            </UBadge>
            <UBadge color="success" variant="soft">
              {{ $t('user.status.normal') }}
            </UBadge>
          </div>
          <p class="mt-1 truncate text-sm text-slate-500 dark:text-slate-400">{{ user?.user.email || '-' }}</p>
        </div>
      </div>

      <dl class="grid gap-2 text-sm sm:grid-cols-2 lg:w-[min(380px,36vw)]">
        <div v-for="item in identityItems" :key="item.label" class="rounded-md bg-slate-50 px-3 py-2 dark:bg-slate-950/60">
          <dt class="text-xs text-slate-500 dark:text-slate-400">{{ item.label }}</dt>
          <dd class="mt-1 truncate font-medium text-slate-950 dark:text-white">{{ item.value }}</dd>
        </div>
      </dl>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { AuthUser } from '~/composables/useAuth'
import { formatDateTime } from '~/utils/format'

const props = defineProps<{
  user?: AuthUser | null
}>()

const { t, locale } = useI18n()

const displayName = computed(() => props.user?.user.username || t('user.fallbackName'))
const roleDisplayName = computed(() => props.user?.role.name || (props.user?.role.isStaff ? t('user.staff') : t('user.member')))
const identityItems = computed(() => [
  {
    label: t('user.fields.id'),
    value: props.user?.user.id ? `#${props.user.user.id}` : '-'
  },
  {
    label: t('user.fields.createdAt'),
    value: formatDateTime(props.user?.user.createdAt, locale.value)
  }
])
</script>
