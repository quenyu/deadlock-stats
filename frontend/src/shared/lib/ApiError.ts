import type { AxiosError } from 'axios'

export interface ApiErrorData {
  message?: string
  code?: string
  details?: Record<string, unknown>
}

export class ApiError extends Error {
  public readonly status: number
  public readonly code: string
  public readonly details?: Record<string, unknown>
  public readonly originalError: unknown

  constructor(
    message: string,
    status: number = 500,
    code: string = 'UNKNOWN_ERROR',
    details?: Record<string, unknown>,
    originalError?: unknown
  ) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.details = details
    this.originalError = originalError

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ApiError)
    }
  }

  static fromAxiosError(error: AxiosError<ApiErrorData>): ApiError {
    const status = error.response?.status || 500
    const data = error.response?.data
    const message = data?.message || error.message || 'An unexpected error occurred'
    const code = data?.code || `HTTP_${status}`
    const details = data?.details

    return new ApiError(message, status, code, details, error)
  }

  static isApiError(error: unknown): error is ApiError {
    return error instanceof ApiError
  }

  toJSON() {
    return {
      name: this.name,
      message: this.message,
      status: this.status,
      code: this.code,
      details: this.details,
    }
  }
}

