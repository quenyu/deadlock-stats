/**
 * Crosshair-related validation schemas
 */

import { z } from 'zod'
import { primitiveSchemas } from '../base'

/**
 * Crosshair settings schema
 * Matches backend: backend/internal/domain/crosshair.go CrosshairSettings
 */
export const crosshairSettingsSchema = z.object({
  color: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format'),
  thickness: z.number().int().min(0),
  length: z.number().int().min(0),
  gap: z.number().int().min(0),
  dot: z.boolean(),
  opacity: z.number().min(0).max(1),
  pipOpacity: z.number().min(0).max(1),
  dotOutlineOpacity: z.number().min(0).max(1),
  hitMarkerDuration: z.number().min(0),
  pipBorder: z.boolean(),
  pipGapStatic: z.boolean(),
})

/**
 * Crosshair settings type
 */
export type CrosshairSettings = z.infer<typeof crosshairSettingsSchema>

/**
 * Crosshair schema (published)
 * Matches backend: backend/internal/domain/crosshair.go Crosshair
 */
export const crosshairSchema = z.object({
  id: primitiveSchemas.uuid,
  author_id: primitiveSchemas.uuid,
  author_name: z.string().optional(),
  author_avatar: z.string().optional(),
  title: primitiveSchemas.nonEmptyString.max(100, 'Title too long'),
  description: z.string().max(500, 'Description too long').optional(),
  settings: crosshairSettingsSchema,
  is_public: z.boolean(),
  likes_count: primitiveSchemas.nonNegativeInt,
  view_count: primitiveSchemas.nonNegativeInt,
  created_at: primitiveSchemas.isoDate,
  updated_at: primitiveSchemas.isoDate,
})

/**
 * Crosshair type
 */
export type Crosshair = z.infer<typeof crosshairSchema>

/**
 * Crosshair creation schema
 */
export const createCrosshairSchema = z.object({
  title: primitiveSchemas.nonEmptyString.max(100, 'Title too long'),
  description: z.string().max(500, 'Description too long').optional(),
  settings: crosshairSettingsSchema,
  is_public: z.boolean().default(true),
})

/**
 * Crosshair creation type
 */
export type CreateCrosshair = z.infer<typeof createCrosshairSchema>

/**
 * Crosshair update schema
 */
export const updateCrosshairSchema = createCrosshairSchema.partial()

/**
 * Published crosshairs response schema
 */
export const publishedCrosshairsResponseSchema = z.object({
  crosshairs: z.array(crosshairSchema),
  total: primitiveSchemas.nonNegativeInt,
})

/**
 * Published crosshairs response type
 */
export type PublishedCrosshairsResponse = z.infer<typeof publishedCrosshairsResponseSchema>

