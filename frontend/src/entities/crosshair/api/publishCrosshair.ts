import { api } from '@/shared/api/api'
import { CrosshairSettings } from '../types/types'

interface PublishCrosshairRequest {
  title: string
  description: string
  settings: CrosshairSettings
  is_public: boolean
}

interface PublishCrosshairResponse {
  id: string
  title: string
  description: string
  settings: CrosshairSettings
  likes_count: number
  is_public: boolean
  view_count: number
  created_at: string
  updated_at: string
}

export const publishCrosshair = async (data: PublishCrosshairRequest): Promise<PublishCrosshairResponse> => {
  const response = await api.post('/crosshairs', data)
  return response.data
} 