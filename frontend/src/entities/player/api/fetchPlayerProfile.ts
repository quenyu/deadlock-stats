import { api } from '@/shared/api/api'
import { API_ENDPOINTS } from '@/shared/constants/api'
import { type PlayerProfile } from '../types/types'

export const fetchPlayerProfile = async (
  steamId: string
): Promise<PlayerProfile> => {
  const response = await api.get<PlayerProfile>(
    API_ENDPOINTS.player.profile(steamId)
  )
  return response.data
} 