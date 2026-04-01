import { api } from '@/shared/api/api'
import { CrosshairListItem } from '../types/types'

interface GetCrosshairsResponse {
  crosshairs: CrosshairListItem[]
  total: number
  page: number
  limit: number
}

export const fetchPublishedCrosshairs = async (page = 1, limit = 20): Promise<GetCrosshairsResponse> => {
  const response = await api.get(`/crosshairs?page=${page}&limit=${limit}`)
  return response.data
} 