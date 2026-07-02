<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <form class="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1fr)_300px] xl:items-start" @submit.prevent="handleSubmit">
        <main class="min-w-0 space-y-3">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <div v-if="nodesPending" class="grid gap-3 md:grid-cols-[minmax(0,240px)_minmax(0,1fr)]">
              <div class="h-10 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
              <div class="h-10 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>
            <template v-else>
              <div class="grid gap-3 md:grid-cols-[minmax(0,240px)_minmax(0,1fr)]">
                <UFormField :label="$t('forum.create.fields.category')">
                  <select
                    v-model.number="selectedCategoryId"
                    class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                    :disabled="pending || categories.length === 0"
                    @change="handleCategoryChange"
                  >
                    <option v-for="category in categories" :key="category.id" :value="category.id">
                      {{ categoryDisplayName(category) }}
                    </option>
                  </select>
                </UFormField>

                <UFormField :label="$t('forum.create.fields.node')" required>
                  <select
                    v-model="form.nodeId"
                    class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                    :disabled="pending || selectedCategoryNodes.length === 0"
                  >
                    <option value="0">{{ $t('forum.create.fields.nodePlaceholder') }}</option>
                    <option v-for="node in selectedCategoryNodes" :key="node.id" :value="String(node.id)">
                      {{ nodeDisplayName(node) }}
                    </option>
                  </select>
                </UFormField>
              </div>

              <div v-if="!selectedNode" class="mt-3 flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
                <UIcon name="i-lucide-circle-slash" class="size-4" />
                <span>{{ $t('forum.empty.nodes') }}</span>
              </div>
            </template>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <div class="grid grid-cols-1 gap-4">
              <UFormField :label="$t('forum.create.fields.subject')" required>
                <UInput v-model="form.subject" class="w-full" :disabled="pending" :placeholder="$t('forum.create.placeholders.subject')" />
              </UFormField>

              <UFormField :label="$t('forum.create.fields.content')" required>
                <UTextarea
                  v-if="editorMode === 'write'"
                  v-model="form.content"
                  class="w-full"
                  :rows="14"
                  :disabled="pending"
                  :placeholder="$t('forum.create.placeholders.content')"
                />
                <div v-else class="min-h-[22rem] rounded-md border border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950">
                  <div v-if="renderedContentPreview" class="rich-text" v-html="renderedContentPreview" />
                  <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ $t('forum.create.preview.empty') }}</p>
                </div>

                <div class="mt-3 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  <div class="inline-flex w-fit rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950">
                    <button type="button" :class="editorModeButtonClass('write')" @click="editorMode = 'write'">
                      {{ $t('forum.create.preview.write') }}
                    </button>
                    <button type="button" :class="editorModeButtonClass('preview')" @click="editorMode = 'preview'">
                      {{ $t('forum.create.preview.preview') }}
                    </button>
                  </div>
                  <span class="text-xs text-slate-500 dark:text-slate-400">
                    {{ $t('forum.create.summary.contentLength', { count: numberFormatter.format(form.content.trim().length) }) }}
                  </span>
                </div>
              </UFormField>
            </div>
          </section>
        </main>

        <aside class="space-y-3 xl:sticky xl:top-20">
          <section class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <UButton type="submit" color="primary" icon="i-lucide-square-pen" block :loading="pending" :disabled="!canSubmit">
              {{ $t('forum.create.submit') }}
            </UButton>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('forum.create.guide.title') }}</h2>
            <ul class="mt-3 space-y-2 text-sm leading-6 text-slate-600 dark:text-slate-300">
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('forum.create.guide.search') }}</span>
              </li>
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('forum.create.guide.titleTip') }}</span>
              </li>
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('forum.create.guide.kindness') }}</span>
              </li>
              <li class="flex gap-2 rounded-md bg-amber-50 px-2 py-1.5 text-amber-800 dark:bg-amber-950/40 dark:text-amber-200">
                <UIcon name="i-lucide-triangle-alert" class="mt-1 size-4 shrink-0" />
                <span>{{ $t('forum.create.guide.finality') }}</span>
              </li>
            </ul>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <dl class="space-y-2 text-sm">
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.create.summary.subjectLength') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(form.subject.trim().length) }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.create.summary.content') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(form.content.trim().length) }}</dd>
              </div>
            </dl>
          </section>

          <section v-if="nodesError" class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
            <p>{{ nodesError }}</p>
            <UButton type="button" class="mt-3" color="error" variant="soft" size="sm" icon="i-lucide-refresh-cw" :loading="nodesPending" @click="loadNodes">
              {{ $t('common.retry') }}
            </UButton>
          </section>
        </aside>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useForum, type ForumNode, type ForumNodeCategory } from '~/composables/useForum'
import { localizeI18nName } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

type EditorMode = 'write' | 'preview'

interface ForumNodeWithCategory {
  category: ForumNodeCategory
  node: ForumNode
}

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
const selectedCategoryId = ref(0)
const editorMode = ref<EditorMode>('write')
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))

const form = reactive({
  nodeId: String(readNodeIdQuery()),
  subject: '',
  content: ''
})

const flatNodeItems = computed<ForumNodeWithCategory[]>(() => {
  return categories.value.flatMap((category) => category.nodes.map((node) => ({ category, node })))
})
const selectedNodeItem = computed(() => {
  return flatNodeItems.value.find((item) => item.node.id === Number(form.nodeId)) || null
})
const selectedCategory = computed(() => {
  return categories.value.find((category) => category.id === selectedCategoryId.value)
    || selectedNodeItem.value?.category
    || categories.value[0]
    || null
})
const selectedCategoryNodes = computed(() => selectedCategory.value?.nodes || [])
const selectedNode = computed(() => selectedNodeItem.value?.node || null)
const renderedContentPreview = computed(() => renderUserMarkdown(form.content).trim())
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
    const data = await forum.listNodes({ scope: 'create' })
    categories.value = (data.list || []).filter((category) => category.nodes.length > 0)
    resolveSelectedNode()
  } catch (error) {
    categories.value = []
    selectedCategoryId.value = 0
    form.nodeId = '0'
    nodesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    nodesPending.value = false
  }
}

function resolveSelectedNode() {
  const matched = selectedNodeItem.value
  if (matched) {
    selectedCategoryId.value = matched.category.id
    return
  }

  const firstCategory = categories.value[0]
  selectedCategoryId.value = firstCategory?.id || 0
  form.nodeId = String(firstCategory?.nodes[0]?.id || 0)
}

function handleCategoryChange() {
  if (pending.value) return

  const category = categories.value.find((item) => item.id === selectedCategoryId.value)
  if (!category) {
    form.nodeId = '0'
    return
  }

  const currentNodeInCategory = category.nodes.some((node) => node.id === Number(form.nodeId))
  if (!currentNodeInCategory) {
    form.nodeId = String(category.nodes[0]?.id || 0)
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

function editorModeButtonClass(mode: EditorMode) {
  const active = editorMode.value === mode
  return [
    'h-7 rounded px-3 text-xs font-medium transition-colors',
    active
      ? 'bg-white text-slate-950 shadow-sm ring-1 ring-slate-200 dark:bg-slate-800 dark:text-white dark:ring-slate-700'
      : 'text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100'
  ].join(' ')
}
</script>
