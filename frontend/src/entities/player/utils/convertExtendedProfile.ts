import { ExtendedPlayerProfileDTO } from '@/entities/deadlock/types/types'
import { PlayerProfile, Match, HeroStat, PerformanceDynamics } from '@/entities/player/types/types'

export const convertExtendedToPlayerProfile = (dto: ExtendedPlayerProfileDTO): PlayerProfile => {
  const recentMatches: Match[] = (dto.match_history ?? []).map(m => ({
    id: String(m.match_id),
    hero_name: m.hero_name,
    hero_avatar: m.hero_avatar ?? '',
    result: m.match_result === 1 ? 'Win' : 'Loss',
    player_kills: m.player_kills,
    player_deaths: m.player_deaths,
    player_assists: m.player_assists,
    match_duration_s: m.match_duration_s,
    player_rank_change: m.player_rank_change,
    player_rank_after_match: m.player_rank_after_match,
    rank_name: m.rank_name,
    sub_rank: m.sub_rank,
    rank_image: m.rank_image,
    match_time: Number.isFinite(m.start_time) && m.start_time > 0 ? new Date(m.start_time * 1000).toISOString() : '',
    souls: m.net_worth,
  }))

  const heroStats: HeroStat[] = (dto.hero_stats ?? []).map(h => ({
    hero_name: h.hero_name,
    matches: h.matches_played,
    win_rate: h.win_rate,
    kda: h.kda,
    hero_avatar: h.hero_avatar,
  }))

  const dynamics: PerformanceDynamics = dto.performance_dynamics

  const totalKills = recentMatches.reduce((sum, match) => sum + match.player_kills, 0);
  const totalDeaths = recentMatches.reduce((sum, match) => sum + match.player_deaths, 0);
  const totalAssists = recentMatches.reduce((sum, match) => sum + match.player_assists, 0);

  const profile: PlayerProfile = {
    steam_id: String(dto.card?.account_id ?? ''),
    nickname: dto.nickname,
    avatar_url: dto.avatar_url,
    last_match_time: recentMatches.length ? recentMatches[0].match_time : '',
    last_updated_at: new Date().toISOString(),
    player_rank: dto.card?.ranked_rank ?? 0,
    rank_name: dto.hero_stats.length > 0 ? dto.match_history[0].rank_name : 'Unranked',
    sub_rank: dto.hero_stats.length > 0 ? dto.match_history[0].sub_rank : 0,
    rank_image: dto.rank_image,
    win_rate: dto.win_rate,
    kd_ratio: dto.kd_ratio,
    total_matches: dto.total_matches,
    total_kills: totalKills,
    total_deaths: totalDeaths,
    total_assists: totalAssists,
    max_kills_in_match: recentMatches.length > 0 ? Math.max(...recentMatches.map(m => m.player_kills)) : 0,
    avg_souls_per_min: dto.avg_souls_per_min,
    recent_matches: recentMatches,
    hero_stats: heroStats,
    performance_dynamics: dynamics,
  }
  return profile
} 