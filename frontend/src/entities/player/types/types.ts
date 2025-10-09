// Matches backend: backend/internal/domain/match.go Match
export interface Match {
  match_id: string
  hero_id: number
  player_kills: number
  player_deaths: number
  player_assists: number
  net_worth: number
  match_duration_s: number
  match_result: number // 0 = loss, 1 = win
  player_team: number
  start_time: number
  hero_name: string
  hero_avatar?: string
  player_rank_after_match: number
  rank_name: string
  sub_rank: number
  rank_image: string
  player_rank_change: number
  kills?: number
  deaths?: number
  assists?: number
  duration_minutes?: number
  match_time?: string
  result: string // "Win" or "Loss"
}

// Matches backend: backend/internal/domain/hero_stat.go HeroStat
export interface HeroStat {
  hero_id: number
  hero_name: string
  matches_played: number
  win_rate: number
  kda: number
  hero_avatar?: string
}

export interface Trend {
  trend: 'up' | 'down' | 'stable';
  value: string;
  sparkline: number[];
}

export interface PerformanceDynamics {
  win_loss: Trend;
  kda: Trend;
  rank: Trend;
}

// Matches backend: backend/internal/domain/player_profile.go PlayerProfile
export interface PlayerProfile {
  steam_id: string
  nickname: string
  avatar_url: string
  profile_url: string
  created_at: string
  last_match_time: string
  player_rank: number
  rank_name: string
  sub_rank: number
  rank_image: string
  win_rate: number
  kd_ratio: number
  avg_matches_per_day: number
  favorite_hero: string
  last_updated_at: string
  total_matches: number
  total_kills: number
  total_deaths: number
  total_assists: number
  max_kills_in_match: number
  avg_damage_per_match: number
  avg_objectives_per_match: number
  avg_souls_per_min: number
  recent_matches: Match[]
  hero_stats: HeroStat[]
  performance_dynamics: PerformanceDynamics
}

export interface RankPoint {
  match_id: string;
  rank: number;
  timestamp: string;
  rank_name: string;
  sub_rank?: number;
} 