<template>
  <div class="flex h-full flex-col overflow-hidden bg-white dark:bg-slate-950">
    <div
      class="flex h-16 shrink-0 items-center border-b border-slate-200 dark:border-slate-800"
      :class="collapsed ? 'justify-center px-2' : 'justify-between px-4'"
    >
      <NuxtLink :to="localePath('/')" class="flex min-w-0 items-center gap-3" @click="$emit('navigate')">
        <span class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-slate-950 text-white dark:bg-white dark:text-slate-950">
          <UIcon name="i-lucide-radio-tower" class="size-5" />
        </span>
        <span v-if="!collapsed" class="truncate text-base font-semibold text-slate-950 dark:text-white">NextPT</span>
      </NuxtLink>

      <UButton
        v-if="showClose"
        color="neutral"
        variant="ghost"
        icon="i-lucide-x"
        :aria-label="$t('common.cancel')"
        @click="$emit('close')"
      />
    </div>

    <div
      v-if="returnAction"
      class="shrink-0 border-b border-slate-200 dark:border-slate-800"
      :class="collapsed ? 'p-2' : 'p-3'"
    >
      <NuxtLink
        :to="localePath(returnAction.to)"
        class="flex min-h-10 items-center gap-3 rounded-md border border-slate-200 bg-slate-50 text-sm font-medium text-slate-600 transition-colors hover:border-slate-300 hover:bg-white hover:text-slate-950 dark:border-slate-800 dark:bg-slate-900/60 dark:text-slate-300 dark:hover:border-slate-700 dark:hover:bg-slate-900 dark:hover:text-white"
        :class="collapsed ? 'justify-center px-0' : 'px-3'"
        :title="collapsed ? returnAction.label : undefined"
        :aria-label="returnAction.label"
        @click="$emit('navigate')"
      >
        <UIcon :name="returnAction.icon" class="size-4 shrink-0" />
        <span v-if="!collapsed" class="truncate">{{ returnAction.label }}</span>
      </NuxtLink>
    </div>

    <div class="min-h-0 flex-1 overflow-y-auto py-4" :class="collapsed ? 'px-2' : 'px-3'">
      <nav class="space-y-6">
        <section v-for="section in sections" :key="section.key">
          <h2 v-if="!collapsed && section.label" class="px-2 text-xs font-semibold uppercase text-slate-400 dark:text-slate-500">{{ section.label }}</h2>
          <div class="space-y-1" :class="!collapsed && section.label ? 'mt-2' : ''">
            <NuxtLink
              v-for="item in section.items"
              :key="item.to"
              :to="localePath(item.to)"
              class="flex min-h-10 items-center gap-3 rounded-md text-sm font-medium transition-colors"
              :class="[
                item.active ? 'bg-slate-950 text-white dark:bg-white dark:text-slate-950' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-900 dark:hover:text-white',
                collapsed ? 'justify-center px-0' : 'px-3'
              ]"
              :title="collapsed ? item.label : undefined"
              @click="$emit('navigate')"
            >
              <UIcon :name="item.icon" class="size-4 shrink-0" />
              <span v-if="!collapsed" class="truncate">{{ item.label }}</span>
            </NuxtLink>
          </div>
        </section>
      </nav>
    </div>

    <div class="shrink-0 border-t border-slate-200 dark:border-slate-800" :class="collapsed ? 'p-2' : 'p-3'">
      <UDropdownMenu
        v-if="userMenuItems?.length"
        :items="userMenuItems"
        :content="{ align: collapsed ? 'center' : 'end' }"
      >
        <button
          type="button"
          class="flex w-full min-w-0 items-center gap-3 rounded-md py-2 text-left transition-colors hover:bg-slate-100 disabled:cursor-wait disabled:opacity-70 dark:hover:bg-slate-900"
          :class="collapsed ? 'justify-center px-0' : 'px-2'"
          :title="collapsed ? user?.username || 'NextPT' : undefined"
          :aria-label="user?.username || 'NextPT'"
          :disabled="loggingOut"
        >
          <UAvatar :src="user?.avatar || undefined" :alt="user?.username || 'NextPT'" size="sm" />
          <div v-if="!collapsed" class="min-w-0 flex-1">
            <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ user?.username || 'NextPT' }}</p>
            <p class="truncate text-xs text-slate-500 dark:text-slate-400">{{ user?.roleName || user?.email || '-' }}</p>
          </div>
          <UIcon v-if="!collapsed" name="i-lucide-chevron-up" class="size-4 shrink-0 text-slate-400" />
        </button>
      </UDropdownMenu>

      <div
        v-else
        class="flex min-w-0 items-center gap-3 rounded-md py-2"
        :class="collapsed ? 'justify-center px-0' : 'px-2'"
        :title="collapsed ? user?.username || 'NextPT' : undefined"
      >
        <UAvatar :src="user?.avatar || undefined" :alt="user?.username || 'NextPT'" size="sm" />
        <div v-if="!collapsed" class="min-w-0">
          <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ user?.username || 'NextPT' }}</p>
          <p class="truncate text-xs text-slate-500 dark:text-slate-400">{{ user?.roleName || user?.email || '-' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface ShellNavItem {
  label: string
  to: string
  icon: string
  active: boolean
}

interface ShellNavSection {
  key: string
  label: string
  items: ShellNavItem[]
}

interface ShellSidebarAction {
  label: string
  to: string
  icon: string
}

interface ShellUserMenuItem {
  label: string
  icon?: string
  onSelect?: () => void | Promise<void>
}

interface ShellUser {
  username?: string
  email?: string
  roleName?: string
  avatar?: string
}

defineProps<{
  sections: ShellNavSection[]
  user?: ShellUser | null
  userMenuItems?: ShellUserMenuItem[][]
  returnAction?: ShellSidebarAction | null
  loggingOut?: boolean
  collapsed?: boolean
  showClose?: boolean
}>()

defineEmits<{
  close: []
  navigate: []
}>()

const localePath = useLocalePath()
</script>
