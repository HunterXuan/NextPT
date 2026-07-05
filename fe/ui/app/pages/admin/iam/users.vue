<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
        <form class="flex flex-col gap-2 lg:flex-row lg:items-center" @submit.prevent="submitSearch">
          <div class="relative min-w-0 flex-1">
            <UIcon name="i-lucide-search" class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-slate-400" />
            <input
              v-model="searchInput"
              type="search"
              class="h-10 w-full rounded-md border border-slate-200 bg-white pl-9 pr-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
              :placeholder="$t('admin.iam.users.searchPlaceholder')"
            >
          </div>
          <select v-model="query.order" class="h-10 shrink-0 rounded-md border border-slate-200 bg-white px-3 text-sm leading-10 text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 lg:w-52 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" @change="reloadFromFirstPage">
            <option v-for="option in orderOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
          <UButton type="submit" color="primary" icon="i-lucide-search" :loading="pending" class="h-10 shrink-0 justify-center px-4">
            {{ $t('admin.iam.users.search') }}
          </UButton>
        </form>
      </div>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_clamp(620px,36vw,760px)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-users" class="size-5 text-sky-600 dark:text-sky-300" />
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.list') }}</h2>
            </div>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[900px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[32%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.user') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.role') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.status') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.lastLogin') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.createdAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 6" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-3 py-3"><div class="h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-800" /></td>
                  <td class="px-3 py-3"><div class="h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="h-4 w-32 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>

          <div v-else-if="users.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.iam.users.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[900px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[32%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.user') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.role') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.status') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.lastLogin') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.iam.users.table.createdAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in users"
                  :key="item.id"
                  class="cursor-pointer border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === item.id ? 'bg-sky-50/70 dark:bg-sky-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                  role="button"
                  tabindex="0"
                  @click="selectUser(item)"
                  @keydown.enter.prevent="selectUser(item)"
                  @keydown.space.prevent="selectUser(item)"
                >
                  <td class="px-3 py-2.5 align-middle">
                    <div class="flex min-w-0 items-center gap-3">
                      <IamUserAvatar :id="item.id" :username="item.username" :avatar="item.avatar" size="sm" />
                      <div class="block min-w-0 text-left">
                        <span class="block truncate text-sm font-semibold text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300">
                          {{ item.username }}
                        </span>
                        <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">#{{ item.id }} · {{ item.email }}</span>
                      </div>
                    </div>
                  </td>
                  <td class="px-3 py-2.5 align-middle text-sm text-slate-600 dark:text-slate-300">
                    <span class="inline-flex max-w-full rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                      <span class="truncate">{{ roleName(item.role) }}</span>
                    </span>
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <UBadge :color="statusColor(item.status)" variant="soft">{{ statusLabel(item.status) }}</UBadge>
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <p class="truncate text-sm text-slate-600 dark:text-slate-300">{{ formatDateTime(item.lastLogin, locale) }}</p>
                    <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ item.lastIp || '-' }}</p>
                  </td>
                  <td class="px-3 py-2.5 align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ formatDateTime(item.createdAt, locale) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <AppPager
            class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
            size="sm"
            :page="query.page"
            :total="total"
            :page-size="query.size"
            :page-size-options="pageSizes"
            :disabled="pending"
            @page-change="changePage"
            @page-size-change="changePageSize"
          />
        </section>

        <section class="space-y-4">
          <div class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div v-if="selectedUser" class="p-4">
              <div class="flex min-w-0 items-start gap-3">
                <IamUserAvatar :id="selectedUser.id" :username="selectedUser.username" :avatar="selectedUser.avatar" size="lg" />
                <div class="min-w-0 flex-1">
                  <div class="flex flex-wrap items-center gap-2">
                    <p class="truncate text-base font-semibold text-slate-950 dark:text-white">{{ selectedUser.username }}</p>
                    <UBadge :color="statusColor(selectedUser.status)" variant="soft">{{ statusLabel(selectedUser.status) }}</UBadge>
                  </div>
                  <p class="mt-1 truncate text-sm text-slate-500 dark:text-slate-400">#{{ selectedUser.id }} · {{ selectedUser.email }}</p>
                  <div class="mt-3 flex flex-wrap items-center gap-2 text-xs">
                    <span class="inline-flex max-w-full items-center gap-1.5 rounded-md bg-slate-100 px-2 py-1 font-medium text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                      <UIcon name="i-lucide-shield" class="size-3.5 text-slate-400" />
                      <span class="truncate">{{ roleName(selectedUser.role) }}</span>
                    </span>
                    <span
                      class="inline-flex items-center gap-1.5 rounded-md px-2 py-1 font-medium"
                      :class="activeUserMods.length > 0
                        ? 'bg-red-50 text-red-700 dark:bg-red-950/40 dark:text-red-300'
                        : 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300'"
                    >
                      <UIcon :name="activeUserMods.length > 0 ? 'i-lucide-shield-alert' : 'i-lucide-check-circle-2'" class="size-3.5" />
                      {{ $t('admin.iam.users.mod.activeSummary', { count: numberFormatter.format(activeUserMods.length) }) }}
                    </span>
                  </div>
                </div>
              </div>

              <div class="mt-4 grid grid-cols-2 gap-1 rounded-md bg-slate-100 p-1 dark:bg-slate-800/80">
                <button
                  v-for="panel in userPanelOptions"
                  :key="panel.value"
                  type="button"
                  class="inline-flex h-8 items-center justify-center gap-1.5 rounded px-2 text-xs font-medium transition"
                  :class="activeUserPanel === panel.value
                    ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-950 dark:text-white'
                    : 'text-slate-600 hover:text-slate-950 dark:text-slate-300 dark:hover:text-white'"
                  @click="activeUserPanel = panel.value"
                >
                  <UIcon :name="panel.icon" class="size-3.5" />
                  <span>{{ panel.label }}</span>
                </button>
              </div>
            </div>

            <div v-else class="flex flex-col items-center justify-center px-4 py-14 text-center">
              <div class="flex size-12 items-center justify-center rounded-lg bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-500">
                <UIcon name="i-lucide-user-round-check" class="size-6" />
              </div>
              <p class="mt-3 text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.form.empty') }}</p>
              <p class="mt-1 max-w-64 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.panels.selectHint') }}</p>
            </div>
          </div>

          <div v-if="selectedUser && activeUserPanel === 'profile'" class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
                {{ $t('admin.iam.users.panels.profile') }}
              </h2>
            </div>

            <form class="space-y-4 p-4" @submit.prevent="saveUser">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.form.status') }}</span>
                <select v-model.number="form.status" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || saving">
                  <option v-for="option in userStatusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
                </select>
              </label>

              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.form.role') }}</span>
                <select v-model.number="form.role" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || saving || roles.length === 0">
                  <option v-for="role in roles" :key="role.id" :value="role.id">{{ roleNameWithLevel(role) }}</option>
                </select>
              </label>

              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.form.passkey') }}</span>
                <input v-model="form.passkey" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 font-mono text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || saving">
              </label>

              <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
                {{ formError }}
              </p>

              <div class="flex flex-col gap-2 sm:flex-row sm:justify-end">
                <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" :disabled="!selectedUser || saving || !isProfileDirty" @click="resetProfileForm">
                  {{ $t('admin.actions.reset') }}
                </UButton>
                <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving" :disabled="!selectedUser || !isProfileDirty">
                  {{ $t('common.save') }}
                </UButton>
                <UPopover
                  :content="{ side: 'top', align: 'end', sideOffset: 8 }"
                  :ui="{ content: 'w-72 p-3' }"
                >
                  <UButton type="button" color="neutral" variant="outline" icon="i-lucide-log-out" :loading="sessionDeleting" :disabled="!selectedUser">
                    {{ $t('admin.iam.users.form.kick') }}
                  </UButton>

                  <template #content="{ close }">
                    <div class="space-y-3">
                      <p class="text-sm font-medium text-slate-950 dark:text-white">
                        {{ $t('admin.iam.users.form.confirmKickTitle') }}
                      </p>
                      <p class="text-xs text-slate-500 dark:text-slate-400">
                        {{ $t('admin.iam.users.form.confirmKickDescription') }}
                      </p>
                      <div class="flex justify-end gap-2">
                        <UButton color="neutral" variant="ghost" size="xs" type="button" @click="close()">
                          {{ $t('common.cancel') }}
                        </UButton>
                        <UButton color="warning" size="xs" type="button" icon="i-lucide-log-out" :loading="sessionDeleting" @click="deleteSessions(close)">
                          {{ $t('admin.iam.users.form.kick') }}
                        </UButton>
                      </div>
                    </div>
                  </template>
                </UPopover>
              </div>
            </form>
          </div>

          <div v-if="selectedUser && activeUserPanel === 'stat'" class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.stat.title') }}</h2>
            </div>

            <div class="space-y-4 p-4">
              <div v-if="statPending" class="grid grid-cols-2 gap-2">
                <div v-for="item in 4" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
              </div>
              <p v-else-if="statError" class="rounded-md border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-800 dark:border-amber-900 dark:bg-amber-950/40 dark:text-amber-200">
                {{ statError }}
              </p>
              <div v-else class="grid grid-cols-2 gap-2">
                <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stat.uploaded') }}</p>
                  <p class="mt-1 truncate text-sm font-semibold text-slate-950 dark:text-white">{{ formatBytes(userStat?.uploaded || 0) }}</p>
                </div>
                <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stat.downloaded') }}</p>
                  <p class="mt-1 truncate text-sm font-semibold text-slate-950 dark:text-white">{{ formatBytes(userStat?.downloaded || 0) }}</p>
                </div>
                <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stat.seedTime') }}</p>
                  <p class="mt-1 truncate text-sm font-semibold text-slate-950 dark:text-white">{{ durationLabel(userStat?.seedTime || 0) }}</p>
                </div>
                <div class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stat.bonus') }}</p>
                  <p class="mt-1 truncate text-sm font-semibold text-amber-600 dark:text-amber-300">{{ numberFormatter.format(userStat?.bonus || 0) }}</p>
                </div>
              </div>

              <form class="space-y-3" @submit.prevent="saveStatDiff">
                <div class="grid gap-2 sm:grid-cols-3">
                  <label class="block">
                    <span class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stat.uploadedDiff') }}</span>
                    <input v-model="statForm.uploadedDiff" type="number" step="1" class="mt-1 h-9 w-full rounded-md border border-slate-200 bg-white px-2 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="!selectedUser || statSaving">
                  </label>
                  <label class="block">
                    <span class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stat.downloadedDiff') }}</span>
                    <input v-model="statForm.downloadedDiff" type="number" step="1" class="mt-1 h-9 w-full rounded-md border border-slate-200 bg-white px-2 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="!selectedUser || statSaving">
                  </label>
                  <label class="block">
                    <span class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stat.bonusDiff') }}</span>
                    <input v-model="statForm.bonusDiff" type="number" step="0.01" class="mt-1 h-9 w-full rounded-md border border-slate-200 bg-white px-2 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="!selectedUser || statSaving">
                  </label>
                </div>
                <div class="flex justify-end">
                  <UButton type="submit" color="neutral" variant="outline" icon="i-lucide-chart-no-axes-combined" :loading="statSaving" :disabled="!selectedUser || !hasStatDiff" class="w-full justify-center sm:w-auto">
                    {{ $t('admin.iam.users.stat.submit') }}
                  </UButton>
                </div>
              </form>
            </div>
          </div>

          <template v-if="selectedUser && activeUserPanel === 'mod'">
            <div class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
              <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
                <div>
                  <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.mod.title') }}</h2>
                  <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                    {{ $t('admin.iam.users.mod.activeSummary', { count: numberFormatter.format(activeUserMods.length) }) }}
                  </p>
                </div>
              </div>

              <div class="space-y-4 p-4">
                <form class="space-y-3" @submit.prevent="applyUserMod">
                  <label class="block">
                    <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.mod.type') }}</span>
                    <select v-model.number="modForm.type" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || modApplying">
                      <option v-for="option in modTypeOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
                    </select>
                  </label>

                  <label class="block">
                    <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.mod.durationDays') }}</span>
                    <input v-model.number="modForm.durationDays" type="number" min="0" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || modApplying">
                    <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.mod.permanentHint') }}</p>
                  </label>

                  <label class="block">
                    <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.mod.reason') }}</span>
                    <textarea v-model="modForm.reason" rows="3" class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || modApplying" />
                  </label>

                  <p v-if="modsError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
                    {{ modsError }}
                  </p>

                  <UButton type="submit" color="primary" icon="i-lucide-shield-alert" :loading="modApplying" :disabled="!canApplyUserMod">
                    {{ $t('admin.iam.users.mod.submit') }}
                  </UButton>
                </form>

              <div class="border-t border-slate-200 pt-4 dark:border-slate-800">
                <div v-if="modsPending && userMods.length === 0" class="space-y-2">
                  <div v-for="item in 3" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
                </div>
                <div v-else-if="userMods.length === 0" class="rounded-md border border-dashed border-slate-200 px-3 py-8 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
                  {{ $t('admin.iam.users.mod.empty') }}
                </div>
                <div v-else class="space-y-2">
                  <article v-for="item in sortedUserMods" :key="item.id" class="rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
                    <div class="flex items-start justify-between gap-3">
                      <div class="min-w-0">
                        <div class="flex flex-wrap items-center gap-2">
                          <UBadge color="neutral" variant="soft">{{ modTypeLabel(item.mod_type) }}</UBadge>
                          <UBadge :color="item.is_active ? 'error' : 'neutral'" variant="soft">
                            {{ item.is_active ? $t('admin.iam.users.mod.status.active') : $t('admin.iam.users.mod.status.inactive') }}
                          </UBadge>
                        </div>
                        <p class="mt-2 whitespace-pre-wrap break-words text-sm text-slate-700 dark:text-slate-200">{{ item.reason || '-' }}</p>
                        <p class="mt-2 text-xs text-slate-500 dark:text-slate-400">
                          {{ formatDateTime(item.created_at, locale) }} · {{ $t('admin.iam.users.mod.expireAt') }} {{ modExpireLabel(item) }}
                        </p>
                      </div>
                      <UPopover
                        v-if="item.is_active"
                        :content="{ side: 'top', align: 'end', sideOffset: 8 }"
                        :ui="{ content: 'w-72 p-3' }"
                      >
                        <UButton
                          color="neutral"
                          variant="outline"
                          size="xs"
                          icon="i-lucide-undo-2"
                          :loading="modRemovingId === item.id"
                          :disabled="modRemovingId > 0"
                        >
                          {{ $t('admin.iam.users.mod.remove') }}
                        </UButton>

                        <template #content="{ close }">
                          <div class="space-y-3">
                            <p class="text-sm font-medium text-slate-950 dark:text-white">
                              {{ $t('admin.iam.users.mod.confirmRemoveTitle') }}
                            </p>
                            <p class="text-xs text-slate-500 dark:text-slate-400">
                              {{ $t('admin.iam.users.mod.confirmRemoveDescription', { type: modTypeLabel(item.mod_type) }) }}
                            </p>
                            <div class="flex justify-end gap-2">
                              <UButton color="neutral" variant="ghost" size="xs" type="button" @click="close()">
                                {{ $t('common.cancel') }}
                              </UButton>
                              <UButton color="warning" size="xs" type="button" icon="i-lucide-undo-2" :loading="modRemovingId === item.id" :disabled="modRemovingId > 0" @click="removeUserMod(item, close)">
                                {{ $t('admin.iam.users.mod.remove') }}
                              </UButton>
                            </div>
                          </div>
                        </template>
                      </UPopover>
                    </div>
                  </article>
                </div>
              </div>
            </div>
          </div>
          </template>

          <div v-if="selectedUser && activeUserPanel === 'permission'" class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.permission.title') }}</h2>
            </div>

            <div class="space-y-4 p-4">
              <div v-if="permissionPending" class="space-y-3">
                <div class="grid gap-3 2xl:grid-cols-2">
                  <div v-for="item in 2" :key="item" class="h-36 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
                </div>
                <div class="h-48 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
              </div>

              <template v-else>
                <p v-if="permissionError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
                  {{ permissionError }}
                </p>

                <div class="grid gap-3 2xl:grid-cols-2">
                  <section class="min-w-0 rounded-md border border-slate-200 dark:border-slate-800">
                    <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">
                      <div class="min-w-0">
                        <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.permission.roleTitle') }}</p>
                        <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ roleName(selectedUser.role) }}</p>
                      </div>
                      <UBadge color="neutral" variant="soft">{{ rolePermissions.length }}</UBadge>
                    </div>
                    <div v-if="rolePermissions.length" class="max-h-44 overflow-y-auto p-2">
                      <code
                        v-for="permission in rolePermissions"
                        :key="permission"
                        class="mb-1 block truncate rounded bg-slate-50 px-2 py-1.5 text-xs font-semibold text-slate-600 last:mb-0 dark:bg-slate-950 dark:text-slate-300"
                      >
                        {{ permission }}
                      </code>
                    </div>
                    <div v-else class="px-3 py-8 text-center text-sm text-slate-500 dark:text-slate-400">
                      {{ $t('admin.iam.users.permission.noRolePermissions') }}
                    </div>
                  </section>

                  <section class="min-w-0 rounded-md border border-slate-200 dark:border-slate-800">
                    <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">
                      <div>
                        <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.permission.userTitle') }}</p>
                        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.permission.manualHint') }}</p>
                      </div>
                      <UBadge color="neutral" variant="soft">{{ userAcls.length }}</UBadge>
                    </div>
                    <div v-if="userAcls.length" class="max-h-44 overflow-y-auto p-2">
                      <label
                        v-for="acl in userAcls"
                        :key="acl.id"
                        class="mb-1 flex items-center gap-2 rounded px-2 py-1.5 last:mb-0"
                        :class="acl.sourceType === manualPermissionSource ? 'hover:bg-slate-50 dark:hover:bg-slate-950' : 'opacity-70'"
                      >
                        <input
                          v-model="selectedUserAclIds"
                          type="checkbox"
                          class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                          :value="acl.id"
                          :disabled="permissionSaving || acl.sourceType !== manualPermissionSource"
                        >
                        <code class="min-w-0 flex-1 truncate text-xs font-semibold" :class="acl.isDeny ? 'text-red-600 dark:text-red-300' : 'text-slate-700 dark:text-slate-200'">
                          {{ acl.rawPermKey }}
                        </code>
                        <UBadge :color="permissionSourceColor(acl.sourceType)" variant="soft">{{ permissionSourceLabel(acl.sourceType) }}</UBadge>
                      </label>
                    </div>
                    <div v-else class="px-3 py-8 text-center text-sm text-slate-500 dark:text-slate-400">
                      {{ $t('admin.iam.users.permission.noUserPermissions') }}
                    </div>
                  </section>
                </div>

                <section class="rounded-md border border-slate-200 dark:border-slate-800">
                  <div class="border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">
                    <div class="flex items-center justify-between gap-3">
                      <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.permission.availableTitle') }}</p>
                      <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.permission.selectedCount', { count: selectedPermissionKeys.length }) }}</p>
                    </div>
                    <div class="relative mt-2">
                      <UIcon name="i-lucide-search" class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-slate-400" />
                      <input
                        v-model="permissionForm.permKey"
                        class="h-9 w-full rounded-md border border-slate-200 bg-white pl-9 pr-9 font-mono text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                        :placeholder="$t('admin.iam.users.permission.placeholder')"
                        :disabled="!selectedUser || permissionSaving"
                      >
                      <button
                        v-if="permissionForm.permKey"
                        type="button"
                        class="absolute right-1 top-1/2 flex size-7 -translate-y-1/2 items-center justify-center rounded text-slate-400 transition hover:bg-slate-100 hover:text-slate-700 dark:hover:bg-slate-800 dark:hover:text-slate-200"
                        :disabled="permissionSaving"
                        @click="clearPermissionSearch"
                      >
                        <UIcon name="i-lucide-x" class="size-4" />
                      </button>
                    </div>
                  </div>

                  <div v-if="filteredPermissions.length" class="max-h-52 overflow-y-auto p-2">
                    <label
                      v-for="permission in filteredPermissions"
                      :key="permission"
                      class="mb-1 flex items-center gap-2 rounded px-2 py-1.5 transition last:mb-0 hover:bg-slate-50 dark:hover:bg-slate-950"
                    >
                      <input v-model="selectedPermissionKeys" type="checkbox" class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600" :value="permission" :disabled="permissionSaving">
                      <code class="min-w-0 flex-1 truncate text-xs font-semibold text-slate-700 dark:text-slate-200">{{ permission }}</code>
                    </label>
                  </div>
                  <div v-else class="px-3 py-8 text-center text-sm text-slate-500 dark:text-slate-400">
                    {{ $t('admin.iam.users.permission.empty') }}
                  </div>
                </section>

                <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  <label class="flex items-center gap-2 text-sm font-medium text-slate-700 dark:text-slate-200">
                    <input v-model="permissionForm.isDeny" type="checkbox" class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600" :disabled="!selectedUser || permissionSaving">
                    <span>{{ $t('admin.iam.users.permission.deny') }}</span>
                  </label>

                  <div class="grid grid-cols-2 gap-2 sm:flex sm:justify-end">
                    <UButton type="button" color="primary" icon="i-lucide-plus" :loading="permissionSaving" :disabled="!selectedUser || grantPermissionKeys.length === 0" @click="applyPermission('grant')">
                      {{ $t('admin.iam.users.permission.grantSelected', { count: grantPermissionKeys.length }) }}
                    </UButton>
                    <UButton type="button" color="neutral" variant="outline" icon="i-lucide-minus" :loading="permissionSaving" :disabled="!selectedUser || selectedUserAclItems.length === 0" @click="applyPermission('revoke')">
                      {{ $t('admin.iam.users.permission.revokeSelected', { count: selectedUserAclItems.length }) }}
                    </UButton>
                  </div>
                </div>
              </template>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminIamRole, AdminIamUser, AdminIamUserPermissionDetailOut, AdminIamUserStat, AdminIamUserStatUpdateInput, AdminIamUserUpdateInput, AdminModUserItem } from '~/composables/useAdmin'
import { formatBytes, formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

useHead({ title: t('admin.iam.users.title') })

const users = ref<AdminIamUser[]>([])
const roles = ref<AdminIamRole[]>([])
const permissions = ref<string[]>([])
const userMods = ref<AdminModUserItem[]>([])
const total = ref(0)
const pending = ref(false)
const rolesPending = ref(false)
const saving = ref(false)
const modsPending = ref(false)
const modApplying = ref(false)
const modRemovingId = ref(0)
const sessionDeleting = ref(false)
const permissionPending = ref(false)
const permissionSaving = ref(false)
const statPending = ref(false)
const statSaving = ref(false)
const selectedId = ref<number | null>(null)
const userStat = ref<AdminIamUserStat | null>(null)
const userPermissionDetail = ref<AdminIamUserPermissionDetailOut | null>(null)
const errorMessage = ref('')
const formError = ref('')
const modsError = ref('')
const permissionError = ref('')
const statError = ref('')
const searchInput = ref('')
const activeUserPanel = ref<'profile' | 'stat' | 'mod' | 'permission'>('profile')
const selectedPermissionKeys = ref<string[]>([])
const selectedUserAclIds = ref<number[]>([])

const query = reactive({
  search: '',
  order: 'id desc',
  page: 1,
  size: 30
})

const form = reactive({
  status: 1,
  role: 0,
  passkey: ''
})
const originalPasskey = ref('')

const permissionForm = reactive({
  permKey: '',
  isDeny: false
})

const modForm = reactive({
  type: 1,
  durationDays: 7,
  reason: ''
})

const statForm = reactive({
  uploadedDiff: '',
  downloadedDiff: '',
  bonusDiff: ''
})

const manualPermissionSource = 1
const userModPermissionSource = 2
const pageSizes = [20, 30, 50, 100]
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))
const roleMap = computed(() => new Map(roles.value.map((item) => [item.id, item])))
const { roleNameWithLevel } = useAdminIamRoleLevels(roles)
const selectedUser = computed(() => users.value.find((item) => item.id === selectedId.value) || null)
const isProfileDirty = computed(() => {
  const user = selectedUser.value
  return Boolean(user && Object.keys(buildUserUpdateInput(user)).length > 0)
})
const hasStatDiff = computed(() => Object.keys(buildStatDiffInput()).length > 0)
const activeUserMods = computed(() => userMods.value.filter((item) => item.is_active))
const sortedUserMods = computed(() => {
  return userMods.value.slice().sort((a, b) => {
    if (a.is_active !== b.is_active) return a.is_active ? -1 : 1
    return Date.parse(b.created_at || '') - Date.parse(a.created_at || '')
  })
})
const canApplyUserMod = computed(() => {
  return Boolean(selectedUser.value && modForm.type > 0 && modForm.reason.trim() && !modApplying.value)
})
const rolePermissions = computed(() => {
  const role = selectedUser.value ? roleMap.value.get(selectedUser.value.role) : null
  return normalizePermissionList(role?.permissions).sort()
})
const userAcls = computed(() => userPermissionDetail.value?.userAcls || [])
const selectedUserAclItems = computed(() => {
  const ids = new Set(selectedUserAclIds.value)
  return userAcls.value.filter((acl) => ids.has(acl.id) && acl.sourceType === manualPermissionSource)
})
const grantPermissionKeys = computed(() => {
  const keys = new Set(selectedPermissionKeys.value.map((item) => item.trim()).filter(Boolean))
  const customKey = permissionForm.permKey.trim()
  if (keys.size === 0 && customKey) keys.add(customKey)
  return Array.from(keys).filter(isWildcardPermission).sort()
})
const filteredPermissions = computed(() => {
  const keyword = permissionForm.permKey.trim().toLowerCase()
  const source = permissions.value
  if (!keyword) return source.slice(0, 80)
  return source.filter((permission) => permission.toLowerCase().includes(keyword)).slice(0, 80)
})
const orderOptions = computed(() => [
  { value: 'id desc', label: t('admin.iam.users.orders.newest') },
  { value: 'id asc', label: t('admin.iam.users.orders.oldest') },
  { value: 'last_login desc', label: t('admin.iam.users.orders.lastLogin') },
  { value: 'role desc', label: t('admin.iam.users.orders.roleDesc') }
])
const userStatusOptions = computed(() => [
  { value: 0, label: t('admin.iam.users.status.pending') },
  { value: 1, label: t('admin.iam.users.status.confirmed') },
  { value: 2, label: t('admin.iam.users.status.disabled') }
])
const modTypeOptions = computed(() => [
  { value: 1, label: t('admin.iam.users.mod.types.warned') },
  { value: 2, label: t('admin.iam.users.mod.types.banned') },
  { value: 3, label: t('admin.iam.users.mod.types.leechWarned') },
  { value: 4, label: t('admin.iam.users.mod.types.uploadBanned') },
  { value: 5, label: t('admin.iam.users.mod.types.downloadBanned') },
  { value: 6, label: t('admin.iam.users.mod.types.forumBanned') }
])
const userPanelOptions = computed(() => [
  { value: 'profile' as const, label: t('admin.iam.users.panels.profile'), icon: 'i-lucide-id-card' },
  { value: 'stat' as const, label: t('admin.iam.users.panels.stat'), icon: 'i-lucide-chart-no-axes-combined' },
  { value: 'mod' as const, label: t('admin.iam.users.panels.mod'), icon: 'i-lucide-shield-alert' },
  { value: 'permission' as const, label: t('admin.iam.users.panels.permission'), icon: 'i-lucide-key-round' }
])

onMounted(async () => {
  await Promise.all([loadRolesAndPermissions(), loadUsers()])
})

async function loadRolesAndPermissions() {
  rolesPending.value = true
  try {
    const data = await adminApi.listIamRoles()
    roles.value = (data.roles || []).sort((a, b) => b.level - a.level || a.id - b.id)
    const roleId = roles.value[0]?.id || 1
    const permissionData = await adminApi.listIamPermissions(roleId)
    permissions.value = (permissionData.permissions || []).slice().sort()
  } catch {
    permissions.value = []
  } finally {
    rolesPending.value = false
  }
}

async function loadUsers() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listIamUsers({
      search: query.search,
      order: query.order,
      page: query.page,
      size: query.size
    })
    users.value = data.users || []
    total.value = data.total || 0
    if (selectedId.value) {
      const next = users.value.find((item) => item.id === selectedId.value)
      if (next) {
        selectUser(next)
      } else {
        selectedId.value = null
        userMods.value = []
        userStat.value = null
        userPermissionDetail.value = null
        resetPermissionSelection()
      }
    }
  } catch (error: unknown) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function submitSearch() {
  query.search = searchInput.value.trim()
  reloadFromFirstPage()
}

function reloadFromFirstPage() {
  query.page = 1
  loadUsers()
}

function changePageSize(size: number) {
  query.size = size
  reloadFromFirstPage()
}

function changePage(page: number) {
  query.page = Math.min(Math.max(1, page), totalPages.value)
  loadUsers()
}

function roleName(roleId: number) {
  const role = roleMap.value.get(roleId)
  if (!role) return `#${roleId}`
  return roleNameWithLevel(role)
}

function statusLabel(status: number) {
  return userStatusOptions.value.find((item) => item.value === status)?.label || String(status)
}

function statusColor(status: number) {
  if (status === 1) return 'success'
  if (status === 2) return 'error'
  return 'warning'
}

function selectUser(user: AdminIamUser) {
  selectedId.value = user.id
  fillProfileForm(user)
  formError.value = ''
  modsError.value = ''
  permissionError.value = ''
  statError.value = ''
  resetStatForm()
  resetPermissionSelection()
  loadSelectedUserStat()
  loadSelectedUserMods()
  loadSelectedUserPermissions()
}

function fillProfileForm(user: AdminIamUser) {
  form.status = user.status
  form.role = user.role
  form.passkey = user.passkey || ''
  originalPasskey.value = user.passkey || ''
}

function resetProfileForm() {
  const user = selectedUser.value
  if (!user) return
  fillProfileForm(user)
  formError.value = ''
}

function clearPermissionSearch() {
  if (permissionSaving.value) return
  permissionForm.permKey = ''
}

function buildUserUpdateInput(user: AdminIamUser): AdminIamUserUpdateInput {
  const input: AdminIamUserUpdateInput = {}
  if (form.status !== user.status) input.status = Number(form.status)
  if (form.role !== user.role) input.role = Number(form.role)
  const passkey = form.passkey.trim()
  if (passkey !== originalPasskey.value) input.passkey = passkey
  return input
}

async function saveUser() {
  const user = selectedUser.value
  if (!user) return
  const input = buildUserUpdateInput(user)
  if (input.passkey !== undefined && !input.passkey) {
    formError.value = t('admin.iam.users.form.passkeyRequired')
    return
  }
  if (Object.keys(input).length === 0) {
    toast.add({ title: t('admin.iam.users.form.noChanges') })
    return
  }

  saving.value = true
  formError.value = ''
  try {
    await adminApi.updateIamUser(user.id, input)
    toast.add({ title: t('admin.iam.users.form.saved') })
    await loadUsers()
  } catch (error: unknown) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

async function deleteSessions(close?: () => void) {
  const user = selectedUser.value
  if (!user) return

  sessionDeleting.value = true
  try {
    await adminApi.deleteIamUserSessions(user.id)
    toast.add({ title: t('admin.iam.users.form.kicked') })
    close?.()
  } catch (error: unknown) {
    toast.add({
      color: 'error',
      title: error instanceof ApiError ? error.message : t('common.requestFailed')
    })
  } finally {
    sessionDeleting.value = false
  }
}

async function loadSelectedUserMods() {
  const user = selectedUser.value
  if (!user) return

  modsPending.value = true
  modsError.value = ''
  try {
    const data = await adminApi.listUserMods(user.id)
    userMods.value = data.list || []
  } catch (error: unknown) {
    userMods.value = []
    modsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    modsPending.value = false
  }
}

async function loadSelectedUserPermissions() {
  const user = selectedUser.value
  if (!user) return

  permissionPending.value = true
  permissionError.value = ''
  try {
    userPermissionDetail.value = await adminApi.getIamUserPermissions(user.id, {
      page: 1,
      size: 1000,
      sourceType: manualPermissionSource,
      wildcardOnly: true
    })
    selectedUserAclIds.value = selectedUserAclIds.value.filter((id) => {
      return userAcls.value.some((acl) => acl.id === id && acl.sourceType === manualPermissionSource)
    })
  } catch (error: unknown) {
    userPermissionDetail.value = null
    permissionError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    permissionPending.value = false
  }
}

async function applyUserMod() {
  const user = selectedUser.value
  if (!user || !canApplyUserMod.value) return

  const durationDays = Math.max(0, Number(modForm.durationDays || 0))
  modApplying.value = true
  modsError.value = ''
  try {
    await adminApi.applyUserMod(user.id, {
      type: Number(modForm.type),
      reason: modForm.reason.trim(),
      duration: Math.round(durationDays * 24 * 60 * 60)
    })
    toast.add({ title: t('admin.iam.users.mod.applied'), color: 'success', icon: 'i-lucide-check-circle' })
    modForm.reason = ''
    await loadSelectedUserMods()
    await loadSelectedUserPermissions()
  } catch (error: unknown) {
    modsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    modApplying.value = false
  }
}

async function removeUserMod(item: AdminModUserItem, close?: () => void) {
  const user = selectedUser.value
  if (!user || !item.is_active) return

  modRemovingId.value = item.id
  modsError.value = ''
  try {
    await adminApi.removeUserMod(user.id, item.id)
    toast.add({ title: t('admin.iam.users.mod.removed'), color: 'success', icon: 'i-lucide-check-circle' })
    close?.()
    await loadSelectedUserMods()
    await loadSelectedUserPermissions()
  } catch (error: unknown) {
    modsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    modRemovingId.value = 0
  }
}

async function applyPermission(action: 'grant' | 'revoke') {
  const user = selectedUser.value
  if (!user) return

  const grantKeys = grantPermissionKeys.value
  const revokeAcls = selectedUserAclItems.value
  if (action === 'grant' && grantKeys.length === 0) return
  if (action === 'revoke' && revokeAcls.length === 0) return

  permissionSaving.value = true
  permissionError.value = ''
  try {
    if (action === 'grant') {
      await adminApi.grantIamUserPermission(user.id, {
        permKeys: grantKeys,
        isDeny: permissionForm.isDeny
      })
      toast.add({ title: t('admin.iam.users.permission.granted', { count: grantKeys.length }) })
      selectedPermissionKeys.value = []
      permissionForm.permKey = ''
    } else {
      await adminApi.revokeIamUserPermission(user.id, {
        ids: revokeAcls.map((acl) => acl.id)
      })
      toast.add({ title: t('admin.iam.users.permission.revoked', { count: revokeAcls.length }) })
      selectedUserAclIds.value = []
    }
    await loadSelectedUserPermissions()
  } catch (error: unknown) {
    permissionError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    permissionSaving.value = false
  }
}

async function loadSelectedUserStat() {
  const user = selectedUser.value
  if (!user) return

  statPending.value = true
  statError.value = ''
  try {
    userStat.value = await adminApi.getIamUserStat(user.id)
  } catch (error: unknown) {
    userStat.value = null
    statError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    statPending.value = false
  }
}

function buildStatDiffInput(): AdminIamUserStatUpdateInput {
  const input: AdminIamUserStatUpdateInput = {}
  const uploadedDiff = Number(statForm.uploadedDiff)
  const downloadedDiff = Number(statForm.downloadedDiff)
  const bonusDiff = Number(statForm.bonusDiff)
  if (Number.isFinite(uploadedDiff) && uploadedDiff !== 0) input.uploadedDiff = uploadedDiff
  if (Number.isFinite(downloadedDiff) && downloadedDiff !== 0) input.downloadedDiff = downloadedDiff
  if (Number.isFinite(bonusDiff) && bonusDiff !== 0) input.bonusDiff = bonusDiff
  return input
}

async function saveStatDiff() {
  const user = selectedUser.value
  if (!user) return
  const input = buildStatDiffInput()
  if (Object.keys(input).length === 0) return

  statSaving.value = true
  statError.value = ''
  try {
    await adminApi.incrementIamUserStat(user.id, input)
    toast.add({ title: t('admin.iam.users.stat.saved'), color: 'success', icon: 'i-lucide-check-circle' })
    resetStatForm()
    await loadSelectedUserStat()
  } catch (error: unknown) {
    statError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    statSaving.value = false
  }
}

function resetStatForm() {
  statForm.uploadedDiff = ''
  statForm.downloadedDiff = ''
  statForm.bonusDiff = ''
}

function resetPermissionSelection() {
  selectedPermissionKeys.value = []
  selectedUserAclIds.value = []
}

function normalizePermissionList(value: unknown): string[] {
  if (Array.isArray(value)) {
    return value.filter((item): item is string => typeof item === 'string' && item.trim() !== '')
  }
  if (value && typeof value === 'object') {
    return Object.values(value).filter((item): item is string => typeof item === 'string' && item.trim() !== '')
  }
  return []
}

function isWildcardPermission(permission: string) {
  return permission === '*' || permission.endsWith(':*')
}

function permissionSourceLabel(sourceType: number) {
  if (sourceType === manualPermissionSource) return t('admin.iam.users.permission.sources.manual')
  if (sourceType === userModPermissionSource) return t('admin.iam.users.permission.sources.mod')
  return `#${sourceType}`
}

function permissionSourceColor(sourceType: number) {
  if (sourceType === manualPermissionSource) return 'primary'
  if (sourceType === userModPermissionSource) return 'warning'
  return 'neutral'
}

function durationLabel(seconds: number) {
  const hours = seconds / 3600
  if (hours >= 1) return `${numberFormatter.value.format(Math.round(hours * 10) / 10)} h`
  return `${numberFormatter.value.format(Math.round(seconds / 60))} min`
}

function modTypeLabel(type: number) {
  return modTypeOptions.value.find((item) => item.value === type)?.label || `#${type}`
}

function modExpireLabel(item: AdminModUserItem) {
  return item.expire_at ? formatDateTime(item.expire_at, locale.value) : t('admin.iam.users.mod.permanent')
}
</script>
