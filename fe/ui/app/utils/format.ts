import type { I18nName } from '~/types/i18n'

export function formatBytes(value?: number | null) {
  const bytes = Number(value || 0)
  if (!Number.isFinite(bytes) || bytes < 1) return '0 B'

  const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB']
  const exponent = Math.min(Math.max(0, Math.floor(Math.log(bytes) / Math.log(1024))), units.length - 1)
  const amount = bytes / Math.pow(1024, exponent)

  return `${amount.toFixed(amount >= 10 || exponent === 0 ? 0 : 1)} ${units[exponent]}`
}

export function formatDateTime(value?: string | null, locale?: string) {
  if (!value) return '-'

  try {
    return new Intl.DateTimeFormat(locale || undefined, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    }).format(new Date(value))
  } catch {
    return value
  }
}

export function formatDateOnly(value?: string | null, locale?: string) {
  if (!value) return '-'

  try {
    return new Intl.DateTimeFormat(locale || undefined, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    }).format(new Date(value))
  } catch {
    return value
  }
}

export function formatRelativeDateTime(value?: string | null, locale?: string) {
  const date = parseDateTime(value)
  if (!date) return '-'

  const formatter = new Intl.RelativeTimeFormat(locale || undefined, { numeric: 'auto' })
  const diffSeconds = Math.round((date.getTime() - Date.now()) / 1000)
  const absSeconds = Math.abs(diffSeconds)

  if (absSeconds < 45) return formatter.format(0, 'second')
  if (absSeconds < 45 * 60) return formatter.format(Math.round(diffSeconds / 60), 'minute')
  if (absSeconds < 22 * 60 * 60) return formatter.format(Math.round(diffSeconds / 60 / 60), 'hour')
  if (absSeconds < 30 * 24 * 60 * 60) return formatter.format(Math.round(diffSeconds / 60 / 60 / 24), 'day')
  if (absSeconds < 12 * 30 * 24 * 60 * 60) return formatter.format(Math.round(diffSeconds / 60 / 60 / 24 / 30), 'month')
  return formatter.format(Math.round(diffSeconds / 60 / 60 / 24 / 365), 'year')
}

export function parseDateTime(value?: string | null) {
  if (!value) return null

  const date = new Date(value)
  if (!Number.isNaN(date.getTime())) return date

  const normalizedDate = new Date(value.replace(' ', 'T'))
  return Number.isNaN(normalizedDate.getTime()) ? null : normalizedDate
}

export function localizeI18nName(name: I18nName | null | undefined, locale: string, fallback = '-') {
  if (!name) return fallback

  const preferred = name[locale]
  if (typeof preferred === 'string' && preferred) return preferred

  const zhCN = name['zh-CN']
  if (typeof zhCN === 'string' && zhCN) return zhCN

  for (const value of Object.values(name)) {
    if (typeof value === 'string' && value) return value
  }

  return fallback
}
