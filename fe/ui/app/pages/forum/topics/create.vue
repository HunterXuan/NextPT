<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mb-6 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-slate-500 dark:text-slate-400">{{ $t('forum.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('forum.create.title') }}</h1>
        </div>
        <UButton color="neutral" variant="outline" icon="i-lucide-arrow-left" :to="localePath('/forum')">
          {{ $t('forum.create.back') }}
        </UButton>
      </div>

      <form class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_340px]" @submit.prevent="handleSubmit">
        <div class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.create.sections.topic') }}</h2>
            </template>

            <div class="grid grid-cols-1 gap-4">
              <UFormField :label="$t('forum.create.fields.node')" required>
                <select
                  v-model="form.nodeId"
                  class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="nodesPending || pending || categories.length === 0"
                >
                  <option value="0">{{ $t('forum.create.fields.nodePlaceholder') }}</option>
                  <optgroup v-for="category in categories" :key="category.id" :label="categoryDisplayName(category)">
                    <option v-for="node in category.nodes" :key="node.id" :value="String(node.id)">
                      {{ nodeDisplayName(node) }}
                    </option>
                  </optgroup>
                </select>
              </UFormField>

              <UFormField :label="$t('forum.create.fields.subject')" required>
                <UInput v-model="form.subject" class="w-full" :disabled="pending" :placeholder="$t('forum.create.placeholders.subject')" />
              </UFormField>

              <UFormField :label="$t('forum.create.fields.content')" required>
                <UTextarea
                  v-model="form.content"
                  class="w-full"
                  :rows="12"
                  :disabled="pending"
                  :placeholder="$t('forum.create.placeholders.content')"
                />
              </UFormField>
            </div>
          </UCard>
        </div>

        <aside class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.create.sections.publish') }}</h2>
            </template>

            <div class="space-y-4">
              <dl class="space-y-3 text-sm">
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.create.summary.node') }}</dt>
                  <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ selectedNodeName }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.create.summary.subjectLength') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(form.subject.trim().length) }}</dd>
                </div>
              </dl>

              <div v-if="nodesError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
                {{ nodesError }}
              </div>

              <UButton type="submit" color="primary" icon="i-lucide-square-pen" block :loading="pending" :disabled="!canSubmit">
                {{ $t('forum.create.submit') }}
              </UButton>
            </div>
          </UCard>
        </aside>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useForum, type ForumNode, type ForumNodeCategory } from '~/composables/useForum'
import { localizeI18nName } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const forum = useForum()

const categories = ref<ForumNodeCategory[]>([])
const nodesPending = ref(true)
const nodesError = ref('')
const pending = ref(false)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))

const form = reactive({
  nodeId: String(readNodeIdQuery()),
  subject: '',
  content: ''
})

const flatNodes = computed(() => categories.value.flatMap((category) => category.nodes))
const selectedNode = computed(() => flatNodes.value.find((node) => node.id === Number(form.nodeId)) || null)
const selectedNodeName = computed(() => selectedNode.value ? nodeDisplayName(selectedNode.value) : '-')
const canSubmit = computed(() => {
  return Number(form.nodeId) > 0 && form.subject.trim().length >= 2 && form.content.trim().length >= 2 && !pending.value
})

useHead(() => ({
  title: `${t('forum.create.metaTitle')} - NextPT`
}))

onMounted(loadNodes)

function readNodeIdQuery() {
  const raw = Array.isArray(route.query.nodeId) ? route.query.nodeId[0] : route.query.nodeId
  const parsed = Number(raw)
  return Number.isInteger(parsed) && parsed > 0 ? parsed : 0
}

async function loadNodes() {
  nodesPending.value = true
  nodesError.value = ''
  try {
    const data = await forum.listNodes()
    categories.value = (data.list || []).filter((category) => category.nodes.length > 0)
    if (Number(form.nodeId) > 0 && !flatNodes.value.some((node) => node.id === Number(form.nodeId))) {
      form.nodeId = '0'
    }
  } catch (error) {
    categories.value = []
    nodesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    nodesPending.value = false
  }
}

async function handleSubmit() {
  if (!canSubmit.value) return

  pending.value = true
  try {
    const out = await forum.createTopic({
      nodeId: Number(form.nodeId),
      subject: form.subject.trim(),
      content: form.content.trim()
    })

    toast.add({
      title: t('forum.create.success', { id: out.id }),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await navigateTo(localePath(`/forum/topics/${out.id}`))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    pending.value = false
  }
}

function categoryDisplayName(category: ForumNodeCategory) {
  return localizeI18nName(category.nameI18n, locale.value, t('forum.fallback.category'))
}

function nodeDisplayName(node: ForumNode) {
  return localizeI18nName(node.nameI18n, locale.value, node.slug || `#${node.id}`)
}
</script>
