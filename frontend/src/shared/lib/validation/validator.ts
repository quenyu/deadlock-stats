/**
 * Validation utilities and helpers
 */

import { z } from 'zod'
import { createLogger } from '../logger'

const log = createLogger('Validator')

/**
 * Validation result type
 */
export type ValidationResult<T> =
  | { success: true; data: T }
  | { success: false; errors: z.ZodError }

/**
 * Safe parse with logging
 */
export function safeParse<T extends z.ZodTypeAny>(
  schema: T,
  data: unknown,
  context?: string
): ValidationResult<z.infer<T>> {
  const result = schema.safeParse(data)

  if (!result.success) {
    log.warn('Validation failed', {
      context,
      errors: result.error.errors,
      data,
    })
    return { success: false, errors: result.error }
  }

  return { success: true, data: result.data }
}

/**
 * Parse and throw on error
 */
export function parse<T extends z.ZodTypeAny>(
  schema: T,
  data: unknown,
  context?: string
): z.infer<T> {
  try {
    return schema.parse(data)
  } catch (error) {
    if (error instanceof z.ZodError) {
      log.error('Validation failed', {
        context,
        errors: error.errors,
        data,
      })
    }
    throw error
  }
}

/**
 * Validate API response
 */
export function validateApiResponse<T extends z.ZodTypeAny>(
  schema: T,
  data: unknown,
  endpoint?: string
): z.infer<T> {
  return parse(schema, data, `API Response: ${endpoint || 'unknown'}`)
}

/**
 * Validate form data
 */
export function validateFormData<T extends z.ZodTypeAny>(
  schema: T,
  data: unknown,
  formName?: string
): ValidationResult<z.infer<T>> {
  return safeParse(schema, data, `Form: ${formName || 'unknown'}`)
}

/**
 * Create validator function
 */
export function createValidator<T extends z.ZodTypeAny>(schema: T) {
  return {
    parse: (data: unknown, context?: string) => parse(schema, data, context),
    safeParse: (data: unknown, context?: string) => safeParse(schema, data, context),
    validate: (data: unknown) => schema.safeParse(data).success,
  }
}

/**
 * Format validation errors for display
 */
export function formatValidationErrors(error: z.ZodError): Record<string, string> {
  const errors: Record<string, string> = {}

  error.errors.forEach((err) => {
    const path = err.path.join('.')
    errors[path] = err.message
  })

  return errors
}

/**
 * Get first validation error message
 */
export function getFirstError(error: z.ZodError): string {
  return error.errors[0]?.message || 'Validation error'
}

