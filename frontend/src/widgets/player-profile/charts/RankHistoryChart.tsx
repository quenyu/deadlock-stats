import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { ResponsiveContainer, LineChart, CartesianGrid, XAxis, YAxis, Tooltip, Legend, Line } from 'recharts'
import { RankBadge } from '@/shared/ui/RankBadge/RankBadge'

interface RankPoint {
  match_id: string
  rank: number
  timestamp: string
  rank_name: string
  sub_rank?: number
  rank_image?: string
}

interface RankHistoryChartProps {
  rankHistory: RankPoint[]
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload
    return (
      <div className="bg-background border p-3 rounded-lg shadow-lg flex flex-col items-center space-y-2">
        <RankBadge rank={data.rankName} subrank={data.subRank} rankImage={data.rankImage} />
        <p className="text-sm text-muted-foreground">Date: {label}</p>
      </div>
    )
  }
  return null
}

export const RankHistoryChart = ({ rankHistory }: RankHistoryChartProps) => {
  if (!rankHistory || rankHistory.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Rank History</CardTitle>
        </CardHeader>
        <CardContent>
          <p>Not enough match history to display rank changes.</p>
        </CardContent>
      </Card>
    )
  }

  const data = rankHistory.map(point => ({
    name: point.timestamp ? new Date(point.timestamp).toLocaleDateString() : 'Unknown',
    Rank: point.rank,
    rankName: point.rank_name,
    subRank: point.sub_rank,
    rankImage: point.rank_image,
  })).reverse()

  return (
    <Card>
      <CardHeader>
        <CardTitle>Rank History</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Line type="monotone" dataKey="Rank" stroke="#8884d8" />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
} 