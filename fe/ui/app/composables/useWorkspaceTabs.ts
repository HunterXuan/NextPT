export type WorkspaceTabMode = 'app' | 'admin'

export interface WorkspaceTabItem {
  id: string
  to: string
  title: string
  createdAt: number
  updatedAt: number
}

const maxWorkspaceTabs = 16

function normalizeWorkspacePathOnly(path: string) {
  const normalized = path.trim() || '/'
  if (normalized === '/') return normalized
  return normalized.replace(/\/+$/, '') || '/'
}

function normalizeWorkspaceTabPath(to: string) {
  const [withoutHash] = String(to || '').split('#')
  const [path] = String(withoutHash || '').split('?')
  return normalizeWorkspacePathOnly(path || '/')
}

function normalizeWorkspaceTabTo(to: string) {
  const value = String(to || '').trim() || '/'
  const hashIndex = value.indexOf('#')
  const beforeHash = hashIndex >= 0 ? value.slice(0, hashIndex) : value
  const hash = hashIndex >= 0 ? value.slice(hashIndex) : ''
  const queryIndex = beforeHash.indexOf('?')
  const path = queryIndex >= 0 ? beforeHash.slice(0, queryIndex) : beforeHash
  const query = queryIndex >= 0 ? beforeHash.slice(queryIndex) : ''

  return `${normalizeWorkspacePathOnly(path || '/')}${query}${hash}`
}

function storageKey(mode: WorkspaceTabMode) {
  return `nextpt_workspace_tabs_${mode}`
}

export function useWorkspaceTabs(mode: WorkspaceTabMode = 'app') {
  const tabs = useState<WorkspaceTabItem[]>(`workspace-tabs:${mode}`, () => [])
  const activeId = useState<string>(`workspace-tabs-active:${mode}`, () => '')
  const loaded = useState<boolean>(`workspace-tabs-loaded:${mode}`, () => false)

  function loadTabs() {
    if (!import.meta.client) return
    if (loaded.value) {
      normalizeCurrentTabs()
      return
    }

    try {
      const rawValue = localStorage.getItem(storageKey(mode))
      const parsed = rawValue ? JSON.parse(rawValue) : []
      tabs.value = Array.isArray(parsed)
        ? normalizeLoadedTabs(parsed.filter(isWorkspaceTabItem)).slice(0, maxWorkspaceTabs)
        : []
    } catch {
      tabs.value = []
    } finally {
      loaded.value = true
    }
  }

  function normalizeCurrentTabs() {
    tabs.value = normalizeLoadedTabs(tabs.value).slice(0, maxWorkspaceTabs)
    if (activeId.value) {
      activeId.value = normalizeWorkspaceTabPath(activeId.value)
    }
  }

  function saveTabs() {
    if (!import.meta.client || !loaded.value) return
    localStorage.setItem(storageKey(mode), JSON.stringify(tabs.value))
  }

  function upsertTab(to: string, title: string) {
    loadTabs()

    const normalizedTo = normalizeWorkspaceTabTo(to)
    const tabId = normalizeWorkspaceTabPath(normalizedTo)
    const now = Date.now()
    const nextTitle = title.trim() || tabId
    const existing = tabs.value.find((tab) => tab.id === tabId)

    if (existing) {
      tabs.value = tabs.value.map((tab) => tab.id === tabId
        ? { ...tab, to: normalizedTo, title: nextTitle, updatedAt: now }
        : tab
      )
    } else {
      tabs.value = [...tabs.value, {
        id: tabId,
        to: normalizedTo,
        title: nextTitle,
        createdAt: now,
        updatedAt: now
      }]
    }

    activeId.value = tabId
    trimTabs(tabId)
    saveTabs()
  }

  function updateTabTitle(to: string, title: string) {
    const normalizedTo = normalizeWorkspaceTabPath(to)
    const nextTitle = title.trim()
    if (!nextTitle) return

    loadTabs()
    const tab = tabs.value.find((item) => item.id === normalizedTo)
    if (!tab) return

    tabs.value = tabs.value.map((item) => item.id === normalizedTo
      ? { ...item, title: nextTitle, updatedAt: Date.now() }
      : item
    )
    saveTabs()
  }

  function removeTab(id: string) {
    loadTabs()

    const index = tabs.value.findIndex((tab) => tab.id === id)
    if (index < 0) return null

    const wasActive = activeId.value === id
    const nextTabs = tabs.value.filter((tab) => tab.id !== id)

    let nextTab: WorkspaceTabItem | null = null
    if (wasActive) {
      nextTab = nextTabs[index] || nextTabs[index - 1] || null
      activeId.value = nextTab?.id || ''
    }

    tabs.value = nextTabs
    saveTabs()
    return nextTab
  }

  function closeOtherTabs(id: string) {
    loadTabs()
    tabs.value = tabs.value.filter((tab) => tab.id === id)
    activeId.value = tabs.value[0]?.id || ''
    saveTabs()
  }

  function closeRightTabs(id: string) {
    loadTabs()
    const index = tabs.value.findIndex((tab) => tab.id === id)
    if (index < 0) return

    const removedActive = tabs.value.slice(index + 1).some((tab) => tab.id === activeId.value)
    tabs.value = tabs.value.slice(0, index + 1)
    if (removedActive) {
      activeId.value = id
    }
    saveTabs()
  }

  function closeAllTabs() {
    tabs.value = []
    activeId.value = ''
    saveTabs()
  }

  function moveTab(fromIndex: number, toIndex: number) {
    loadTabs()
    if (fromIndex === toIndex || fromIndex < 0 || toIndex < 0 || fromIndex >= tabs.value.length || toIndex >= tabs.value.length) return

    const nextTabs = [...tabs.value]
    const [tab] = nextTabs.splice(fromIndex, 1)
    if (!tab) return

    nextTabs.splice(toIndex, 0, tab)
    tabs.value = nextTabs
    saveTabs()
  }

  function setActiveTab(to: string) {
    activeId.value = normalizeWorkspaceTabPath(to)
  }

  function trimTabs(preferredId: string) {
    let nextTabs = tabs.value
    while (nextTabs.length > maxWorkspaceTabs) {
      const removableIndex = nextTabs.findIndex((tab) => tab.id !== preferredId)
      nextTabs = nextTabs.filter((_, index) => index !== (removableIndex >= 0 ? removableIndex : 0))
    }
    tabs.value = nextTabs
  }

  return {
    tabs,
    activeId,
    loadTabs,
    upsertTab,
    updateTabTitle,
    removeTab,
    closeOtherTabs,
    closeRightTabs,
    closeAllTabs,
    moveTab,
    setActiveTab,
    normalizeTabPath: normalizeWorkspaceTabPath
  }
}

function isWorkspaceTabItem(value: unknown): value is WorkspaceTabItem {
  if (!value || typeof value !== 'object') return false

  const item = value as Partial<WorkspaceTabItem>
  return typeof item.id === 'string'
    && typeof item.to === 'string'
    && typeof item.title === 'string'
    && typeof item.createdAt === 'number'
    && typeof item.updatedAt === 'number'
}

function normalizeLoadedTabs(items: WorkspaceTabItem[]) {
  const normalizedTabs: WorkspaceTabItem[] = []

  for (const item of items) {
    const normalizedTo = normalizeWorkspaceTabTo(item.to || item.id)
    const id = normalizeWorkspaceTabPath(normalizedTo)
    const normalizedItem = {
      ...item,
      id,
      to: normalizedTo,
      title: item.title.trim() || id
    }
    const existingIndex = normalizedTabs.findIndex((tab) => tab.id === id)
    if (existingIndex >= 0) {
      normalizedTabs[existingIndex] = normalizedItem
    } else {
      normalizedTabs.push(normalizedItem)
    }
  }

  return normalizedTabs
}
