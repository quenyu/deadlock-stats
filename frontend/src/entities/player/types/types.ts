export interface Match {
  id: string;
  hero_name: string;
  hero_avatar: string;
  result: 'Win' | 'Loss';
  player_kills: number;
  player_deaths: number;
  player_assists: number;
  match_duration_s: number;
  player_rank_change: number;
  player_rank_after_match: number;
  rank_name: string;
  sub_rank: number | undefined;
  rank_image: string | undefined;
  match_time: string;
  souls: number;
  player_score: number;
}

export interface HeroStat {
  hero_name: string;
  matches: number;
  win_rate: number;
  kda: number;
  hero_avatar?: string;
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

export interface PlayerProfile {
  steam_id: string
  nickname: string
  avatar_url: string
  last_match_time: string
  last_updated_at: string
  player_rank: number
  rank_name: string
  sub_rank: number | undefined
  rank_image: string
  win_rate: number
  kd_ratio: number
  total_matches: number
  total_kills: number
  total_deaths: number
  total_assists: number
  max_kills_in_match: number
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