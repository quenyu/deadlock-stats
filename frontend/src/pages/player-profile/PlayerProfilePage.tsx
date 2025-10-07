import { useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { useExtendedProfileStore } from '@/entities/deadlock/model/store'
import { convertExtendedToPlayerProfile } from '@/entities/player/utils/convertExtendedProfile'
import { PlayerInfoCard } from '@/widgets/player-profile/PlayerInfoCard'
import { PerformanceSnapshot } from '@/widgets/player-profile/PerformanceSnapshot'
import { RecentMatchesTimeline } from '@/widgets/player-profile/RecentMatchesTimeline'
import { HeroStats } from '@/widgets/player-profile/HeroStats'
import { RankHistoryChart } from '@/widgets/player-profile/charts/RankHistoryChart'
import { FeaturedHeroes } from '@/widgets/player-profile/FeaturedHeroes'
import { PersonalRecords } from '@/widgets/player-profile/PersonalRecords'
import { BestMates } from '@/widgets/player-profile/BestMates'
import { HeroMMRChart } from '@/widgets/player-profile/charts/HeroMMRChart'

export const PlayerProfilePage = () => {
  const { steamId } = useParams<{ steamId: string }>()
  const { profile, loading, error, fetchProfile } = useExtendedProfileStore()

  useEffect(() => {
    if (steamId) {
      fetchProfile(steamId)
    }
  }, [steamId, fetchProfile])

  if (loading) {
    return <div className="text-center py-10">Loading profile...</div>
  }

  if (error) {
    return <div className="text-center py-10 text-red-500">Error: {error}</div>
  }

  if (!profile) {
    return <div className="text-center py-10">Player not found.</div>
  }

  const converted = convertExtendedToPlayerProfile(profile)

  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8">
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-8">
          <PlayerInfoCard profile={{
            nickname: profile.nickname || 'Unknown',
            avatar_url: profile.avatar_url || '',
            player_rank: profile.player_rank || 0,
            rank_name: profile.rank_name || 'Unranked',
            rank_image: profile.rank_image || '',
            sub_rank: profile.sub_rank || 0,
            peak_rank: profile.peak_rank,
            peak_rank_name: profile.peak_rank_name,
            peak_rank_image: profile.peak_rank_image
          }} />

          {profile.featured_heroes && profile.featured_heroes.length > 0 && (
            <FeaturedHeroes featuredHeroes={profile.featured_heroes} />
          )}
          
          <RecentMatchesTimeline />
          
          <RankHistoryChart rankHistory={Array.isArray(converted.recent_matches) ? converted.recent_matches.map(match => ({
            match_id: match.id,
            rank: match.player_rank_after_match,
            mmr_score: match.player_score,
            timestamp: match.match_time,
            rank_name: match.rank_name,
            sub_rank: match.sub_rank,
            rank_image: match.rank_image,
          })) : []} />
        </div>
        
        <div className="space-y-8">
          <PerformanceSnapshot stats={{
            win_rate: profile.win_rate || 0,
            kd_ratio: profile.kd_ratio || 0,
            total_matches: profile.total_matches || 0,
            performance_dynamics: profile.performance_dynamics,
            avg_kills_per_match: profile.avg_kills_per_match ?? 0,
            avg_deaths_per_match: profile.avg_deaths_per_match ?? 0,
            avg_assists_per_match: profile.avg_assists_per_match ?? 0,
            avg_match_duration: profile.avg_match_duration ?? 0
          }} />
          
          {profile.personal_records && (
            <PersonalRecords records={profile.personal_records} />
          )}

          {profile.mate_stats && profile.mate_stats.length > 0 && (
            <BestMates mates={profile.mate_stats} />
          )}
          
          <HeroStats heroStats={Array.isArray(converted.hero_stats) ? converted.hero_stats : []} />
          
          {profile.hero_mmr_history && Array.isArray(profile.hero_mmr_history) && profile.hero_mmr_history.length > 0 && (
            <HeroMMRChart heroMMRHistory={profile.hero_mmr_history.map(hero => ({
              ...hero,
              history: Array.isArray(hero.history) ? hero.history.map(point => ({
                ...point,
                match_id: point.match_id || 0,
                start_time: point.start_time || 0,
                player_score: point.player_score || 0,
                rank: point.rank || 0
              })) : []
            }))} />
          )}
        </div>
      </div>
    </div>
  )
} 