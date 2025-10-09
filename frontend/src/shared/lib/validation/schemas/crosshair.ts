/**
 * Crosshair-related validation schemas
 */

import { z } from 'zod'
import { primitiveSchemas } from '../base'

/**
 * Crosshair settings schema
 */
export const crosshairSettingsSchema = z.object({
  // Dot settings
  dotEnabled: z.boolean(),
  dotSize: z.number().min(0).max(100),
  dotColor: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format'),
  dotOpacity: z.number().min(0).max(1),
  dotOutlineEnabled: z.boolean(),
  dotOutlineSize: z.number().min(0).max(10),
  dotOutlineColor: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format'),
  dotOutlineOpacity: z.number().min(0).max(1),

  // Cross settings
  crossEnabled: z.boolean(),
  crossLength: z.number().min(0).max(100),
  crossThickness: z.number().min(0).max(20),
  crossGap: z.number().min(0).max(50),
  crossColor: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format'),
  crossOpacity: z.number().min(0).max(1),
  crossOutlineEnabled: z.boolean(),
  crossOutlineSize: z.number().min(0).max(10),
  crossOutlineColor: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format'),
  crossOutlineOpacity: z.number().min(0).max(1),

  // Circle settings
  circleEnabled: z.boolean(),
  circleRadius: z.number().min(0).max(100),
  circleThickness: z.number().min(0).max(20),
  circleColor: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format'),
  circleOpacity: z.number().min(0).max(1),
  circleOutlineEnabled: z.boolean(),
  circleOutlineSize: z.number().min(0).max(10),
  circleOutlineColor: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format'),
  circleOutlineOpacity: z.number().min(0).max(1),

  // T-shape settings
  tShapeEnabled: z.boolean(),
  tShapeLength: z.number().min(0).max(100),
  tShapeThickness: z.number().min(0).max(20),
  tShapeGap: z.number().min(0).max(50),
  tShapeColor: z.string().regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format'),
  tShapeOpacity: z.number().min(0).max(1),
})

/**
 * Crosshair settings type
 */
export type CrosshairSettings = z.infer<typeof crosshairSettingsSchema>

/**
 * Crosshair schema (published)
 */
export const crosshairSchema = z.object({
  id: primitiveSchemas.uuid,
  userId: primitiveSchemas.uuid,
  title: primitiveSchemas.nonEmptyString.max(100, 'Title too long'),
  description: z.string().max(500, 'Description too long').optional(),
  settings: crosshairSettingsSchema,
  isPublic: z.boolean(),
  likes: primitiveSchemas.nonNegativeInt,
  createdAt: primitiveSchemas.timestamp,
  updatedAt: primitiveSchemas.timestamp,
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
  isPublic: z.boolean().default(true),
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

