<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
      <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ t('user.sessions.title') }}</h2>
      <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ t('user.sessions.subtitle') }}</p>
    </div>

    <div v-if="pending" class="space-y-2 p-4">
      <div v-for="item in 2" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>

    <div v-else-if="errorMessage" class="flex items-center gap-3 px-4 py-6 text-sm text-red-600 dark:text-red-300">
      <UIcon name="i-lucide-circle-alert" class="size-5 shrink-0" />
      <span>{{ errorMessage }}</span>
    </div>

    <div v-else-if="sessions.length === 0" class="flex flex-col items-center justify-center px-4 py-8 text-center">
      <UIcon name="i-lucide-monitor-off" class="size-7 text-slate-400" />
      <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">{{ t('user.sessions.empty') }}</p>
    </div>

    <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
      <article v-for="item in sessions" :key="item.id" class="flex items-center gap-3 px-4 py-3">
        <span
          class="flex size-9 shrink-0 items-center justify-center rounded-md"
          :class="item.current
            ? 'bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
            : 'bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-400'"
        >
          <UIcon :name="deviceIcon(item.userAgent)" class="size-4" />
        </span>

        <div class="min-w-0 flex-1">
          <div class="flex flex-wrap items-center gap-2">
            <span class="text-sm font-medium text-slate-950 dark:text-white">{{ sessionLabel(item.userAgent) }}</span>
            <UBadge v-if="item.current" color="primary" variant="subtle" size="sm">{{ t('user.sessions.current') }}</UBadge>
            <span class="font-mono text-xs text-slate-500 dark:text-slate-400">{{ item.ip || '-' }}</span>
          </div>
          <UTooltip :text="item.userAgent || '-'" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
            <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ item.userAgent || '-' }}</p>
          </UTooltip>
        </div>

        <div class="hidden shrink-0 text-right sm:block">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ t('user.sessions.lastSeen') }}</p>
          <time class="mt-1 block text-xs text-slate-600 dark:text-slate-300">{{ formatDateTime(item.lastSeenAt, locale) }}</time>
        </div>

        <UPopover
          :open="confirmSessionId === item.id"
          :content="{ side: 'top', align: 'end', sideOffset: 8 }"
          :ui="{ content: 'w-64 p-3' }"
          @update:open="setConfirmOpen(item.id, $event)"
        >
          <UTooltip :text="item.current ? t('user.sessions.logoutCurrent') : t('user.sessions.revoke')">
            <UButton
              color="error"
              variant="ghost"
              size="sm"
              icon="i-lucide-log-out"
              :aria-label="item.current ? t('user.sessions.logoutCurrent') : t('user.sessions.revoke')"
              :loading="deletingSessionId === item.id"
            />
          </UTooltip>
          <template #content="{ close }">
            <div class="space-y-3">
              <p class="text-sm font-medium text-slate-950 dark:text-white">{{ item.current ? t('user.sessions.logoutCurrentConfirm') : t('user.sessions.revokeConfirm') }}</p>
              <div class="flex justify-end gap-2">
                <UButton color="neutral" variant="ghost" size="xs" @click="closeConfirm(close)">{{ t('common.cancel') }}</UButton>
                <UButton color="error" size="xs" :loading="deletingSessionId === item.id" @click="removeSession(item, close)">{{ t('common.confirm') }}</UButton>
              </div>
            </div>
          </template>
        </UPopover>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { AuthSessionItem } from '~/composables/useAuth'
import { ApiError } from '~/composables/useApi'
import { formatDateTime } from '~/utils/format'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const { listSessions, deleteSession, clearLocalSession } = useAuth()

const sessions = ref<AuthSessionItem[]>([])
const pending = ref(false)
const errorMessage = ref('')
const confirmSessionId = ref('')
const deletingSessionId = ref('')

onMounted(loadSessions)

async function loadSessions() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await listSessions()
    sessions.value = data.list || []
  } catch (error) {
    sessions.value = []
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function setConfirmOpen(sessionId: string, open: boolean) {
  if (deletingSessionId.value) return
  confirmSessionId.value = open ? sessionId : ''
}

function closeConfirm(close?: () => void) {
  confirmSessionId.value = ''
  close?.()
}

async function removeSession(item: AuthSessionItem, close?: () => void) {
  deletingSessionId.value = item.id
  try {
    await deleteSession(item.id)
    closeConfirm(close)
    if (item.current) {
      clearLocalSession()
      await navigateTo(localePath('/login'))
      return
    }
    sessions.value = sessions.value.filter(session => session.id !== item.id)
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    deletingSessionId.value = ''
  }
}

function sessionLabel(userAgent: string) {
  const browser = browserName(userAgent)
  const system = systemName(userAgent)
  return [browser, system].filter(Boolean).join(' · ') || t('user.sessions.unknownDevice')
}

function browserName(userAgent: string) {
  if (/Edg\//i.test(userAgent)) return 'Edge'
  if (/Firefox\//i.test(userAgent)) return 'Firefox'
  if (/Chrome\//i.test(userAgent)) return 'Chrome'
  if (/Safari\//i.test(userAgent)) return 'Safari'
  return ''
}

function systemName(userAgent: string) {
  if (/Android/i.test(userAgent)) return 'Android'
  if (/iPhone|iPad|iPod/i.test(userAgent)) return 'iOS'
  if (/Windows/i.test(userAgent)) return 'Windows'
  if (/Macintosh|Mac OS X/i.test(userAgent)) return 'macOS'
  if (/Linux/i.test(userAgent)) return 'Linux'
  return ''
}

function deviceIcon(userAgent: string) {
  return /Android|iPhone|iPad|iPod|Mobile/i.test(userAgent) ? 'i-lucide-smartphone' : 'i-lucide-monitor'
}
</script>
