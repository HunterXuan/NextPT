export const useSiteBranding = () => {
  const config = useRuntimeConfig()

  return {
    siteName: computed(() => config.public.siteName.trim() || 'NextPT'),
    siteLogo: computed(() => config.public.siteLogo.trim()),
    siteLogoDark: computed(() => config.public.siteLogoDark.trim()),
    siteFavicon: computed(() => config.public.siteFavicon.trim() || '/favicon.ico')
  }
}
