<template>
  <UPopover
    v-if="userId > 0"
    v-model:open="open"
    mode="hover"
    :open-delay="600"
    :close-delay="300"
    :content="{ side: 'top', align: 'start', sideOffset: 8 }"
    :ui="{ content: 'w-80 max-w-[calc(100vw-1rem)] overflow-hidden p-0' }"
    :class="attrs.class"
  >
    <slot :label="displayName" :user="triggerUser">
      <span class="inline-flex min-w-0 cursor-default items-center gap-2 text-inherit transition-colors hover:text-sky-700 dark:hover:text-sky-300">
        <IamUserAvatar v-if="showAvatar" :user="triggerUser" :size="avatarSize" />
        <span class="truncate">{{ displayName }}</span>
      </span>
    </slot>

    <template #content>
      <div v-if="pending" class="p-3.5">
        <div class="flex items-center gap-3">
          <div class="size-10 animate-pulse rounded-md bg-slate-200 dark:bg-slate-800" />
          <div class="min-w-0 flex-1 space-y-2">
            <div class="h-4 w-28 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
            <div class="h-3 w-40 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
          </div>
        </div>
      </div>

      <div v-else-if="errorMessage" class="flex items-center gap-2 p-4 text-sm text-red-600 dark:text-red-300">
        <UIcon name="i-lucide-circle-alert" class="size-4 shrink-0" />
        <span>{{ errorMessage }}</span>
      </div>

      <div v-else-if="publicUser">
        <div class="flex items-center gap-3 px-4 py-3.5">
          <IamUserAvatar
            :id="publicUser.user.id"
            :username="publicUser.user.username"
            :avatar="publicUser.profile.avatar"
            size="md"
          />
          <div class="min-w-0 flex-1">
            <div class="flex min-w-0 items-center gap-1.5">
              <p class="truncate text-base font-semibold text-slate-950 dark:text-white">{{ publicUser.user.username }}</p>
              <UBadge :color="publicUser.role.isStaff ? 'warning' : 'primary'" variant="soft" size="xs" class="shrink-0">{{ roleName }}</UBadge>
            </div>
            <div class="mt-1 flex items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400">
              <span class="tabular-nums">#{{ publicUser.user.id }}</span>
              <span class="text-slate-300 dark:text-slate-700">&middot;</span>
              <span class="truncate">{{ $t('user.publicCard.joined', { date: joinedAt }) }}</span>
            </div>
            <p v-if="publicUser.profile.signature" class="mt-1.5 line-clamp-1 text-xs leading-5 text-slate-600 dark:text-slate-300">
              {{ publicUser.profile.signature }}
            </p>
          </div>
        </div>

        <dl class="grid grid-cols-4 border-t border-slate-200 bg-slate-50/70 dark:border-slate-800 dark:bg-slate-950/40">
          <div
            v-for="(item, index) in statItems"
            :key="item.label"
            class="min-w-0 border-slate-200 px-2 py-2.5 text-center dark:border-slate-800"
            :class="index < statItems.length - 1 ? 'border-r' : ''"
          >
            <dt class="truncate text-[10px] leading-4 text-slate-500 dark:text-slate-400">{{ item.label }}</dt>
            <dd class="mt-0.5 truncate text-[13px] font-semibold tabular-nums text-slate-950 dark:text-white" :title="item.value">{{ item.value }}</dd>
          </div>
        </dl>

        <div v-if="renderedBio" class="border-t border-slate-200 px-4 py-2.5 dark:border-slate-800">
          <div class="rich-text rich-text-compact max-h-16 overflow-hidden text-xs leading-5" v-html="renderedBio" />
        </div>
      </div>
    </template>
  </UPopover>

  <slot v-else :label="displayName" :user="triggerUser">
    <span class="inline-flex min-w-0 items-center gap-2 text-inherit" :class="attrs.class">
      <IamUserAvatar v-if="showAvatar" :user="triggerUser" :size="avatarSize" />
      <span class="truncate">{{ displayName }}</span>
    </span>
  </slot>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { PublicUser } from '~/composables/useIamUsers'
import { formatBytes, formatDateOnly } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

defineOptions({ inheritAttrs: false })

interface UserReference {
  id?: number | string
  username?: string
  avatar?: string | null
}

const props = withDefaults(defineProps<{
  user?: UserReference | null
  id?: number | string
  username?: string
  avatar?: string | null
  fallback?: string
  showAvatar?: boolean
  avatarSize?: 'xs' | 'sm' | 'md' | 'lg'
}>(), {
  user: null,
  id: 0,
  username: '',
  avatar: '',
  fallback: '-',
  showAvatar: false,
  avatarSize: 'xs'
})

const { t, locale } = useI18n()
const route = useRoute()
const attrs = useAttrs()
const iamUsers = useIamUsers()
const open = ref(false)
const pending = ref(false)
const errorMessage = ref('')
const publicUser = ref<PublicUser | null>(null)

const userId = computed(() => Number(props.id || props.user?.id || 0))
const knownUsername = computed(() => String(props.username || props.user?.username || '').trim())
const knownAvatar = computed(() => props.avatar || props.user?.avatar || '')
const displayName = computed(() => publicUser.value?.user.username || knownUsername.value || (userId.value > 0 ? `#${userId.value}` : props.fallback))
const triggerUser = computed(() => ({
  id: userId.value,
  username: displayName.value,
  avatar: publicUser.value?.profile.avatar || knownAvatar.value
}))
const roleName = computed(() => publicUser.value?.role.name || (publicUser.value?.role.isStaff ? t('user.staff') : t('user.member')))
const joinedAt = computed(() => formatDateOnly(publicUser.value?.user.createdAt, locale.value))
const renderedBio = computed(() => renderUserMarkdown(publicUser.value?.profile.info || '').trim())
const ratio = computed(() => {
  const stat = publicUser.value?.stat
  if (!stat) return '-'
  if (stat.downloaded === 0 && stat.uploaded > 0) return '∞'
  return Number(stat.shareRatio || 0).toFixed(2)
})
const statItems = computed(() => [
  { label: t('user.publicCard.stats.uploaded'), value: formatBytes(publicUser.value?.stat.uploaded) },
  { label: t('user.publicCard.stats.downloaded'), value: formatBytes(publicUser.value?.stat.downloaded) },
  { label: t('user.publicCard.stats.ratio'), value: ratio.value },
  { label: t('user.publicCard.stats.seedTime'), value: formatDuration(publicUser.value?.stat.seedTime || 0) }
])

watch([userId, locale], () => {
  open.value = false
  publicUser.value = userId.value > 0 ? iamUsers.getCachedUser(userId.value) : null
  pending.value = false
  errorMessage.value = ''
}, { immediate: true })

watch(() => route.fullPath, () => {
  open.value = false
})

watch(open, (value) => {
  if (value) loadUser()
})

async function loadUser() {
  if (userId.value <= 0) return

  const cached = iamUsers.getCachedUser(userId.value)
  if (cached) {
    publicUser.value = cached
    return
  }

  const requestedId = userId.value
  const requestedLocale = locale.value
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await iamUsers.getUser(requestedId)
    if (userId.value === requestedId && locale.value === requestedLocale) publicUser.value = data
  } catch (error) {
    if (userId.value === requestedId && locale.value === requestedLocale) {
      errorMessage.value = error instanceof ApiError ? error.message : t('user.publicCard.loadFailed')
    }
  } finally {
    if (userId.value === requestedId && locale.value === requestedLocale) pending.value = false
  }
}

function formatDuration(seconds: number) {
  const normalized = Math.max(0, Math.floor(Number(seconds) || 0))
  if (normalized >= 86400) return t('user.duration.days', { count: Math.floor(normalized / 86400) })
  if (normalized >= 3600) return t('user.duration.hours', { count: Math.floor(normalized / 3600) })
  if (normalized >= 60) return t('user.duration.minutes', { count: Math.floor(normalized / 60) })
  return t('user.duration.seconds', { count: normalized })
}
</script>
