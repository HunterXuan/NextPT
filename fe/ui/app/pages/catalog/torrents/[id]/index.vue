<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div class="min-w-0">
          <UButton color="neutral" variant="ghost" icon="i-lucide-arrow-left" :to="localePath('/catalog/torrents')">
            {{ $t('catalog.torrents.detail.back') }}
          </UButton>
          <p class="mt-4 text-sm font-medium text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.eyebrow') }}</p>
          <h1 class="mt-1 break-words text-2xl font-semibold text-slate-950 dark:text-white">
            {{ torrent?.name || $t('catalog.torrents.detail.titleFallback', { id: torrentId }) }}
          </h1>
          <p v-if="torrent?.subTitle" class="mt-2 break-words text-sm text-slate-500 dark:text-slate-400">
            {{ torrent.subTitle }}
          </p>
        </div>

        <div class="grid grid-cols-1 gap-2 sm:grid-cols-3 lg:flex lg:items-center">
          <UButton
            color="primary"
            icon="i-lucide-download"
            :loading="downloadPending"
            :disabled="!torrent || pending"
            @click="handleDownload"
          >
            {{ $t('catalog.torrents.detail.actions.download') }}
          </UButton>
          <UButton
            color="neutral"
            :variant="torrent?.isLiked ? 'soft' : 'outline'"
            :icon="torrent?.isLiked ? 'i-lucide-heart' : 'i-lucide-heart-plus'"
            :loading="likePending"
            :disabled="!torrent || pending"
            @click="handleToggleLike"
          >
            {{ torrent?.isLiked ? $t('catalog.torrents.detail.actions.liked') : $t('catalog.torrents.detail.actions.like') }}
          </UButton>
          <UButton
            color="neutral"
            :variant="torrent?.isBookmarked ? 'soft' : 'outline'"
            :icon="torrent?.isBookmarked ? 'i-lucide-bookmark-check' : 'i-lucide-bookmark-plus'"
            :loading="bookmarkPending"
            :disabled="!torrent || pending"
            @click="handleToggleBookmark"
          >
            {{ torrent?.isBookmarked ? $t('catalog.torrents.detail.actions.bookmarked') : $t('catalog.torrents.detail.actions.bookmark') }}
          </UButton>
        </div>
      </div>

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

      <div v-else-if="torrent" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px] lg:items-start">
        <main class="space-y-6">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <div class="flex flex-wrap items-center gap-2">
              <UBadge color="neutral" variant="soft">{{ categoryName }}</UBadge>
              <UBadge :color="torrent.type === 1 ? 'primary' : 'neutral'" variant="subtle">
                {{ torrent.type === 1 ? $t('catalog.torrents.types.multi') : $t('catalog.torrents.types.single') }}
              </UBadge>
              <UBadge v-if="torrent.anonymous" color="neutral" variant="outline">
                {{ torrentOwnerName(torrent) }}
              </UBadge>
            </div>

            <div class="mt-5 grid grid-cols-2 gap-3 md:grid-cols-4">
              <button
                type="button"
                :class="peerStatCardClass('seeders')"
                :aria-pressed="activePeerView === 'seeders'"
                @click="handlePeerStatClick('seeders')"
              >
                <div class="flex items-center justify-between gap-2">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.stats.seeders') }}</p>
                  <UIcon name="i-lucide-users" class="size-4 shrink-0 text-emerald-500" />
                </div>
                <p class="mt-1 text-lg font-semibold text-emerald-600 dark:text-emerald-400">{{ numberFormatter.format(torrent.seeders) }}</p>
              </button>
              <button
                type="button"
                :class="peerStatCardClass('leechers')"
                :aria-pressed="activePeerView === 'leechers'"
                @click="handlePeerStatClick('leechers')"
              >
                <div class="flex items-center justify-between gap-2">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.stats.leechers') }}</p>
                  <UIcon name="i-lucide-loader-circle" class="size-4 shrink-0 text-sky-500" />
                </div>
                <p class="mt-1 text-lg font-semibold text-sky-600 dark:text-sky-400">{{ numberFormatter.format(torrent.leechers) }}</p>
              </button>
              <button
                type="button"
                :class="peerStatCardClass('completed')"
                :aria-pressed="activePeerView === 'completed'"
                @click="handlePeerStatClick('completed')"
              >
                <div class="flex items-center justify-between gap-2">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.stats.completed') }}</p>
                  <UIcon name="i-lucide-circle-check" class="size-4 shrink-0 text-slate-500 dark:text-slate-400" />
                </div>
                <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(torrent.snatched) }}</p>
              </button>
              <div class="rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.stats.size') }}</p>
                <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ formatBytes(torrent.size) }}</p>
              </div>
            </div>

            <div v-if="activePeerView" class="mt-4 overflow-hidden rounded-md border border-slate-200 dark:border-slate-800">
              <div class="flex items-start justify-between gap-3 border-b border-slate-200 bg-slate-50 px-4 py-3 dark:border-slate-800 dark:bg-slate-950">
                <div class="min-w-0">
                  <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ activePeerPanelTitle }}</h2>
                  <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ activePeerPanelSummary }}</p>
                </div>
                <div class="flex shrink-0 items-center gap-1">
                  <UButton
                    v-if="activePeerView !== 'completed'"
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    icon="i-lucide-refresh-cw"
                    :aria-label="$t('common.refresh')"
                    :loading="peersPending"
                    @click="loadPeers(true)"
                  />
                  <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-x" :aria-label="$t('common.cancel')" @click="activePeerView = null" />
                </div>
              </div>

              <div v-if="activePeerView === 'completed'" class="px-4 py-8 text-center">
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
                <div v-for="peer in visiblePeers" :key="`${peer.userId}-${peer.startedAt}-${peer.isSeeder}`" class="grid gap-3 px-4 py-3 text-sm sm:grid-cols-[minmax(0,1fr)_180px] sm:items-center">
                  <div class="min-w-0">
                    <div class="flex min-w-0 items-center gap-2">
                      <p class="truncate font-medium text-slate-950 dark:text-white">{{ peer.username || `#${peer.userId}` }}</p>
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

          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.description.title') }}</h2>
            <p
              v-if="torrent.description"
              class="mt-4 whitespace-pre-wrap break-words text-sm leading-6 text-slate-700 dark:text-slate-200"
            >
              {{ torrent.description }}
            </p>
            <p v-else class="mt-4 text-sm text-slate-500 dark:text-slate-400">
              {{ $t('catalog.torrents.detail.description.empty') }}
            </p>
          </section>

          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
              <div>
                <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.files.title') }}</h2>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                  {{ $t('catalog.torrents.fileCount', { count: numberFormatter.format(torrent.fileCount) }) }}
                </p>
              </div>
              <div v-if="filesLoaded && hasFileTreeDirectories" class="flex items-center gap-2">
                <UButton
                  color="neutral"
                  variant="outline"
                  size="xs"
                  icon="i-lucide-folder-open"
                  :loading="filesPending"
                  :disabled="allFileTreeExpanded"
                  @click="expandAllFileTree"
                >
                  {{ $t('catalog.torrents.detail.files.expandAll') }}
                </UButton>
                <UButton
                  color="neutral"
                  variant="outline"
                  size="xs"
                  icon="i-lucide-folder"
                  :disabled="filesPending || expandedFileNodeIds.size === 0"
                  @click="collapseAllFileTree"
                >
                  {{ $t('catalog.torrents.detail.files.collapseAll') }}
                </UButton>
              </div>
            </div>

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

            <div v-else class="max-h-[560px] overflow-auto">
              <div
                v-for="row in visibleFileTreeRows"
                :key="row.node.id"
                class="grid min-h-11 gap-2 border-b border-slate-100 px-4 py-2 text-sm last:border-b-0 sm:grid-cols-[minmax(0,1fr)_120px] sm:items-center dark:border-slate-800"
              >
                <button
                  v-if="row.node.type === 'directory'"
                  type="button"
                  class="flex min-w-0 items-center gap-2 rounded-md py-1 pr-2 text-left text-slate-700 transition-colors hover:bg-slate-50 dark:text-slate-200 dark:hover:bg-slate-950"
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
                  <span class="shrink-0 text-xs text-slate-500 dark:text-slate-400">
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
                <span class="text-slate-500 sm:text-right dark:text-slate-400">{{ formatBytes(row.node.size) }}</span>
              </div>
            </div>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <div class="flex flex-col gap-3 border-b border-slate-200 pb-4 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
              <div>
                <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.subtitles.title') }}</h2>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                  {{ $t('catalog.torrents.detail.subtitles.summary', { count: numberFormatter.format(subtitleTotal) }) }}
                </p>
              </div>
              <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" :loading="subtitlesPending" @click="loadSubtitles">
                {{ $t('common.refresh') }}
              </UButton>
            </div>

            <form class="mt-4 grid grid-cols-1 gap-3 lg:grid-cols-[minmax(0,1fr)_180px_auto]" @submit.prevent="handleSubtitleUpload">
              <label class="flex min-h-10 cursor-pointer items-center gap-2 rounded-md border border-slate-200 bg-slate-50 px-3 text-sm text-slate-700 transition hover:border-slate-300 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-200 dark:hover:border-slate-700">
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
            <div v-else class="mt-4 divide-y divide-slate-100 dark:divide-slate-800">
              <div v-for="subtitle in subtitles" :key="subtitle.id" class="grid gap-3 py-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center">
                <div class="min-w-0">
                  <div class="flex min-w-0 items-center gap-2">
                    <UIcon name="i-lucide-captions" class="size-4 shrink-0 text-sky-500" />
                    <p class="truncate text-sm font-medium text-slate-950 dark:text-white">{{ subtitle.fileName }}</p>
                    <UBadge color="neutral" variant="soft">{{ subtitle.language }}</UBadge>
                  </div>
                  <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                    {{ subtitle.username || `#${subtitle.userId}` }} · {{ formatBytes(subtitle.size) }} · {{ formatDateTime(subtitle.createdAt, locale) }}
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
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <div class="flex flex-col gap-3 border-b border-slate-200 pb-4 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
              <div>
                <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.comments.title') }}</h2>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                  {{ $t('catalog.torrents.detail.comments.summary', { count: numberFormatter.format(commentTotal) }) }}
                </p>
              </div>
              <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" :loading="commentsPending" @click="loadComments">
                {{ $t('common.refresh') }}
              </UButton>
            </div>

            <form class="mt-4 space-y-3" @submit.prevent="handleCommentSubmit">
              <UTextarea v-model="commentForm.content" class="w-full" :rows="4" :placeholder="$t('catalog.torrents.detail.comments.placeholder')" :disabled="commentSubmitPending" />
              <div class="flex justify-end">
                <UButton type="submit" color="primary" icon="i-lucide-send" :loading="commentSubmitPending" :disabled="commentForm.content.trim().length < 3">
                  {{ $t('catalog.torrents.detail.comments.submit') }}
                </UButton>
              </div>
            </form>

            <div v-if="commentsError" class="mt-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
              {{ commentsError }}
            </div>

            <div v-if="commentsPending && comments.length === 0" class="mt-5 space-y-3">
              <div v-for="item in 3" :key="item" class="h-24 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>
            <div v-else-if="comments.length === 0" class="mt-5 rounded-md border border-dashed border-slate-200 px-4 py-10 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
              {{ $t('catalog.torrents.detail.comments.empty') }}
            </div>
            <div v-else class="mt-5 divide-y divide-slate-100 dark:divide-slate-800">
              <article v-for="comment in comments" :key="comment.id" class="py-4">
                <div class="flex items-start gap-3">
                  <UAvatar :alt="comment.username || `#${comment.userId}`" size="sm" />
                  <div class="min-w-0 flex-1">
                    <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                      <p class="text-sm font-medium text-slate-950 dark:text-white">{{ comment.username || `#${comment.userId}` }}</p>
                      <span class="text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(comment.createdAt, locale) }}</span>
                    </div>
                    <p class="mt-2 whitespace-pre-wrap break-words text-sm leading-6 text-slate-700 dark:text-slate-200">{{ comment.content }}</p>
                    <div class="mt-3 flex flex-wrap items-center gap-2">
                      <UButton
                        color="neutral"
                        :variant="comment.isLiked ? 'soft' : 'ghost'"
                        size="xs"
                        :icon="comment.isLiked ? 'i-lucide-heart' : 'i-lucide-heart-plus'"
                        :loading="commentLikePendingId === comment.id"
                        @click="handleCommentLike(comment)"
                      >
                        {{ numberFormatter.format(comment.likeCount) }}
                      </UButton>
                      <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-coins" @click="startCommentReward(comment.id)">
                        {{ $t('catalog.torrents.detail.comments.reward', { count: numberFormatter.format(comment.rewardCount) }) }}
                      </UButton>
                      <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-flag" @click="startCommentReport(comment.id)">
                        {{ $t('catalog.torrents.detail.actions.report') }}
                      </UButton>
                    </div>
                    <form v-if="activeCommentRewardId === comment.id" class="mt-3 flex flex-col gap-2 rounded-md bg-slate-50 p-3 sm:flex-row dark:bg-slate-950" @submit.prevent="handleCommentReward(comment)">
                      <UInput v-model="commentRewardAmount" class="sm:w-40" type="number" min="1" step="1" :placeholder="$t('catalog.torrents.detail.reward.amount')" :disabled="commentRewardPendingId === comment.id" />
                      <div class="flex gap-2">
                        <UButton color="neutral" variant="ghost" size="sm" type="button" @click="activeCommentRewardId = 0">{{ $t('common.cancel') }}</UButton>
                        <UButton color="primary" variant="soft" size="sm" type="submit" :loading="commentRewardPendingId === comment.id" :disabled="Number(commentRewardAmount) < 1">
                          {{ $t('catalog.torrents.detail.reward.submit') }}
                        </UButton>
                      </div>
                    </form>
                    <form v-if="activeCommentReportId === comment.id" class="mt-3 grid gap-2 rounded-md bg-slate-50 p-3 dark:bg-slate-950" @submit.prevent="handleCommentReport(comment.id)">
                      <UTextarea v-model="commentReportReason" :rows="2" :placeholder="$t('catalog.torrents.detail.report.reason')" :disabled="reportPending" />
                      <div class="flex justify-end gap-2">
                        <UButton color="neutral" variant="ghost" size="xs" type="button" @click="activeCommentReportId = 0">{{ $t('common.cancel') }}</UButton>
                        <UButton color="error" variant="soft" size="xs" type="submit" :loading="reportPending" :disabled="commentReportReason.trim().length < 5">
                          {{ $t('catalog.torrents.detail.report.submit') }}
                        </UButton>
                      </div>
                    </form>
                  </div>
                </div>
              </article>
            </div>
          </section>
        </main>

        <aside class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900 lg:sticky lg:top-24">
          <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.info.title') }}</h2>
          <dl class="mt-4 space-y-4 text-sm">
            <div class="flex items-center justify-between gap-4">
              <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.category') }}</dt>
              <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ categoryName }}</dd>
            </div>
            <div class="flex items-center justify-between gap-4">
              <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.publisher') }}</dt>
              <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ publisherName }}</dd>
            </div>
            <div class="flex items-center justify-between gap-4">
              <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.publishedAt') }}</dt>
              <dd class="text-right font-medium text-slate-950 dark:text-white">{{ formatDateTime(torrent.createdAt, locale) }}</dd>
            </div>
            <div class="flex items-center justify-between gap-4">
              <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.type') }}</dt>
              <dd class="font-medium text-slate-950 dark:text-white">
                {{ torrent.type === 1 ? $t('catalog.torrents.types.multi') : $t('catalog.torrents.types.single') }}
              </dd>
            </div>
            <div class="flex items-center justify-between gap-4">
              <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.fileCount') }}</dt>
              <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(torrent.fileCount) }}</dd>
            </div>
            <div class="flex items-center justify-between gap-4">
              <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.detail.info.torrentId') }}</dt>
              <dd class="font-medium text-slate-950 dark:text-white">#{{ torrent.id }}</dd>
            </div>
          </dl>

          <div class="mt-6 border-t border-slate-200 pt-5 dark:border-slate-800">
            <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.actions.title') }}</h3>
            <div class="mt-3 grid gap-2">
              <UButton
                v-if="canEditTorrent"
                color="neutral"
                variant="outline"
                icon="i-lucide-pen-line"
                block
                :to="localePath(`/catalog/torrents/${torrent.id}/edit`)"
              >
                {{ $t('catalog.torrents.detail.actions.edit') }}
              </UButton>
              <UButton color="neutral" variant="outline" icon="i-lucide-flag" block @click="showTorrentReport = !showTorrentReport">
                {{ $t('catalog.torrents.detail.actions.report') }}
              </UButton>
              <UButton
                v-if="isStaff"
                color="error"
                variant="soft"
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
          </div>

          <div class="mt-6 border-t border-slate-200 pt-5 dark:border-slate-800">
            <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.detail.reward.title') }}</h3>
            <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
              {{ $t('catalog.torrents.detail.reward.summary', { count: numberFormatter.format(rewardTotal) }) }}
            </p>
            <form class="mt-3 flex gap-2" @submit.prevent="handleTorrentReward">
              <UInput v-model="torrentRewardAmount" type="number" min="1" step="1" :placeholder="$t('catalog.torrents.detail.reward.amount')" :disabled="torrentRewardPending" />
              <UButton type="submit" color="primary" icon="i-lucide-coins" :loading="torrentRewardPending" :disabled="Number(torrentRewardAmount) < 1">
                {{ $t('catalog.torrents.detail.reward.submit') }}
              </UButton>
            </form>

            <div v-if="rewardsError" class="mt-3 text-sm text-red-600 dark:text-red-300">{{ rewardsError }}</div>
            <div v-else-if="rewardsPending && rewards.length === 0" class="mt-4 space-y-2">
              <div v-for="item in 3" :key="item" class="h-9 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>
            <div v-else-if="rewards.length === 0" class="mt-4 text-sm text-slate-500 dark:text-slate-400">
              {{ $t('catalog.torrents.detail.reward.empty') }}
            </div>
            <div v-else class="mt-4 space-y-3">
              <div v-for="reward in rewards" :key="reward.id" class="flex items-center justify-between gap-3 text-sm">
                <span class="min-w-0 truncate text-slate-600 dark:text-slate-300">#{{ reward.userId }}</span>
                <span class="shrink-0 font-medium text-amber-600 dark:text-amber-300">{{ formatBonus(reward.amount) }}</span>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { CatalogCategory, CommentItem, SubtitleItem, TorrentDetail, TorrentFileItem, TorrentPeerItem, TorrentRewardItem } from '~/composables/useCatalogTorrents'

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

type PeerView = 'seeders' | 'leechers' | 'completed'

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
const rewards = ref<TorrentRewardItem[]>([])
const comments = ref<CommentItem[]>([])
const subtitles = ref<SubtitleItem[]>([])
const categories = ref<CatalogCategory[]>([])
const pending = ref(false)
const downloadPending = ref(false)
const likePending = ref(false)
const bookmarkPending = ref(false)
const peersPending = ref(false)
const rewardsPending = ref(false)
const commentsPending = ref(false)
const subtitlesPending = ref(false)
const torrentRewardPending = ref(false)
const commentSubmitPending = ref(false)
const subtitleUploadPending = ref(false)
const peersLoaded = ref(false)
const errorMessage = ref('')
const peersError = ref('')
const rewardsError = ref('')
const commentsError = ref('')
const subtitlesError = ref('')
const expandedFileNodeIds = ref(new Set<string>())
const filesLoaded = ref(false)
const filesPending = ref(false)
const filesError = ref('')
const rewardTotal = ref(0)
const commentTotal = ref(0)
const subtitleTotal = ref(0)
const selectedSubtitleFile = ref<File | null>(null)
const subtitleFileInputKey = ref(0)
const subtitleDownloadPendingId = ref(0)
const commentLikePendingId = ref(0)
const commentRewardPendingId = ref(0)
const activeCommentRewardId = ref(0)
const activeCommentReportId = ref(0)
const activeSubtitleReportId = ref(0)
const activePeerView = ref<PeerView | null>(null)
const reportPending = ref(false)
const adminActionPending = ref('')
const showTorrentReport = ref(false)
const torrentRewardAmount = ref('10')
const commentRewardAmount = ref('5')
const torrentReportReason = ref('')
const commentReportReason = ref('')
const subtitleReportReason = ref('')

const commentForm = reactive({
  content: ''
})

const subtitleForm = reactive({
  language: 'zh-CN'
})

const subtitleLanguageOptions = [
  { value: 'zh-CN', label: '简体中文' },
  { value: 'zh-TW', label: '繁體中文' },
  { value: 'en-US', label: 'English' },
  { value: 'ja-JP', label: '日本語' },
  { value: 'ko-KR', label: '한국어' },
  { value: 'other', label: 'Other' }
]

const fileTreeRootId = 'dir:/'
const torrentId = computed(() => readRouteId())
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const fileTree = computed(() => filesLoaded.value ? buildFileTree(files.value) : buildFileTreePlaceholder())
const fileTreeDirectoryIds = computed(() => collectFileTreeDirectoryIds(fileTree.value))
const hasFileTreeDirectories = computed(() => fileTreeDirectoryIds.value.length > 0)
const allFileTreeExpanded = computed(() => hasFileTreeDirectories.value && fileTreeDirectoryIds.value.every((id) => expandedFileNodeIds.value.has(id)))
const visibleFileTreeRows = computed(() => flattenFileTree(fileTree.value, expandedFileNodeIds.value))
const canUploadSubtitle = computed(() => Boolean(selectedSubtitleFile.value && subtitleForm.language && !subtitleUploadPending.value))
const canEditTorrent = computed(() => Boolean(torrent.value && (isStaff.value || user.value?.id === torrent.value.ownerId)))
const peerSeederCount = computed(() => peers.value.filter((item) => item.isSeeder).length)
const peerLeecherCount = computed(() => peers.value.filter((item) => !item.isSeeder).length)
const visiblePeers = computed(() => {
  if (activePeerView.value === 'seeders') return peers.value.filter((item) => item.isSeeder)
  if (activePeerView.value === 'leechers') return peers.value.filter((item) => !item.isSeeder)
  return []
})
const activePeerPanelTitle = computed(() => {
  if (activePeerView.value === 'seeders') return t('catalog.torrents.detail.peers.seedersTitle')
  if (activePeerView.value === 'leechers') return t('catalog.torrents.detail.peers.leechersTitle')
  if (activePeerView.value === 'completed') return t('catalog.torrents.detail.peers.completedTitle')
  return ''
})
const activePeerPanelSummary = computed(() => {
  if (!torrent.value) return ''
  if (activePeerView.value === 'completed') {
    return t('catalog.torrents.detail.peers.completedSummary', {
      count: numberFormatter.value.format(torrent.value.snatched)
    })
  }

  const seeders = peersLoaded.value ? peerSeederCount.value : torrent.value.seeders
  const leechers = peersLoaded.value ? peerLeecherCount.value : torrent.value.leechers
  return t('catalog.torrents.detail.peers.summary', {
    seeders: numberFormatter.value.format(seeders),
    leechers: numberFormatter.value.format(leechers)
  })
})
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
  const ownerName = torrent.ownerName || `#${torrent.ownerId}`
  if (torrent.anonymous) {
    return torrent.ownerId > 0 ? t('catalog.torrents.anonymousOwner', { name: ownerName }) : t('catalog.torrents.anonymous')
  }
  return torrent.ownerId > 0 ? ownerName : '-'
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

function readRouteId() {
  const raw = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id
  const id = Number(raw)
  return Number.isInteger(id) && id > 0 ? id : 0
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
    void Promise.all([loadRewards(), loadComments(), loadSubtitles()])
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
  }
}

function resetInteractionState() {
  peers.value = []
  peersLoaded.value = false
  rewards.value = []
  comments.value = []
  subtitles.value = []
  rewardTotal.value = 0
  commentTotal.value = 0
  subtitleTotal.value = 0
  peersError.value = ''
  rewardsError.value = ''
  commentsError.value = ''
  subtitlesError.value = ''
  selectedSubtitleFile.value = null
  subtitleFileInputKey.value += 1
  activeCommentRewardId.value = 0
  activeCommentReportId.value = 0
  activeSubtitleReportId.value = 0
  activePeerView.value = null
  showTorrentReport.value = false
  commentForm.content = ''
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
  activePeerView.value = view
  if (view !== 'completed') {
    await loadPeers()
  }
}

function peerStatCardClass(view: PeerView) {
  const base = 'rounded-md border px-3 py-3 text-left transition focus:outline-none focus:ring-2 focus:ring-sky-200 dark:focus:ring-sky-900'
  if (activePeerView.value === view) {
    return `${base} border-sky-400 bg-sky-50 dark:border-sky-700 dark:bg-sky-950/40`
  }
  return `${base} border-slate-200 hover:border-sky-300 hover:bg-slate-50 dark:border-slate-800 dark:hover:border-sky-800 dark:hover:bg-slate-950`
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

async function loadRewards() {
  if (torrentId.value <= 0 || rewardsPending.value) return

  rewardsPending.value = true
  rewardsError.value = ''
  try {
    const data = await catalogTorrents.listRewards(torrentId.value, 1, 20)
    rewards.value = data.list || []
    rewardTotal.value = data.total || 0
  } catch (error) {
    rewards.value = []
    rewardTotal.value = 0
    rewardsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    rewardsPending.value = false
  }
}

async function handleTorrentReward() {
  if (!torrent.value || torrentRewardPending.value) return

  const amount = Number(torrentRewardAmount.value)
  if (!Number.isFinite(amount) || amount < 1) return

  torrentRewardPending.value = true
  try {
    await catalogTorrents.rewardTorrent(torrent.value.id, amount)
    toast.add({
      title: t('catalog.torrents.detail.reward.success'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await loadRewards()
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    torrentRewardPending.value = false
  }
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
    const data = await catalogTorrents.listComments(torrentId.value, 1, 20)
    comments.value = data.list || []
    commentTotal.value = data.total || 0
  } catch (error) {
    comments.value = []
    commentTotal.value = 0
    commentsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    commentsPending.value = false
  }
}

async function handleCommentSubmit() {
  const content = commentForm.content.trim()
  if (!torrent.value || content.length < 3 || commentSubmitPending.value) return

  commentSubmitPending.value = true
  try {
    await catalogTorrents.createComment(torrent.value.id, content)
    commentForm.content = ''
    toast.add({
      title: t('catalog.torrents.detail.comments.created'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
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

function startCommentReward(commentId: number) {
  activeCommentRewardId.value = activeCommentRewardId.value === commentId ? 0 : commentId
  activeCommentReportId.value = 0
  commentRewardAmount.value = '5'
}

async function handleCommentReward(comment: CommentItem) {
  if (!torrent.value || commentRewardPendingId.value > 0) return

  const amount = Number(commentRewardAmount.value)
  if (!Number.isFinite(amount) || amount < 1) return

  commentRewardPendingId.value = comment.id
  try {
    await catalogTorrents.rewardComment(torrent.value.id, comment.id, amount)
    comment.rewardCount += 1
    activeCommentRewardId.value = 0
    toast.add({
      title: t('catalog.torrents.detail.reward.success'),
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
    commentRewardPendingId.value = 0
  }
}

function startCommentReport(commentId: number) {
  activeCommentReportId.value = activeCommentReportId.value === commentId ? 0 : commentId
  activeCommentRewardId.value = 0
  commentReportReason.value = ''
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

async function loadSubtitles() {
  if (torrentId.value <= 0 || subtitlesPending.value) return

  subtitlesPending.value = true
  subtitlesError.value = ''
  try {
    const data = await catalogTorrents.listSubtitles(torrentId.value, 1, 20)
    subtitles.value = data.list || []
    subtitleTotal.value = data.total || 0
  } catch (error) {
    subtitles.value = []
    subtitleTotal.value = 0
    subtitlesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    subtitlesPending.value = false
  }
}

function handleSubtitleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  selectedSubtitleFile.value = input.files?.[0] || null
}

async function handleSubtitleUpload() {
  if (!torrent.value || !selectedSubtitleFile.value || !subtitleForm.language || subtitleUploadPending.value) return

  subtitleUploadPending.value = true
  try {
    await catalogTorrents.uploadSubtitle(torrent.value.id, selectedSubtitleFile.value, subtitleForm.language)
    selectedSubtitleFile.value = null
    subtitleFileInputKey.value += 1
    toast.add({
      title: t('catalog.torrents.detail.subtitles.uploaded'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await loadSubtitles()
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

function formatBonus(value: number) {
  return numberFormatter.value.format(Number(value || 0))
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
