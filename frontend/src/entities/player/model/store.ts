import { create } from 'zustand'
import { type PlayerProfile } from '../types/types'
import { fetchExtendedPlayerProfile } from '../api/fetchExtendedPlayerProfile'
import { convertExtendedToPlayerProfile } from '../utils/convertExtendedProfile'
import { extractErrorMessage } from '@/shared/lib/errors'

interface PlayerProfileState {
  profile: PlayerProfile | null
  loading: boolean
  error: string | null
  fetchProfile: (steamId: string) => Promise<void>
}

export const usePlayerProfileStore = create<PlayerProfileState>((set) => ({
  profile: null,
  loading: true,
  error: null,
  fetchProfile: async (steamId: string) => {
    try {
      set({ loading: true, error: null })
      const dto = await fetchExtendedPlayerProfile(steamId)
      const profile = convertExtendedToPlayerProfile(dto)
      set({ profile, loading: false })
    } catch (error) {
      set({ error: extractErrorMessage(error, 'Failed to fetch player profile.'), loading: false })
    }
  },
})) 