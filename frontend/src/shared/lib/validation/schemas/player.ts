/**
 * Player-related validation schemas
 */

import { z } from 'zod'
import { primitiveSchemas, objectSchemas } from '../base'

/**
 * Player stats schema
 */
export const playerStatsSchema = z.object({
  steamId: primitiveSchemas.nonEmptyString,
  personaName: primitiveSchemas.nonEmptyString,
  avatar: primitiveSchemas.url.optional(),
  totalMatches: primitiveSchemas.nonNegativeInt,
  wins: primitiveSchemas.nonNegativeInt,
  losses: primitiveSchemas.nonNegativeInt,
  winrate: z.number().min(0).max(100),
  kills: primitiveSchemas.nonNegativeInt,
  deaths: primitiveSchemas.nonNegativeInt,
  assists: primitiveSchemas.nonNegativeInt,
  kda: z.number().nonnegative(),
  avgKills: z.number().nonnegative(),
  avgDeaths: z.number().nonnegative(),
  avgAssists: z.number().nonnegative(),
  lastMatchDate: primitiveSchemas.timestamp.optional(),
  rank: primitiveSchemas.nonEmptyString.optional(),
  mmr: primitiveSchemas.nonNegativeInt.optional(),
})

/**
 * Player stats type
 */
export type PlayerStats = z.infer<typeof playerStatsSchema>

/**
 * Player search result schema
 */
export const playerSearchResultSchema = z.object({
  steamId: primitiveSchemas.nonEmptyString,
  personaName: primitiveSchemas.nonEmptyString,
  avatar: primitiveSchemas.url.optional(),
  totalMatches: primitiveSchemas.nonNegativeInt,
  winrate: z.number().min(0).max(100).optional(),
  rank: primitiveSchemas.nonEmptyString.optional(),
  lastMatchDate: primitiveSchemas.timestamp.optional(),
})

/**
 * Player search result type
 */
export type PlayerSearchResult = z.infer<typeof playerSearchResultSchema>

/**
 * Player search response schema
 */
export const playerSearchResponseSchema = z.object({
  players: z.array(playerSearchResultSchema),
  meta: objectSchemas.paginationMeta,
})

/**
 * Player search response type
 */
export type PlayerSearchResponse = z.infer<typeof playerSearchResponseSchema>

/**
 * Player profile schema
 */
export const playerProfileSchema = z.object({
  steamId: primitiveSchemas.nonEmptyString,
  personaName: primitiveSchemas.nonEmptyString,
  profileUrl: primitiveSchemas.url.optional(),
  avatar: primitiveSchemas.url.optional(),
  stats: playerStatsSchema,
  recentMatches: z.array(z.unknown()).optional(), // Define match schema separately
})

/**
 * Player profile type
 */
export type PlayerProfile = z.infer<typeof playerProfileSchema>

