import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { API_BASE_URL } from '../constants/api'
import { toast } from 'sonner'
import { ApiError, type ApiErrorData } from '../lib/ApiError'

export const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  timeout: 30000, // 30 seconds
})

api.interceptors.request.use(
  (config) => {
    config.headers['X-Request-ID'] = crypto.randomUUID()
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const status = error.response?.status
    const config = error.config as InternalAxiosRequestConfig & { _retry?: boolean; _retryCount?: number }

    // Handle 401 Unauthorized
    if (status === 401) {
      localStorage.removeItem('token')
      toast.error('Session expired. Please log in again.')
      // Optionally redirect to login
      // window.location.href = '/login'
    }

    // Handle 403 Forbidden
    if (status === 403) {
      toast.error('Access denied. You do not have permission to perform this action.')
    }

    // Handle 429 Rate Limit
    if (status === 429) {
      const retryAfter = error.response?.headers['retry-after']
      toast.warning(`Too many requests. Please try again ${retryAfter ? `in ${retryAfter} seconds` : 'later'}.`)
    }

    // Handle 500+ Server Errors
    if (status && status >= 500) {
      toast.error('Server error. Our team has been notified. Please try again later.')
    }

    // Retry logic for GET requests (idempotent)
    if (
      config &&
      config.method === 'get' &&
      !config._retry &&
      status &&
      status >= 500 &&
      status < 600
    ) {
      const retryCount = config._retryCount || 0
      const maxRetries = 2

      if (retryCount < maxRetries) {
        config._retry = true
        config._retryCount = retryCount + 1

        // Exponential backoff: 1s, 2s, 4s
        const delay = Math.pow(2, retryCount) * 1000
        await new Promise((resolve) => setTimeout(resolve, delay))

        return api.request(config)
      }
    }

    // Convert to ApiError for consistent error handling
    const apiError = ApiError.fromAxiosError(error as AxiosError<ApiErrorData>)
    return Promise.reject(apiError)
  }
)