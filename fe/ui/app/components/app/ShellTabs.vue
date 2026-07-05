<template>
  <div class="border-t border-slate-200 bg-slate-50/92 dark:border-slate-800 dark:bg-slate-950/92">
    <div class="flex h-10 items-center gap-1 px-2 sm:px-3 lg:px-5">
      <UTooltip
        :text="$t('workspaceTabs.scrollLeft')"
        :content="{ side: 'bottom', sideOffset: 8 }"
        :delay-duration="300"
      >
        <UButton
          class="hidden sm:inline-flex"
          color="neutral"
          variant="ghost"
          size="xs"
          icon="i-lucide-chevron-left"
          :aria-label="$t('workspaceTabs.scrollLeft')"
          @click="scrollTabs(-260)"
        />
      </UTooltip>

      <div
        ref="scroller"
        class="min-w-0 flex-1 overflow-x-auto [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
        :aria-label="$t('workspaceTabs.title')"
      >
        <div class="flex min-w-max items-center py-1">
          <template
            v-for="(tab, index) in tabs"
            :key="tab.id"
          >
            <span
              v-if="index > 0"
              class="mx-1 h-5 w-px shrink-0 bg-slate-200/90 transition-opacity dark:bg-slate-800"
              :class="tab.id === activeId || tabs[index - 1]?.id === activeId ? 'opacity-0' : 'opacity-100'"
              aria-hidden="true"
            />
            <UContextMenu
              :items="tabContextMenuItems(tab)"
              :content="{ sideOffset: 6 }"
              :modal="false"
            >
              <div
                :ref="(el) => setTabElement(tab.id, el)"
                class="group flex h-8 max-w-[220px] shrink-0 items-center overflow-hidden rounded-md border transition"
                :class="tab.id === activeId
                  ? 'border-slate-300 bg-white text-slate-950 shadow-sm dark:border-slate-700 dark:bg-slate-900 dark:text-white'
                  : 'border-transparent bg-transparent text-slate-600 hover:border-slate-200 hover:bg-white/80 dark:text-slate-300 dark:hover:border-slate-800 dark:hover:bg-slate-900/80'"
                draggable="true"
                @dragstart="handleDragStart(index)"
                @dragover.prevent
                @drop="handleDrop(index)"
                @dragend="draggingIndex = -1"
              >
                <UTooltip
                  :text="tab.title"
                  :content="{ side: 'bottom', sideOffset: 8 }"
                  :delay-duration="600"
                >
                  <button
                    type="button"
                    class="flex min-w-0 flex-1 items-center gap-1.5 px-2 py-1 text-left text-xs"
                    :aria-current="tab.id === activeId ? 'page' : undefined"
                    @click="activateTab(tab)"
                  >
                    <UIcon :name="tabIcon(tab)" class="size-3.5 shrink-0 text-slate-400" />
                    <span class="truncate">{{ tab.title }}</span>
                  </button>
                </UTooltip>

                <UTooltip
                  :text="$t('workspaceTabs.close')"
                  :content="{ side: 'bottom', sideOffset: 8 }"
                  :delay-duration="300"
                >
                  <button
                    type="button"
                    draggable="false"
                    class="mr-1 flex size-5 shrink-0 items-center justify-center rounded text-slate-400 opacity-80 transition hover:bg-slate-100 hover:text-slate-700 group-hover:opacity-100 dark:hover:bg-slate-800 dark:hover:text-slate-200"
                    :aria-label="$t('workspaceTabs.close')"
                    @pointerdown.stop
                    @mousedown.stop
                    @dragstart.stop.prevent
                    @click.stop.prevent="closeTab(tab)"
                  >
                    <UIcon name="i-lucide-x" class="size-3.5" />
                  </button>
                </UTooltip>
              </div>
            </UContextMenu>
          </template>
        </div>
      </div>

      <UTooltip
        :text="$t('workspaceTabs.scrollRight')"
        :content="{ side: 'bottom', sideOffset: 8 }"
        :delay-duration="300"
      >
        <UButton
          class="hidden sm:inline-flex"
          color="neutral"
          variant="ghost"
          size="xs"
          icon="i-lucide-chevron-right"
          :aria-label="$t('workspaceTabs.scrollRight')"
          @click="scrollTabs(260)"
        />
      </UTooltip>

      <div class="h-5 w-px bg-slate-200 dark:bg-slate-800" />

      <UTooltip
        :text="$t('workspaceTabs.refresh')"
        :content="{ side: 'bottom', sideOffset: 8 }"
        :delay-duration="300"
      >
        <UButton
          color="neutral"
          variant="ghost"
          size="xs"
          icon="i-lucide-refresh-cw"
          :aria-label="$t('workspaceTabs.refresh')"
          :disabled="!activeTab || refreshing"
          :loading="refreshing"
          @click="refreshActiveTab"
        />
      </UTooltip>

      <UDropdownMenu :items="tabMenuItems" :content="{ align: 'end' }">
        <UButton
          color="neutral"
          variant="ghost"
          size="xs"
          icon="i-lucide-more-horizontal"
          :aria-label="$t('workspaceTabs.title')"
          :disabled="!activeTab"
        />
      </UDropdownMenu>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ContextMenuItem, DropdownMenuItem } from '@nuxt/ui'
import type { WorkspaceTabItem, WorkspaceTabMode } from '~/composables/useWorkspaceTabs'

const props = defineProps<{
  mode: WorkspaceTabMode
  currentTitle: string
}>()

const emit = defineEmits<{
  refresh: []
}>()

const route = useRoute()
const localePath = useLocalePath()
const { t } = useI18n()
const { tabs, activeId, loadTabs, upsertTab, removeTab, closeOtherTabs, closeRightTabs, closeAllTabs, moveTab, setActiveTab, normalizeTabPath } = useWorkspaceTabs(props.mode)

const scroller = ref<HTMLElement | null>(null)
const draggingIndex = ref(-1)
const suppressedRouteId = ref('')
const tabAutoAddPaused = ref(false)
const refreshing = ref(false)
const tabElements = new Map<string, HTMLElement>()

const activeTab = computed(() => tabs.value.find((tab) => tab.id === activeId.value) || null)
const fallbackPath = computed(() => localePath(props.mode === 'admin' ? '/admin' : '/'))
const fallbackTitle = computed(() => props.mode === 'admin' ? t('admin.nav.overview') : t('nav.home'))
const tabMenuItems = computed<DropdownMenuItem[][]>(() => [
  [
    {
      label: t('workspaceTabs.refresh'),
      icon: 'i-lucide-refresh-cw',
      disabled: refreshing.value,
      onSelect: () => {
        void refreshActiveTab()
      }
    }
  ],
  [
    {
      label: t('workspaceTabs.closeCurrent'),
      icon: 'i-lucide-x',
      onSelect: () => {
        if (activeTab.value) void closeTab(activeTab.value)
      }
    },
    {
      label: t('workspaceTabs.closeOther'),
      icon: 'i-lucide-copy-x',
      onSelect: () => {
        void closeOtherActiveTabs()
      }
    },
    {
      label: t('workspaceTabs.closeRight'),
      icon: 'i-lucide-panel-right-close',
      onSelect: () => {
        void closeRightActiveTabs()
      }
    },
    {
      label: t('workspaceTabs.closeAll'),
      icon: 'i-lucide-circle-x',
      onSelect: () => {
        void closeAllAndNavigate()
      }
    }
  ]
])

function tabContextMenuItems(tab: WorkspaceTabItem): ContextMenuItem[][] {
  return [
    [
      {
        label: t('workspaceTabs.refresh'),
        icon: 'i-lucide-refresh-cw',
        disabled: refreshing.value,
        onSelect: () => {
          void refreshTab(tab)
        }
      }
    ],
    [
      {
        label: t('workspaceTabs.close'),
        icon: 'i-lucide-x',
        onSelect: () => {
          void closeTab(tab)
        }
      },
      {
        label: t('workspaceTabs.closeOther'),
        icon: 'i-lucide-copy-x',
        onSelect: () => {
          void closeOtherTabsFor(tab)
        }
      },
      {
        label: t('workspaceTabs.closeRight'),
        icon: 'i-lucide-panel-right-close',
        onSelect: () => {
          void closeRightTabsFor(tab)
        }
      },
      {
        label: t('workspaceTabs.closeAll'),
        icon: 'i-lucide-circle-x',
        onSelect: () => {
          void closeAllAndNavigate()
        }
      }
    ]
  ]
}

onMounted(() => {
  loadTabs()
  touchCurrentTab()
})

watch(
  () => [route.fullPath, props.currentTitle],
  () => {
    if (!import.meta.client) return
    const routeId = normalizeTabPath(route.fullPath)
    if (suppressedRouteId.value && suppressedRouteId.value !== routeId) {
      suppressedRouteId.value = ''
    }
    touchCurrentTab()
  }
)

watch(
  activeId,
  () => {
    void nextTick(scrollActiveTabIntoView)
  },
  { flush: 'post' }
)

function touchCurrentTab() {
  if (tabAutoAddPaused.value) return
  if (suppressedRouteId.value === normalizeTabPath(route.fullPath)) return

  upsertTab(route.fullPath, props.currentTitle)
  void nextTick(scrollActiveTabIntoView)
}

async function activateTab(tab: WorkspaceTabItem) {
  setActiveTab(tab.to)
  if (normalizeTabPath(route.fullPath) === tab.id) return
  await navigateTo(tab.to)
}

async function closeTab(tab: WorkspaceTabItem) {
  const wasActive = tab.id === activeId.value
  const nextTab = removeTab(tab.id)

  if (wasActive) {
    const destination = nextTab?.to || fallbackPath.value
    if (nextTab) {
      suppressedRouteId.value = tab.id
      await navigateTo(destination)
      suppressedRouteId.value = ''
      return
    }

    tabAutoAddPaused.value = true
    try {
      await navigateTo(destination)
      upsertTab(destination, fallbackTitle.value)
      await nextTick()
      scrollActiveTabIntoView()
    } finally {
      tabAutoAddPaused.value = false
    }
    return
  }

  if (suppressedRouteId.value === tab.id) {
    suppressedRouteId.value = ''
  }
}

async function refreshActiveTab() {
  if (!activeTab.value) return
  await refreshTab(activeTab.value)
}

async function refreshTab(tab: WorkspaceTabItem) {
  if (!import.meta.client || refreshing.value) return

  refreshing.value = true
  try {
    if (route.fullPath !== tab.to || normalizeTabPath(route.fullPath) !== tab.id) {
      setActiveTab(tab.to)
      await navigateTo(tab.to)
    }
    await nextTick()
    emit('refresh')
  } finally {
    refreshing.value = false
  }
}

function closeOtherActiveTabs() {
  if (!activeTab.value) return
  closeOtherTabsFor(activeTab.value)
}

async function closeRightActiveTabs() {
  if (!activeTab.value) return
  await closeRightTabsFor(activeTab.value)
}

async function closeOtherTabsFor(tab: WorkspaceTabItem) {
  closeOtherTabs(tab.id)
  if (normalizeTabPath(route.fullPath) !== tab.id) {
    await navigateTo(tab.to)
  } else {
    setActiveTab(tab.to)
  }
}

async function closeRightTabsFor(tab: WorkspaceTabItem) {
  const targetIndex = tabs.value.findIndex((item) => item.id === tab.id)
  const activeIndex = tabs.value.findIndex((item) => item.id === activeId.value)
  const shouldActivateTarget = targetIndex >= 0 && activeIndex > targetIndex

  closeRightTabs(tab.id)
  if (shouldActivateTarget && normalizeTabPath(route.fullPath) !== tab.id) {
    await navigateTo(tab.to)
  } else if (normalizeTabPath(route.fullPath) === tab.id) {
    setActiveTab(tab.to)
  }
}

async function closeAllAndNavigate() {
  tabAutoAddPaused.value = true
  try {
    closeAllTabs()
    await navigateTo(fallbackPath.value)
    upsertTab(fallbackPath.value, fallbackTitle.value)
    await nextTick()
    scrollActiveTabIntoView()
  } finally {
    tabAutoAddPaused.value = false
  }
}

function scrollTabs(left: number) {
  scroller.value?.scrollBy({ left, behavior: 'smooth' })
}

function handleDragStart(index: number) {
  draggingIndex.value = index
}

function handleDrop(index: number) {
  if (draggingIndex.value < 0) return
  moveTab(draggingIndex.value, index)
  draggingIndex.value = -1
}

function setTabElement(id: string, element: Element | null) {
  if (element instanceof HTMLElement) {
    tabElements.set(id, element)
  } else {
    tabElements.delete(id)
  }
}

function scrollActiveTabIntoView() {
  const element = tabElements.get(activeId.value)
  element?.scrollIntoView({ block: 'nearest', inline: 'nearest', behavior: 'smooth' })
}

function tabIcon(tab: WorkspaceTabItem) {
  const to = tab.to
  if (to.includes('/admin')) return 'i-lucide-shield-check'
  if (to.includes('/catalog/torrents')) return 'i-lucide-library'
  if (to.includes('/catalog/subtitles')) return 'i-lucide-captions'
  if (to.includes('/forum')) return 'i-lucide-messages-square'
  if (to.includes('/iam/users')) return 'i-lucide-user-round'
  return 'i-lucide-house'
}
</script>
