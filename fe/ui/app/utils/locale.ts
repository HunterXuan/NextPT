export function normalizeLocaleCode(code: string | null | undefined, supportedCodes: string[], defaultCode: string) {
  const fallback = supportedCodes.includes(defaultCode) ? defaultCode : supportedCodes[0] || defaultCode

  if (!code) return fallback
  if (supportedCodes.includes(code)) return code
  return fallback
}
