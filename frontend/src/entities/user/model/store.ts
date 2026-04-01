import { create } from 'zustand'
import { devtools } from 'zustand/middleware'
import { User } from '../types/types'
import { fetchCurrentUser } from '../api/fetchCurrentUser'
import { API_ENDPOINTS } from '@/shared/constants/api'
import { extractErrorMessage } from '@/shared/lib/errors'

interface UserState {
  user: User | null
  isLoading: boolean
  error: string | null
  fetchUser: () => Promise<void>
  logout: () => void
}

const useUserStore = create<UserState>()(
  devtools(
    (set) => ({
      user: null,
      isLoading: true,
      error: null,
      fetchUser: async () => {
        set({ isLoading: true, error: null })
        try {
          const response = await fetchCurrentUser()

          if (!response) {
            // User not authenticated - this is normal, not an error
            set({ user: null, isLoading: false, error: null })
            return
          }

          set({ user: response as User, isLoading: false })
        } catch (error) {
          // Only set error for unexpected failures (not 401)
          const errorMessage = extractErrorMessage(error, 'Failed to fetch user')
          set({ user: null, isLoading: false, error: errorMessage })
        }
      },
      logout: async () => {
        try {
          await (await import('@/shared/api/api')).api.get(API_ENDPOINTS.auth.logout)
        } catch (e) {
          // ignore
        } finally {
          set({ user: null })
        }
      },
    }),
    { name: 'UserStore', enabled: import.meta.env.MODE === 'development' }
  )
)

export default useUserStore