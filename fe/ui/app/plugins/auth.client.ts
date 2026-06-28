export default defineNuxtPlugin(async () => {
  const { isLoggedIn, fetchUser } = useAuth()

  if (isLoggedIn.value) {
    try {
      await fetchUser()
    } catch {
      // The API layer clears the token on 401. Other failures should not block hydration.
    }
  }
})
