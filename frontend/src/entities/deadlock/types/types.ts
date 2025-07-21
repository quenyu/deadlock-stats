export interface DeadlockCardSlot {
  hero: { id: number | null }
}

export interface DeadlockCard {
  account_id: number | null
  ranked_badge_level: number | null
  ranked_rank: number | null
  ranked_subrank: number | null
  slots: DeadlockCardSlot[]
}

export interface DeadlockMMR {
  match_id: number
  start_time: number
  rank: number
  division: number
  division_tier: number
}

export interface MatchRaw {
  match_id: number
  hero_id: number
  player_kills: number
  player_deaths: number
  player_assists: number
  net_worth: number
  match_duration_s: number
  match_result: number
  start_time: number
  hero_name: string
  hero_avatar?: string
  player_rank_after_match: number
  rank_name: string
  sub_rank?: number
  rank_image?: string
  player_rank_change: number

  kills?: number
  deaths?: number
  assists?: number
  duration_minutes?: number
  match_time?: string
  souls?: number
}

export interface HeroStatRaw {
  hero_id: number
  matches_played: number
  wins: number
  kills: number
  deaths: number
  assists: number
  hero_name: string
  hero_avatar?: string
  win_rate: number
  kda: number
}

export interface PerformanceTrend {
  trend: 'up' | 'down' | 'stable'
  value: string
  sparkline: number[]
}

export interface PerformanceDynamics {
  win_loss: PerformanceTrend
  kda: PerformanceTrend
  rank: PerformanceTrend
}

export interface ExtendedPlayerProfileDTO {
  card: DeadlockCard | null
  match_history: MatchRaw[]
  hero_stats: HeroStatRaw[]
  mmr_history: DeadlockMMR[]
  total_matches: number
  win_rate: number
  kd_ratio: number
  performance_dynamics: PerformanceDynamics
  avg_souls_per_min: number

  nickname: string
  avatar_url: string
  rank_image: string
} 