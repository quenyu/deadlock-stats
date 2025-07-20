import { api } from '@/shared/api/api';
import { User } from '../types/types';
import axios from 'axios';

export const fetchCurrentUser = async () => {
    try {
      const response = await api.get<User>('/users/me');
      
      if (response.status === 200 && response.data) {
        return response.data;
      }
      
      return null;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        console.error('API Error:', error.response?.status, error.response?.data);
        
        if (error.response?.status === 401) {
          console.log('User not authenticated');
          return null;
        }
      }
      
      const errorMessage =
        error instanceof Error ? error.message : 'An unknown error occurred';
      console.error('Error fetching user:', errorMessage);
      throw new Error(errorMessage);
    }
}