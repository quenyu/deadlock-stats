import { api } from '@/shared/api/api';
import { User } from '../types/types';

export const fetchCurrentUser = async () => {
    try {
      const response = await api.get<User>('/users/me')

      if (response.status === 401) {
        return
      }

        if (response.status !== 200) {
        throw new Error('Failed to fetch user data')
      }

      return response.data
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : 'An unknown error occurred'
      throw new Error(errorMessage)
    }
}