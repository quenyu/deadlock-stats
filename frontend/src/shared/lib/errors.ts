import axios, { type AxiosError } from 'axios'
import { ApiError, type ApiErrorData } from './ApiError'

export function extractErrorMessage(error: unknown, defaultMessage: string = 'Something went wrong'): string {
  // Handle ApiError instances
  if (ApiError.isApiError(error)) {
    return error.message
  }

  // Handle Axios errors
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<ApiErrorData>
    const dataMessage = axiosError.response?.data?.message
    return dataMessage || error.message || defaultMessage
  }

  // Handle standard Error instances
  if (error instanceof Error) {
    return error.message || defaultMessage
  }

  return defaultMessage
}

export function toApiError(error: unknown): ApiError {
  if (ApiError.isApiError(error)) {
    return error
  }

  if (axios.isAxiosError(error)) {
    return ApiError.fromAxiosError(error as AxiosError<ApiErrorData>)
  }

  if (error instanceof Error) {
    return new ApiError(error.message, 500, 'UNKNOWN_ERROR', undefined, error)
  }

  return new ApiError(
    typeof error === 'string' ? error : 'An unexpected error occurred',
    500,
    'UNKNOWN_ERROR',
    undefined,
    error
  )
}


