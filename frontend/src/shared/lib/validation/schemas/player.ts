/**
 * Player-related validation schemas
 */

import { z } from 'zod'
import { primitiveSchemas, objectSchemas } from '../base'
import {
  matchSchema,
  heroStatSchema,
  heroStatRawSchema,
  performanceDynamicsSchema,
  deadlockMMRSchema,
  featuredHeroSchema,
  personalRecordsSchema,
  mateStatSchema,
  heroMMRHistorySchema,
} from './match'

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
 * Extended player profile schema (from /players/:steamId endpoint)
 * Matches backend: backend/internal/dto/extended_player_profile.go
 */
export const extendedPlayerProfileSchema = z.object({
  match_history: z.array(matchSchema),
  hero_stats: z.array(heroStatRawSchema),
  mmr_history: z.array(deadlockMMRSchema),
  total_matches: primitiveSchemas.nonNegativeInt,
  win_rate: z.number().min(0).max(100),
  kd_ratio: z.number().nonnegative(),
  performance_dynamics: performanceDynamicsSchema,
  avg_souls_per_min: z.number().nonnegative(),
  player_rank: primitiveSchemas.nonNegativeInt,
  nickname: primitiveSchemas.nonEmptyString,
  avatar_url: z.string(),
  rank_image: z.string(),
  rank_name: z.string(), 
  sub_rank: primitiveSchemas.nonNegativeInt,
  featured_heroes: z.array(featuredHeroSchema),
  peak_rank: primitiveSchemas.nonNegativeInt,
  peak_rank_name: z.string(),
  peak_rank_image: z.string(),
  personal_records: personalRecordsSchema,
  mate_stats: z.array(mateStatSchema),
  hero_mmr_history: z.array(heroMMRHistorySchema),
  last_updated_at: primitiveSchemas.isoDate,
  avg_assists_per_match: z.number().nonnegative(),
  avg_deaths_per_match: z.number().nonnegative(),
  avg_kills_per_match: z.number().nonnegative(),
  avg_match_duration: z.number().nonnegative(),
})

/**
 * Extended player profile type
 */
export type ExtendedPlayerProfile = z.infer<typeof extendedPlayerProfileSchema>

/**
 * Player profile schema (domain model)
 * Matches backend: backend/internal/domain/player_profile.go
 */
export const playerProfileSchema = z.object({
  steam_id: z.string(),
  nickname: primitiveSchemas.nonEmptyString,
  avatar_url: z.string(),
  profile_url: z.string(),
  created_at: primitiveSchemas.isoDate,
  last_match_time: primitiveSchemas.isoDate,
  player_rank: primitiveSchemas.nonNegativeInt,
  rank_name: primitiveSchemas.nonEmptyString,
  sub_rank: primitiveSchemas.nonNegativeInt,
  rank_image: z.string(),
  win_rate: z.number().nonnegative(),
  kd_ratio: z.number().nonnegative(),
  avg_matches_per_day: z.number().nonnegative(),
  favorite_hero: z.string(),
  last_updated_at: primitiveSchemas.isoDate,
  total_matches: primitiveSchemas.nonNegativeInt,
  total_kills: primitiveSchemas.nonNegativeInt,
  total_deaths: primitiveSchemas.nonNegativeInt,
  total_assists: primitiveSchemas.nonNegativeInt,
  max_kills_in_match: primitiveSchemas.nonNegativeInt,
  avg_damage_per_match: z.number().nonnegative(),
  avg_objectives_per_match: z.number().nonnegative(),
  avg_souls_per_min: z.number().nonnegative(),
  recent_matches: z.array(matchSchema),
  hero_stats: z.array(heroStatSchema),
  performance_dynamics: performanceDynamicsSchema,
})

/**
 * Player profile type
 */
export type PlayerProfile = z.infer<typeof playerProfileSchema>

