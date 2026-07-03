<template>
  <UCard class="rounded-lg">
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.invites.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
            {{ $t('user.invites.summary', { available: numberFormatter.format(availableInvites.length), total: numberFormatter.format(inviteTotal) }) }}
          </p>
        </div>
        <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" :loading="invitesPending" @click="loadInvites">
          {{ $t('common.refresh') }}
        </UButton>
      </div>
    </template>

    <form class="grid gap-3 md:grid-cols-[minmax(0,1fr)_220px_auto] md:items-end" @submit.prevent="handleInviteSend">
      <UFormField :label="$t('user.invites.email')" required>
        <UInput v-model="inviteEmail" class="w-full" type="email" icon="i-lucide-mail" :disabled="inviteSendPending || availableInvites.length === 0" />
      </UFormField>
      <UFormField :label="$t('user.invites.code')" required>
        <select
          v-model="selectedInviteHash"
          class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
          :disabled="inviteSendPending || availableInvites.length === 0"
        >
          <option value="">{{ $t('user.invites.selectCode') }}</option>
          <option v-for="invite in availableInvites" :key="invite.id" :value="invite.hash">
            {{ invite.hash }}
          </option>
        </select>
      </UFormField>
      <UButton type="submit" color="primary" icon="i-lucide-send" :loading="inviteSendPending" :disabled="!canSendInvite">
        {{ $t('user.invites.send') }}
      </UButton>
    </form>

    <div v-if="invitesError" class="mt-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
      {{ invitesError }}
    </div>
    <div v-else-if="invitesPending && invites.length === 0" class="mt-4 space-y-2">
      <div v-for="item in 6" :key="item" class="h-12 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>
    <div v-else-if="invites.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-3 py-12 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
      {{ $t('user.invites.empty') }}
    </div>
    <div v-else class="mt-4 divide-y divide-slate-100 dark:divide-slate-800">
      <div v-for="invite in invites" :key="invite.id" class="grid gap-3 py-3 text-sm md:grid-cols-[minmax(0,1fr)_180px_100px] md:items-center">
        <p class="min-w-0 truncate font-mono text-xs text-slate-700 dark:text-slate-300">{{ invite.hash }}</p>
        <p class="min-w-0 truncate text-xs text-slate-500 dark:text-slate-400">
          {{ invite.inviteeEmail || invite.inviteeName || '-' }}
        </p>
        <UBadge class="justify-self-start md:justify-self-end" :color="inviteStatusColor(invite.status)" variant="soft">
          {{ inviteStatusLabel(invite.status) }}
        </UBadge>
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useInvites, type InviteItem } from '~/composables/useInvites'

const { t, locale } = useI18n()
const toast = useToast()
const { fetchUser } = useAuth()
const inviteService = useInvites()

const invites = ref<InviteItem[]>([])
const inviteTotal = ref(0)
const invitePage = ref(1)
const inviteSize = 20
const invitesPending = ref(true)
const invitesError = ref('')
const inviteSendPending = ref(false)
const inviteEmail = ref('')
const selectedInviteHash = ref('')

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const availableInvites = computed(() => invites.value.filter((invite) => invite.status === 0))
const canSendInvite = computed(() => {
  return Boolean(selectedInviteHash.value && inviteEmail.value.trim() && !inviteSendPending.value && availableInvites.value.length > 0)
})

onMounted(loadInvites)

async function loadInvites() {
  if (invitesPending.value && invites.value.length > 0) return

  invitesPending.value = true
  invitesError.value = ''

  try {
    const data = await inviteService.listInvites({
      page: invitePage.value,
      size: inviteSize
    })
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
