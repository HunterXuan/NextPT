<template>
  <div class="min-h-[calc(100vh-4rem)] bg-neutral-50 py-6 dark:bg-neutral-950">
    <div class="grid w-full gap-4 px-3 sm:px-4 lg:px-5 xl:grid-cols-[minmax(0,1fr)_minmax(360px,480px)] xl:items-start">
      <section class="overflow-hidden rounded-lg border border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-900">
        <header class="flex flex-wrap items-center justify-between gap-3 border-b border-neutral-200 px-4 py-3 dark:border-neutral-800">
          <h2 class="text-sm font-semibold">{{ t('trackerAdmin.title') }}</h2>
          <UInput v-model="search" icon="i-lucide-search" :placeholder="t('trackerAdmin.search')" :aria-label="t('trackerAdmin.search')" class="w-full sm:w-56" />
        </header>
        <p v-if="pending" class="px-4 py-12 text-center text-sm text-muted">{{ t('trackerAdmin.loading') }}</p>
        <div v-else-if="loadError" class="space-y-3 px-4 py-8 text-center">
          <p class="text-sm text-error">{{ loadError }}</p>
          <UButton color="neutral" variant="soft" icon="i-lucide-refresh-cw" @click="load">{{ t('trackerAdmin.retry') }}</UButton>
        </div>
        <p v-else-if="!visibleRules.length" class="px-4 py-12 text-center text-sm text-muted">{{ t('trackerAdmin.empty') }}</p>
        <div v-else class="overflow-x-auto">
          <table class="w-full min-w-[440px] table-fixed text-left text-sm">
            <thead class="border-b border-neutral-200 bg-neutral-50 text-xs text-muted dark:border-neutral-800 dark:bg-neutral-950">
              <tr><th class="w-[48%] px-4 py-3">{{ t('trackerAdmin.family') }}</th><th class="w-[30%] px-4 py-3">{{ t('trackerAdmin.peerPattern') }}</th><th class="w-[22%] px-4 py-3 text-center">{{ t('trackerAdmin.status') }}</th></tr>
            </thead>
            <tbody>
              <tr v-for="rule in visibleRules" :key="rule.id" class="cursor-pointer border-b border-neutral-200 last:border-0 hover:bg-neutral-50 dark:border-neutral-800 dark:hover:bg-neutral-950" :class="{ 'app-selected': selectedId === rule.id }" role="button" tabindex="0" :aria-selected="selectedId === rule.id" @click="select(rule)" @keydown.enter.prevent="select(rule)" @keydown.space.prevent="select(rule)">
                <td class="px-4 py-3"><span class="block truncate font-medium hover:text-primary">{{ rule.family }}</span></td>
                <td class="px-4 py-3"><span class="block truncate font-mono text-muted">{{ rule.peerIdPattern || '-' }}</span></td>
                <td class="px-4 py-3 text-center"><UBadge :color="rule.enabled ? 'success' : 'neutral'" variant="soft">{{ t(rule.enabled ? 'trackerAdmin.enabled' : 'trackerAdmin.disabled') }}</UBadge></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="rounded-lg border border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-900">
        <header class="flex items-center justify-between gap-3 border-b border-neutral-200 px-4 py-3 dark:border-neutral-800">
          <h2 class="text-sm font-semibold">{{ t(selectedId ? 'trackerAdmin.edit' : 'trackerAdmin.create') }}</h2>
          <UTooltip :text="t('trackerAdmin.create')"><UButton icon="i-lucide-plus" color="neutral" variant="soft" :aria-label="t('trackerAdmin.create')" :disabled="saving" @click="startCreate" /></UTooltip>
        </header>
        <form class="space-y-4 p-4" @submit.prevent="save">
          <UFormField :label="t('trackerAdmin.family')" required><UInput v-model="form.family" class="w-full" size="lg" :maxlength="60" required :disabled="saving" /></UFormField>
          <UFormField :label="t('trackerAdmin.peerPattern')"><UInput v-model="form.peerIdPattern" class="w-full" size="lg" :maxlength="200" placeholder="^-qB" :disabled="saving" /></UFormField>
          <UFormField :label="t('trackerAdmin.pattern')"><UInput v-model="form.agentPattern" class="w-full" size="lg" :maxlength="200" placeholder="^qBittorrent/" :disabled="saving" /></UFormField>
          <UFormField :label="t('trackerAdmin.comment')"><UTextarea v-model="form.comment" class="w-full" :maxlength="200" :rows="3" :disabled="saving" /></UFormField>
          <USwitch v-model="form.enabled" :label="t('trackerAdmin.enabled')" :disabled="saving" />
          <AdminTrackerClientExceptionsEditor v-model="form.exceptions" :disabled="saving" />

          <details class="border-t border-neutral-200 pt-4 dark:border-neutral-800">
            <summary class="flex cursor-pointer list-none items-center gap-2 text-sm font-medium"><UIcon name="i-lucide-flask-conical" class="size-4" />{{ t('trackerAdmin.test') }}<UIcon name="i-lucide-chevron-down" class="ml-auto size-4 text-muted" /></summary>
            <div class="mt-4 space-y-3">
              <UFormField label="Peer ID"><UInput v-model="sample.peerId" class="w-full" :maxlength="200" placeholder="-qB4500-abcdefghijkl" /></UFormField>
              <UFormField label="User-Agent"><UInput v-model="sample.userAgent" class="w-full" :maxlength="500" placeholder="qBittorrent/4.5.0" /></UFormField>
              <div class="flex flex-wrap items-center gap-3"><UButton type="button" color="neutral" variant="soft" icon="i-lucide-play" :loading="testing" @click="check">{{ t('trackerAdmin.runTest') }}</UButton><UBadge v-if="testResult !== null" :color="testResult.matched ? 'success' : 'error'" variant="soft">{{ t(testResult.matched ? 'trackerAdmin.matched' : 'trackerAdmin.notMatched') }}</UBadge></div>
              <p v-if="testResult?.blockedBy" class="text-sm text-error">{{ t('trackerAdmin.blockedBy', { name: testResult.blockedBy }) }}</p>
            </div>
          </details>
          <footer class="flex items-center justify-between gap-3 border-t border-neutral-200 pt-4 dark:border-neutral-800">
            <UButton v-if="selectedId" type="button" color="error" variant="soft" icon="i-lucide-trash-2" :disabled="saving" @click="() => { deleteOpen = true }">{{ t('trackerAdmin.delete') }}</UButton><span v-else />
            <UButton type="submit" color="primary" variant="soft" icon="i-lucide-save" :loading="saving" :disabled="!form.family.trim() || (!form.peerIdPattern.trim() && !form.agentPattern.trim())">{{ t('trackerAdmin.save') }}</UButton>
          </footer>
        </form>
      </section>
    </div>
    <UModal v-model:open="deleteOpen" :title="t('trackerAdmin.delete')" :ui="{ content: 'max-w-sm' }">
      <template #body><p class="text-sm">{{ t('trackerAdmin.deleteConfirm', { name: form.family }) }}</p></template>
      <template #footer><div class="flex w-full justify-end gap-2"><UButton color="neutral" variant="soft" :disabled="saving" @click="() => { deleteOpen = false }">{{ t('trackerAdmin.cancel') }}</UButton><UButton color="error" variant="soft" :loading="saving" @click="remove">{{ t('trackerAdmin.delete') }}</UButton></div></template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { AdminTrackerClientWhitelist, AdminTrackerClientWhitelistInput } from '~/composables/useAdmin'
import { ApiError } from '~/composables/useApi'

definePageMeta({ layout: 'admin', middleware: ['auth', 'admin'] })
const { t } = useI18n()
const api = useAdmin()
const toast = useToast()
const rules = ref<AdminTrackerClientWhitelist[]>([])
const selectedId = ref<number | null>(null)
const search = ref('')
const pending = ref(false)
const saving = ref(false)
const testing = ref(false)
const loadError = ref('')
const deleteOpen = ref(false)
const testResult = ref<{ matched: boolean; blockedBy: string } | null>(null)
const form = reactive<AdminTrackerClientWhitelistInput>({ family: '', peerIdPattern: '', agentPattern: '', exceptions: [], enabled: true, comment: '' })
const sample = reactive({ peerId: '', userAgent: '' })
const visibleRules = computed(() => rules.value.filter(rule => `${rule.family} ${rule.peerIdPattern} ${rule.agentPattern}`.toLowerCase().includes(search.value.trim().toLowerCase())))
watch(() => [form.peerIdPattern, form.agentPattern, form.exceptions, sample.peerId, sample.userAgent], () => { testResult.value = null }, { deep: true })
onMounted(async () => {
  await load()
  if (rules.value[0]) select(rules.value[0])
})
useHead(() => ({ title: t('trackerAdmin.title') }))

async function load() {
  pending.value = true
  loadError.value = ''
  try { rules.value = (await api.listTrackerClientWhitelists()).list || [] }
  catch (error) { loadError.value = error instanceof ApiError ? error.message : t('common.requestFailed') }
  finally { pending.value = false }
}
function select(rule: AdminTrackerClientWhitelist) {
  if (saving.value) return
  selectedId.value = rule.id
  Object.assign(form, { family: rule.family, peerIdPattern: rule.peerIdPattern, agentPattern: rule.agentPattern, exceptions: (rule.exceptions || []).map(exception => ({ ...exception })), enabled: rule.enabled, comment: rule.comment })
  testResult.value = null
}
function startCreate() {
  selectedId.value = null
  Object.assign(form, { family: '', peerIdPattern: '', agentPattern: '', exceptions: [], enabled: true, comment: '' })
  testResult.value = null
}
function failed(error: unknown) { toast.add({ title: error instanceof ApiError ? error.message : t('common.requestFailed'), color: 'error' }) }
async function save() {
  if (saving.value) return
  saving.value = true
  try {
    const input = { ...form, exceptions: form.exceptions.map(exception => ({ ...exception })) }
    if (selectedId.value) await api.updateTrackerClientWhitelist(selectedId.value, input)
    else selectedId.value = (await api.createTrackerClientWhitelist(input)).id
    toast.add({ title: t('trackerAdmin.saved'), color: 'success' })
    await load()
    const saved = rules.value.find(rule => rule.id === selectedId.value)
    if (saved) Object.assign(form, { ...saved, exceptions: (saved.exceptions || []).map(exception => ({ ...exception })) })
  } catch (error) { failed(error) }
  finally { saving.value = false }
}
async function remove() {
  if (!selectedId.value || saving.value) return
  saving.value = true
  try {
    await api.deleteTrackerClientWhitelist(selectedId.value)
    deleteOpen.value = false
    startCreate()
    toast.add({ title: t('trackerAdmin.deleted'), color: 'success' })
    await load()
  } catch (error) { failed(error) }
  finally { saving.value = false }
}
async function check() {
  if (testing.value) return
  testing.value = true
  testResult.value = null
  const input = { peerIdPattern: form.peerIdPattern, agentPattern: form.agentPattern, exceptions: form.exceptions.map(exception => ({ ...exception })), ...sample }
  try {
    const result = await api.checkTrackerClientWhitelist(input)
    if (JSON.stringify(input) === JSON.stringify({ peerIdPattern: form.peerIdPattern, agentPattern: form.agentPattern, exceptions: form.exceptions, ...sample })) testResult.value = result
  } catch (error) { failed(error) }
  finally { testing.value = false }
}
</script>
