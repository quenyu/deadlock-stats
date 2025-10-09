/**
 * Match-related validation schemas
 */

import { z } from 'zod'
import { primitiveSchemas } from '../base'

/**
 * Match result enum
 */
export const matchResultSchema = z.enum(['win', 'loss'])

/**
 * Match result type
 */
export type MatchResult = z.infer<typeof matchResultSchema>

/**
 * Player match stats schema
 */
export const playerMatchStatsSchema = z.object({
  matchId: primitiveSchemas.nonEmptyString,
  steamId: primitiveSchemas.nonEmptyString,
  heroName: primitiveSchemas.nonEmptyString,
  kills: primitiveSchemas.nonNegativeInt,
  deaths: primitiveSchemas.nonNegativeInt,
  assists: primitiveSchemas.nonNegativeInt,
  netWorth: primitiveSchemas.nonNegativeInt,
  denies: primitiveSchemas.nonNegativeInt,
  lastHits: primitiveSchemas.nonNegativeInt,
  damage: primitiveSchemas.nonNegativeInt,
  healing: primitiveSchemas.nonNegativeInt,
  result: matchResultSchema,
  duration: primitiveSchemas.positiveInt,
  createdAt: primitiveSchemas.timestamp,
})

/**
 * Player match stats type
 */
export type PlayerMatchStats = z.infer<typeof playerMatchStatsSchema>

/**
 * Match schema
 */
export const matchSchema = z.object({
  matchId: primitiveSchemas.nonEmptyString,
  startTime: primitiveSchemas.timestamp,
  duration: primitiveSchemas.positiveInt,
  winningTeam: z.enum(['team1', 'team2']),
  players: z.array(playerMatchStatsSchema),
  createdAt: primitiveSchemas.timestamp,
})

/**
 * Match type
 */
export type Match = z.infer<typeof matchSchema>

/**
 * Recent matches response schema
 */
export const recentMatchesResponseSchema = z.object({
  matches: z.array(playerMatchStatsSchema),
  total: primitiveSchemas.nonNegativeInt,
})

/**
 * Recent matches response type
 */
export type RecentMatchesResponse = z.infer<typeof recentMatchesResponseSchema>

