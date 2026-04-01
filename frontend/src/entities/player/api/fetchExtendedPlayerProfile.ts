import { api } from '@/shared/api/api'
import { ExtendedPlayerProfileDTO } from '@/entities/deadlock/types/types'

export const fetchExtendedPlayerProfile = async (steamId: string): Promise<ExtendedPlayerProfileDTO> => {
  const response = await api.get<{profile: ExtendedPlayerProfileDTO, loadTime: number, steamID: string}>(`/players/${steamId}`)
  return response.data.profile
} 