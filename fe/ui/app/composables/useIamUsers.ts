export interface PublicUser {
  user: PublicUserAccount
  role: PublicUserRole
  profile: PublicUserProfile
  stat: PublicUserStat
}

export interface PublicUserAccount {
  id: number
  username: string
  createdAt?: string | null
}

export interface PublicUserRole {
  id: number
  name: string
  isStaff: boolean
}

export interface PublicUserProfile {
  avatar?: string
  info?: string
  signature?: string
}

export interface PublicUserStat {
  uploaded: number
  downloaded: number
  shareRatio: number
  seedTime: number
}

interface PublicUserCacheEntry {
  data: PublicUser
  expiresAt: number
}

const publicUserCacheTtl = 5 * 60 * 1000
const publicUserRequests = new Map<string, Promise<PublicUser>>()

export function useIamUsers() {
  const { locale } = useI18n()
  const cache = useState<Record<string, PublicUserCacheEntry>>('iam:public-users', () => ({}))

  function cacheKey(id: number) {
    return `${locale.value}:${id}`
  }

  function getCachedUser(id: number) {
    const entry = cache.value[cacheKey(id)]
    if (!entry || entry.expiresAt <= Date.now()) return null
    return entry.data
  }

  async function getUser(id: number) {
    const cached = getCachedUser(id)
    if (cached) return cached

    const key = cacheKey(id)
    const pending = publicUserRequests.get(key)
    if (pending) return await pending

    const request = fetchApi<PublicUser>(`/api/iam/users/${id}`)
      .then((data) => {
        cache.value = {
          ...cache.value,
          [key]: {
            data,
            expiresAt: Date.now() + publicUserCacheTtl
          }
        }
        return data
      })
      .finally(() => publicUserRequests.delete(key))

    publicUserRequests.set(key, request)
    return await request
  }

  return {
    getCachedUser,
    getUser
  }
}
