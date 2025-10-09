import { api } from '@/shared/api/api';
import { User } from '../types/types';
import axios from 'axios';
import { logger } from '@/shared/lib/logger';

export const fetchCurrentUser = async () => {
    try {
      const response = await api.get<User>('/users/me');
      
      if (response.status === 200 && response.data) {
        return response.data;
      }
      
      return null;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        logger.error('API Error during user fetch', error, {
          status: error.response?.status,
          data: error.response?.data
        });
        
        if (error.response?.status === 401) {
          logger.info('User not authenticated');
          return null;
        }
      }
      
      const errorMessage =
        error instanceof Error ? error.message : 'An unknown error occurred';
      logger.error('Error fetching user', error);
      throw new Error(errorMessage);
    }
}