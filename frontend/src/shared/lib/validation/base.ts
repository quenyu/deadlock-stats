import { z } from 'zod'

export const primitiveSchemas = {
  /** Non-empty string */
  nonEmptyString: z.string().min(1, 'Field is required'),
  
  /** Email */
  email: z.email('Invalid email format'),
  
  /** URL */
  url: z.url('Invalid URL format'),
  
  /** UUID */
  uuid: z.uuid('Invalid UUID format'),
  
  /** Positive integer */
  positiveInt: z.number().int().positive('Must be positive integer'),
  
  /** Non-negative integer */
  nonNegativeInt: z.number().int().nonnegative('Must be non-negative integer'),
  
  /** Timestamp */
  timestamp: z.number().int().nonnegative(),
  
  /** ISO 8601 datetime string */
  isoDate: z.string().datetime(),
  
  /** ISO date only (YYYY-MM-DD) */
  isoDateOnly: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, 'Invalid date format (YYYY-MM-DD)'),
} as const


export const objectSchemas = {
  /** Pagination params */
  pagination: z.object({
    page: z.number().int().positive().default(1),
    limit: z.number().int().positive().max(100).default(10),
  }),
  
  /** Pagination response metadata */
  // Matches backend: backend/internal/dto/search_result.go SearchResult
  paginationMeta: z.object({
    total_count: z.number().int().nonnegative(),
    page: z.number().int().positive(),
    page_size: z.number().int().positive(),
    total_pages: z.number().int().nonnegative(),
  }),
  
  /** Error response */
  error: z.object({
    error: z.string(),
    message: z.string().optional(),
    details: z.unknown().optional(),
  }),
} as const

/**
 * Nullable helper
 */
export const nullable = <T extends z.ZodTypeAny>(schema: T) => 
  z.union([schema, z.null()])

/**
 * Optional nullable helper
 */
export const optionalNullable = <T extends z.ZodTypeAny>(schema: T) =>
  z.union([schema, z.null(), z.undefined()])

/**
 * Create enum from array
 */
export const createEnum = <T extends readonly [string, ...string[]]>(values: T) =>
  z.enum(values)

/**
 * Create record with specific value type
 */
export const createRecord = <V extends z.ZodTypeAny>(
  valueSchema: V
) => z.record(z.string(), valueSchema)

