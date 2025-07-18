import axios, { AxiosRequestHeaders } from 'axios';

const baseURL = import.meta.env.VITE_API_URL;

export const api = axios.create({
	baseURL,
	withCredentials: true,
});

api.interceptors.request.use((config) => {
	const token = localStorage.getItem('token')
	if (token) {
		const headers = (config.headers ?? {}) as AxiosRequestHeaders
		headers.Authorization = token
		config.headers = headers
	}
	return config
})

api.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response?.status === 401) {
			localStorage.removeItem('token')
		}
		return Promise.reject(error)
	}
)