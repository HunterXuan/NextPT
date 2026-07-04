import { computed, type Ref } from 'vue'
import type { AdminIamRole } from '~/composables/useAdmin'
import type { I18nName } from '~/types/i18n'
import { localizeI18nName } from '~/utils/format'

export function useAdminIamRoleLevels(roles: Ref<AdminIamRole[]>) {
  const { t, locale } = useI18n()
  const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
  const roleOptions = computed(() => {
    return [...roles.value].sort((a, b) => a.level - b.level || a.id - b.id)
  })
  const roleMapByLevel = computed(() => {
    return new Map(roleOptions.value.map((role) => [role.level, role]))
  })

  function roleName(role: AdminIamRole) {
    return localizeI18nName(normalizeRoleName(role.nameI18N), locale.value, roleLevelLabel(role.level))
  }

  function roleNameWithLevel(role: AdminIamRole) {
    return t('admin.forum.roles.withLevel', {
      name: roleName(role),
      level: numberFormatter.value.format(Number(role.level || 0))
    })
  }

  function roleNameByLevel(level: number) {
    const role = roleMapByLevel.value.get(Number(level || 0))
    return role ? roleNameWithLevel(role) : roleLevelLabel(level)
  }

  function roleLevelLabel(level: number) {
    return t('admin.forum.roles.level', { level: numberFormatter.value.format(Number(level || 0)) })
  }

  function hasRoleLevel(level: number) {
    return roleMapByLevel.value.has(Number(level || 0))
  }

  return {
    roleOptions,
    roleName,
    roleNameWithLevel,
    roleNameByLevel,
    roleLevelLabel,
    hasRoleLevel
  }
}

function normalizeRoleName(value: AdminIamRole['nameI18N']) {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return null
  return value as I18nName
}
