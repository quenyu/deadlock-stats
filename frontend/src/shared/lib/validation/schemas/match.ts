/**
 * Match-related validation schemas
 * Matches backend: backend/internal/domain/match.go
 */

import { z } from 'zod'
import { primitiveSchemas } from '../base'

/**
 * Match schema - matches backend Match model
 */
export const matchSchema = z.object({
  match_id: z.string(),
  hero_id: z.number().int(),
  player_kills: z.number().int().nonnegative(),
  player_deaths: z.number().int().nonnegative(),
  player_assists: z.number().int().nonnegative(),
  net_worth: z.number().int().nonnegative(),
  match_duration_s: z.number().int().nonnegative(),
  match_result: z.number().int(),
  player_team: z.number().int(),
  start_time: z.number().int(),
  hero_name: z.string(),
  hero_avatar: z.string().optional(),
  player_rank_after_match: z.number().int(),
  rank_name: z.string(),
  sub_rank: z.number().int(),
  rank_image: z.string(),
  player_rank_change: z.number().int(),
  kills: z.number().int().optional(),
  deaths: z.number().int().optional(),
  assists: z.number().int().optional(),
  duration_minutes: z.number().int().optional(),
  match_time: primitiveSchemas.isoDate.optional(),
  result: z.string(),
})

export type Match = z.infer<typeof matchSchema>

/**
 * Hero stat schema - matches backend HeroStat model
 */
export const heroStatSchema = z.object({
  hero_id: z.number().int(),
  hero_name: z.string(),
  matches_played: z.number().int().nonnegative(),
  win_rate: z.number().nonnegative(),
  kda: z.number().nonnegative(),
  hero_avatar: z.string().optional(),
})

export type HeroStat = z.infer<typeof heroStatSchema>

/**
 * Hero stat raw schema (from API with additional fields)
 * Note: wins, kills, deaths, assists are optional because backend might not send them
 */
export const heroStatRawSchema = heroStatSchema.extend({
  wins: z.number().int().nonnegative().optional(),
  kills: z.number().int().nonnegative().optional(),
  deaths: z.number().int().nonnegative().optional(),
  assists: z.number().int().nonnegative().optional(),
})

export type HeroStatRaw = z.infer<typeof heroStatRawSchema>

/**
 * Performance trend schema
 */
export const performanceTrendSchema = z.object({
  trend: z.enum(['up', 'down', 'stable']),
  value: z.string(),
  sparkline: z.array(z.number()),
})

export type PerformanceTrend = z.infer<typeof performanceTrendSchema>

/**
 * Performance dynamics schema
 */
export const performanceDynamicsSchema = z.object({
  win_loss: performanceTrendSchema,
  kda: performanceTrendSchema,
  rank: performanceTrendSchema,
})

export type PerformanceDynamics = z.infer<typeof performanceDynamicsSchema>

/**
 * Deadlock MMR schema
 */
export const deadlockMMRSchema = z.object({
  match_id: z.number().int(),
  rank: z.number().int(),
  start_time: z.number().int(),
  player_score: z.number(),
  division: z.number().int(),
  division_tier: z.number().int(),
})

export type DeadlockMMR = z.infer<typeof deadlockMMRSchema>

/**
 * Featured hero schema
 */
export const featuredHeroSchema = z.object({
  hero_id: z.number().int(),
  hero_name: z.string(),
  hero_image: z.string(),
  kills: z.number().int().optional(),
  wins: z.number().int().optional(),
  stat_id: z.number().int().optional(),
  stat_score: z.number().int().optional(),
})

export type FeaturedHero = z.infer<typeof featuredHeroSchema>

/**
 * Personal records schema
 */
export const personalRecordsSchema = z.object({
  max_kills: z.number().int().nonnegative(),
  max_assists: z.number().int().nonnegative(),
  max_net_worth: z.number().int().nonnegative(),
  best_kda: z.number().nonnegative(),
  max_kills_match_id: z.string(),
  max_assists_match_id: z.string(),
  max_net_worth_match_id: z.string(),
  best_kda_match_id: z.string(),
})

export type PersonalRecords = z.infer<typeof personalRecordsSchema>

/**
 * Mate stat schema
 */
export const mateStatSchema = z.object({
  steam_id: z.string(),
  nickname: z.string(),
  avatar_url: z.string(),
  games: z.number().int().nonnegative(),
  wins: z.number().int().nonnegative(),
  win_rate: z.number().nonnegative(),
})

export type MateStat = z.infer<typeof mateStatSchema>

/**
 * Hero MMR history schema
 */
export const heroMMRHistorySchema = z.object({
  hero_id: z.number().int(),
  hero_name: z.string(),
  history: z.array(deadlockMMRSchema),
})

export type HeroMMRHistory = z.infer<typeof heroMMRHistorySchema>

/**
 * Recent matches response schema
 */
export const recentMatchesResponseSchema = z.object({
  matches: z.array(matchSchema),
  total: z.number().int().nonnegative(),
  page: z.number().int().positive(),
  limit: z.number().int().positive(),
})

export type RecentMatchesResponse = z.infer<typeof recentMatchesResponseSchema>
