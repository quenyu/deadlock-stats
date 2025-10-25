import { useParams } from 'react-router-dom'
import { usePlayerProfile } from '@/shared/lib/react-query/hooks'
import { PlayerInfoCard } from '@/widgets/player-profile/PlayerInfoCard'
import { PerformanceSnapshot } from '@/widgets/player-profile/PerformanceSnapshot'
import { RecentMatchesTimeline } from '@/widgets/player-profile/RecentMatchesTimeline'
import { HeroStats } from '@/widgets/player-profile/HeroStats'
import { RankHistoryChart } from '@/widgets/player-profile/charts/RankHistoryChart'
import { FeaturedHeroes } from '@/widgets/player-profile/FeaturedHeroes'
import { PersonalRecords } from '@/widgets/player-profile/PersonalRecords'
import { BestMates } from '@/widgets/player-profile/BestMates'
import { HeroMMRChart } from '@/widgets/player-profile/charts/HeroMMRChart'
import { Skeleton } from '@/shared/ui/skeleton'

export const PlayerProfilePage = () => {
  const { steamId } = useParams<{ steamId: string }>()
  const { data: profile, isLoading: loading, error } = usePlayerProfile(steamId || '')

  if (loading) {
    return (
      <div className="container mx-auto p-4 sm:p-6 lg:p-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-8">
            <Skeleton width="100%" height={200} variant="rounded" />
            <Skeleton width="100%" height={300} variant="rounded" />
          </div>
          <div className="space-y-8">
            <Skeleton width="100%" height={150} variant="rounded" />
            <Skeleton width="100%" height={200} variant="rounded" />
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return <div className="text-center py-10 text-red-500">Error: {error.message}</div>
  }

  if (!profile) {
    return <div className="text-center py-10">Player not found.</div>
  }

  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8">
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-8">
          <PlayerInfoCard profile={{
            nickname: profile.nickname,
            avatar_url: profile.avatar_url,
            player_rank: profile.player_rank,
            rank_name: profile.rank_name,
            rank_image: profile.rank_image,
            sub_rank: profile.sub_rank,
            peak_rank: profile.peak_rank,
            peak_rank_name: profile.peak_rank_name,
            peak_rank_image: profile.peak_rank_image
          }} />

          {profile.featured_heroes.length > 0 && (
            <FeaturedHeroes featuredHeroes={profile.featured_heroes} />
          )}
          
          <RecentMatchesTimeline />
          
          <RankHistoryChart rankHistory={profile.match_history.map(match => ({
            match_id: match.match_id,
            rank: match.player_rank_after_match,
            mmr_score: 0,
            timestamp: match.match_time || '',
            rank_name: match.rank_name,
            sub_rank: match.sub_rank,
            rank_image: match.rank_image,
          }))} />
        </div>
        
        <div className="space-y-8">
          <PerformanceSnapshot stats={{
            win_rate: profile.win_rate,
            kd_ratio: profile.kd_ratio,
            total_matches: profile.total_matches,
            performance_dynamics: profile.performance_dynamics,
            avg_kills_per_match: profile.avg_kills_per_match,
            avg_deaths_per_match: profile.avg_deaths_per_match,
            avg_assists_per_match: profile.avg_assists_per_match,
            avg_match_duration: profile.avg_match_duration
          }} />
          
          <PersonalRecords records={profile.personal_records} />

          {profile.mate_stats.length > 0 && (
            <BestMates mates={profile.mate_stats} />
          )}
          
          <HeroStats heroStats={profile.hero_stats.map(h => ({
            hero_id: h.hero_id,
            hero_name: h.hero_name,
            matches_played: h.matches_played,
            win_rate: h.win_rate,
            kda: h.kda,
            hero_avatar: h.hero_avatar,
          }))} />
          
          {profile.hero_mmr_history.length > 0 && (
            <HeroMMRChart heroMMRHistory={profile.hero_mmr_history} />
          )}
        </div>
      </div>
    </div>
  )
} 