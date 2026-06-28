<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-8 dark:bg-slate-950">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mb-6 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-slate-500 dark:text-slate-400">{{ $t('user.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold text-slate-950 dark:text-white">{{ displayName }}</h1>
        </div>
        <div class="flex items-center gap-2">
          <UBadge :color="user?.isStaff ? 'primary' : 'neutral'" variant="soft">
            {{ roleDisplayName }}
          </UBadge>
          <UBadge color="success" variant="soft">
            {{ $t('user.status.normal') }}
          </UBadge>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        <UCard v-for="item in statCards" :key="item.label" class="rounded-lg">
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-sm text-slate-500 dark:text-slate-400">{{ item.label }}</p>
              <p class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ item.value }}</p>
            </div>
            <span class="flex size-10 items-center justify-center rounded-md" :class="item.iconClass">
              <UIcon :name="item.icon" class="size-5" />
            </span>
          </div>
        </UCard>
      </div>

      <div class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-[360px_1fr]">
        <div class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <div class="flex items-center gap-3">
                <UAvatar :src="user?.avatar || undefined" :alt="displayName" size="lg" />
                <div class="min-w-0">
                  <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ displayName }}</p>
                  <p class="truncate text-xs text-slate-500 dark:text-slate-400">{{ user?.email }}</p>
                </div>
              </div>
            </template>

            <dl class="space-y-3 text-sm">
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('user.fields.id') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">#{{ user?.id }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('user.fields.role') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ roleDisplayName }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('user.fields.createdAt') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ formatDateTime(user?.createdAt) }}</dd>
              </div>
            </dl>
          </UCard>

          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.passkey.title') }}</h2>
            </template>

            <div class="rounded-md border border-slate-200 bg-slate-50 p-3 dark:border-slate-800 dark:bg-slate-950">
              <p class="break-all font-mono text-xs text-slate-700 dark:text-slate-300">{{ user?.passkey || '-' }}</p>
            </div>
            <UButton
              class="mt-3"
              color="neutral"
              variant="outline"
              icon="i-lucide-copy"
              block
              :disabled="!user?.passkey"
              @click="copyPasskey"
            >
              {{ $t('common.copy') }}
            </UButton>
          </UCard>

          <UCard class="rounded-lg">
            <template #header>
              <div class="flex items-start justify-between gap-3">
                <div>
                  <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.invites.title') }}</h2>
                  <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                    {{ $t('user.invites.summary', { available: numberFormatter.format(availableInvites.length), total: numberFormatter.format(inviteTotal) }) }}
                  </p>
                </div>
                <UButton color="neutral" variant="ghost" size="sm" icon="i-lucide-refresh-cw" :loading="invitesPending" :aria-label="$t('common.refresh')" @click="loadInvites" />
              </div>
            </template>

            <form class="space-y-3" @submit.prevent="handleInviteSend">
              <UFormField :label="$t('user.invites.email')" required>
                <UInput v-model="inviteEmail" class="w-full" type="email" icon="i-lucide-mail" :disabled="inviteSendPending || availableInvites.length === 0" />
              </UFormField>
              <UFormField :label="$t('user.invites.code')" required>
                <select
                  v-model="selectedInviteHash"
                  class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="inviteSendPending || availableInvites.length === 0"
                >
                  <option value="">{{ $t('user.invites.selectCode') }}</option>
                  <option v-for="invite in availableInvites" :key="invite.id" :value="invite.hash">
                    {{ invite.hash }}
                  </option>
                </select>
              </UFormField>
              <UButton type="submit" color="primary" icon="i-lucide-send" block :loading="inviteSendPending" :disabled="!canSendInvite">
                {{ $t('user.invites.send') }}
              </UButton>
            </form>

            <div v-if="invitesError" class="mt-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
              {{ invitesError }}
            </div>
            <div v-else-if="invitesPending && invites.length === 0" class="mt-4 space-y-2">
              <div v-for="item in 3" :key="item" class="h-12 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>
            <div v-else-if="invites.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-3 py-8 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
              {{ $t('user.invites.empty') }}
            </div>
            <div v-else class="mt-4 divide-y divide-slate-100 dark:divide-slate-800">
              <div v-for="invite in invites" :key="invite.id" class="py-3 text-sm">
                <div class="flex items-center justify-between gap-3">
                  <p class="min-w-0 truncate font-mono text-xs text-slate-700 dark:text-slate-300">{{ invite.hash }}</p>
                  <UBadge :color="inviteStatusColor(invite.status)" variant="soft">
                    {{ inviteStatusLabel(invite.status) }}
                  </UBadge>
                </div>
                <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
                  {{ invite.inviteeEmail || invite.inviteeName || '-' }}
                </p>
              </div>
            </div>
          </UCard>
        </div>

        <div class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.trafficHistory.title') }}</h2>
                  <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.trafficHistory.subtitle') }}</p>
                </div>
                <div class="flex flex-wrap items-center gap-2">
                  <div class="flex rounded-md border border-slate-200 bg-white p-1 dark:border-slate-800 dark:bg-slate-950">
                    <button
                      v-for="period in trafficHistoryPeriods"
                      :key="period.value"
                      type="button"
                      class="inline-flex h-8 items-center gap-1.5 rounded px-3 text-sm font-medium transition"
                      :class="trafficHistoryPeriod === period.value
                        ? 'bg-slate-950 text-white dark:bg-white dark:text-slate-950'
                        : 'text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-900'"
                      :disabled="trafficHistoryPending"
                      @click="setTrafficHistoryPeriod(period.value)"
                    >
                      <UIcon :name="period.icon" class="size-4" />
                      {{ period.label }}
                    </button>
                  </div>
                  <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" :loading="trafficHistoryPending" @click="loadTrafficHistory">
                    {{ $t('common.refresh') }}
                  </UButton>
                </div>
              </div>
            </template>

            <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
              <div v-for="item in trafficHistorySummaryCards" :key="item.label" class="rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ item.label }}</p>
                <p class="mt-1 truncate text-sm font-semibold text-slate-950 dark:text-white">{{ item.value }}</p>
              </div>
            </div>

            <div class="mt-4 flex flex-wrap gap-2">
              <UButton
                v-for="mode in trafficHistoryModes"
                :key="mode.value"
                color="neutral"
                :variant="trafficChartMode === mode.value ? 'soft' : 'ghost'"
                size="sm"
                :icon="mode.icon"
                @click="trafficChartMode = mode.value"
              >
                {{ mode.label }}
              </UButton>
            </div>

            <div v-if="trafficHistoryError" class="mt-4 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
              <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
              <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ trafficHistoryError }}</p>
              <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadTrafficHistory">
                {{ $t('common.retry') }}
              </UButton>
            </div>
            <div v-else-if="trafficHistoryPending && trafficHistory.length === 0" class="mt-4 h-72 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            <div v-else-if="trafficChartItems.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-10 text-center dark:border-slate-800">
              <UIcon name="i-lucide-chart-no-axes-column" class="mx-auto size-8 text-slate-400" />
              <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.trafficHistory.empty') }}</p>
            </div>
            <div v-else class="mt-4">
              <div class="mb-3 flex flex-wrap items-center gap-x-5 gap-y-2 text-xs text-slate-500 dark:text-slate-400">
                <span v-for="series in trafficChartSeries" :key="series.key" class="inline-flex items-center gap-2">
                  <span class="size-2 rounded-full" :style="{ backgroundColor: series.color }" />
                  {{ series.label }}
                </span>
              </div>

              <div class="overflow-hidden rounded-md border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-950">
                <svg class="h-72 w-full" viewBox="0 0 1000 280" role="img" :aria-label="$t('user.trafficHistory.title')">
                  <g v-for="line in trafficChartGridLines" :key="line.y">
                    <line :x1="trafficChartPadding.left" :x2="1000 - trafficChartPadding.right" :y1="line.y" :y2="line.y" class="stroke-slate-100 dark:stroke-slate-800" stroke-width="1" />
                    <text :x="trafficChartPadding.left - 10" :y="line.y + 4" text-anchor="end" class="fill-slate-400 text-[11px]">{{ line.label }}</text>
                  </g>

                  <g v-for="series in trafficChartSeries" :key="series.key">
                    <polyline :points="series.polyline" fill="none" :stroke="series.color" stroke-width="4" stroke-linecap="round" stroke-linejoin="round" />
                    <circle v-for="point in series.points" :key="`${series.key}-${point.date}`" :cx="point.x" :cy="point.y" r="4" :fill="series.color" />
                  </g>

                  <g v-for="label in trafficChartXAxisLabels" :key="label.date">
                    <text :x="label.x" y="264" text-anchor="middle" class="fill-slate-400 text-[11px]">{{ label.date }}</text>
                  </g>
                </svg>
              </div>
            </div>
          </UCard>

          <UCard class="rounded-lg">
            <template #header>
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.bonusLogs.title') }}</h2>
                  <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                    {{ $t('user.bonusLogs.summary', { total: numberFormatter.format(bonusLogTotal), hourly: formatBonus(hourlyBonus) }) }}
                  </p>
                </div>
                <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" :loading="bonusLogsPending" @click="loadBonusLogs">
                  {{ $t('common.refresh') }}
                </UButton>
              </div>
            </template>

            <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div class="rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.hourly') }}</p>
                <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ formatBonus(hourlyBonus) }}</p>
              </div>
              <div class="rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.balance') }}</p>
                <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ formatBonus(user?.bonus || 0) }}</p>
              </div>
            </div>

            <div v-if="bonusLogsError" class="mt-4 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
              <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
              <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ bonusLogsError }}</p>
              <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadBonusLogs">
                {{ $t('common.retry') }}
              </UButton>
            </div>
            <div v-else-if="bonusLogsPending && bonusLogs.length === 0" class="mt-4 space-y-3">
              <div v-for="item in 4" :key="item" class="h-14 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>
            <div v-else-if="bonusLogs.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-10 text-center dark:border-slate-800">
              <UIcon name="i-lucide-coins" class="mx-auto size-8 text-slate-400" />
              <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.empty') }}</p>
            </div>
            <div v-else class="mt-4 divide-y divide-slate-100 dark:divide-slate-800">
              <article v-for="log in bonusLogs" :key="log.id" class="grid gap-3 py-3 text-sm sm:grid-cols-[minmax(0,1fr)_120px] sm:items-center">
                <div class="min-w-0">
                  <div class="flex min-w-0 items-center gap-2">
                    <UBadge :color="log.amount >= 0 ? 'success' : 'error'" variant="soft">
                      {{ log.amount >= 0 ? '+' : '' }}{{ formatBonus(log.amount) }}
                    </UBadge>
                    <p class="truncate font-medium text-slate-950 dark:text-white">{{ bonusActionLabel(log.action) }}</p>
                  </div>
                  <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
                    {{ log.remark || `${log.targetType || '-'} #${log.targetId || '-'}` }}
                  </p>
                </div>
                <div class="text-xs text-slate-500 sm:text-right dark:text-slate-400">
                  <p>{{ formatDateTime(log.createdAt, locale) }}</p>
                  <p class="mt-1">{{ $t('user.bonusLogs.after', { balance: formatBonus(log.balanceAfter) }) }}</p>
                </div>
              </article>
            </div>

            <div v-if="bonusLogTotalPages > 1" class="mt-4 flex flex-col gap-3 border-t border-slate-200 pt-4 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
              <p class="text-sm text-slate-500 dark:text-slate-400">
                {{ $t('user.bonusLogs.pageSummary', { page: bonusLogPage, pages: bonusLogTotalPages }) }}
              </p>
              <div class="flex items-center gap-2">
                <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-chevron-left" :disabled="bonusLogPage <= 1 || bonusLogsPending" @click="goToBonusLogPage(bonusLogPage - 1)">
                  {{ $t('common.previous') }}
                </UButton>
                <UButton color="neutral" variant="outline" size="sm" trailing-icon="i-lucide-chevron-right" :disabled="bonusLogPage >= bonusLogTotalPages || bonusLogsPending" @click="goToBonusLogPage(bonusLogPage + 1)">
                  {{ $t('common.next') }}
                </UButton>
              </div>
            </div>
          </UCard>

          <UCard class="rounded-lg">
            <template #header>
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.snatches.title') }}</h2>
                  <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                    {{ $t('user.snatches.summary', { total: numberFormatter.format(snatchTotal) }) }}
                  </p>
                </div>
                <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" :loading="snatchesPending" @click="loadSnatches">
                  {{ $t('common.refresh') }}
                </UButton>
              </div>
            </template>

            <div class="flex flex-wrap gap-2">
              <UButton
                v-for="filter in snatchFilters"
                :key="filter.value"
                color="neutral"
                :variant="snatchStatus === filter.value ? 'soft' : 'ghost'"
                size="sm"
                :icon="filter.icon"
                :disabled="snatchesPending"
                @click="setSnatchStatus(filter.value)"
              >
                {{ filter.label }}
              </UButton>
            </div>

            <div v-if="snatchesError" class="mt-4 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
              <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
              <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ snatchesError }}</p>
              <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadSnatches">
                {{ $t('common.retry') }}
              </UButton>
            </div>
            <div v-else-if="snatchesPending && snatches.length === 0" class="mt-4 space-y-3">
              <div v-for="item in 4" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>
            <div v-else-if="snatches.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-10 text-center dark:border-slate-800">
              <UIcon name="i-lucide-inbox" class="mx-auto size-8 text-slate-400" />
              <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.snatches.empty') }}</p>
            </div>
            <div v-else class="mt-4 divide-y divide-slate-100 dark:divide-slate-800">
              <article v-for="snatch in snatches" :key="snatch.id" class="py-4">
                <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_280px] xl:items-start">
                  <div class="min-w-0">
                    <div class="flex flex-wrap items-center gap-2">
                      <UBadge :color="snatch.isFinished ? 'success' : 'primary'" variant="soft">
                        {{ snatch.isFinished ? $t('user.snatches.finished') : $t('user.snatches.unfinished') }}
                      </UBadge>
                      <span class="text-xs text-slate-500 dark:text-slate-400">#{{ snatch.torrentId }}</span>
                    </div>
                    <h3 class="mt-2 truncate text-sm font-semibold">
                      <NuxtLink
                        :to="localePath(`/catalog/torrents/${snatch.torrentId}`)"
                        class="text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300"
                      >
                        {{ snatch.torrentName || $t('catalog.torrents.detail.titleFallback', { id: snatch.torrentId }) }}
                      </NuxtLink>
                    </h3>
                    <p class="mt-2 text-xs text-slate-500 dark:text-slate-400">
                      {{ $t('user.snatches.lastAction') }} {{ formatDateTime(snatch.lastActionAt, locale) }}
                    </p>
                  </div>

                  <div class="grid grid-cols-2 gap-x-4 gap-y-3 text-xs sm:grid-cols-4 xl:grid-cols-2">
                    <div>
                      <p class="text-slate-500 dark:text-slate-400">{{ $t('user.snatches.uploaded') }}</p>
                      <p class="mt-1 font-medium text-slate-950 dark:text-white">{{ formatBytes(snatch.uploaded) }}</p>
                    </div>
                    <div>
                      <p class="text-slate-500 dark:text-slate-400">{{ $t('user.snatches.downloaded') }}</p>
                      <p class="mt-1 font-medium text-slate-950 dark:text-white">{{ formatBytes(snatch.downloaded) }}</p>
                    </div>
                    <div>
                      <p class="text-slate-500 dark:text-slate-400">{{ $t('user.snatches.seedTime') }}</p>
                      <p class="mt-1 font-medium text-slate-950 dark:text-white">{{ formatDuration(snatch.seedTime) }}</p>
                    </div>
                    <div>
                      <p class="text-slate-500 dark:text-slate-400">{{ $t('user.snatches.leechTime') }}</p>
                      <p class="mt-1 font-medium text-slate-950 dark:text-white">{{ formatDuration(snatch.leechTime) }}</p>
                    </div>
                  </div>
                </div>
              </article>
            </div>

            <div v-if="snatchTotalPages > 1" class="mt-4 flex flex-col gap-3 border-t border-slate-200 pt-4 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
              <p class="text-sm text-slate-500 dark:text-slate-400">
                {{ $t('user.snatches.pageSummary', { page: snatchPage, pages: snatchTotalPages }) }}
              </p>
              <div class="flex items-center gap-2">
                <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-chevron-left" :disabled="snatchPage <= 1 || snatchesPending" @click="goToSnatchPage(snatchPage - 1)">
                  {{ $t('common.previous') }}
                </UButton>
                <UButton color="neutral" variant="outline" size="sm" trailing-icon="i-lucide-chevron-right" :disabled="snatchPage >= snatchTotalPages || snatchesPending" @click="goToSnatchPage(snatchPage + 1)">
                  {{ $t('common.next') }}
                </UButton>
              </div>
            </div>
          </UCard>

          <UCard class="rounded-lg">
            <template #header>
              <div>
                <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.profile.title') }}</h2>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.profile.subtitle') }}</p>
              </div>
            </template>

            <form class="grid grid-cols-1 gap-4" @submit.prevent="handleProfileSave">
              <UFormField :label="$t('user.profile.avatar')">
                <UInput
                  v-model="profileForm.avatar"
                  class="w-full"
                  icon="i-lucide-image"
                  :disabled="profilePending || loading"
                />
              </UFormField>

              <UFormField :label="$t('user.profile.signature')">
                <UInput
                  v-model="profileForm.signature"
                  class="w-full"
                  icon="i-lucide-pen-line"
                  :disabled="profilePending || loading"
                />
              </UFormField>

              <UFormField :label="$t('user.profile.info')">
                <UTextarea
                  v-model="profileForm.info"
                  class="w-full"
                  :rows="5"
                  :disabled="profilePending || loading"
                />
              </UFormField>

              <div class="flex justify-end">
                <UButton type="submit" color="primary" :loading="profilePending" :disabled="loading">
                  {{ $t('common.save') }}
                </UButton>
              </div>
            </form>
          </UCard>

          <UCard class="rounded-lg">
            <template #header>
              <div>
                <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.security.title') }}</h2>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.security.subtitle') }}</p>
              </div>
            </template>

            <form class="grid grid-cols-1 gap-4 md:grid-cols-2" @submit.prevent="handlePasswordChange">
              <UFormField :label="$t('user.security.oldPassword')" required>
                <UInput
                  v-model="passwordForm.oldPassword"
                  class="w-full"
                  type="password"
                  autocomplete="current-password"
                  :disabled="passwordPending"
                />
              </UFormField>

              <UFormField :label="$t('user.security.newPassword')" required>
                <UInput
                  v-model="passwordForm.newPassword"
                  class="w-full"
                  type="password"
                  autocomplete="new-password"
                  :disabled="passwordPending"
                />
              </UFormField>

              <div class="md:col-span-2 flex justify-end">
                <UButton type="submit" color="neutral" variant="outline" :loading="passwordPending" :disabled="!canChangePassword">
                  {{ $t('user.security.submit') }}
                </UButton>
              </div>
            </form>
          </UCard>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAccounting } from '~/composables/useAccounting'
import { ApiError } from '~/composables/useApi'
import { useEconomy, type BonusLogItem } from '~/composables/useEconomy'
import { useInvites, type InviteItem } from '~/composables/useInvites'
import type { SnatchItem, TrafficHistoryItem, TrafficHistoryPeriod, TrafficSummary } from '~/composables/useAccounting'

type SnatchStatus = 'all' | 'finished' | 'unfinished'
type TrafficChartMode = 'traffic' | 'time' | 'bonus'

interface TrafficChartPoint {
  date: string
  x: number
  y: number
  value: number
}

interface TrafficChartLabel {
  date: string
  x: number
}

interface TrafficChartSeries {
  key: string
  label: string
  color: string
  points: TrafficChartPoint[]
  polyline: string
}

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const toast = useToast()
const { user, fetchUser, updateProfile, changePassword } = useAuth()
const accounting = useAccounting()
const economy = useEconomy()
const inviteService = useInvites()

const loading = ref(true)
const traffic = ref<TrafficSummary | null>(null)
const trafficHistory = ref<TrafficHistoryItem[]>([])
const trafficHistoryPeriod = ref<TrafficHistoryPeriod>('daily')
const trafficChartMode = ref<TrafficChartMode>('traffic')
const snatches = ref<SnatchItem[]>([])
const bonusLogs = ref<BonusLogItem[]>([])
const invites = ref<InviteItem[]>([])
const snatchTotal = ref(0)
const snatchPage = ref(1)
const snatchSize = 10
const bonusLogTotal = ref(0)
const bonusLogPage = ref(1)
const bonusLogSize = 10
const inviteTotal = ref(0)
const invitePage = ref(1)
const inviteSize = 10
const hourlyBonus = ref(0)
const snatchStatus = ref<SnatchStatus>('all')
const profilePending = ref(false)
const passwordPending = ref(false)
const trafficHistoryPending = ref(false)
const trafficHistoryError = ref('')
const snatchesPending = ref(false)
const snatchesError = ref('')
const bonusLogsPending = ref(false)
const bonusLogsError = ref('')
const invitesPending = ref(false)
const invitesError = ref('')
const inviteSendPending = ref(false)
const inviteEmail = ref('')
const selectedInviteHash = ref('')

const profileForm = reactive({
  avatar: '',
  signature: '',
  info: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: ''
})

const displayName = computed(() => user.value?.username || t('user.fallbackName'))
const roleDisplayName = computed(() => user.value?.roleName || (user.value?.isStaff ? t('user.staff') : t('user.member')))
const trafficUploaded = computed(() => traffic.value?.uploaded ?? user.value?.uploaded ?? 0)
const trafficDownloaded = computed(() => traffic.value?.downloaded ?? user.value?.downloaded ?? 0)
const ratio = computed(() => traffic.value?.shareRatio ?? user.value?.shareRatio ?? 0)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const snatchTotalPages = computed(() => Math.max(1, Math.ceil(snatchTotal.value / snatchSize)))
const bonusLogTotalPages = computed(() => Math.max(1, Math.ceil(bonusLogTotal.value / bonusLogSize)))
const availableInvites = computed(() => invites.value.filter((invite) => invite.status === 0))
const canSendInvite = computed(() => {
  return Boolean(selectedInviteHash.value && inviteEmail.value.trim() && !inviteSendPending.value && availableInvites.value.length > 0)
})
const trafficChartPadding = {
  top: 24,
  right: 28,
  bottom: 44,
  left: 64
}
const trafficChartWidth = 1000
const trafficChartHeight = 280
const trafficChartInnerWidth = trafficChartWidth - trafficChartPadding.left - trafficChartPadding.right
const trafficChartInnerHeight = trafficChartHeight - trafficChartPadding.top - trafficChartPadding.bottom

const trafficHistoryPeriods = computed<Array<{ value: TrafficHistoryPeriod, label: string, icon: string }>>(() => [
  {
    value: 'daily',
    label: t('user.trafficHistory.periods.daily'),
    icon: 'i-lucide-calendar-days'
  },
  {
    value: 'monthly',
    label: t('user.trafficHistory.periods.monthly'),
    icon: 'i-lucide-calendar-range'
  }
])

const trafficHistoryModes = computed<Array<{ value: TrafficChartMode, label: string, icon: string }>>(() => [
  {
    value: 'traffic',
    label: t('user.trafficHistory.modes.traffic'),
    icon: 'i-lucide-arrow-up-down'
  },
  {
    value: 'time',
    label: t('user.trafficHistory.modes.time'),
    icon: 'i-lucide-timer'
  },
  {
    value: 'bonus',
    label: t('user.trafficHistory.modes.bonus'),
    icon: 'i-lucide-coins'
  }
])

const trafficChartItems = computed(() => {
  const limit = trafficHistoryPeriod.value === 'daily' ? 30 : 12
  return [...trafficHistory.value]
    .sort((a, b) => a.date.localeCompare(b.date))
    .slice(-limit)
})

const trafficHistoryTotals = computed(() => {
  return trafficChartItems.value.reduce(
    (total, item) => {
      total.uploaded += Number(item.uploaded || 0)
      total.downloaded += Number(item.downloaded || 0)
      total.activeTime += Number(item.seedTime || 0) + Number(item.leechTime || 0)
      total.bonus += Number(item.bonus || 0)
      return total
    },
    {
      uploaded: 0,
      downloaded: 0,
      activeTime: 0,
      bonus: 0
    }
  )
})

const trafficHistorySummaryCards = computed(() => [
  {
    label: t('user.trafficHistory.uploaded'),
    value: formatBytes(trafficHistoryTotals.value.uploaded)
  },
  {
    label: t('user.trafficHistory.downloaded'),
    value: formatBytes(trafficHistoryTotals.value.downloaded)
  },
  {
    label: t('user.trafficHistory.activeTime'),
    value: formatDuration(trafficHistoryTotals.value.activeTime)
  },
  {
    label: t('user.trafficHistory.bonus'),
    value: formatBonus(trafficHistoryTotals.value.bonus)
  }
])

const trafficChartMax = computed(() => {
  const values = trafficChartItems.value.flatMap((item) => {
    if (trafficChartMode.value === 'time') {
      return [Number(item.seedTime || 0), Number(item.leechTime || 0)]
    }
    if (trafficChartMode.value === 'bonus') {
      return [Number(item.bonus || 0)]
    }
    return [Number(item.uploaded || 0), Number(item.downloaded || 0)]
  })
  return Math.max(1, ...values)
})
const trafficChartSeries = computed<TrafficChartSeries[]>(() => {
  if (trafficChartMode.value === 'time') {
    return [
      buildTrafficChartSeries('seedTime', t('user.trafficHistory.seedTime'), '#10b981'),
      buildTrafficChartSeries('leechTime', t('user.trafficHistory.leechTime'), '#8b5cf6')
    ]
  }

  if (trafficChartMode.value === 'bonus') {
    return [
      buildTrafficChartSeries('bonus', t('user.trafficHistory.bonus'), '#f59e0b')
    ]
  }

  return [
    buildTrafficChartSeries('uploaded', t('user.trafficHistory.uploaded'), '#10b981'),
    buildTrafficChartSeries('downloaded', t('user.trafficHistory.downloaded'), '#0ea5e9')
  ]
})
const trafficChartGridLines = computed(() => {
  return [0, 0.25, 0.5, 0.75, 1].map((ratio) => {
    const value = trafficChartMax.value * (1 - ratio)
    return {
      y: trafficChartPadding.top + ratio * trafficChartInnerHeight,
      label: formatTrafficChartValue(value)
    }
  })
})
const trafficChartXAxisLabels = computed<TrafficChartLabel[]>(() => {
  const items = trafficChartItems.value
  if (items.length === 0) return []

  const indexes = new Set([0, Math.floor((items.length - 1) / 2), items.length - 1])
  return [...indexes].map((index) => ({
    date: items[index].date,
    x: trafficChartX(index, items.length)
  }))
})

const snatchFilters = computed<Array<{ value: SnatchStatus, label: string, icon: string }>>(() => [
  {
    value: 'all',
    label: t('user.snatches.filters.all'),
    icon: 'i-lucide-list'
  },
  {
    value: 'finished',
    label: t('user.snatches.filters.finished'),
    icon: 'i-lucide-circle-check'
  },
  {
    value: 'unfinished',
    label: t('user.snatches.filters.unfinished'),
    icon: 'i-lucide-clock'
  }
])

const statCards = computed(() => [
  {
    label: t('user.stats.uploaded'),
    value: formatBytes(trafficUploaded.value),
    icon: 'i-lucide-upload',
    iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300'
  },
  {
    label: t('user.stats.downloaded'),
    value: formatBytes(trafficDownloaded.value),
    icon: 'i-lucide-download',
    iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
  },
  {
    label: t('user.stats.ratio'),
    value: Number.isFinite(ratio.value) ? ratio.value.toFixed(2) : '-',
    icon: 'i-lucide-scale',
    iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300'
  },
  {
    label: t('user.stats.bonus'),
    value: String(user.value?.bonus ?? 0),
    icon: 'i-lucide-coins',
    iconClass: 'bg-violet-100 text-violet-700 dark:bg-violet-950 dark:text-violet-300'
  }
])

const canChangePassword = computed(() => passwordForm.oldPassword.length > 0 && passwordForm.newPassword.length >= 6)

watch(
  user,
  (value) => {
    profileForm.avatar = value?.avatar || ''
    profileForm.signature = value?.signature || ''
    profileForm.info = value?.info || ''
  },
  { immediate: true }
)

onMounted(async () => {
  loading.value = true

  try {
    await fetchUser()
    await Promise.all([
      loadTraffic(),
      loadTrafficHistory(),
      loadSnatches(),
      loadHourlyBonus(),
      loadBonusLogs(),
      loadInvites()
    ])
  } catch {
    await navigateTo(localePath('/login'))
  } finally {
    loading.value = false
  }
})

async function loadTraffic() {
  try {
    traffic.value = await accounting.getTraffic()
  } catch {
    traffic.value = null
  }
}

async function loadTrafficHistory() {
  if (trafficHistoryPending.value) return

  trafficHistoryPending.value = true
  trafficHistoryError.value = ''

  try {
    const data = await accounting.listTrafficHistory({
      period: trafficHistoryPeriod.value
    })
    trafficHistory.value = data.list || []
  } catch (error) {
    trafficHistory.value = []
    trafficHistoryError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    trafficHistoryPending.value = false
  }
}

async function loadHourlyBonus() {
  try {
    const data = await economy.getHourlyBonus()
    hourlyBonus.value = data.hourlyBonus || 0
  } catch {
    hourlyBonus.value = 0
  }
}

async function loadBonusLogs() {
  if (bonusLogsPending.value) return

  bonusLogsPending.value = true
  bonusLogsError.value = ''

  try {
    let data = await economy.listBonusLogs({
      page: bonusLogPage.value,
      size: bonusLogSize
    })

    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / bonusLogSize))
    if (bonusLogPage.value > nextTotalPages) {
      bonusLogPage.value = nextTotalPages
      data = await economy.listBonusLogs({
        page: bonusLogPage.value,
        size: bonusLogSize
      })
    }

    bonusLogs.value = data.list || []
    bonusLogTotal.value = data.total || 0
  } catch (error) {
    bonusLogs.value = []
    bonusLogTotal.value = 0
    bonusLogsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    bonusLogsPending.value = false
  }
}

function goToBonusLogPage(page: number) {
  bonusLogPage.value = Math.min(Math.max(1, page), bonusLogTotalPages.value)
  loadBonusLogs()
}

async function loadInvites() {
  if (invitesPending.value) return

  invitesPending.value = true
  invitesError.value = ''

  try {
    const data = await inviteService.listInvites({
      page: invitePage.value,
      size: inviteSize
    })
    invites.value = data.list || []
    inviteTotal.value = data.total || 0

    if (!selectedInviteHash.value || !availableInvites.value.some((invite) => invite.hash === selectedInviteHash.value)) {
      selectedInviteHash.value = availableInvites.value[0]?.hash || ''
    }
  } catch (error) {
    invites.value = []
    inviteTotal.value = 0
    selectedInviteHash.value = ''
    invitesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    invitesPending.value = false
  }
}

async function handleInviteSend() {
  if (!canSendInvite.value) return

  inviteSendPending.value = true
  try {
    await inviteService.sendInvite(selectedInviteHash.value, inviteEmail.value.trim())
    toast.add({
      title: t('user.invites.sent'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    inviteEmail.value = ''
    await Promise.all([loadInvites(), fetchUser()])
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    inviteSendPending.value = false
  }
}

function setTrafficHistoryPeriod(period: TrafficHistoryPeriod) {
  if (trafficHistoryPeriod.value === period) return

  trafficHistoryPeriod.value = period
  loadTrafficHistory()
}

async function loadSnatches() {
  if (snatchesPending.value) return

  snatchesPending.value = true
  snatchesError.value = ''

  try {
    let data = await accounting.listSnatches({
      page: snatchPage.value,
      size: snatchSize,
      isFinished: snatchStatus.value === 'all' ? undefined : snatchStatus.value === 'finished'
    })

    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / snatchSize))
    if (snatchPage.value > nextTotalPages) {
      snatchPage.value = nextTotalPages
      data = await accounting.listSnatches({
        page: snatchPage.value,
        size: snatchSize,
        isFinished: snatchStatus.value === 'all' ? undefined : snatchStatus.value === 'finished'
      })
    }

    snatches.value = data.list || []
    snatchTotal.value = data.total || 0
  } catch (error) {
    snatches.value = []
    snatchTotal.value = 0
    snatchesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    snatchesPending.value = false
  }
}

function setSnatchStatus(status: SnatchStatus) {
  if (snatchStatus.value === status) return

  snatchStatus.value = status
  snatchPage.value = 1
  loadSnatches()
}

function goToSnatchPage(page: number) {
  snatchPage.value = Math.min(Math.max(1, page), snatchTotalPages.value)
  loadSnatches()
}

async function handleProfileSave() {
  profilePending.value = true

  try {
    await updateProfile({
      avatar: profileForm.avatar.trim(),
      signature: profileForm.signature.trim(),
      info: profileForm.info.trim()
    })
    toast.add({
      title: t('user.profile.saved'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    profilePending.value = false
  }
}

async function handlePasswordChange() {
  if (!canChangePassword.value) return

  passwordPending.value = true

  try {
    await changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })
    toast.add({
      title: t('user.security.changed'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await navigateTo(localePath('/login'))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    passwordPending.value = false
  }
}

async function copyPasskey() {
  if (!user.value?.passkey || !navigator?.clipboard) return

  await navigator.clipboard.writeText(user.value.passkey)
  toast.add({
    title: t('user.passkey.copied'),
    color: 'success',
    icon: 'i-lucide-check-circle'
  })
}

function formatDuration(value?: number | null) {
  const seconds = Math.max(0, Number(value || 0))
  if (seconds < 60) {
    return t('user.duration.seconds', { count: numberFormatter.value.format(Math.floor(seconds)) })
  }

  const minutes = seconds / 60
  if (minutes < 60) {
    return t('user.duration.minutes', { count: numberFormatter.value.format(Math.floor(minutes)) })
  }

  const hours = minutes / 60
  if (hours < 24) {
    return t('user.duration.hours', { count: numberFormatter.value.format(Number(hours.toFixed(hours >= 10 ? 0 : 1))) })
  }

  const days = hours / 24
  return t('user.duration.days', { count: numberFormatter.value.format(Number(days.toFixed(days >= 10 ? 0 : 1))) })
}

function formatBonus(value?: number | null) {
  return numberFormatter.value.format(Number(value || 0))
}

function bonusActionLabel(action?: string | null) {
  if (!action) return '-'
  const key = `user.bonusLogs.actions.${action}`
  const translated = t(key)
  return translated === key ? action : translated
}

function inviteStatusLabel(status: number) {
  const key = `user.invites.status.${status}`
  const translated = t(key)
  return translated === key ? String(status) : translated
}

function inviteStatusColor(status: number) {
  if (status === 0) return 'success'
  if (status === 1) return 'primary'
  if (status === 2) return 'neutral'
  if (status === 3) return 'error'
  return 'warning'
}

function formatTrafficChartValue(value: number) {
  if (trafficChartMode.value === 'time') return formatDuration(value)
  if (trafficChartMode.value === 'bonus') return formatBonus(value)
  return formatBytes(value)
}

function trafficChartX(index: number, total: number) {
  if (total <= 1) {
    return trafficChartPadding.left + trafficChartInnerWidth / 2
  }

  return trafficChartPadding.left + (index / (total - 1)) * trafficChartInnerWidth
}

function trafficChartY(value: number) {
  const ratio = Math.min(1, Math.max(0, Number(value || 0) / trafficChartMax.value))
  return trafficChartPadding.top + (1 - ratio) * trafficChartInnerHeight
}

function buildTrafficChartSeries(key: keyof Pick<TrafficHistoryItem, 'uploaded' | 'downloaded' | 'seedTime' | 'leechTime' | 'bonus'>, label: string, color: string): TrafficChartSeries {
  const points = buildTrafficChartPoints(key)
  return {
    key,
    label,
    color,
    points,
    polyline: points.map((point) => `${point.x},${point.y}`).join(' ')
  }
}

function buildTrafficChartPoints(key: keyof Pick<TrafficHistoryItem, 'uploaded' | 'downloaded' | 'seedTime' | 'leechTime' | 'bonus'>): TrafficChartPoint[] {
  return trafficChartItems.value.map((item, index) => {
    const value = Number(item[key] || 0)
    return {
      date: item.date,
      value,
      x: trafficChartX(index, trafficChartItems.value.length),
      y: trafficChartY(value)
    }
  })
}

useSeoMeta({
  title: t('user.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
