import { ExtendedPlayerProfileDTO } from '@/entities/deadlock/types/types'
import { PlayerProfile, Match, HeroStat } from '@/entities/player/types/types'

export const convertExtendedToPlayerProfile = (dto: ExtendedPlayerProfileDTO): PlayerProfile => {
  const dynamics = dto.performance_dynamics || {
    win_loss: { trend: 'stable', value: '0/0', sparkline: [] },
    kda: { trend: 'stable', value: '0.00 KDA', sparkline: [] },
    rank: { trend: 'stable', value: '0 Rank', sparkline: [] }
  }

  const recentMatches: Match[] = (dto.match_history || []).map(m => ({
    match_id: String(m.match_id),
    hero_id: m.hero_id,
    hero_name: m.hero_name,
    hero_avatar: m.hero_avatar || '',
    result: m.match_result === 1 ? 'Win' : 'Loss',
    player_kills: m.player_kills,
    player_deaths: m.player_deaths,
    player_assists: m.player_assists,
    net_worth: m.net_worth,
    match_duration_s: m.match_duration_s,
    match_result: m.match_result,
    player_team: m.player_team || 0,
    start_time: m.start_time,
    player_rank_change: m.player_rank_change,
    player_rank_after_match: m.player_rank_after_match,
    rank_name: m.rank_name,
    sub_rank: m.sub_rank || 0,
    rank_image: m.rank_image || '',
    kills: m.kills,
    deaths: m.deaths,
    assists: m.assists,
    duration_minutes: m.duration_minutes,
    match_time: Number.isFinite(m.start_time) && m.start_time > 0 ? new Date(m.start_time * 1000).toISOString() : undefined,
  }))

  const heroStats: HeroStat[] = (dto.hero_stats || []).map(h => ({
    hero_id: h.hero_id,
    hero_name: h.hero_name,
    matches_played: h.matches_played,
    win_rate: h.win_rate,
    kda: h.kda,
    hero_avatar: h.hero_avatar,
  }))

  const totalKills = recentMatches.reduce((sum, match) => sum + match.player_kills, 0)
  const totalDeaths = recentMatches.reduce((sum, match) => sum + match.player_deaths, 0)
  const totalAssists = recentMatches.reduce((sum, match) => sum + match.player_assists, 0)

  return {
    steam_id: '',
    nickname: dto.nickname || 'Unknown',
    avatar_url: dto.avatar_url || '',
    profile_url: '',
    created_at: '',
    player_rank: dto.player_rank || 0,
    rank_name: dto.rank_name || 'Unranked',
    rank_image: dto.rank_image || '',
    sub_rank: dto.sub_rank || 0,
    total_matches: dto.total_matches || 0,
    win_rate: dto.win_rate || 0,
    kd_ratio: dto.kd_ratio || 0,
    avg_matches_per_day: 0,
    favorite_hero: '',
    total_kills: totalKills,
    total_deaths: totalDeaths,
    total_assists: totalAssists,
    max_kills_in_match: recentMatches.length > 0 ? Math.max(...recentMatches.map(m => m.player_kills)) : 0,
    avg_damage_per_match: 0,
    avg_objectives_per_match: 0,
    avg_souls_per_min: dto.avg_souls_per_min || 0,
    performance_dynamics: dynamics,
    recent_matches: recentMatches,
    hero_stats: heroStats,
    last_match_time: recentMatches.length > 0 && recentMatches[0].match_time ? recentMatches[0].match_time : '',
    last_updated_at: new Date().toISOString()
  }
} 