<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between gap-3">
      <p class="text-xs leading-5 text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.tasks.hint') }}</p>
      <UButton type="button" color="neutral" variant="outline" size="sm" icon="i-lucide-plus" :disabled="disabled" @click="addTask">
        {{ $t('admin.site.configs.tasks.add') }}
      </UButton>
    </div>

    <div v-if="tasks.length === 0" class="rounded-md border border-dashed border-slate-300 px-4 py-10 text-center text-sm text-slate-500 dark:border-slate-700 dark:text-slate-400">
      {{ $t('admin.site.configs.tasks.empty') }}
    </div>

    <section v-for="(task, index) in tasks" :key="task.key" class="overflow-hidden rounded-md border border-slate-200 dark:border-slate-800">
      <header class="flex items-center justify-between gap-3 border-b border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950/60">
        <div class="min-w-0">
          <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ task.nameI18n['zh-CN'] || task.key }}</p>
          <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ task.key }}</p>
        </div>
        <div class="flex items-center gap-2">
          <USwitch v-model="task.enabled" :disabled="disabled" />
          <UTooltip :text="$t('common.delete')">
            <UButton type="button" color="error" variant="ghost" size="xs" icon="i-lucide-trash-2" :disabled="disabled" :aria-label="$t('common.delete')" @click="removeTask(index)" />
          </UTooltip>
        </div>
      </header>

      <div class="space-y-3 p-3">
        <div class="grid gap-3 sm:grid-cols-2">
          <UFormField :label="$t('admin.site.configs.tasks.ruleType')">
            <USelect v-model="task.rule.type" class="w-full" size="lg" :items="ruleTypes" value-key="value" :disabled="disabled" />
          </UFormField>
          <UFormField :label="$t('admin.site.configs.tasks.cycle')">
            <USelect v-model="task.cycle" class="w-full" size="lg" :items="cycleOptions" value-key="value" :disabled="disabled" />
          </UFormField>
        </div>

        <div class="grid gap-3 sm:grid-cols-3">
          <UFormField v-for="locale in locales" :key="locale" :label="`${$t('admin.site.configs.tasks.name')} (${locale})`">
            <UInput v-model="task.nameI18n[locale]" class="w-full" size="lg" :disabled="disabled" />
          </UFormField>
        </div>
        <div class="grid gap-3 sm:grid-cols-3">
          <UFormField v-for="locale in locales" :key="locale" :label="`${$t('admin.site.configs.tasks.description')} (${locale})`">
            <UInput v-model="task.descriptionI18n[locale]" class="w-full" size="lg" :disabled="disabled" />
          </UFormField>
        </div>

        <UFormField :label="targetLabel(task.rule.type)">
          <AdminIamRoleLevelSelect
            v-if="task.rule.type === 'iam.role_level_reached'"
            v-model="task.rule.target"
            :roles="normalRoles"
            :disabled="disabled"
          />
          <UInput v-else v-model="task.rule.target" type="number" min="1" step="1" class="w-full" size="lg" :disabled="disabled" />
        </UFormField>

        <div class="border-t border-slate-100 pt-3 dark:border-slate-800">
          <div class="mb-2 flex items-center justify-between gap-3">
            <p class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.tasks.rewards') }}</p>
            <UButton type="button" color="neutral" variant="ghost" size="xs" icon="i-lucide-plus" :disabled="disabled || task.rewards.length >= rewardTypes.length" @click="addReward(task)">
              {{ $t('admin.site.configs.tasks.addReward') }}
            </UButton>
          </div>
          <div class="space-y-2">
            <div v-for="(reward, rewardIndex) in task.rewards" :key="`${reward.type}-${rewardIndex}`" class="grid grid-cols-[minmax(0,1fr)_120px_auto] items-end gap-2">
              <UFormField :label="$t('admin.site.configs.tasks.rewardType')">
                <USelect v-model="reward.type" class="w-full" size="md" :items="availableRewardTypes(task, rewardIndex)" value-key="value" :disabled="disabled" />
              </UFormField>
              <UFormField :label="rewardAmountLabel(reward.type)">
                <UInput v-model="reward.amount" type="number" :min="rewardMinimum(reward.type)" :step="reward.type === 'bonus' ? 0.1 : 1" class="w-full" size="md" :disabled="disabled" />
              </UFormField>
              <UTooltip :text="$t('common.delete')">
                <UButton type="button" color="error" variant="ghost" size="sm" icon="i-lucide-trash-2" :disabled="disabled" :aria-label="$t('common.delete')" @click="task.rewards.splice(rewardIndex, 1)" />
              </UTooltip>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { AdminIamRole } from '~/composables/useAdmin'

type RuleType = 'catalog.torrent_published' | 'tracker.seed_duration' | 'tracker.uploaded' | 'iam.role_level_reached'
type RewardType = 'bonus' | 'vip' | 'invite'
interface TaskForm {
  key: string
  enabled: boolean
  cycle: 'once' | 'weekly' | 'monthly'
  nameI18n: Record<string, string>
  descriptionI18n: Record<string, string>
  rule: { type: RuleType, target: number }
  rewards: { type: RewardType, amount: number }[]
}

const props = withDefaults(defineProps<{ modelValue: unknown, roles: AdminIamRole[], disabled?: boolean }>(), { disabled: false })
const emit = defineEmits<{ 'update:modelValue': [value: unknown] }>()
const { t } = useI18n()
const locales = ['zh-CN', 'zh-TW', 'en-US']
const tasks = ref<TaskForm[]>([])
let resetting = false
let lastSnapshot = ''
const normalRoles = computed(() => props.roles.filter(role => !role.isStaff))
const ruleTypes = computed(() => [
  { value: 'catalog.torrent_published', label: t('admin.site.configs.tasks.rules.torrentPublished') },
  { value: 'tracker.seed_duration', label: t('admin.site.configs.tasks.rules.seedDuration') },
  { value: 'tracker.uploaded', label: t('admin.site.configs.tasks.rules.uploaded') },
  { value: 'iam.role_level_reached', label: t('admin.site.configs.tasks.rules.roleLevel') }
])
const rewardTypes = computed(() => [
  { value: 'bonus', label: t('admin.site.configs.tasks.rewardTypes.bonus') },
  { value: 'vip', label: t('admin.site.configs.tasks.rewardTypes.vip') },
  { value: 'invite', label: t('admin.site.configs.tasks.rewardTypes.invite') }
])
const cycleOptions = computed(() => [
  { value: 'once', label: t('admin.site.configs.tasks.cycles.once') },
  { value: 'weekly', label: t('admin.site.configs.tasks.cycles.weekly') },
  { value: 'monthly', label: t('admin.site.configs.tasks.cycles.monthly') }
])

watch(() => props.modelValue, reset, { deep: true, immediate: true })
watch(tasks, emitValue, { deep: true })

function reset() {
  const snapshot = JSON.stringify(props.modelValue || [])
  if (snapshot === lastSnapshot) return
  resetting = true
  tasks.value = Array.isArray(props.modelValue) ? props.modelValue.filter(isSupportedTask).map(toTask) : []
  nextTick(() => { resetting = false })
}
function toTask(value: any): TaskForm {
  const task: TaskForm = {
    key: String(value?.key || createKey()), enabled: value?.enabled !== false, cycle: value?.cycle || 'once',
    nameI18n: { ...value?.nameI18n }, descriptionI18n: { ...value?.descriptionI18n },
    rule: { type: value?.rule?.type || 'catalog.torrent_published', target: Number(value?.rule?.target || 1) },
    rewards: Array.isArray(value?.rewards) ? value.rewards.map((reward: any) => ({ type: reward.type, amount: Number(reward.amount || 1) })) : [{ type: 'bonus', amount: 100 }]
  }
  return task
}
function isSupportedTask(value: any) {
  return ['catalog.torrent_published', 'tracker.seed_duration', 'tracker.uploaded', 'iam.role_level_reached'].includes(value?.rule?.type)
}
function addTask() { tasks.value.push(toTask({})) }
function removeTask(index: number) { tasks.value.splice(index, 1) }
function targetLabel(type: RuleType) {
  if (type === 'catalog.torrent_published') return t('admin.site.configs.tasks.targets.torrents')
  if (type === 'tracker.seed_duration') return t('admin.site.configs.tasks.targets.hours')
  if (type === 'tracker.uploaded') return t('admin.site.configs.tasks.targets.gib')
  return t('admin.site.configs.tasks.targets.role')
}
function addReward(task: TaskForm) {
  const type = rewardTypes.value.map(item => item.value as RewardType).find(type => !task.rewards.some(reward => reward.type === type))
  if (type) task.rewards.push({ type, amount: type === 'bonus' ? 100 : 1 })
}
function availableRewardTypes(task: TaskForm, index: number) { return rewardTypes.value.filter(option => option.value === task.rewards[index].type || !task.rewards.some((reward, rewardIndex) => rewardIndex !== index && reward.type === option.value)) }
function rewardAmountLabel(type: RewardType) { return t(`admin.site.configs.tasks.amounts.${type}`) }
function rewardMinimum(type: RewardType) { return type === 'bonus' ? 0.1 : 1 }
function createKey() { return `task_${Math.random().toString(36).slice(2, 10)}` }
function emitValue() {
  if (resetting) return
  const value = tasks.value.map(task => ({ ...task, nameI18n: trimI18n(task.nameI18n), descriptionI18n: trimI18n(task.descriptionI18n), rule: { ...task.rule, target: Number(task.rule.target) }, rewards: task.rewards.map(reward => ({ ...reward, amount: Number(reward.amount) })) }))
  lastSnapshot = JSON.stringify(value)
  emit('update:modelValue', value)
}
function trimI18n(values: Record<string, string>) { return Object.fromEntries(Object.entries(values).map(([locale, value]) => [locale, String(value || '').trim()]).filter(([, value]) => value)) }
function validate() {
  if (tasks.value.some(task => !task.nameI18n['zh-CN'] || task.rewards.length === 0 || task.rewards.some(reward => !Number.isFinite(reward.amount) || reward.amount <= 0))) return { valid: false, message: t('admin.site.configs.tasks.invalid') }
  return { valid: true, value: tasks.value }
}
defineExpose({ validate })
</script>
