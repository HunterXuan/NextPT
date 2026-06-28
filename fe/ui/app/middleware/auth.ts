export default defineNuxtRouteMiddleware((to) => {
  const localePath = useLocalePath()
  const { isLoggedIn } = useAuth()

  if (!isLoggedIn.value) {
    return navigateTo({
      path: localePath('/login'),
      query: { redirect: to.fullPath }
    })
  }
})
