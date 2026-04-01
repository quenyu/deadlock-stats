/**
 * React hook for form validation
 */

import { useState, useCallback } from 'react'
import { z } from 'zod'
import { safeParse, formatValidationErrors, type ValidationResult } from '../validator'

/**
 * Validation hook for forms
 */
export function useValidation<T extends z.ZodTypeAny>(schema: T) {
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validate = useCallback(
    (data: unknown): ValidationResult<z.infer<T>> => {
      const result = safeParse(schema, data)

      if (!result.success) {
        setErrors(formatValidationErrors(result.errors))
        return result
      }

      setErrors({})
      return result
    },
    [schema]
  )

  const clearErrors = useCallback(() => {
    setErrors({})
  }, [])

  const setFieldError = useCallback((field: string, message: string) => {
    setErrors((prev) => ({ ...prev, [field]: message }))
  }, [])

  const clearFieldError = useCallback((field: string) => {
    setErrors((prev) => {
      const next = { ...prev }
      delete next[field]
      return next
    })
  }, [])

  return {
    validate,
    errors,
    clearErrors,
    setFieldError,
    clearFieldError,
    hasErrors: Object.keys(errors).length > 0,
  }
}

