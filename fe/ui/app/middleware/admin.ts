export default defineNuxtRouteMiddleware(async () => {
  const localePath = useLocalePath()
  const { isLoggedIn, user, fetchUser } = useAuth()

  if (!isLoggedIn.value) {
    return navigateTo(localePath('/login'))
  }

  if (!user.value) {
    try {
      await fetchUser()
    } catch {
      return navigateTo(localePath('/login'))
    }
  }

  if (!user.value?.isStaff) {
    return navigateTo(localePath('/iam/users/me'))
  }
})
