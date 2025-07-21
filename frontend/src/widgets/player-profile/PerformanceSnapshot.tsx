import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import { type PlayerProfile } from '@/entities/player'
import { Progress } from '@/shared/ui/progress'
import { Badge } from '@/shared/ui/badge'

export const PerformanceSnapshot = ({ stats }: { stats: PlayerProfile }) => {
  const avgKills = (stats.total_kills / stats.total_matches).toFixed(1)
  const avgDeaths = (stats.total_deaths / stats.total_matches).toFixed(1)
  const avgAssists = (stats.total_assists / stats.total_matches).toFixed(1)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Performance Snapshot</CardTitle>
        <CardDescription>Key performance indicators.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center justify-between">
          <p className="text-sm font-medium text-muted-foreground">
            Total Matches
          </p>
          <Badge variant="secondary">{stats.total_matches}</Badge>
        </div>
        <div className="space-y-1">
          <p className="text-sm font-medium text-muted-foreground">Win Rate</p>
          <div className="flex items-center gap-2">
            <Progress value={stats.win_rate} className="w-full" />
            <span className="text-sm font-semibold">
              {stats.win_rate.toFixed(1)}%
            </span>
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4 text-center">
          <div>
            <p className="text-sm text-muted-foreground">KDA Ratio</p>
            <p className="text-2xl font-bold text-green-500">
              {stats.kd_ratio.toFixed(2)}
            </p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Avg. K/D/A</p>
            <p className="text-lg font-semibold">
              {avgKills} / {avgDeaths} / {avgAssists}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
} 