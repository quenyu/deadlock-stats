import { create } from 'zustand'
import { ExtendedPlayerProfileDTO } from '@/entities/deadlock'
import { fetchExtendedPlayerProfile } from '@/entities/player/api/fetchExtendedPlayerProfile'

interface ExtendedProfileState {
  profile: ExtendedPlayerProfileDTO | null
  loading: boolean
  error: string | null
  fetchProfile: (steamId: string) => Promise<void>
}

export const useExtendedProfileStore = create<ExtendedProfileState>((set) => ({
  profile: null,
  loading: true,
  error: null,
  fetchProfile: async (steamId: string) => {
    try {
      set({ loading: true, error: null })
      const data = await fetchExtendedPlayerProfile(steamId)
      set({ profile: data, loading: false })
    } catch (err) {
      set({ error: 'Failed to fetch player profile.', loading: false })
    }
  },
})) 