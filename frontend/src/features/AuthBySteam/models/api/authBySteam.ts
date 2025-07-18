import { api } from '@/shared/api/api'
import { routes } from '@/shared/constants/routes'

export const AuthBySteam = async () => {
    try {
        const response = await api.post(routes.auth.steam)

        return response.data
    } catch (error) {
        console.error(error)
    }
}