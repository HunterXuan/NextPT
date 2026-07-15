<template>
  <aside class="app-sticky-offset space-y-6 lg:sticky">
    <UCard class="rounded-lg">
      <template #header>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ t('forum.detail.info.title') }}</h2>
      </template>

      <dl class="space-y-3 text-sm">
        <div class="flex items-center justify-between gap-3">
          <dt class="text-slate-500 dark:text-slate-400">{{ t('forum.detail.info.topic') }}</dt>
          <dd class="font-medium text-slate-950 dark:text-white">#{{ topic.id }}</dd>
        </div>
        <div class="flex items-center justify-between gap-3">
          <dt class="text-slate-500 dark:text-slate-400">{{ t('forum.detail.info.views') }}</dt>
          <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(topic.views) }}</dd>
        </div>
        <div class="flex items-center justify-between gap-3">
          <dt class="text-slate-500 dark:text-slate-400">{{ t('forum.detail.info.replies') }}</dt>
          <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(topic.replyCount) }}</dd>
        </div>
        <div class="flex items-center justify-between gap-3">
          <dt class="text-slate-500 dark:text-slate-400">{{ t('forum.detail.info.status') }}</dt>
          <dd class="font-medium text-slate-950 dark:text-white">{{ topicStatusText }}</dd>
        </div>
      </dl>
    </UCard>

    <RewardPanel
      :title="t('forum.detail.reward.title')"
      :source-key="topic.id"
      success-key="forum.detail.reward.success"
      :load-rewards="loadRewards"
      :submit-reward="submitReward"
    />

    <UCard v-if="canManageTopic" class="rounded-lg">
      <template #header>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ t('forum.detail.admin.title') }}</h2>
      </template>

      <div class="grid grid-cols-2 gap-2">
        <UButton
          color="neutral"
          variant="outline"
          class="justify-center"
          :icon="topic.isLocked ? 'i-lucide-lock-open' : 'i-lucide-lock'"
          :loading="adminActionPending === 'lock'"
          @click="$emit('adminAction', topic.isLocked ? 'unlock' : 'lock')"
        >
          {{ topic.isLocked ? t('forum.detail.admin.unlock') : t('forum.detail.admin.lock') }}
        </UButton>
        <UButton
          color="neutral"
          variant="outline"
          class="justify-center"
          :icon="topic.isSticky ? 'i-lucide-pin-off' : 'i-lucide-pin'"
          :loading="adminActionPending === 'pin'"
          @click="$emit('adminAction', topic.isSticky ? 'unpin' : 'pin')"
        >
          {{ topic.isSticky ? t('forum.detail.admin.unpin') : t('forum.detail.admin.pin') }}
        </UButton>
      </div>

      <form class="mt-4 grid gap-2 border-t border-slate-200 pt-4 dark:border-slate-800" @submit.prevent="$emit('move')">
        <UFormField :label="t('forum.detail.admin.moveTo')">
          <USelect
            :model-value="moveNodeIdValue"
            class="w-full"
            size="lg"
            :ui="{ base: 'h-10 w-full' }"
            :items="adminNodeOptions"
            value-key="value"
            :disabled="adminNodesPending || adminActionPending === 'move'"
            @update:model-value="moveNodeIdValue = Number($event || 0)"
          />
        </UFormField>
        <UButton type="submit" color="primary" variant="soft" icon="i-lucide-move-right" :loading="adminActionPending === 'move'" :disabled="!moveNodeIdValue || moveNodeIdValue === topic.nodeId">
          {{ t('forum.detail.admin.move') }}
        </UButton>
        <p v-if="adminError" class="text-sm text-red-600 dark:text-red-300">{{ adminError }}</p>
      </form>

      <div class="mt-4 border-t border-slate-200 pt-4 dark:border-slate-800">
        <UPopover
          :open="deleteConfirmOpen"
          :content="{ side: 'top', align: 'center', sideOffset: 8 }"
          :ui="{ content: 'w-72 p-3' }"
          @update:open="setDeleteConfirmOpen"
        >
          <UButton color="error" variant="soft" icon="i-lucide-trash-2" block :loading="adminActionPending === 'delete'" :disabled="adminActionPending === 'delete'">
            {{ t('forum.detail.admin.delete') }}
          </UButton>

          <template #content="{ close }">
            <div class="space-y-3">
              <p class="text-sm font-medium text-slate-950 dark:text-white">
                {{ t('forum.detail.admin.confirmDelete') }}
              </p>
              <div class="flex justify-end gap-2">
                <UButton color="neutral" variant="ghost" size="xs" type="button" @click="closeDeletePopover(close)">
                  {{ t('common.cancel') }}
                </UButton>
                <UButton color="error" size="xs" type="button" icon="i-lucide-trash-2" :loading="adminActionPending === 'delete'" :disabled="adminActionPending === 'delete'" @click="confirmDelete(close)">
                  {{ t('forum.detail.admin.delete') }}
                </UButton>
              </div>
            </div>
          </template>
        </UPopover>
      </div>
    </UCard>
  </aside>
</template>

<script setup lang="ts">
import type { AdminForumCategory, AdminForumNode } from '~/composables/useAdmin'
import type { ForumRewardListOut, ForumTopicDetail } from '~/composables/useForum'
import { localizeI18nName } from '~/utils/format'

type TopicAdminAction = 'lock' | 'unlock' | 'pin' | 'unpin'
type RewardLoader = (page: number, size: number) => Promise<ForumRewardListOut>
type RewardSubmitter = (amount: number) => Promise<void>

const props = withDefaults(defineProps<{
  topic: ForumTopicDetail
  canManageTopic?: boolean
  moveNodeId: number
  adminNodes: AdminForumNode[]
  adminCategories: AdminForumCategory[]
  adminNodesPending?: boolean
  adminActionPending?: string
  adminError?: string
  loadRewards: RewardLoader
  submitReward: RewardSubmitter
}>(), {
  canManageTopic: false,
  adminNodesPending: false,
  adminActionPending: '',
  adminError: ''
})

const emit = defineEmits<{
  'update:moveNodeId': [value: number]
  adminAction: [action: TopicAdminAction]
  move: []
  delete: []
}>()

const { t, locale } = useI18n()

const deleteConfirmOpen = ref(false)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const moveNodeIdValue = computed({
  get: () => props.moveNodeId,
  set: (value) => emit('update:moveNodeId', value)
})
const topicStatusText = computed(() => {
  if (props.topic.isLocked) return t('forum.detail.info.statusLocked')
  if (props.topic.isSticky) return t('forum.detail.info.statusSticky')
  return t('forum.detail.info.statusNormal')
})
const adminNodeGroups = computed(() => {
  const nodesByCategory = new Map<number, AdminForumNode[]>()
  for (const node of props.adminNodes) {
    const nodes = nodesByCategory.get(node.categoryId) || []
    nodes.push(node)
    nodesByCategory.set(node.categoryId, nodes)
  }

  const groups = props.adminCategories
    .map((category) => ({
      key: `category:${category.id}`,
      name: adminCategoryName(category),
      nodes: nodesByCategory.get(category.id) || []
    }))
    .filter((group) => group.nodes.length > 0)

  const groupedNodeIds = new Set(groups.flatMap((group) => group.nodes.map((node) => node.id)))
  const uncategorizedNodes = props.adminNodes.filter((node) => !groupedNodeIds.has(node.id))
  if (uncategorizedNodes.length > 0) {
    groups.push({
      key: 'category:unknown',
      name: t('forum.fallback.category'),
      nodes: uncategorizedNodes
    })
  }

  return groups
})
const adminNodeOptions = computed(() => adminNodeGroups.value.flatMap(group => group.nodes.map(node => ({
  value: node.id,
  label: `${group.name} / ${forumNodeName(node)}`
}))))

function forumNodeName(node: AdminForumNode) {
  return localizeI18nName(node.nameI18N as any, locale.value, node.slug || `#${node.id}`)
}

function adminCategoryName(category: AdminForumCategory) {
  return localizeI18nName(category.nameI18N as any, locale.value, t('forum.fallback.category'))
}

function setDeleteConfirmOpen(value: boolean) {
  if (props.adminActionPending === 'delete') return
  deleteConfirmOpen.value = value
}

function closeDeletePopover(close?: () => void) {
  deleteConfirmOpen.value = false
  close?.()
}

function confirmDelete(close?: () => void) {
  emit('delete')
  closeDeletePopover(close)
}
</script>
