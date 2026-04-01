import { create } from 'zustand'
import { fetchExtendedPlayerProfile } from '@/entities/player/api/fetchExtendedPlayerProfile'
import { ExtendedPlayerProfileDTO } from '@/entities/deadlock/types/types'
import { extractErrorMessage } from '@/shared/lib/errors'

interface ExtendedProfileState {
  profile: ExtendedPlayerProfileDTO | null
  loading: boolean
  error: string | null
  fetchProfile: (steamId: string) => Promise<void>
}

export const useExtendedProfileStore = create<ExtendedProfileState>((set) => ({
  profile: null,
  loading: false,
  error: null,
  fetchProfile: async (steamId: string) => {
    set({ loading: true, error: null })
    try {
      const data = await fetchExtendedPlayerProfile(steamId)
      set({ profile: data, loading: false })
    } catch (error) {
      const errorMessage = extractErrorMessage(error, 'Failed to fetch extended player profile')
      set({ error: errorMessage, loading: false })
    }
  },
})) 