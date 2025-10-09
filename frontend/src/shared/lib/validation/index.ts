/**
 * Validation - Modular Zod validation system
 * 
 * @example
 * ```ts
 * import { userSchema, validateApiResponse } from '@/shared/lib/validation'
 * 
 * // Validate API response
 * const user = validateApiResponse(userSchema, response.data, '/api/users/me')
 * 
 * // Safe parse with error handling
 * const result = safeParse(userSchema, data)
 * if (!result.success) {
 *   console.log(formatValidationErrors(result.errors))
 * }
 * ```
 */

// Export base schemas and utilities
export * from './base'

// Export all schemas and types
export * from './schemas/user'
export * from './schemas/player'
export * from './schemas/crosshair'
export * from './schemas/match'

// Export hooks
export * from './hooks'

// Export validator utilities
export * from './validator'

// Re-export zod for convenience
export { z } from 'zod'

