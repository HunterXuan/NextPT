<template>
  <div class="min-h-[calc(100vh-4rem)] overflow-x-hidden bg-slate-50 py-6 lg:overflow-x-visible dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div v-if="pending" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <div class="space-y-4">
          <div class="h-44 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
          <div class="h-64 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
        </div>
        <div class="h-80 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
      </div>

      <div v-else-if="errorMessage" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadPage">
          {{ $t('common.retry') }}
        </UButton>
      </div>

      <div v-else-if="torrent" class="grid min-w-0 gap-4 lg:grid-cols-[minmax(0,1fr)_44px_340px] lg:items-start">
        <main class="min-w-0 space-y-6">
          <section id="torrent-top" class="scroll-mt-24 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="p-4 sm:p-5">
              <div class="space-y-4">
                <div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
                  <div class="flex min-w-0 flex-wrap items-center gap-2">
                    <span class="inline-flex h-6 max-w-full items-center rounded border border-slate-200 bg-slate-50 px-2 text-xs font-medium text-slate-600 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-300">
                      <span class="truncate">{{ categoryName }}</span>
                    </span>
                    <span
                      v-for="badge in torrentStatusBadges(torrent)"
                      :key="badge.key"
                      class="inline-flex h-6 shrink-0 items-center rounded px-2 text-xs font-semibold leading-none"
                      :class="badge.class"
                      :title="badge.title"
                    >
                      {{ badge.label }}
                    </span>
                  </div>

                  <div class="flex shrink-0 flex-wrap items-center gap-2 lg:justify-end">
                    <UTooltip
                      :text="$t('catalog.torrents.detail.actions.download')"
                      :content="{ side: 'top', sideOffset: 8 }"
                      :delay-duration="120"
                    >
                      <UButton
                        class="shrink-0"
                        color="primary"
                        icon="i-lucide-download"
                        :loading="downloadPending"
                        :disabled="downloadPending"
                        :aria-label="$t('catalog.torrents.detail.actions.download')"
                        @click="handleDownload"
                      />
                    </UTooltip>
                    <UTooltip
                      :text="$t('catalog.torrents.detail.actions.like')"
                      :content="{ side: 'top', sideOffset: 8 }"
                      :delay-duration="120"
                    >
                      <UButton
                        :color="torrent.isLiked ? 'error' : 'neutral'"
                        :variant="torrent.isLiked ? 'soft' : 'outline'"
                        icon="i-lucide-heart"
                        :loading="likePending"
                        :disabled="likePending"
                        :aria-label="$t('catalog.torrents.detail.actions.like')"
                        @click="handleToggleLike"
                      >
                        {{ numberFormatter.format(torrent.likeCount || 0) }}
                      </UButton>
                    </UTooltip>
                    <UTooltip
                      :text="$t('catalog.torrents.detail.actions.bookmark')"
                      :content="{ side: 'top', sideOffset: 8 }"
                      :delay-duration="120"
                    >
                      <UButton
                        color="neutral"
                        :variant="torrent.isBookmarked ? 'soft' : 'outline'"
                        :icon="torrent.isBookmarked ? 'i-lucide-bookmark-check' : 'i-lucide-bookmark'"
                        :loading="bookmarkPending"
                        :disabled="bookmarkPending"
                        :aria-label="$t('catalog.torrents.detail.actions.bookmark')"
                        @click="handleToggleBookmark"
                      />
                    </UTooltip>
                  </div>
                </div>

                <div class="min-w-0">
                  <h1 class="break-words [overflow-wrap:anywhere] text-xl font-semibold leading-7 text-slate-950 md:text-2xl dark:text-white">
                    {{ torrent.name || $t('catalog.torrents.detail.titleFallback', { id: torrent.id }) }}
                  </h1>
                  <p v-if="torrent.subTitle" class="mt-2 break-words [overflow-wrap:anywhere] text-sm leading-6 text-slate-600 dark:text-slate-300">
                    {{ torrent.subTitle }}
                  </p>
                </div>
              </div>
            </div>

            <div class="border-t border-slate-200 bg-slate-50/70 dark:border-slate-800 dark:bg-slate-950/40">
              <div class="grid grid-cols-2 divide-x divide-y divide-slate-200 sm:grid-cols-4 sm:divide-y-0 dark:divide-slate-800">
                <button
                  type="button"
                  :class="fileStatCardClass()"
                  :aria-pressed="filePanelOpen"
                  @click="handleFileStatClick"
                >
                  <div class="flex min-w-0 items-center justify-between gap-3">
                    <div class="flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                      <UIcon name="i-lucide-folder-tree" class="size-3.5 shrink-0 text-amber-500" />
                      <span class="truncate">{{ $t('catalog.torrents.fileCount', { count: numberFormatter.format(torrent.fileCount) }) }}</span>
                    </div>
                    <p class="shrink-0 truncate text-right text-sm font-semibold text-slate-950 dark:text-white">
                      {{ formatBytes(torrent.size) }}
                    </p>
                  </div>
                </button>
                <button
                  type="button"
                  :class="peerStatCardClass('seeders')"
                  :aria-pressed="peerView === 'seeders'"
                  @click="handlePeerStatClick('seeders')"
                >
                  <div class="flex min-w-0 items-center justify-between gap-3">
                    <div class="flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                      <UIcon name="i-lucide-users" class="size-3.5 shrink-0 text-emerald-500" />
                      <span class="truncate">{{ $t('catalog.torrents.detail.stats.seeders') }}</span>
                    </div>
                    <p class="shrink-0 text-right text-sm font-semibold text-emerald-600 dark:text-emerald-400">{{ numberFormatter.format(torrent.seeders) }}</p>
                  </div>
                </button>
                <button
                  type="button"
                  :class="peerStatCardClass('leechers')"
                  :aria-pressed="peerView === 'leechers'"
                  @click="handlePeerStatClick('leechers')"
                >
                  <div class="flex min-w-0 items-center justify-between gap-3">
                    <div class="flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                      <UIcon name="i-lucide-loader-circle" class="size-3.5 shrink-0 text-sky-500" />
                      <span class="truncate">{{ $t('catalog.torrents.detail.stats.leechers') }}</span>
                    </div>
                    <p class="shrink-0 text-right text-sm font-semibold text-sky-600 dark:text-sky-400">{{ numberFormatter.format(torrent.leechers) }}</p>
                  </div>
                </button>
                <button
                  type="button"
                  :class="peerStatCardClass('completed')"
                  :aria-pressed="peerView === 'completed'"
                  @click="handlePeerStatClick('completed')"
                >
                  <div class="flex min-w-0 items-center justify-between gap-3">
                    <div class="flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                      <UIcon name="i-lucide-circle-check" class="size-3.5 shrink-0 text-slate-500 dark:text-slate-400" />
                      <span class="truncate">{{ $t('catalog.torrents.detail.stats.completed') }}</span>
                    </div>
                    <p class="shrink-0 text-right text-sm font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(torrent.snatched) }}</p>
                  </div>
                </button>
              </div>
            </div>

            <div v-if="filePanelOpen" class="border-t border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
              <div v-if="filesError" class="flex flex-col items-center justify-center px-4 py-10 text-center">
                <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
                <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ filesError }}</p>
                <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadFiles(true)">
                  {{ $t('common.retry') }}
                </UButton>
              </div>

              <div v-else-if="fileTree.length === 0" class="px-4 py-10 text-center text-sm text-slate-500 dark:text-slate-400">
                {{ $t('catalog.torrents.detail.files.empty') }}
              </div>

              <div v-else class="max-h-[520px] overflow-auto">
                <div
                  v-for="row in visibleFileTreeRows"
                  :key="row.node.id"
                  class="grid min-h-11 gap-2 border-b border-slate-100 px-4 py-2 text-sm last:border-b-0 sm:grid-cols-[minmax(0,1fr)_120px] sm:items-center dark:border-slate-800"
                  :class="row.node.id === fileTreeRootId ? 'bg-slate-50/70 dark:bg-slate-950/40' : ''"
                >
                  <button
                    v-if="row.node.type === 'directory'"
                    type="button"
                    class="flex min-w-0 items-center gap-2 rounded-md py-1 pr-2 text-left text-slate-700 transition-colors hover:bg-white dark:text-slate-200 dark:hover:bg-slate-950"
                    :style="fileTreeIndent(row.depth)"
                    :aria-expanded="expandedFileNodeIds.has(row.node.id)"
                    :disabled="filesPending"
                    @click="toggleFileTreeNode(row.node)"
                  >
                    <UIcon
                      v-if="filesPending && !filesLoaded && row.node.id === fileTreeRootId"
                      name="i-lucide-loader-circle"
                      class="size-4 shrink-0 animate-spin text-slate-400"
                    />
                    <UIcon
                      v-else
                      :name="expandedFileNodeIds.has(row.node.id) ? 'i-lucide-chevron-down' : 'i-lucide-chevron-right'"
                      class="size-4 shrink-0 text-slate-400"
                    />
                    <UIcon
                      :name="expandedFileNodeIds.has(row.node.id) ? 'i-lucide-folder-open' : 'i-lucide-folder'"
                      class="size-4 shrink-0 text-amber-500"
                    />
                    <span class="min-w-0 truncate font-medium">{{ row.node.name }}</span>
                    <span v-if="row.node.id !== fileTreeRootId" class="shrink-0 text-xs text-slate-500 dark:text-slate-400">
                      {{ $t('catalog.torrents.fileCount', { count: numberFormatter.format(row.node.fileCount) }) }}
                    </span>
                  </button>
                  <div
                    v-else
                    class="flex min-w-0 items-center gap-2 py-1 pr-2 text-slate-700 dark:text-slate-200"
                    :style="fileTreeIndent(row.depth)"
                  >
                    <span class="size-4 shrink-0" />
                    <UIcon name="i-lucide-file" class="size-4 shrink-0 text-slate-400" />
                    <span class="min-w-0 truncate font-mono text-xs">{{ row.node.name }}</span>
                  </div>
                  <div v-if="row.node.id === fileTreeRootId" class="flex items-center justify-end gap-1">
                    <UTooltip
                      :text="$t('catalog.torrents.detail.files.expandAll')"
                      :content="{ side: 'top', sideOffset: 8 }"
                      :delay-duration="120"
                    >
                      <UButton
                        color="neutral"
                        variant="ghost"
                        size="xs"
                        icon="i-lucide-folder-open"
                        :loading="filesPending"
                        :disabled="filesPending || (filesLoaded && allFileTreeExpanded)"
                        :aria-label="$t('catalog.torrents.detail.files.expandAll')"
                        @click.stop="expandAllFileTree"
                      />
                    </UTooltip>
                    <UTooltip
                      :text="$t('catalog.torrents.detail.files.collapseAll')"
                      :content="{ side: 'top', sideOffset: 8 }"
                      :delay-duration="120"
                    >
                      <UButton
                        color="neutral"
                        variant="ghost"
                        size="xs"
                        icon="i-lucide-folder"
                        :disabled="filesPending || expandedFileNodeIds.size === 0"
                        :aria-label="$t('catalog.torrents.detail.files.collapseAll')"
                        @click.stop="collapseAllFileTree"
                      />
                    </UTooltip>
                  </div>
                  <span v-else class="text-slate-500 sm:text-right dark:text-slate-400">{{ formatBytes(row.node.size) }}</span>
                </div>
              </div>
            </div>

            <div v-if="peerView" class="border-t border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
              <div v-if="peerView === 'completed'" class="px-4 py-8 text-center">
                <UIcon name="i-lucide-circle-check" class="mx-auto size-8 text-slate-400" />
                <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">
                  {{ $t('catalog.torrents.detail.peers.completedUnavailable') }}
                </p>
              </div>
              <div v-else-if="peersError" class="flex flex-col items-center justify-center px-4 py-8 text-center">
                <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
                <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ peersError }}</p>
                <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadPeers(true)">
                  {{ $t('common.retry') }}
                </UButton>
              </div>
              <div v-else-if="peersPending && !peersLoaded" class="space-y-2 px-4 py-4">
                <div v-for="item in 3" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
              </div>
              <div v-else-if="visiblePeers.length === 0" class="px-4 py-8 text-center text-sm text-slate-500 dark:text-slate-400">
                {{ $t('catalog.torrents.detail.peers.empty') }}
              </div>
              <div v-else class="max-h-80 divide-y divide-slate-100 overflow-auto dark:divide-slate-800">
                <div v-for="peer in visiblePeers" :key="`${peer.user?.id || 0}-${peer.startedAt}-${peer.isSeeder}`" class="grid gap-3 px-4 py-3 text-sm sm:grid-cols-[minmax(0,1fr)_180px] sm:items-center">
                  <div class="min-w-0">
                    <div class="flex min-w-0 items-center gap-2">
                      <p class="truncate font-medium text-slate-950 dark:text-white">{{ peerUserName(peer) }}</p>
                      <UBadge :color="peer.isSeeder ? 'success' : 'primary'" variant="soft">
                        {{ peer.isSeeder ? $t('catalog.torrents.detail.peers.seeder') : $t('catalog.torrents.detail.peers.leecher') }}
                      </UBadge>
                    </div>
                    <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(peer.startedAt, locale) }}</p>
                  </div>
                  <div class="grid grid-cols-2 gap-2 text-xs sm:text-right">
                    <div>
                      <p class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.peers.uploaded') }}</p>
                      <p class="mt-1 font-medium text-slate-950 dark:text-white">{{ formatBytes(peer.uploaded) }}</p>
                    </div>
                    <div>
                      <p class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.peers.downloaded') }}</p>
                      <p class="mt-1 font-medium text-slate-950 dark:text-white">{{ formatBytes(peer.downloaded) }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </section>

          <section id="torrent-subtitles" class="scroll-mt-24 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <button
              type="button"
              class="flex w-full items-center justify-between gap-4 px-4 py-3 text-left transition hover:bg-slate-50 focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-sky-200 dark:hover:bg-slate-950 dark:focus-visible:ring-sky-900"
              :aria-expanded="subtitlesPanelOpen"
              @click="toggleSubtitlesPanel"
            >
              <div class="min-w-0">
                <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.subtitles.title') }}</h2>
              </div>
              <div class="flex shrink-0 items-center gap-2">
                <UBadge v-if="subtitlesLoaded" color="neutral" variant="soft">
                  {{ $t('catalog.torrents.detail.subtitles.summary', { count: numberFormatter.format(subtitleTotal) }) }}
                </UBadge>
                <UIcon :name="subtitlesPanelOpen ? 'i-lucide-chevron-up' : 'i-lucide-chevron-down'" class="size-4 text-slate-400" />
              </div>
            </button>

            <div v-if="subtitlesPanelOpen" class="border-t border-slate-200 p-4 dark:border-slate-800">
              <form class="grid grid-cols-1 gap-3 rounded-md border border-slate-200 bg-slate-50 p-3 lg:grid-cols-[minmax(0,1fr)_180px_auto_auto] lg:items-center dark:border-slate-800 dark:bg-slate-950" @submit.prevent="handleSubtitleUpload">
                <label class="flex min-h-10 cursor-pointer items-center gap-2 rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-700 transition hover:border-slate-300 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-200 dark:hover:border-slate-600">
                  <UIcon name="i-lucide-file-up" class="size-4 shrink-0 text-slate-400" />
                  <span class="min-w-0 truncate">{{ selectedSubtitleFile?.name || $t('catalog.torrents.detail.subtitles.choose') }}</span>
                  <input :key="subtitleFileInputKey" class="sr-only" type="file" :disabled="subtitleUploadPending" @change="handleSubtitleFileChange">
                </label>
                <select
                  v-model="subtitleForm.language"
                  class="h-10 rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="subtitleUploadPending"
                >
                  <option v-for="option in subtitleLanguageOptions" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
                <label class="flex h-10 items-center justify-between gap-3 rounded-md border border-slate-200 bg-white px-3 text-sm font-medium text-slate-700 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-200">
                  <span>{{ $t('catalog.torrents.detail.subtitles.anonymous') }}</span>
                  <input
                    v-model="subtitleForm.anonymous"
                    type="checkbox"
                    class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                    :disabled="subtitleUploadPending"
                  >
                </label>
                <UButton type="submit" color="primary" icon="i-lucide-upload" :loading="subtitleUploadPending" :disabled="!canUploadSubtitle">
                  {{ $t('catalog.torrents.detail.subtitles.upload') }}
                </UButton>
              </form>

              <div v-if="subtitlesError" class="mt-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
                {{ subtitlesError }}
              </div>

              <div v-if="subtitlesPending && subtitles.length === 0" class="mt-4 space-y-2">
                <div v-for="item in 3" :key="item" class="h-14 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
              </div>
              <div v-else-if="subtitles.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-8 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
                {{ $t('catalog.torrents.detail.subtitles.empty') }}
              </div>
              <div v-else class="mt-4 overflow-hidden rounded-md border border-slate-200 dark:border-slate-800">
                <div v-for="subtitle in subtitles" :key="subtitle.id" class="grid gap-3 border-b border-slate-100 px-3 py-3 last:border-b-0 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center dark:border-slate-800">
                  <div class="min-w-0">
                    <div class="flex min-w-0 items-center gap-2">
                      <UIcon name="i-lucide-captions" class="size-4 shrink-0 text-sky-500" />
                      <p class="truncate text-sm font-medium text-slate-950 dark:text-white">{{ subtitle.fileName }}</p>
                      <UBadge color="neutral" variant="soft">{{ subtitleLanguageLabel(subtitle.language) }}</UBadge>
                    </div>
                    <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                      {{ subtitleUploaderName(subtitle) }} · {{ formatBytes(subtitle.size) }} · {{ formatDateTime(subtitle.createdAt, locale) }}
                    </p>
                  </div>
                  <div class="flex items-center gap-2">
                    <UButton color="neutral" variant="outline" size="xs" icon="i-lucide-download" :loading="subtitleDownloadPendingId === subtitle.id" @click="handleSubtitleDownload(subtitle)">
                      {{ $t('catalog.torrents.detail.subtitles.download') }}
                    </UButton>
                    <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-flag" @click="startSubtitleReport(subtitle.id)">
                      {{ $t('catalog.torrents.detail.actions.report') }}
                    </UButton>
                  </div>
                  <form v-if="activeSubtitleReportId === subtitle.id" class="sm:col-span-2 grid gap-2 rounded-md bg-slate-50 p-3 dark:bg-slate-950" @submit.prevent="handleSubtitleReport(subtitle.id)">
                    <UTextarea v-model="subtitleReportReason" :rows="2" :placeholder="$t('catalog.torrents.detail.report.reason')" :disabled="reportPending" />
                    <div class="flex justify-end gap-2">
                      <UButton color="neutral" variant="ghost" size="xs" type="button" @click="activeSubtitleReportId = 0">{{ $t('common.cancel') }}</UButton>
                      <UButton color="error" variant="soft" size="xs" type="submit" :loading="reportPending" :disabled="subtitleReportReason.trim().length < 5">
                        {{ $t('catalog.torrents.detail.report.submit') }}
                      </UButton>
                    </div>
                  </form>
                </div>
              </div>

              <AppPager
                v-if="subtitleTotal > subtitleSize"
                class="mt-4 border-t border-slate-200 pt-4 dark:border-slate-800"
                size="xs"
                :page="subtitlePage"
                :total="subtitleTotal"
                :page-size="subtitleSize"
                :disabled="subtitlesPending"
                @page-change="goToSubtitlePage"
              />
            </div>
          </section>

          <section id="torrent-description" class="scroll-mt-24 rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 pb-4 dark:border-slate-800">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.description.title') }}</h2>
            </div>
            <div class="mt-4">
              <div
                v-if="renderedDescription"
                class="rich-text"
                v-html="renderedDescription"
              />
              <p v-else class="text-sm text-slate-500 dark:text-slate-400">
                {{ $t('catalog.torrents.detail.description.empty') }}
              </p>
            </div>
          </section>

          <section id="torrent-comments" class="scroll-mt-24 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.comments.title') }}</h2>
              <UBadge color="neutral" variant="soft">
                {{ $t('catalog.torrents.detail.comments.summary', { count: numberFormatter.format(commentTotal) }) }}
              </UBadge>
            </div>

            <form id="torrent-comment-composer" class="border-b border-slate-200 bg-slate-50/60 p-4 dark:border-slate-800 dark:bg-slate-950/40" @submit.prevent="handleCommentSubmit">
              <UTextarea
                v-if="commentEditorMode === 'write'"
                v-model="commentForm.content"
                class="w-full"
                :rows="4"
                :placeholder="$t('catalog.torrents.detail.comments.placeholder')"
                :disabled="commentSubmitPending"
              />
              <div v-else class="min-h-28 rounded-md border border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950">
                <div v-if="renderedCommentPreview" class="rich-text rich-text-compact" v-html="renderedCommentPreview" />
                <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.comments.previewEmpty') }}</p>
              </div>
              <div class="mt-3 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div class="inline-flex w-fit rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950">
                  <button
                    type="button"
                    :class="commentEditorTabClass('write')"
                    @click="commentEditorMode = 'write'"
                  >
                    {{ $t('catalog.torrents.detail.comments.edit') }}
                  </button>
                  <button
                    type="button"
                    :class="commentEditorTabClass('preview')"
                    @click="commentEditorMode = 'preview'"
                  >
                    {{ $t('catalog.torrents.detail.comments.preview') }}
                  </button>
                </div>
                <UButton type="submit" color="primary" icon="i-lucide-send" :loading="commentSubmitPending" :disabled="commentForm.content.trim().length < 3">
                  {{ $t('catalog.torrents.detail.comments.submit') }}
                </UButton>
              </div>
            </form>

            <CatalogCommentList
              v-model:active-report-id="activeCommentReportId"
              v-model:report-reason="commentReportReason"
              :comments="comments"
              :total="commentTotal"
              :page="commentPage"
              :page-size="commentSize"
              :pending="commentsPending"
              :error="commentsError"
              :like-pending-id="commentLikePendingId"
              :report-pending="reportPending"
              :submit-reward="submitCommentReward"
              @retry="loadComments"
              @page-change="goToCommentPage"
              @toggle-like="handleCommentLike"
              @reward-success="handleCommentRewardSuccess"
              @quote="insertCommentQuote"
              @report="handleCommentReport"
            />
          </section>
          <div id="torrent-bottom" class="h-px scroll-mt-24" />
        </main>

        <nav
          class="hidden lg:sticky lg:top-[5.5rem] lg:flex lg:h-[calc(100vh-7rem)] lg:items-center lg:justify-center"
          :aria-label="$t('catalog.torrents.detail.navigation.title')"
        >
          <div class="relative flex min-h-56 w-full justify-center">
            <div class="absolute inset-y-0 left-1/2 w-px -translate-x-1/2 bg-slate-200 dark:bg-slate-800" />
            <div class="relative z-10 flex flex-col items-center gap-2">
              <UTooltip
                v-for="item in sectionNavItems"
                :key="item.id"
                :text="item.label"
                :content="{ side: 'left', sideOffset: 10 }"
                :delay-duration="120"
              >
                <button
                  type="button"
                  :class="sectionRailButtonClass(item.id)"
                  :aria-label="item.label"
                  :aria-current="activeSectionId === item.id ? 'location' : undefined"
                  @click="scrollToSection(item.id)"
                >
                  <UIcon :name="item.icon" class="size-4" />
                </button>
              </UTooltip>
            </div>
          </div>
        </nav>

        <aside class="min-w-0 space-y-3 lg:sticky lg:top-[5.5rem]">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <div class="flex items-center justify-between gap-3">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.info.title') }}</h2>
              <span class="shrink-0 rounded border border-slate-200 px-2 py-0.5 text-xs font-medium text-slate-500 dark:border-slate-700 dark:text-slate-400">
                #{{ torrent.id }}
              </span>
            </div>
            <dl class="mt-3 divide-y divide-slate-100 text-sm dark:divide-slate-800">
              <div class="grid grid-cols-[88px_minmax(0,1fr)] gap-3 py-2.5">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.category') }}</dt>
                <dd class="min-w-0 truncate text-right font-medium text-slate-950 dark:text-white">{{ categoryName }}</dd>
              </div>
              <div class="grid grid-cols-[88px_minmax(0,1fr)] gap-3 py-2.5">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.publisher') }}</dt>
                <dd class="min-w-0 truncate text-right font-medium text-slate-950 dark:text-white">{{ publisherName }}</dd>
              </div>
              <div class="grid grid-cols-[88px_minmax(0,1fr)] gap-3 py-2.5">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.publishedAt') }}</dt>
                <dd class="text-right font-medium text-slate-950 dark:text-white">{{ formatDateTime(torrent.createdAt, locale) }}</dd>
              </div>
            </dl>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.actions.title') }}</h2>
            <div class="mt-3 grid grid-cols-2 gap-2">
              <UButton
                v-if="canEditTorrent"
                color="neutral"
                variant="outline"
                size="sm"
                icon="i-lucide-pen-line"
                block
                :to="localePath(`/catalog/torrents/${torrent.id}/edit`)"
              >
                {{ $t('catalog.torrents.detail.actions.edit') }}
              </UButton>
              <UButton
                color="neutral"
                variant="outline"
                size="sm"
                icon="i-lucide-flag"
                block
                :class="canEditTorrent ? '' : 'col-span-2'"
                @click="showTorrentReport = !showTorrentReport"
              >
                {{ $t('catalog.torrents.detail.actions.report') }}
              </UButton>
              <UButton
                v-if="isStaff"
                class="col-span-2"
                color="error"
                variant="soft"
                size="sm"
                icon="i-lucide-trash-2"
                block
                :loading="adminActionPending === 'delete'"
                @click="handleAdminDeleteTorrent"
              >
                {{ $t('catalog.torrents.detail.actions.adminDelete') }}
              </UButton>
            </div>

            <form v-if="showTorrentReport" class="mt-3 grid gap-2 rounded-md bg-slate-50 p-3 dark:bg-slate-950" @submit.prevent="handleTorrentReport">
              <UTextarea v-model="torrentReportReason" :rows="3" :placeholder="$t('catalog.torrents.detail.report.reason')" :disabled="reportPending" />
              <div class="flex justify-end gap-2">
                <UButton color="neutral" variant="ghost" size="xs" type="button" @click="showTorrentReport = false">{{ $t('common.cancel') }}</UButton>
                <UButton color="error" variant="soft" size="xs" type="submit" :loading="reportPending" :disabled="torrentReportReason.trim().length < 5">
                  {{ $t('catalog.torrents.detail.report.submit') }}
                </UButton>
              </div>
            </form>
          </section>

          <RewardPanel
            :title="$t('catalog.torrents.detail.reward.title')"
            :source-key="torrent.id"
            summary-key="catalog.torrents.detail.reward.summary"
            all-title-key="catalog.torrents.detail.reward.allTitle"
            all-description-key="catalog.torrents.detail.reward.allDescription"
            success-key="catalog.torrents.detail.reward.success"
            :load-rewards="loadTorrentRewards"
            :submit-reward="submitTorrentReward"
          />
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { CatalogCategory, CommentItem, SubtitleItem, TorrentDetail, TorrentFileItem, TorrentPeerItem } from '~/composables/useCatalogTorrents'
import { renderUserMarkdown } from '~/utils/richText'

interface FileTreeNode {
  id: string
  name: string
  path: string
  type: 'directory' | 'file'
  size: number
  fileCount: number
  children: FileTreeNode[]
}

interface FileTreeRow {
  node: FileTreeNode
  depth: number
}

interface TorrentStatusBadge {
  key: string
  label: string
  class: string
  title: string
}

interface SectionNavItem {
  id: string
  icon: string
  label: string
}

type PeerView = 'seeders' | 'leechers' | 'completed'
type CommentEditorMode = 'write' | 'preview'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()
const adminApi = useAdmin()
const { user, isStaff } = useAuth()

const torrent = ref<TorrentDetail | null>(null)
const files = ref<TorrentFileItem[]>([])
const peers = ref<TorrentPeerItem[]>([])
const comments = ref<CommentItem[]>([])
const subtitles = ref<SubtitleItem[]>([])
const categories = ref<CatalogCategory[]>([])
const pending = ref(false)
const downloadPending = ref(false)
const likePending = ref(false)
const bookmarkPending = ref(false)
const peersPending = ref(false)
const commentsPending = ref(false)
const subtitlesPending = ref(false)
const commentSubmitPending = ref(false)
const subtitleUploadPending = ref(false)
const peersLoaded = ref(false)
const errorMessage = ref('')
const peersError = ref('')
const commentsError = ref('')
const subtitlesError = ref('')
const expandedFileNodeIds = ref(new Set<string>())
const filesLoaded = ref(false)
const filesPending = ref(false)
const filesError = ref('')
const commentTotal = ref(0)
const subtitleTotal = ref(0)
const selectedSubtitleFile = ref<File | null>(null)
const subtitleFileInputKey = ref(0)
const subtitleDownloadPendingId = ref(0)
const commentLikePendingId = ref(0)
const activeCommentReportId = ref(0)
const activeSubtitleReportId = ref(0)
const peerView = ref<PeerView | null>(null)
const filePanelOpen = ref(false)
const subtitlesPanelOpen = ref(false)
const subtitlesLoaded = ref(false)
const reportPending = ref(false)
const adminActionPending = ref('')
const showTorrentReport = ref(false)
const torrentReportReason = ref('')
const commentReportReason = ref('')
const subtitleReportReason = ref('')
const commentEditorMode = ref<CommentEditorMode>('write')
const activeSectionId = ref('torrent-top')
let sectionObserver: IntersectionObserver | null = null

const commentForm = reactive({
  content: ''
})

const commentSize = 20
const subtitleSize = 10
const commentPage = ref(1)
const subtitlePage = ref(1)
const subtitleForm = reactive({
  language: 'zh-CN',
  anonymous: false
})

const subtitleLanguageOptions = [
  { value: 'zh-CN', label: '简体中文' },
  { value: 'zh-TW', label: '繁體中文' },
  { value: 'en-US', label: 'English' },
  { value: 'ja-JP', label: '日本語' },
  { value: 'ko-KR', label: '한국어' },
  { value: 'other', label: 'Other' }
]

const subtitleLanguageLabelMap = new Map(subtitleLanguageOptions.map((option) => [option.value, option.label]))

const fileTreeRootId = 'dir:/'
const torrentId = computed(() => readRouteId())
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const commentTotalPages = computed(() => Math.max(1, Math.ceil(commentTotal.value / commentSize)))
const fileTree = computed(() => filesLoaded.value ? buildFileTree(files.value) : buildFileTreePlaceholder())
const fileTreeDirectoryIds = computed(() => collectFileTreeDirectoryIds(fileTree.value))
const hasFileTreeDirectories = computed(() => fileTreeDirectoryIds.value.length > 0)
const allFileTreeExpanded = computed(() => hasFileTreeDirectories.value && fileTreeDirectoryIds.value.every((id) => expandedFileNodeIds.value.has(id)))
const visibleFileTreeRows = computed(() => flattenFileTree(fileTree.value, expandedFileNodeIds.value))
const renderedDescription = computed(() => renderRichText(torrent.value?.description || ''))
const renderedCommentPreview = computed(() => renderRichText(commentForm.content))
const canUploadSubtitle = computed(() => Boolean(selectedSubtitleFile.value && subtitleForm.language && !subtitleUploadPending.value))
const canEditTorrent = computed(() => Boolean(torrent.value && (isStaff.value || user.value?.id === torrent.value.owner?.id)))
const visiblePeers = computed(() => {
  if (peerView.value === 'seeders') return peers.value.filter((item) => item.isSeeder)
  if (peerView.value === 'leechers') return peers.value.filter((item) => !item.isSeeder)
  return []
})
const sectionNavItems = computed<SectionNavItem[]>(() => [
  {
    id: 'torrent-top',
    icon: 'i-lucide-arrow-up',
    label: t('catalog.torrents.detail.navigation.top')
  },
  {
    id: 'torrent-subtitles',
    icon: 'i-lucide-captions',
    label: t('catalog.torrents.detail.subtitles.title')
  },
  {
    id: 'torrent-description',
    icon: 'i-lucide-align-left',
    label: t('catalog.torrents.detail.description.title')
  },
  {
    id: 'torrent-comments',
    icon: 'i-lucide-message-square',
    label: t('catalog.torrents.detail.comments.title')
  },
  {
    id: 'torrent-bottom',
    icon: 'i-lucide-arrow-down',
    label: t('catalog.torrents.detail.navigation.bottom')
  }
])
const categoryName = computed(() => {
  if (!torrent.value) return '-'

  const category = categories.value.find((item) => item.id === torrent.value?.categoryId)
  return category ? localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`) : t('catalog.torrents.unknownCategory')
})
const publisherName = computed(() => {
  if (!torrent.value) return '-'
  return torrentOwnerName(torrent.value)
})

function torrentOwnerName(torrent: TorrentDetail) {
  const ownerName = rawTorrentOwnerName(torrent)
  if (torrent.anonymous) {
    return torrent.owner?.id > 0 ? t('catalog.torrents.anonymousOwner', { name: ownerName }) : t('catalog.torrents.anonymous')
  }
  return torrent.owner?.id > 0 ? ownerName : '-'
}

function rawTorrentOwnerName(torrent: TorrentDetail) {
  return torrent.owner?.username || (torrent.owner?.id > 0 ? `#${torrent.owner.id}` : '')
}

function renderRichText(content: string) {
  return renderUserMarkdown(content).trim()
}

function commentEditorTabClass(mode: CommentEditorMode) {
  const base = 'h-8 rounded px-3 text-sm font-medium transition focus:outline-none focus-visible:ring-2 focus-visible:ring-sky-200 dark:focus-visible:ring-sky-900'
  if (commentEditorMode.value === mode) {
    return `${base} bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white`
  }
  return `${base} text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white`
}

function sectionRailButtonClass(id: string) {
  const base = 'flex size-9 items-center justify-center rounded-full border text-slate-400 shadow-sm transition focus:outline-none focus-visible:ring-2 focus-visible:ring-sky-200 dark:focus-visible:ring-sky-900'
  if (activeSectionId.value === id) {
    return `${base} border-sky-200 bg-sky-50 text-sky-600 shadow-sky-950/5 dark:border-sky-900 dark:bg-sky-950 dark:text-sky-300`
  }
  return `${base} border-slate-200 bg-white hover:border-sky-200 hover:bg-sky-50 hover:text-sky-600 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-sky-900 dark:hover:bg-sky-950 dark:hover:text-sky-300`
}

async function scrollToSection(id: string) {
  activeSectionId.value = id

  if (id === 'torrent-subtitles') {
    await openSubtitlesPanel()
  }
  if (id === 'torrent-bottom') {
    window.scrollTo({ top: document.documentElement.scrollHeight, behavior: 'smooth' })
    return
  }

  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function setupSectionObserver() {
  if (typeof window === 'undefined' || !('IntersectionObserver' in window)) return

  sectionObserver?.disconnect()
  sectionObserver = new IntersectionObserver((entries) => {
    const visibleEntry = entries
      .filter((entry) => entry.isIntersecting)
      .sort((a, b) => {
        if (b.intersectionRatio !== a.intersectionRatio) return b.intersectionRatio - a.intersectionRatio
        return a.boundingClientRect.top - b.boundingClientRect.top
      })[0]
    if (visibleEntry?.target.id) {
      activeSectionId.value = visibleEntry.target.id
    }
  }, {
    rootMargin: '-24% 0px -62% 0px',
    threshold: [0, 0.1, 0.25, 0.5]
  })

  for (const item of sectionNavItems.value) {
    const element = document.getElementById(item.id)
    if (element) {
      sectionObserver.observe(element)
    }
  }
}

function torrentStatusBadges(torrent: TorrentDetail) {
  const badges: TorrentStatusBadge[] = []
  if (torrent.isPinned) {
    badges.push({
      key: 'pinned',
      label: t('catalog.torrents.status.pinned'),
      class: 'bg-slate-900 text-white dark:bg-white dark:text-slate-950',
      title: t('catalog.torrents.status.pinned')
    })
  }
  if (torrent.isFeatured) {
    badges.push({
      key: 'featured',
      label: t('catalog.torrents.status.featured'),
      class: 'bg-amber-100 text-amber-800 dark:bg-amber-950 dark:text-amber-200',
      title: t('catalog.torrents.status.featured')
    })
  }

  const promotionLabel = torrentPromotionLabel(torrent.spState)
  if (promotionLabel && isPromotionActive(torrent)) {
    badges.push({
      key: `sp-${torrent.spState}`,
      label: promotionLabel,
      class: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-200',
      title: torrent.spExpireAt
        ? t('catalog.torrents.status.expiresAt', { time: formatDateTime(torrent.spExpireAt, locale.value) })
        : promotionLabel
    })
  }
  return badges
}

function torrentPromotionLabel(spState?: number | null) {
  const key = Number(spState || 0)
  return key >= 1 && key <= 7 ? t(`catalog.torrents.status.promotion.${key}`) : ''
}

function isPromotionActive(torrent: TorrentDetail) {
  if (!torrent.spState) return false
  if (!torrent.spExpireAt) return true

  const expireAt = parseDateTime(torrent.spExpireAt)
  return !expireAt || expireAt.getTime() > Date.now()
}

function parseDateTime(value?: string | null) {
  if (!value) return null

  const date = new Date(value)
  if (!Number.isNaN(date.getTime())) return date

  const normalized = new Date(value.replace(' ', 'T'))
  return Number.isNaN(normalized.getTime()) ? null : normalized
}

async function handleAdminDeleteTorrent() {
  if (!torrent.value || !window.confirm(t('catalog.torrents.detail.actions.confirmAdminDelete', { name: torrent.value.name }))) return

  adminActionPending.value = 'delete'
  try {
    await adminApi.deleteCatalogTorrent(torrent.value.id)
    toast.add({ title: t('catalog.torrents.detail.actions.adminDeleted'), color: 'success', icon: 'i-lucide-check-circle' })
    await navigateTo(localePath('/catalog/torrents'))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    adminActionPending.value = ''
  }
}

onMounted(loadPage)
onBeforeUnmount(() => {
  sectionObserver?.disconnect()
  sectionObserver = null
})

watch(() => route.hash, () => {
  void nextTick(() => scrollToCommentHash())
})

function readRouteId() {
  const raw = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id
  const id = Number(raw)
  return Number.isInteger(id) && id > 0 ? id : 0
}

function readCommentPageQuery() {
  const raw = Array.isArray(route.query.commentPage) ? route.query.commentPage[0] : route.query.commentPage
  const page = Number(raw)
  return Number.isInteger(page) && page > 0 ? page : 1
}

async function loadPage() {
  if (torrentId.value <= 0) {
    errorMessage.value = t('catalog.torrents.detail.invalidId')
    return
  }

  pending.value = true
  errorMessage.value = ''

  try {
    const [detail, categoryList] = await Promise.all([
      catalogTorrents.getTorrent(torrentId.value),
      catalogTorrents.listCategories().catch(() => ({ list: [] }))
    ])

    torrent.value = detail
    files.value = []
    filesLoaded.value = false
    filesError.value = ''
    expandedFileNodeIds.value = new Set()
    categories.value = categoryList.list || []
    resetInteractionState()
    void loadComments()
  } catch (error) {
    torrent.value = null
    files.value = []
    filesLoaded.value = false
    filesError.value = ''
    expandedFileNodeIds.value = new Set()
    resetInteractionState()
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
    if (torrent.value) {
      void nextTick(setupSectionObserver)
    }
  }
}

function resetInteractionState() {
  peers.value = []
  peersLoaded.value = false
  comments.value = []
  commentPage.value = readCommentPageQuery()
  subtitles.value = []
  subtitlePage.value = 1
  commentTotal.value = 0
  subtitleTotal.value = 0
  peersError.value = ''
  commentsError.value = ''
  subtitlesError.value = ''
  selectedSubtitleFile.value = null
  subtitleFileInputKey.value += 1
  activeCommentReportId.value = 0
  activeSubtitleReportId.value = 0
  peerView.value = null
  filePanelOpen.value = false
  subtitlesPanelOpen.value = false
  subtitlesLoaded.value = false
  activeSectionId.value = 'torrent-top'
  showTorrentReport.value = false
  commentForm.content = ''
  commentEditorMode.value = 'write'
  torrentReportReason.value = ''
  commentReportReason.value = ''
  subtitleReportReason.value = ''
}

async function loadPeers(force = false) {
  if (torrentId.value <= 0 || peersPending.value) return
  if (peersLoaded.value && !force) return

  peersPending.value = true
  peersError.value = ''
  try {
    const data = await catalogTorrents.listPeers(torrentId.value)
    peers.value = data.list || []
    peersLoaded.value = true
  } catch (error) {
    peers.value = []
    peersLoaded.value = false
    peersError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    peersPending.value = false
  }
}

async function handlePeerStatClick(view: PeerView) {
  if (peerView.value === view) {
    peerView.value = null
    return
  }

  filePanelOpen.value = false
  peerView.value = view
  if (view !== 'completed') {
    await loadPeers()
  }
}

async function handleFileStatClick() {
  filePanelOpen.value = !filePanelOpen.value
  if (!filePanelOpen.value) return

  peerView.value = null
  await loadFiles(true)
}

function fileStatCardClass() {
  const base = 'min-w-0 px-3 py-3 text-left transition focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-sky-200 dark:focus-visible:ring-sky-900'
  if (filePanelOpen.value) {
    return `${base} bg-amber-50 dark:bg-amber-950/30`
  }
  return `${base} hover:bg-white dark:hover:bg-slate-950`
}

function peerStatCardClass(view: PeerView) {
  const base = 'min-w-0 px-3 py-3 text-left transition focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-sky-200 dark:focus-visible:ring-sky-900'
  if (peerView.value === view) {
    return `${base} bg-sky-50 dark:bg-sky-950/40`
  }
  return `${base} hover:bg-white dark:hover:bg-slate-950`
}

async function handleDownload() {
  if (!torrent.value || downloadPending.value) return

  downloadPending.value = true
  try {
    const out = await catalogTorrents.downloadTorrent(torrent.value.id)
    downloadBlob(out.blob, out.filename || fallbackTorrentFilename(torrent.value))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    downloadPending.value = false
  }
}

async function handleToggleLike() {
  if (!torrent.value || likePending.value) return

  likePending.value = true
  try {
    const out = await catalogTorrents.toggleLike(torrent.value.id)
    if (torrent.value.isLiked !== out.isLiked) {
      torrent.value.likeCount = Math.max(0, Number(torrent.value.likeCount || 0) + (out.isLiked ? 1 : -1))
    }
    torrent.value.isLiked = out.isLiked
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    likePending.value = false
  }
}

async function handleToggleBookmark() {
  if (!torrent.value || bookmarkPending.value) return

  bookmarkPending.value = true
  try {
    if (torrent.value.isBookmarked) {
      await catalogTorrents.unbookmark(torrent.value.id)
      torrent.value.isBookmarked = false
    } else {
      await catalogTorrents.bookmark(torrent.value.id)
      torrent.value.isBookmarked = true
    }
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    bookmarkPending.value = false
  }
}

async function loadTorrentRewards(page: number, size: number) {
  if (torrentId.value <= 0) return { list: [], total: 0 }
  return await catalogTorrents.listRewards(torrentId.value, page, size)
}

async function submitTorrentReward(amount: number) {
  if (!torrent.value) return
  await catalogTorrents.rewardTorrent(torrent.value.id, amount)
}

async function handleTorrentReport() {
  if (!torrent.value || reportPending.value) return

  const reason = torrentReportReason.value.trim()
  if (reason.length < 5) return

  reportPending.value = true
  try {
    await catalogTorrents.reportTorrent(torrent.value.id, reason)
    showTorrentReport.value = false
    torrentReportReason.value = ''
    toast.add({
      title: t('catalog.torrents.detail.report.success'),
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
    reportPending.value = false
  }
}

async function loadComments() {
  if (torrentId.value <= 0 || commentsPending.value) return

  commentsPending.value = true
  commentsError.value = ''
  try {
    const data = await catalogTorrents.listComments(torrentId.value, commentPage.value, commentSize)
    comments.value = data.list || []
    commentTotal.value = data.total || 0

    if (commentPage.value > commentTotalPages.value) {
      commentPage.value = commentTotalPages.value
      const nextData = await catalogTorrents.listComments(torrentId.value, commentPage.value, commentSize)
      comments.value = nextData.list || []
      commentTotal.value = nextData.total || 0
    }
    await nextTick()
    scrollToCommentHash('auto')
  } catch (error) {
    comments.value = []
    commentTotal.value = 0
    commentsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    commentsPending.value = false
  }
}

function goToCommentPage(nextPage: number) {
  commentPage.value = Math.min(Math.max(1, nextPage), commentTotalPages.value)
  loadComments()
}

function scrollToCommentAnchor(id: string, behavior: ScrollBehavior = 'smooth') {
  document.getElementById(id)?.scrollIntoView({ behavior, block: 'start' })
}

function scrollToCommentHash(behavior: ScrollBehavior = 'smooth') {
  const id = route.hash.startsWith('#') ? route.hash.slice(1) : ''
  if (!id.startsWith('comment-')) return

  scrollToCommentAnchor(id, behavior)
}

function insertCommentQuote(quote: string) {
  commentForm.content = commentForm.content.trim()
    ? `${commentForm.content.trim()}\n\n${quote}`
    : quote
  commentEditorMode.value = 'write'
  void nextTick(() => {
    document.getElementById('torrent-comment-composer')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  })
}

async function handleCommentSubmit() {
  const content = commentForm.content.trim()
  if (!torrent.value || content.length < 3 || commentSubmitPending.value) return

  commentSubmitPending.value = true
  try {
    await catalogTorrents.createComment(torrent.value.id, content)
    commentForm.content = ''
    commentEditorMode.value = 'write'
    toast.add({
      title: t('catalog.torrents.detail.comments.created'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    commentPage.value = Math.max(1, Math.ceil((commentTotal.value + 1) / commentSize))
    await loadComments()
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    commentSubmitPending.value = false
  }
}

async function handleCommentLike(comment: CommentItem) {
  if (!torrent.value || commentLikePendingId.value > 0) return

  commentLikePendingId.value = comment.id
  try {
    const out = await catalogTorrents.toggleCommentLike(torrent.value.id, comment.id)
    if (comment.isLiked !== out.isLiked) {
      comment.likeCount = Math.max(0, comment.likeCount + (out.isLiked ? 1 : -1))
    }
    comment.isLiked = out.isLiked
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    commentLikePendingId.value = 0
  }
}

async function submitCommentReward(comment: CommentItem, amount: number) {
  if (!torrent.value) return
  await catalogTorrents.rewardComment(torrent.value.id, comment.id, amount)
}

function handleCommentRewardSuccess(comment: CommentItem) {
  comment.rewardCount += 1
}

async function handleCommentReport(commentId: number) {
  if (!torrent.value || reportPending.value) return

  const reason = commentReportReason.value.trim()
  if (reason.length < 5) return

  reportPending.value = true
  try {
    await catalogTorrents.reportComment(torrent.value.id, commentId, reason)
    activeCommentReportId.value = 0
    commentReportReason.value = ''
    toast.add({
      title: t('catalog.torrents.detail.report.success'),
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
    reportPending.value = false
  }
}

async function toggleSubtitlesPanel() {
  if (subtitlesPanelOpen.value) {
    subtitlesPanelOpen.value = false
    return
  }
  await openSubtitlesPanel()
}

async function openSubtitlesPanel() {
  subtitlesPanelOpen.value = true
  if (!subtitlesLoaded.value && !subtitlesPending.value) {
    await loadSubtitles(1)
  }
}

function goToSubtitlePage(page: number) {
  void loadSubtitles(page)
}

async function loadSubtitles(page = subtitlePage.value) {
  if (torrentId.value <= 0 || subtitlesPending.value) return

  subtitlePage.value = Math.max(1, page)
  subtitlesPending.value = true
  subtitlesError.value = ''
  try {
    const data = await catalogTorrents.listSubtitles(torrentId.value, subtitlePage.value, subtitleSize)
    subtitles.value = data.list || []
    subtitleTotal.value = data.total || 0
    subtitlesLoaded.value = true

    const totalPages = Math.max(1, Math.ceil(subtitleTotal.value / subtitleSize))
    if (subtitlePage.value > totalPages) {
      subtitlePage.value = totalPages
      const nextData = await catalogTorrents.listSubtitles(torrentId.value, subtitlePage.value, subtitleSize)
      subtitles.value = nextData.list || []
      subtitleTotal.value = nextData.total || 0
    }
  } catch (error) {
    subtitles.value = []
    subtitleTotal.value = 0
    subtitlesLoaded.value = false
    subtitlesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    subtitlesPending.value = false
  }
}

function handleSubtitleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  selectedSubtitleFile.value = input.files?.[0] || null
}

function subtitleLanguageLabel(language?: string | null) {
  if (!language) return '-'
  return subtitleLanguageLabelMap.get(language) || language
}

function subtitleUploaderName(subtitle: SubtitleItem) {
  if (subtitle.anonymous) {
    const name = rawSubtitleUploaderName(subtitle)
    return name ? t('catalog.torrents.anonymousOwner', { name }) : t('catalog.torrents.anonymous')
  }
  return subtitle.uploader?.id > 0 ? rawSubtitleUploaderName(subtitle) : '-'
}

function rawSubtitleUploaderName(subtitle: SubtitleItem) {
  return subtitle.uploader?.username || (subtitle.uploader?.id > 0 ? `#${subtitle.uploader.id}` : '')
}

function peerUserName(peer: TorrentPeerItem) {
  return peer.user?.username || (peer.user?.id > 0 ? `#${peer.user.id}` : '-')
}

async function handleSubtitleUpload() {
  if (!torrent.value || !selectedSubtitleFile.value || !subtitleForm.language || subtitleUploadPending.value) return

  subtitleUploadPending.value = true
  try {
    await catalogTorrents.uploadSubtitle(torrent.value.id, selectedSubtitleFile.value, subtitleForm.language, subtitleForm.anonymous)
    selectedSubtitleFile.value = null
    subtitleFileInputKey.value += 1
    toast.add({
      title: t('catalog.torrents.detail.subtitles.uploaded'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await loadSubtitles(1)
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    subtitleUploadPending.value = false
  }
}

async function handleSubtitleDownload(subtitle: SubtitleItem) {
  if (subtitleDownloadPendingId.value > 0) return

  subtitleDownloadPendingId.value = subtitle.id
  try {
    const out = await catalogTorrents.downloadSubtitle(subtitle.id)
    downloadBlob(out.blob, out.filename || subtitle.fileName || `subtitle-${subtitle.id}`)
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    subtitleDownloadPendingId.value = 0
  }
}

function startSubtitleReport(subtitleId: number) {
  activeSubtitleReportId.value = activeSubtitleReportId.value === subtitleId ? 0 : subtitleId
  subtitleReportReason.value = ''
}

async function handleSubtitleReport(subtitleId: number) {
  if (reportPending.value) return

  const reason = subtitleReportReason.value.trim()
  if (reason.length < 5) return

  reportPending.value = true
  try {
    await catalogTorrents.reportSubtitle(subtitleId, reason)
    activeSubtitleReportId.value = 0
    subtitleReportReason.value = ''
    toast.add({
      title: t('catalog.torrents.detail.report.success'),
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
    reportPending.value = false
  }
}

function fallbackTorrentFilename(item: TorrentDetail) {
  const name = (item.name || `torrent-${item.id}`).replace(/[\\/:*?"<>|]+/g, '_').trim()
  return `${name || `torrent-${item.id}`}.torrent`
}

function downloadBlob(blob: Blob, filename: string) {
  const href = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = href
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(href)
}

async function loadFiles(expandRoot = false) {
  if (torrentId.value <= 0 || filesPending.value) return false
  if (filesLoaded.value) {
    if (expandRoot) {
      const next = new Set(expandedFileNodeIds.value)
      next.add(fileTreeRootId)
      expandedFileNodeIds.value = next
    }
    return true
  }

  filesPending.value = true
  filesError.value = ''
  try {
    const data = await catalogTorrents.listFiles(torrentId.value)
    files.value = data.list || []
    filesLoaded.value = true

    if (expandRoot) {
      const next = new Set(expandedFileNodeIds.value)
      next.add(fileTreeRootId)
      expandedFileNodeIds.value = next
    }
    return true
  } catch (error) {
    filesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
    return false
  } finally {
    filesPending.value = false
  }
}

async function toggleFileTreeNode(node: FileTreeNode) {
  if (node.type !== 'directory') return
  if (!filesLoaded.value) {
    await loadFiles(true)
    return
  }

  const next = new Set(expandedFileNodeIds.value)
  if (next.has(node.id)) {
    next.delete(node.id)
  } else {
    next.add(node.id)
  }
  expandedFileNodeIds.value = next
}

async function expandAllFileTree() {
  if (!filesLoaded.value) {
    const loaded = await loadFiles()
    if (!loaded) return
  }
  expandedFileNodeIds.value = new Set(fileTreeDirectoryIds.value)
}

function collapseAllFileTree() {
  expandedFileNodeIds.value = new Set()
}

function fileTreeIndent(depth: number) {
  return {
    paddingLeft: `${12 + Math.min(depth, 12) * 18}px`
  }
}

function buildFileTree(items: TorrentFileItem[]) {
  const roots: FileTreeNode[] = []
  const nodeMap = new Map<string, FileTreeNode>()

  for (const item of items) {
    const parts = normalizeFilePath(item.path).split('/').filter(Boolean)
    if (parts.length === 0) continue

    let siblings = roots
    let path = ''

    for (let index = 0; index < parts.length; index += 1) {
      const name = parts[index]
      const isFile = index === parts.length - 1
      path = path ? `${path}/${name}` : name
      const id = `${isFile ? 'file' : 'dir'}:${path}`
      let node = nodeMap.get(id)

      if (!node) {
        node = {
          id,
          name,
          path,
          type: isFile ? 'file' : 'directory',
          size: 0,
          fileCount: 0,
          children: []
        }
        nodeMap.set(id, node)
        siblings.push(node)
      }

      if (isFile) {
        node.size = Number(item.size || 0)
        node.fileCount = 1
      } else {
        node.size += Number(item.size || 0)
        node.fileCount += 1
        siblings = node.children
      }
    }
  }

  sortFileTree(roots)
  if (roots.length === 0) return roots

  return [{
    id: fileTreeRootId,
    name: t('catalog.torrents.detail.files.root'),
    path: '',
    type: 'directory',
    size: roots.reduce((sum, node) => sum + node.size, 0),
    fileCount: roots.reduce((sum, node) => sum + node.fileCount, 0),
    children: roots
  }]
}

function buildFileTreePlaceholder() {
  if (!torrent.value || torrent.value.fileCount <= 0) return []

  return [{
    id: fileTreeRootId,
    name: t('catalog.torrents.detail.files.root'),
    path: '',
    type: 'directory',
    size: torrent.value.size,
    fileCount: torrent.value.fileCount,
    children: []
  }]
}

function normalizeFilePath(path: string) {
  return String(path || '').replace(/\\/g, '/').replace(/^\/+|\/+$/g, '')
}

function sortFileTree(nodes: FileTreeNode[]) {
  nodes.sort((a, b) => {
    if (a.type !== b.type) return a.type === 'directory' ? -1 : 1
    return a.name.localeCompare(b.name, locale.value, { numeric: true, sensitivity: 'base' })
  })

  for (const node of nodes) {
    if (node.children.length > 0) sortFileTree(node.children)
  }
}

function collectFileTreeDirectoryIds(nodes: FileTreeNode[]) {
  const ids: string[] = []
  for (const node of nodes) {
    if (node.type === 'directory') {
      ids.push(node.id)
      ids.push(...collectFileTreeDirectoryIds(node.children))
    }
  }
  return ids
}

function flattenFileTree(nodes: FileTreeNode[], expandedIds: Set<string>, depth = 0): FileTreeRow[] {
  const rows: FileTreeRow[] = []
  for (const node of nodes) {
    rows.push({ node, depth })
    if (node.type === 'directory' && expandedIds.has(node.id)) {
      rows.push(...flattenFileTree(node.children, expandedIds, depth + 1))
    }
  }
  return rows
}

useSeoMeta({
  title: () => torrent.value?.name || t('catalog.torrents.detail.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
