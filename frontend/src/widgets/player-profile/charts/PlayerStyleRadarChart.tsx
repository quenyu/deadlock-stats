import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import {
  PolarGrid,
  PolarAngleAxis,
  Radar,
  RadarChart,
  ResponsiveContainer,
  Tooltip,
  Legend,
} from 'recharts'
import type { Match, PlayerProfile } from '@/entities/player'

interface PlayerStyleRadarChartProps {
  matches: Match[]
  stats: PlayerProfile
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="rounded-lg border bg-background p-2 shadow-sm">
        <div className="grid grid-cols-2 gap-2">
          <div className="flex flex-col">
            <span className="text-[0.70rem] uppercase text-muted-foreground">
              {label}
            </span>
            <span className="font-bold text-muted-foreground">
              {payload[0].value.toFixed(2)}
            </span>
          </div>
        </div>
      </div>
    );
  }

  return null
}

export const PlayerStyleRadarChart = ({
  matches,
  stats,
}: PlayerStyleRadarChartProps) => {
  if (
    !stats ||
    !matches ||
    matches.length === 0 ||
    stats.avg_souls_per_min === 0
  ) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Player Style Analysis</CardTitle>
        </CardHeader>
        <CardContent>
          <p>Not enough data to display player style.</p>
        </CardContent>
      </Card>
    )
  }

  const totalMatches = matches.length
  const avgKills =
    matches.reduce((sum: number, match: Match) => sum + match.player_kills, 0) /
    totalMatches
  const avgAssists =
    matches.reduce(
      (sum: number, match: Match) => sum + match.player_assists,
      0
    ) / totalMatches
  const avgDeaths =
    matches.reduce(
      (sum: number, match: Match) => sum + match.player_deaths,
      0
    ) / totalMatches

  const dataPoints = [
    { subject: 'Aggression', value: (avgKills + avgAssists / 4) * 5.5 },
    { subject: 'Support', value: (avgAssists - avgAssists / 4) * 5 },
    { subject: 'Durability', value: (1 / (avgDeaths || 1)) * 150 },
    { subject: 'Farming', value: stats.avg_souls_per_min / 12 },
  ]

  const maxValue = Math.max(...dataPoints.map(p => p.value), 100)

  const radarData = dataPoints.map(p => ({
    ...p,
    value: (p.value / maxValue) * 100,
  }));


  return (
    <Card>
      <CardHeader>
        <CardTitle>Player Style Analysis</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <RadarChart cx="50%" cy="50%" outerRadius="80%" data={radarData}>
            <PolarGrid />
            <PolarAngleAxis dataKey="subject" />
            <Tooltip content={<CustomTooltip />} />
            <Radar name="Player Style" dataKey="value" stroke="#8884d8" fill="#8884d8" fillOpacity={0.6} />
            <Legend />
          </RadarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
} 