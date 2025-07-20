import { api } from '@/shared/api/api'
import { API_ENDPOINTS } from '@/shared/constants/api'

export const AuthBySteam = async () => {
    try {
        const response = await api.post(API_ENDPOINTS.auth.steamLogin)

        return response.data
    } catch (error) {
        console.error(error)
    }
}