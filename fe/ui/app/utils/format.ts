import type { I18nName } from '~/types/i18n'

export function formatBytes(value?: number | null) {
  const bytes = Number(value || 0)
  if (bytes <= 0) return '0 B'

  const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB']
  const exponent = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1)
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
