/**
 * User-related validation schemas
 */

import { z } from 'zod'
import { primitiveSchemas } from '../base'

/**
 * User schema - matches backend/internal/domain/user.go User
 */
export const userSchema = z.object({
  id: primitiveSchemas.uuid,
  steam_id: primitiveSchemas.nonEmptyString,
  nickname: primitiveSchemas.nonEmptyString,
  profile_url: z.string(),
  avatar_url: z.string(),
  created_at: primitiveSchemas.isoDate,
  updated_at: primitiveSchemas.isoDate,
  
  // Extended fields from UserSearchResult DTO
  account_id: z.number().int().optional(),
  countrycode: z.string().optional(),
  last_updated: z.number().int().optional(),
  realname: z.string().optional(),
  is_deadlock_player: z.boolean().optional(),
  deadlock_status_known: z.boolean().optional(),
})

/**
 * User type inferred from schema (camelCase for new API)
 */
export type User = z.infer<typeof userSchema>
export type UserFromSchema = User // Alias for clarity

/**
 * Partial user schema (for updates)
 */
export const partialUserSchema = userSchema.partial()

/**
 * User creation schema
 */
export const createUserSchema = userSchema.omit({
  id: true,
  created_at: true,
  updated_at: true,
})

/**
 * User update schema
 */
export const updateUserSchema = userSchema
  .omit({
    id: true,
    steam_id: true,
    created_at: true,
    updated_at: true,
  })
  .partial()

