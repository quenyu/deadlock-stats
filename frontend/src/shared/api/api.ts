import axios from 'axios';
import { API_BASE_URL } from '../constants/api';

export const api = axios.create({
	baseURL: API_BASE_URL,
	withCredentials: true,
});

api.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response?.status === 401) {
			localStorage.removeItem('token')
		}
		return Promise.reject(error)
	}
)