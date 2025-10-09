/**
 * User-related validation schemas
 */

import { z } from 'zod'
import { primitiveSchemas } from '../base'

/**
 * User schema
 */
export const userSchema = z.object({
  id: primitiveSchemas.uuid,
  steamId: primitiveSchemas.nonEmptyString,
  personaName: primitiveSchemas.nonEmptyString,
  profileUrl: primitiveSchemas.url.optional(),
  avatar: primitiveSchemas.url.optional(),
  createdAt: primitiveSchemas.timestamp,
  updatedAt: primitiveSchemas.timestamp,
})

/**
 * User type inferred from schema
 */
export type User = z.infer<typeof userSchema>

/**
 * Partial user schema (for updates)
 */
export const partialUserSchema = userSchema.partial()

/**
 * User creation schema
 */
export const createUserSchema = userSchema.omit({
  id: true,
  createdAt: true,
  updatedAt: true,
})

/**
 * User update schema
 */
export const updateUserSchema = userSchema
  .omit({
    id: true,
    steamId: true,
    createdAt: true,
    updatedAt: true,
  })
  .partial()

