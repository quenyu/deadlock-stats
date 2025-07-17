import { create } from 'zustand'
import { User } from '../types'

interface UserState {
  user: User | null
  isLoading: boolean
  error: string | null
  fetchUser: () => Promise<void>
}

const useUserStore = create<UserState>((set) => ({
  user: null,
  isLoading: true,
  error: null,
  fetchUser: async () => {
    set({ isLoading: true, error: null })
    try {
      const response = await fetch('/api/v1/users/me')

      if (response.status === 401) {
        set({ user: null, isLoading: false, error: null })
        return
      }

      if (!response.ok) {
        throw new Error('Failed to fetch user data')
      }

      const user: User = await response.json()
      set({ user, isLoading: false })
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : 'An unknown error occurred'
      set({ user: null, isLoading: false, error: errorMessage })
    }
  },
}))

export default useUserStore