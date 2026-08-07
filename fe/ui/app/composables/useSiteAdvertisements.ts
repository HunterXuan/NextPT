import type { SiteAdvertisement, SiteAdvertisementPlacement } from '~/composables/useSite'

type SiteAdvertisementMap = Partial<Record<SiteAdvertisementPlacement, SiteAdvertisement>>

const siteAdvertisementCacheDuration = 5 * 60 * 1000

export function useSiteAdvertisements() {
  const siteApi = useSite()
  const placements = useState<SiteAdvertisementMap>('site:advertisements', () => ({}))
  const loadedAt = useState<number>('site:advertisements:loaded-at', () => 0)
  const pending = useState<boolean>('site:advertisements:pending', () => false)

  async function load() {
    if (pending.value || Date.now() - loadedAt.value < siteAdvertisementCacheDuration) return

    pending.value = true
    try {
      const data = await siteApi.listAdvertisements()
      placements.value = data.placements || {}
      loadedAt.value = Date.now()
    } catch {
      // Advertising should never prevent the page itself from rendering.
    } finally {
      pending.value = false
    }
  }

  function advertisement(placement: SiteAdvertisementPlacement) {
    return computed(() => placements.value[placement] || null)
  }

  return {
    advertisement,
    load
  }
}
