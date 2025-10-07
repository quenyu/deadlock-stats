import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { TrendingUp, TrendingDown, Minus } from 'lucide-react'

interface PerformanceSnapshotProps {
  stats: {
    win_rate: number
    kd_ratio: number
    total_matches: number
    performance_dynamics?: {
      win_loss: { trend: string; value: string; sparkline: number[] }
      kda: { trend: string; value: string; sparkline: number[] }
      rank: { trend: string; value: string; sparkline: number[] }
    }
    avg_kills_per_match: number
    avg_deaths_per_match: number
    avg_assists_per_match: number
    avg_match_duration: number
  }
}

export const PerformanceSnapshot = ({ stats }: PerformanceSnapshotProps) => {
  const renderTrend = (trend: string | undefined) => {
    if (!trend) return <Minus className="h-4 w-4 text-gray-400" />
    
    switch (trend) {
      case 'up':
        return <TrendingUp className="h-4 w-4 text-green-500" />
      case 'down':
        return <TrendingDown className="h-4 w-4 text-red-500" />
      default:
        return <Minus className="h-4 w-4 text-gray-400" />
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Performance Snapshot</CardTitle>
        <CardDescription>Recent performance trends and statistics.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <p className="text-sm font-medium text-muted-foreground">Win Rate</p>
            <p className="text-2xl font-bold">{stats.win_rate.toFixed(1)}%</p>
          </div>
          <div className="space-y-2">
            <p className="text-sm font-medium text-muted-foreground">K/D Ratio</p>
            <p className="text-2xl font-bold">{stats.kd_ratio.toFixed(2)}</p>
          </div>
          <div className="space-y-2">
            <p className="text-sm font-medium text-muted-foreground">Total Matches</p>
            <p className="text-2xl font-bold">{stats.total_matches}</p>
          </div>
          <div className="space-y-2">
            <p className="text-sm font-medium text-muted-foreground">Avg Kills/Match</p>
            <p className="text-2xl font-bold">{stats.avg_kills_per_match.toFixed(1)}</p>
          </div>
        </div>

        {stats.performance_dynamics && (
          <div className="space-y-4">
            <h4 className="text-sm font-medium">Recent Trends</h4>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm">Win/Loss</span>
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">{stats.performance_dynamics.win_loss.value}</span>
                  {renderTrend(stats.performance_dynamics.win_loss.trend)}
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">KDA</span>
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">{stats.performance_dynamics.kda.value}</span>
                  {renderTrend(stats.performance_dynamics.kda.trend)}
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">Rank</span>
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">{stats.performance_dynamics.rank.value}</span>
                  {renderTrend(stats.performance_dynamics.rank.trend)}
                </div>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
} 