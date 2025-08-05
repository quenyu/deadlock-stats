import { api } from '@/shared/api/api'

export const likeCrosshair = async (crosshairId: string): Promise<void> => {
  await api.post(`/crosshairs/${crosshairId}/like`)
}

export const unlikeCrosshair = async (crosshairId: string): Promise<void> => {
  await api.delete(`/crosshairs/${crosshairId}/like`)
} 