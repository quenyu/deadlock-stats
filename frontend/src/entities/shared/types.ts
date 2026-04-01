/**
 * Shared domain models that match backend domain models
 * These models are used across multiple features
 */

// Matches backend: backend/internal/domain/builds.go Build
export interface Build {
  id: string
  author_id: string
  title: string
  description: string
  game_version: string
  is_public: boolean
  view_count: number
  created_at: string
  updated_at: string
}

// Matches backend: backend/internal/domain/comments.go Comment
export interface Comment {
  id: string
  author_id: string
  parent_id: string
  content_type: string
  content_id: string
  body: string
  created_at: string
}

// Matches backend: backend/internal/domain/content_tags.go ContentTag
export interface ContentTag {
  tag_id: number
  content_type: string
  content_id: string
}

// Matches backend: backend/internal/domain/tags.go Tag
export interface Tag {
  id: number
  name: string
}

// Matches backend: backend/internal/domain/votes.go Vote
export interface Vote {
  user_id: string
  content_type: string
  content_id: string
  vote_value: number
  created_at: string
}

// Matches backend: backend/internal/domain/player_stats.go PlayerStats
export interface PlayerStats {
  user_id: string
  kd_ratio: number
  win_rate: number
  avg_matches_per_day: number
  favorite_hero: string
  last_updated_at: string
}

// Matches backend: backend/internal/domain/mate_stat_api.go MateStatAPI
export interface MateStatAPI {
  mate_id: number
  wins: number
  matches_played: number
}

// Matches backend: backend/internal/domain/steam_profile_search.go SteamProfileSearch
export interface SteamProfileSearch {
  account_id: number
  avatar: string
  countrycode: string
  last_updated: number
  personaname: string
  profileurl: string
  realname: string
}

