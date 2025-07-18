import { create } from 'zustand'
import { devtools } from 'zustand/middleware'
import { User } from '../types/types'
import { fetchCurrentUser } from '../api/fetchCurrentUser'

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
            set({ user: null, isLoading: false, error: null })
            return
          }

          set({ user: response, isLoading: false })
        } catch (error) {
          const errorMessage =
            error instanceof Error ? error.message : 'An unknown error occurred'
          set({ user: null, isLoading: false, error: errorMessage })
        }
      },
      logout: async () => {
        try {
          await fetch('/api/v1/auth/logout', { credentials: 'include' })
        } catch (e) {
          // ignore
        } finally {
          localStorage.removeItem('token')
          set({ user: null })
        }
      },
    }),
    { name: 'UserStore', enabled: import.meta.env.MODE === 'development' }
  )
)

export default useUserStore