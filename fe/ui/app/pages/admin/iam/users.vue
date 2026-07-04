<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto]">
        <form class="flex flex-col gap-2 sm:flex-row" @submit.prevent="submitSearch">
          <div class="relative min-w-0 flex-1">
            <UIcon name="i-lucide-search" class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-slate-400" />
            <input
              v-model="searchInput"
              type="search"
              class="h-10 w-full rounded-md border border-slate-200 bg-white pl-9 pr-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
              :placeholder="$t('admin.iam.users.searchPlaceholder')"
            >
          </div>
          <select v-model="query.order" class="h-10 rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" @change="reloadFromFirstPage">
            <option v-for="option in orderOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
          <UButton type="submit" color="primary" icon="i-lucide-search" :loading="pending">
            {{ $t('admin.iam.users.search') }}
          </UButton>
        </form>

        <div class="grid grid-cols-2 gap-2 sm:flex sm:items-center">
          <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
            <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stats.total') }}</p>
            <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(total) }}</p>
          </div>
          <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
            <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.users.stats.current') }}</p>
            <p class="mt-1 text-lg font-semibold text-sky-700 dark:text-sky-300">{{ numberFormatter.format(users.length) }}</p>
          </div>
        </div>
      </div>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_400px]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.list') }}</h2>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[1080px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[24%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.user') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.role') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.status') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.lastLogin') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.createdAt') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.users.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 6" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-4 py-4"><div class="h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-800" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-36 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-20 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
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
            <table class="min-w-[1080px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[24%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.user') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.role') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.status') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.lastLogin') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.users.table.createdAt') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.users.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in users"
                  :key="item.id"
                  class="border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === item.id ? 'bg-sky-50/70 dark:bg-sky-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-4 py-3 align-middle">
                    <button type="button" class="block max-w-full text-left" @click="selectUser(item)">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300">
                        #{{ item.id }} {{ item.username }}
                      </span>
                      <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">{{ item.email }}</span>
                    </button>
                  </td>
                  <td class="px-4 py-3 align-middle text-sm text-slate-600 dark:text-slate-300">
                    <span class="block truncate">{{ roleName(item.role) }}</span>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <UBadge :color="statusColor(item.status)" variant="soft">{{ statusLabel(item.status) }}</UBadge>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <p class="truncate text-sm text-slate-600 dark:text-slate-300">{{ formatDateTime(item.lastLogin, locale) }}</p>
                    <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ item.lastIp || '-' }}</p>
                  </td>
                  <td class="px-4 py-3 align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ formatDateTime(item.createdAt, locale) }}
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <div class="flex items-center justify-end gap-2">
                      <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-pencil" @click="selectUser(item)">
                        {{ $t('admin.actions.edit') }}
                      </UButton>
                    </div>
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
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
                {{ selectedUser ? selectedUser.username : $t('admin.iam.users.form.empty') }}
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
                  <option v-for="role in roles" :key="role.id" :value="role.id">{{ roleName(role.id) }}</option>
                </select>
              </label>

              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.form.passkey') }}</span>
                <input v-model="form.passkey" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 font-mono text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || saving">
              </label>

              <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
                {{ formError }}
              </p>

              <div class="grid grid-cols-2 gap-2">
                <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving" :disabled="!selectedUser">
                  {{ $t('common.save') }}
                </UButton>
                <UButton type="button" color="neutral" variant="outline" icon="i-lucide-log-out" :loading="sessionDeleting" :disabled="!selectedUser" @click="deleteSessions">
                  {{ $t('admin.iam.users.form.kick') }}
                </UButton>
              </div>
            </form>
          </div>

          <div class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
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
                <UButton type="submit" color="neutral" variant="outline" icon="i-lucide-chart-no-axes-combined" :loading="statSaving" :disabled="!selectedUser || !hasStatDiff">
                  {{ $t('admin.iam.users.stat.submit') }}
                </UButton>
              </form>
            </div>
          </div>

          <div class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.ban.title') }}</h2>
            </div>

            <form class="space-y-4 p-4" @submit.prevent="banUser">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.ban.reason') }}</span>
                <textarea v-model="banForm.reason" rows="3" class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || banning" />
              </label>

              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.ban.durationDays') }}</span>
                <input v-model.number="banForm.durationDays" type="number" min="0" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || banning">
              </label>

              <p v-if="banError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
                {{ banError }}
              </p>

              <UButton type="submit" color="error" icon="i-lucide-ban" :loading="banning" :disabled="!selectedUser || !banForm.reason.trim()">
                {{ $t('admin.iam.users.ban.submit') }}
              </UButton>
            </form>
          </div>

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
                      <UButton
                        v-if="item.is_active"
                        color="neutral"
                        variant="outline"
                        size="xs"
                        icon="i-lucide-undo-2"
                        :loading="modRemovingId === item.id"
                        :disabled="modRemovingId > 0"
                        @click="removeUserMod(item)"
                      >
                        {{ $t('admin.iam.users.mod.remove') }}
                      </UButton>
                    </div>
                  </article>
                </div>
              </div>
            </div>
          </div>

          <div class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.users.permission.title') }}</h2>
            </div>

            <form class="space-y-4 p-4" @submit.prevent="applyPermission('grant')">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.users.permission.key') }}</span>
                <input v-model="permissionForm.permKey" list="admin-iam-permissions" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 font-mono text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="!selectedUser || permissionSaving">
                <datalist id="admin-iam-permissions">
                  <option v-for="permission in permissions" :key="permission" :value="permission" />
                </datalist>
              </label>

              <label class="flex items-center gap-2 text-sm font-medium text-slate-700 dark:text-slate-200">
                <input v-model="permissionForm.isDeny" type="checkbox" class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600" :disabled="!selectedUser || permissionSaving">
                <span>{{ $t('admin.iam.users.permission.deny') }}</span>
              </label>

              <p v-if="permissionError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
                {{ permissionError }}
              </p>

              <div class="grid grid-cols-2 gap-2">
                <UButton type="submit" color="primary" icon="i-lucide-plus" :loading="permissionSaving" :disabled="!selectedUser || !permissionForm.permKey.trim()">
                  {{ $t('admin.iam.users.permission.grant') }}
                </UButton>
                <UButton type="button" color="neutral" variant="outline" icon="i-lucide-minus" :loading="permissionSaving" :disabled="!selectedUser || !permissionForm.permKey.trim()" @click="applyPermission('revoke')">
                  {{ $t('admin.iam.users.permission.revoke') }}
                </UButton>
              </div>
            </form>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminIamRole, AdminIamUser, AdminIamUserStat, AdminIamUserStatUpdateInput, AdminIamUserUpdateInput, AdminModUserItem } from '~/composables/useAdmin'
import { formatBytes, formatDateTime, localizeI18nName } from '~/utils/format'

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
const banning = ref(false)
const modsPending = ref(false)
const modApplying = ref(false)
const modRemovingId = ref(0)
const sessionDeleting = ref(false)
const permissionSaving = ref(false)
const statPending = ref(false)
const statSaving = ref(false)
const selectedId = ref<number | null>(null)
const userStat = ref<AdminIamUserStat | null>(null)
const errorMessage = ref('')
const formError = ref('')
const banError = ref('')
const modsError = ref('')
const permissionError = ref('')
const statError = ref('')
const searchInput = ref('')

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

const banForm = reactive({
  reason: '',
  durationDays: 0
})

const permissionForm = reactive({
  permKey: '',
  isDeny: false
})

const modForm = reactive({
  type: 5,
  durationDays: 7,
  reason: ''
})

const statForm = reactive({
  uploadedDiff: '',
  downloadedDiff: '',
  bonusDiff: ''
})

const pageSizes = [20, 30, 50, 100]
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))
const roleMap = computed(() => new Map(roles.value.map((item) => [item.id, item])))
const selectedUser = computed(() => users.value.find((item) => item.id === selectedId.value) || null)
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
  return localizeI18nName(role.nameI18N as any, locale.value, `#${role.id}`)
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
  form.status = user.status
  form.role = user.role
  form.passkey = user.passkey || ''
  originalPasskey.value = user.passkey || ''
  formError.value = ''
  banError.value = ''
  modsError.value = ''
  permissionError.value = ''
  statError.value = ''
  resetStatForm()
  loadSelectedUserStat()
  loadSelectedUserMods()
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

async function deleteSessions() {
  const user = selectedUser.value
  if (!user) return
  if (!window.confirm(t('admin.iam.users.form.confirmKick', { name: user.username }))) return

  sessionDeleting.value = true
  try {
    await adminApi.deleteIamUserSessions(user.id)
    toast.add({ title: t('admin.iam.users.form.kicked') })
  } catch (error: unknown) {
    toast.add({
      color: 'error',
      title: error instanceof ApiError ? error.message : t('common.requestFailed')
    })
  } finally {
    sessionDeleting.value = false
  }
}

async function banUser() {
  const user = selectedUser.value
  if (!user || !banForm.reason.trim()) return

  banning.value = true
  banError.value = ''
  try {
    await adminApi.banIamUser(user.id, {
      reason: banForm.reason.trim(),
      durationDays: Number(banForm.durationDays || 0)
    })
    toast.add({ title: t('admin.iam.users.ban.success') })
    banForm.reason = ''
    banForm.durationDays = 0
    await loadSelectedUserMods()
  } catch (error: unknown) {
    banError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    banning.value = false
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
  } catch (error: unknown) {
    modsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    modApplying.value = false
  }
}

async function removeUserMod(item: AdminModUserItem) {
  const user = selectedUser.value
  if (!user || !item.is_active) return
  if (!window.confirm(t('admin.iam.users.mod.confirmRemove', { type: modTypeLabel(item.mod_type) }))) return

  modRemovingId.value = item.id
  modsError.value = ''
  try {
    await adminApi.removeUserMod(user.id, item.id)
    toast.add({ title: t('admin.iam.users.mod.removed'), color: 'success', icon: 'i-lucide-check-circle' })
    await loadSelectedUserMods()
  } catch (error: unknown) {
    modsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    modRemovingId.value = 0
  }
}

async function applyPermission(action: 'grant' | 'revoke') {
  const user = selectedUser.value
  const permKey = permissionForm.permKey.trim()
  if (!user || !permKey) return

  permissionSaving.value = true
  permissionError.value = ''
  try {
    const input = {
      permKey,
      isDeny: permissionForm.isDeny
    }
    if (action === 'grant') {
      await adminApi.grantIamUserPermission(user.id, input)
      toast.add({ title: t('admin.iam.users.permission.granted') })
    } else {
      await adminApi.revokeIamUserPermission(user.id, input)
      toast.add({ title: t('admin.iam.users.permission.revoked') })
    }
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
