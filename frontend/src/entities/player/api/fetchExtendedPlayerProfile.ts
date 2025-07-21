import { api } from '@/shared/api/api'
import { ExtendedPlayerProfileDTO } from '@/entities/deadlock/types/types'

export const fetchExtendedPlayerProfile = async (steamId: string): Promise<ExtendedPlayerProfileDTO> => {
  const { data } = await api.get<ExtendedPlayerProfileDTO>(`/players/${steamId}`)
  return data
} 