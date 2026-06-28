export interface AuthUser {
  id: number
  username: string
  email: string
  passkey: string
  status: number
  role: number
  roleName: string
  roleLevel: number
  isStaff: boolean
  vipUntil?: string | null
  avatar?: string
  info?: string
  signature?: string
  uploaded: number
  downloaded: number
  bonus: number
  invites: number
  shareRatio: number
  createdAt?: string | null
}

export interface RegisterInput {
  username: string
  email: string
  password: string
  inviteHash?: string
}

export interface ProfileInput {
  avatar?: string
  info?: string
  signature?: string
}

export interface PasswordInput {
  oldPassword: string
  newPassword: string
}

export function useAuth() {
  const token = useCookie<string | null>('nextpt_token', {
    sameSite: 'lax',
    maxAge: 60 * 60 * 24 * 10,
    watch: true
  })
  const user = useState<AuthUser | null>('auth:user', () => null)
  const isLoggedIn = computed(() => Boolean(token.value))
  const isStaff = computed(() => Boolean(user.value?.isStaff))

  setApiAuthToken(token.value)

  function setToken(value: string | null) {
    token.value = value
    setApiAuthToken(value)
  }

  function clearLocalSession() {
    user.value = null
    setToken(null)
  }

  async function login(username: string, password: string) {
    const session = await fetchApi<{ token: string }>('/api/iam/sessions', {
      method: 'POST',
      body: { username, password }
    })

    setToken(session.token)
    await fetchUser()
    return user.value
  }

  async function register(input: RegisterInput) {
    return await fetchApi<{ id: number }>('/api/iam/users', {
      method: 'POST',
      body: input
    })
  }

  async function fetchUser() {
    if (!token.value) {
      user.value = null
      return null
    }

    const data = await fetchApi<AuthUser>('/api/iam/users/me')
    user.value = data
    return data
  }

  async function updateProfile(input: ProfileInput) {
    await fetchApi('/api/iam/users/me', {
      method: 'PATCH',
      body: input
    })
    await fetchUser()
  }

  async function changePassword(input: PasswordInput) {
    await fetchApi('/api/iam/users/me:changePassword', {
      method: 'POST',
      body: input
    })
    clearLocalSession()
  }

  async function logout(remote = true) {
    if (remote && token.value) {
      try {
        await fetchApi('/api/iam/sessions', { method: 'DELETE' })
      } catch {
        // Local session cleanup still matters if the server token is already gone.
      }
    }

    clearLocalSession()
  }

  return {
    token,
    user,
    isLoggedIn,
    isStaff,
    login,
    register,
    fetchUser,
    updateProfile,
    changePassword,
    logout,
    clearLocalSession
  }
}
