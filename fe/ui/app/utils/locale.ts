export function normalizeLocaleCode<T extends string>(code: string | null | undefined, supportedCodes: readonly T[], defaultCode: T): T {
  const fallback = supportedCodes.includes(defaultCode) ? defaultCode : supportedCodes[0] || defaultCode

  if (!code) return fallback
  return supportedCodes.find(supportedCode => supportedCode === code) ?? fallback
}
