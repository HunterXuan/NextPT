export function formatDate(dt: string | null | undefined): string {
    if (!dt) return ''
    try {
        return new Date(dt + '+08:00').toLocaleString()
    } catch (e) {
        return String(dt)
    }
}
