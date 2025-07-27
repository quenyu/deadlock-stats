import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { type Match } from '@/entities/player'
import { api } from '@/shared/api/api'
import { formatDistanceToNow } from 'date-fns'

const HeroIcon = ({ name }: { name: string }) => {
  return <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center text-xs">{name.substring(0, 2)}</div>
}

export const RecentMatchesTimeline = () => {
  const { steamId } = useParams<{ steamId: string }>()
  const [matches, setMatches] = useState<Match[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchMatches = async () => {
      try {
        setLoading(true)
        const response = await api.get<Match[]>(`/players/${steamId}/matches`)
        setMatches(response.data)
      } catch (err) {
        setError('Failed to load recent matches.')
      } finally {
        setLoading(false)
      }
    }

    if (steamId) {
      fetchMatches()
    }
  }, [steamId])

  if (loading) return <Card><CardHeader><CardTitle>Recent Matches</CardTitle></CardHeader><CardContent>Loading...</CardContent></Card>
  if (error) return <Card><CardHeader><CardTitle>Recent Matches</CardTitle></CardHeader><CardContent>{error}</CardContent></Card>

  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Matches</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-6">
          {Array.isArray(matches) && matches.map((match) => (
            <div key={match.id} className="flex items-center gap-4">
              <HeroIcon name={match.hero_name} />
              <div className="flex-grow">
                <div className="flex justify-between items-center">
                  <p className="font-semibold">{match.hero_name} - <span className={match.result === 'Win' ? 'text-green-500' : 'text-red-500'}>{match.result}</span></p>
                  <p className="text-xs text-muted-foreground">{match.match_time ? formatDistanceToNow(new Date(match.match_time)) + ' ago' : 'Unknown time'}</p>
                </div>
                <p className="text-sm text-muted-foreground">{match.player_kills}/{match.player_deaths}/{match.player_assists} - {match.match_duration_s} min</p>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
} 