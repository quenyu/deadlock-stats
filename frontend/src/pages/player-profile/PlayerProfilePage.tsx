import { useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { useExtendedProfileStore } from '@/entities/deadlock/model/store'
import { convertExtendedToPlayerProfile } from '@/entities/player/utils/convertExtendedProfile'
import { PlayerInfoCard } from '@/widgets/player-profile/PlayerInfoCard'
import { PerformanceSnapshot } from '@/widgets/player-profile/PerformanceSnapshot'
import { RecentMatchesTimeline } from '@/widgets/player-profile/RecentMatchesTimeline'
import { HeroStats } from '@/widgets/player-profile/HeroStats'
import { PlayerStyleRadarChart } from '@/widgets/player-profile/charts/PlayerStyleRadarChart'
import { RankHistoryChart } from '@/widgets/player-profile/charts/RankHistoryChart'

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
  console.log('Player Profile Data:', converted);

  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8">
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-8">
          <PlayerInfoCard profile={converted} />
          <RecentMatchesTimeline />
          <RankHistoryChart rankHistory={converted.recent_matches.map(match => ({
            match_id: match.id,
            rank: match.player_rank_after_match,
            timestamp: match.match_time,
            rank_name: match.rank_name,
            sub_rank: match.sub_rank,
            rank_image: match.rank_image,
          }))} />
        </div>
        <div className="space-y-8">
          <PerformanceSnapshot stats={converted} />
          <HeroStats heroStats={converted.hero_stats} />
          <PlayerStyleRadarChart matches={converted.recent_matches} stats={converted} />
        </div>
      </div>
    </div>
  )
} 