<template>
  <UCard class="rounded-lg">
    <template #header>
      <div>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.invites.title') }}</h2>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
          {{ $t('user.invites.summary', { total: numberFormatter.format(inviteTotal) }) }}
        </p>
      </div>
    </template>

    <section class="rounded-lg border border-slate-200 bg-slate-50 p-4 dark:border-slate-800 dark:bg-slate-950/60">
      <div class="mb-4 flex flex-col gap-1">
        <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('user.invites.sendTitle') }}</h3>
        <p class="text-sm text-slate-500 dark:text-slate-400">
          {{ availableInvites.length > 0 ? $t('user.invites.sendDescription') : $t('user.invites.noAvailable') }}
        </p>
      </div>

      <form class="grid gap-3 md:grid-cols-[minmax(0,1fr)_220px_auto] md:items-end" @submit.prevent="handleInviteSend">
        <UFormField :label="$t('user.invites.email')" required>
          <UInput
            v-model="inviteEmail"
            class="w-full"
            type="email"
            size="lg"
            icon="i-lucide-mail"
            :disabled="inviteSendPending || availableInvites.length === 0"
          />
        </UFormField>
        <UFormField :label="$t('user.invites.code')" required>
          <USelect
            v-model="selectedInviteHash"
            class="w-full"
            size="lg"
            :items="inviteOptions"
            value-key="value"
            :placeholder="$t('user.invites.selectCode')"
            :disabled="inviteSendPending || availableInvites.length === 0"
          />
        </UFormField>
        <UButton type="submit" color="primary" size="lg" icon="i-lucide-send" :loading="inviteSendPending" :disabled="!canSendInvite">
          {{ $t('user.invites.send') }}
        </UButton>
      </form>
    </section>

    <div v-if="invitesError" class="mt-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
      {{ invitesError }}
    </div>
    <div v-else-if="invitesPending && invites.length === 0" class="mt-4 space-y-2">
      <div v-for="item in 6" :key="item" class="h-12 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>
    <div v-else-if="invites.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-3 py-12 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
      {{ $t('user.invites.empty') }}
    </div>
    <div v-else class="mt-4 overflow-hidden rounded-lg border border-slate-200 dark:border-slate-800">
      <div class="hidden grid-cols-[minmax(0,1fr)_minmax(0,1fr)_112px_132px] gap-4 border-b border-slate-200 bg-slate-50 px-4 py-2.5 text-xs font-medium text-slate-500 md:grid dark:border-slate-800 dark:bg-slate-950/60 dark:text-slate-400">
        <span>{{ $t('user.invites.code') }}</span>
        <span>{{ $t('user.invites.invitee') }}</span>
        <span class="text-center">{{ $t('user.invites.statusLabel') }}</span>
        <span class="text-right">{{ $t('user.invites.updatedAt') }}</span>
      </div>

      <div class="divide-y divide-slate-100 dark:divide-slate-800">
        <div v-for="invite in invites" :key="invite.id" class="grid gap-3 px-4 py-3 text-sm md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_112px_132px] md:items-center md:gap-4">
          <div class="min-w-0">
            <p class="font-mono text-xs text-slate-700 dark:text-slate-300">{{ formatInviteHash(invite.hash) }}</p>
            <p class="mt-1 text-xs text-slate-500 md:hidden dark:text-slate-400">{{ $t('user.invites.createdAt') }} {{ formatDateTime(invite.createdAt, locale) }}</p>
          </div>
          <p class="min-w-0 truncate text-sm text-slate-600 dark:text-slate-300">
            {{ invite.inviteeEmail || invite.inviteeName || '-' }}
          </p>
          <UBadge class="w-fit md:justify-self-center" :color="inviteStatusColor(invite.status)" variant="soft">
            {{ inviteStatusLabel(invite.status) }}
          </UBadge>
          <p class="hidden text-right text-xs text-slate-500 md:block dark:text-slate-400">
            {{ inviteTimeLabel(invite) }}
          </p>
        </div>
      </div>
    </div>

    <AppPager
      v-if="inviteTotal > 0"
      class="mt-4"
      size="sm"
      :page="invitePage"
      :total="inviteTotal"
      :page-size="inviteSize"
      :disabled="invitesPending"
      @page-change="goToInvitePage"
    />
  </UCard>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useInvites, type InviteItem } from '~/composables/useInvites'
import { formatDateTime } from '~/utils/format'

const { t, locale } = useI18n()
const toast = useToast()
const { fetchUser } = useAuth()
const inviteService = useInvites()

const invites = ref<InviteItem[]>([])
const inviteTotal = ref(0)
const invitePage = ref(1)
const inviteSize = 10
const invitesPending = ref(true)
const invitesError = ref('')
const inviteSendPending = ref(false)
const inviteEmail = ref('')
const selectedInviteHash = ref('')

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const availableInvites = computed(() => invites.value.filter((invite) => invite.status === 0))
const inviteOptions = computed(() => availableInvites.value.map((invite) => ({
  label: formatInviteHash(invite.hash),
  value: invite.hash
})))
const inviteTotalPages = computed(() => Math.max(1, Math.ceil(inviteTotal.value / inviteSize)))
const canSendInvite = computed(() => {
  return Boolean(selectedInviteHash.value && inviteEmail.value.trim() && !inviteSendPending.value && availableInvites.value.length > 0)
})

onMounted(loadInvites)

async function loadInvites() {
  if (invitesPending.value && invites.value.length > 0) return

  invitesPending.value = true
  invitesError.value = ''

  try {
    let data = await inviteService.listInvites({
      page: invitePage.value,
      size: inviteSize
    })

    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / inviteSize))
    if (invitePage.value > nextTotalPages) {
      invitePage.value = nextTotalPages
      data = await inviteService.listInvites({
        page: invitePage.value,
        size: inviteSize
      })
    }

    invites.value = data.list || []
    inviteTotal.value = data.total || 0

    if (!selectedInviteHash.value || !availableInvites.value.some((invite) => invite.hash === selectedInviteHash.value)) {
      selectedInviteHash.value = availableInvites.value[0]?.hash || ''
    }
  } catch (error) {
    invites.value = []
    inviteTotal.value = 0
    selectedInviteHash.value = ''
    invitesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    invitesPending.value = false
  }
}

function goToInvitePage(page: number) {
  invitePage.value = Math.min(Math.max(1, page), inviteTotalPages.value)
  loadInvites()
}

async function handleInviteSend() {
  if (!canSendInvite.value) return

  inviteSendPending.value = true
  try {
    await inviteService.sendInvite(selectedInviteHash.value, inviteEmail.value.trim())
    toast.add({
      title: t('user.invites.sent'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    inviteEmail.value = ''
    await Promise.all([loadInvites(), fetchUser()])
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    inviteSendPending.value = false
  }
}

function formatInviteHash(hash: string) {
  if (!hash) return '-'
  if (hash.length <= 12) return hash
  return `${hash.slice(0, 6)}...${hash.slice(-6)}`
}

function inviteTimeLabel(invite: InviteItem) {
  if (invite.usedAt) return formatDateTime(invite.usedAt, locale.value)
  if (invite.expireAt) return formatDateTime(invite.expireAt, locale.value)
  return formatDateTime(invite.createdAt, locale.value)
}

function inviteStatusLabel(status: number) {
  const key = `user.invites.status.${status}`
  const translated = t(key)
  return translated === key ? String(status) : translated
}

function inviteStatusColor(status: number) {
  if (status === 0) return 'success'
  if (status === 1) return 'primary'
  if (status === 2) return 'neutral'
  if (status === 3) return 'error'
  return 'warning'
}
</script>
