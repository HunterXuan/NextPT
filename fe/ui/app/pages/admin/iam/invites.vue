<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(0,720px)_minmax(460px,1fr)] xl:items-start">
        <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.invites.form.title') }}</h2>
          </div>

          <form class="grid gap-4 p-4" @submit.prevent="grantInvites">
            <div class="grid gap-4 lg:grid-cols-[180px_minmax(0,1fr)]">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.amount') }}</span>
                <input v-model.number="form.amount" type="number" min="1" step="1" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="saving">
              </label>

              <div>
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.targetMode') }}</span>
                <div class="mt-1 grid gap-2 sm:grid-cols-2">
                  <button
                    type="button"
                    class="flex h-10 items-center justify-center rounded-md border px-3 text-sm font-semibold transition"
                    :class="targetMode === 'site' ? 'border-sky-300 bg-sky-50 text-sky-950 dark:border-sky-700 dark:bg-sky-950/40 dark:text-sky-100' : 'border-slate-200 bg-white text-slate-700 hover:border-slate-300 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-200 dark:hover:border-slate-700'"
                    :disabled="saving"
                    @click="setTargetMode('site')"
                  >
                    {{ $t('admin.iam.invites.form.siteTarget') }}
                  </button>
                  <button
                    type="button"
                    class="flex h-10 items-center justify-center rounded-md border px-3 text-sm font-semibold transition"
                    :class="targetMode === 'roles' ? 'border-sky-300 bg-sky-50 text-sky-950 dark:border-sky-700 dark:bg-sky-950/40 dark:text-sky-100' : 'border-slate-200 bg-white text-slate-700 hover:border-slate-300 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-200 dark:hover:border-slate-700'"
                    :disabled="saving"
                    @click="setTargetMode('roles')"
                  >
                    {{ $t('admin.iam.invites.form.roleTarget') }}
                  </button>
                </div>
              </div>
            </div>

            <div v-if="targetMode === 'roles'" class="grid gap-2">
              <div class="flex items-center justify-between gap-3">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.roleIds') }}</span>
                <button
                  v-if="form.roleIds.length > 0"
                  type="button"
                  class="text-xs font-medium text-slate-500 transition hover:text-slate-950 dark:text-slate-400 dark:hover:text-white"
                  :disabled="saving"
                  @click="form.roleIds = []"
                >
                  {{ $t('common.clear') }}
                </button>
              </div>
              <div class="grid max-h-36 gap-1 overflow-y-auto rounded-md border border-slate-200 bg-slate-50 p-2 sm:grid-cols-2 dark:border-slate-800 dark:bg-slate-950/60">
                <p v-if="rolesPending" class="px-2 py-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.form.rolesLoading') }}</p>
                <p v-else-if="roleOptions.length === 0" class="px-2 py-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.form.rolesEmpty') }}</p>
                <template v-else>
                  <label
                    v-for="role in roleOptions"
                    :key="role.id"
                    class="flex cursor-pointer items-center gap-3 rounded-md px-2 py-2 text-sm transition hover:bg-white dark:hover:bg-slate-900"
                  >
                    <input
                      v-model="form.roleIds"
                      type="checkbox"
                      class="size-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600"
                      :value="role.id"
                      :disabled="saving"
                    >
                    <span class="min-w-0 truncate text-slate-700 dark:text-slate-200">{{ roleNameWithLevel(role) }}</span>
                  </label>
                </template>
              </div>
              <p class="text-xs" :class="targetMode === 'roles' && form.roleIds.length === 0 ? 'text-red-500 dark:text-red-300' : 'text-slate-500 dark:text-slate-400'">
                {{ targetMode === 'roles' && form.roleIds.length === 0 ? $t('admin.iam.invites.form.roleRequired') : $t('admin.iam.invites.form.roleIdsHint') }}
              </p>
            </div>

            <div class="grid gap-4 lg:grid-cols-[180px_minmax(0,1fr)]">
              <label class="flex h-10 items-center justify-between gap-4 rounded-md border border-slate-200 px-3 dark:border-slate-800">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.isTemp') }}</span>
                <input v-model="form.isTemp" type="checkbox" class="size-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600" :disabled="saving">
              </label>

              <label class="block">
                <span class="sr-only">{{ $t('admin.iam.invites.form.expireAt') }}</span>
                <input v-model="form.expireAt" type="datetime-local" class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="saving || !form.isTemp">
                <span class="mt-1 block text-xs text-slate-500 dark:text-slate-400">{{ expireHint }}</span>
              </label>
            </div>

            <div class="rounded-md border border-slate-200 bg-slate-50 px-3 py-2 dark:border-slate-800 dark:bg-slate-950/60">
              <dl class="space-y-2 text-sm">
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.preview.amount') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(normalizedAmount) }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.preview.target') }}</dt>
                  <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ targetPreview }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.preview.type') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">{{ inviteTypeLabel }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.preview.expireAt') }}</dt>
                  <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ expirePreview }}</dd>
                </div>
              </dl>
            </div>

            <p v-if="errorMessage" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ errorMessage }}</p>

            <div class="flex justify-end">
              <UButton type="submit" color="primary" icon="i-lucide-ticket-plus" :loading="saving" :disabled="!canSubmit">
                {{ $t('admin.iam.invites.form.submit') }}
              </UButton>
            </div>
          </form>
        </section>

        <aside>
          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
              <div>
                <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.invites.siteList.title') }}</h2>
                <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.siteList.description') }}</p>
              </div>
              <select v-model="siteInviteStatusFilter" class="h-9 rounded-md border border-slate-200 bg-white px-2 text-sm text-slate-700 outline-none dark:border-slate-700 dark:bg-slate-950 dark:text-slate-200">
                <option v-for="option in siteInviteStatusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
              </select>
            </div>

            <div v-if="siteInvitesPending && siteInvites.length === 0" class="space-y-2 p-4">
              <div v-for="item in 5" :key="item" class="h-14 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>
            <div v-else-if="siteInvitesError" class="m-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
              {{ siteInvitesError }}
            </div>
            <div v-else-if="siteInvites.length === 0" class="px-4 py-12 text-center text-sm text-slate-500 dark:text-slate-400">
              {{ $t('admin.iam.invites.siteList.empty') }}
            </div>
            <div v-else class="divide-y divide-slate-100 dark:divide-slate-800">
              <div v-for="invite in siteInvites" :key="invite.id" class="grid gap-2 px-4 py-3 text-sm sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center">
                <div class="min-w-0">
                  <div class="flex min-w-0 items-center gap-2">
                    <UTooltip :text="invite.hash" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="120">
                      <span class="min-w-0 truncate font-mono text-xs font-medium text-slate-800 dark:text-slate-200">
                        {{ formatInviteHash(invite.hash) }}
                      </span>
                    </UTooltip>
                    <UBadge :color="inviteStatusColor(invite.status)" variant="soft">
                      {{ inviteStatusLabel(invite.status) }}
                    </UBadge>
                  </div>
                  <div class="mt-1 flex flex-wrap gap-x-3 gap-y-1 text-xs text-slate-500 dark:text-slate-400">
                    <span>{{ $t('admin.iam.invites.siteList.expireAt') }} {{ inviteValidUntilText(invite) }}</span>
                    <span>{{ $t('admin.iam.invites.siteList.createdAt') }} {{ formatDateTime(invite.createdAt, locale) }}</span>
                    <span v-if="invite.inviteeEmail || invite.inviteeName">{{ $t('admin.iam.invites.siteList.usedBy') }} {{ invite.inviteeEmail || invite.inviteeName }}</span>
                  </div>
                </div>
                <div class="flex items-center justify-end gap-1">
                  <UTooltip :text="$t('common.copy')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="120">
                    <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-copy" :disabled="!canCopySiteInvite(invite)" :aria-label="$t('common.copy')" @click="copySiteInviteHash(invite)" />
                  </UTooltip>
                  <UPopover :content="{ side: 'top', align: 'end', sideOffset: 8 }" :ui="{ content: 'w-72 p-3' }">
                    <UButton color="error" variant="ghost" size="xs" icon="i-lucide-archive" :disabled="!canRecycleSiteInvite(invite) || recyclingInviteId !== null" :loading="recyclingInviteId === invite.id" :aria-label="$t('admin.iam.invites.siteList.recycle')" />
                    <template #content="{ close }">
                      <div class="space-y-3">
                        <p class="text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.iam.invites.siteList.confirmRecycleTitle') }}</p>
                        <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.siteList.confirmRecycleDescription') }}</p>
                        <div class="flex justify-end gap-2">
                          <UButton color="neutral" variant="ghost" size="xs" type="button" @click="close()">{{ $t('common.cancel') }}</UButton>
                          <UButton color="error" size="xs" type="button" icon="i-lucide-archive" :loading="recyclingInviteId === invite.id" :disabled="recyclingInviteId !== null" @click="recycleSiteInvite(invite, close)">
                            {{ $t('admin.iam.invites.siteList.recycle') }}
                          </UButton>
                        </div>
                      </div>
                    </template>
                  </UPopover>
                </div>
              </div>
            </div>

            <AppPager
              v-if="siteInviteTotal > 0"
              class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
              size="sm"
              :page="siteInvitePage"
              :total="siteInviteTotal"
              :page-size="siteInviteSize"
              :disabled="siteInvitesPending"
              @page-change="goToSiteInvitePage"
            />
          </section>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminIamInvite, AdminIamRole } from '~/composables/useAdmin'
import { IamInviteStatusExpired, IamInviteStatusRecycled, IamInviteStatusSent, IamInviteStatusUnused, IamInviteStatusUsed } from '~/composables/useInvites'
import { formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

type InviteTargetMode = 'site' | 'roles'

const saving = ref(false)
const roles = ref<AdminIamRole[]>([])
const rolesPending = ref(true)
const errorMessage = ref('')
const targetMode = ref<InviteTargetMode>('site')
const siteInvites = ref<AdminIamInvite[]>([])
const siteInviteTotal = ref(0)
const siteInvitePage = ref(1)
const siteInviteSize = 8
const siteInvitesPending = ref(true)
const siteInvitesError = ref('')
const siteInviteStatusFilter = ref('')
const recyclingInviteId = ref<number | null>(null)
const form = reactive({
  amount: 1,
  roleIds: [] as number[],
  isTemp: false,
  expireAt: ''
})
const { roleOptions, roleNameWithLevel } = useAdminIamRoleLevels(roles)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const normalizedAmount = computed(() => Math.max(0, Number(form.amount || 0)))
const hasValidTarget = computed(() => targetMode.value === 'site' || form.roleIds.length > 0)
const canSubmit = computed(() => normalizedAmount.value > 0 && hasValidTarget.value && (!form.isTemp || Boolean(form.expireAt)) && !saving.value)
const inviteTypeLabel = computed(() => form.isTemp ? t('admin.iam.invites.preview.temporary') : t('admin.iam.invites.preview.permanent'))
const selectedRoleLabels = computed(() => {
  const selected = new Set(form.roleIds.map((id) => Number(id)))
  return roleOptions.value.filter((role) => selected.has(role.id)).map((role) => roleNameWithLevel(role))
})
const targetPreview = computed(() => {
  if (targetMode.value === 'site') return t('admin.iam.invites.preview.site')
  return selectedRoleLabels.value.length > 0 ? selectedRoleLabels.value.join(' / ') : t('admin.iam.invites.preview.rolesRequired')
})
const expireHint = computed(() => form.isTemp ? t('admin.iam.invites.form.expireAtRequired') : t('admin.iam.invites.form.expireAtDisabled'))
const expirePreview = computed(() => {
  if (!form.isTemp) return t('admin.iam.invites.preview.never')
  return form.expireAt ? new Date(form.expireAt).toLocaleString(locale.value) : t('admin.iam.invites.preview.notSet')
})
const siteInviteTotalPages = computed(() => Math.max(1, Math.ceil(siteInviteTotal.value / siteInviteSize)))
const siteInviteStatusOptions = computed(() => [
  { label: t('admin.iam.invites.siteList.all'), value: '' },
  { label: inviteStatusLabel(IamInviteStatusUnused), value: String(IamInviteStatusUnused) },
  { label: inviteStatusLabel(IamInviteStatusSent), value: String(IamInviteStatusSent) },
  { label: inviteStatusLabel(IamInviteStatusUsed), value: String(IamInviteStatusUsed) },
  { label: inviteStatusLabel(IamInviteStatusExpired), value: String(IamInviteStatusExpired) },
  { label: inviteStatusLabel(IamInviteStatusRecycled), value: String(IamInviteStatusRecycled) }
])

useHead({ title: t('admin.iam.invites.title') })

onMounted(() => {
  void loadRoles()
  void loadSiteInvites()
})

watch(siteInviteStatusFilter, () => {
  siteInvitePage.value = 1
  void loadSiteInvites()
})

async function loadRoles() {
  rolesPending.value = true
  try {
    const data = await adminApi.listIamRoles()
    roles.value = data.roles || []
  } catch {
    roles.value = []
  } finally {
    rolesPending.value = false
  }
}

async function loadSiteInvites() {
  siteInvitesPending.value = true
  siteInvitesError.value = ''
  try {
    let data = await adminApi.listIamInvites({
      page: siteInvitePage.value,
      size: siteInviteSize,
      status: siteInviteStatusFilter.value === '' ? undefined : Number(siteInviteStatusFilter.value)
    })
    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / siteInviteSize))
    if (siteInvitePage.value > nextTotalPages) {
      siteInvitePage.value = nextTotalPages
      data = await adminApi.listIamInvites({
        page: siteInvitePage.value,
        size: siteInviteSize,
        status: siteInviteStatusFilter.value === '' ? undefined : Number(siteInviteStatusFilter.value)
      })
    }
    siteInvites.value = data.invites || []
    siteInviteTotal.value = data.total || 0
  } catch (error) {
    siteInvites.value = []
    siteInviteTotal.value = 0
    siteInvitesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    siteInvitesPending.value = false
  }
}

function setTargetMode(mode: InviteTargetMode) {
  targetMode.value = mode
  if (mode === 'site') {
    form.roleIds = []
  }
}

async function grantInvites() {
  if (!canSubmit.value) return
  const generatedSiteInvites = targetMode.value === 'site'
  saving.value = true
  errorMessage.value = ''
  try {
    await adminApi.grantIamInvites({
      amount: Number(form.amount || 0),
      targetMode: targetMode.value,
      roleIds: targetMode.value === 'roles' ? [...form.roleIds] : [],
      isTemp: form.isTemp,
      expireAt: form.isTemp && form.expireAt ? new Date(form.expireAt).toISOString() : null
    })
    toast.add({ title: t('admin.iam.invites.form.success'), color: 'success', icon: 'i-lucide-check-circle' })
    form.amount = 1
    targetMode.value = 'site'
    form.roleIds = []
    form.isTemp = false
    form.expireAt = ''
    if (generatedSiteInvites) {
      siteInvitePage.value = 1
      await loadSiteInvites()
    }
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

function goToSiteInvitePage(page: number) {
  siteInvitePage.value = Math.min(Math.max(1, page), siteInviteTotalPages.value)
  void loadSiteInvites()
}

async function copySiteInviteHash(invite: AdminIamInvite) {
  if (!canCopySiteInvite(invite) || typeof navigator === 'undefined' || !navigator.clipboard) return
  await navigator.clipboard.writeText(invite.hash)
  toast.add({ title: t('user.invites.copied'), color: 'success', icon: 'i-lucide-check-circle' })
}

async function recycleSiteInvite(invite: AdminIamInvite, close: () => void) {
  if (!canRecycleSiteInvite(invite)) return
  recyclingInviteId.value = invite.id
  try {
    await adminApi.recycleIamInvite(invite.id)
    close()
    toast.add({ title: t('admin.iam.invites.siteList.recycled'), color: 'success', icon: 'i-lucide-check-circle' })
    await loadSiteInvites()
  } catch (error) {
    toast.add({ title: error instanceof ApiError ? error.message : t('common.requestFailed'), color: 'error', icon: 'i-lucide-circle-alert' })
  } finally {
    recyclingInviteId.value = null
  }
}

function canCopySiteInvite(invite: AdminIamInvite) {
  return invite.status === IamInviteStatusUnused && !isInviteExpired(invite)
}

function canRecycleSiteInvite(invite: AdminIamInvite) {
  return invite.status === IamInviteStatusUnused || invite.status === IamInviteStatusSent || invite.status === IamInviteStatusExpired
}

function isInviteExpired(invite: AdminIamInvite) {
  return Boolean(invite.expireAt && new Date(invite.expireAt).getTime() <= Date.now())
}

function formatInviteHash(hash: string) {
  if (!hash) return '-'
  if (hash.length <= 16) return hash
  return `${hash.slice(0, 8)}...${hash.slice(-8)}`
}

function inviteValidUntilText(invite: AdminIamInvite) {
  return invite.expireAt ? formatDateTime(invite.expireAt, locale.value) : t('admin.iam.invites.preview.never')
}

function inviteStatusLabel(status: number) {
  const key = `user.invites.status.${status}`
  const translated = t(key)
  return translated === key ? String(status) : translated
}

function inviteStatusColor(status: number) {
  if (status === IamInviteStatusUnused) return 'success'
  if (status === IamInviteStatusSent) return 'primary'
  if (status === IamInviteStatusUsed) return 'neutral'
  if (status === IamInviteStatusExpired) return 'error'
  return 'warning'
}
</script>
