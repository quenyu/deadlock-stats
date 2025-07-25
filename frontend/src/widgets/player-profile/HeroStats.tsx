import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import type { HeroStat } from '@/entities/player'
import { Progress } from '@/shared/ui/progress'

interface HeroStatsProps {
  heroStats: HeroStat[]
}

export const HeroStats = ({ heroStats }: HeroStatsProps) => {
  if (!heroStats || heroStats.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Hero Performance</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">No hero data available.</p>
        </CardContent>
      </Card>
    )
  }

  const sortedHeroes = [...heroStats].sort((a, b) => b.matches - a.matches).slice(0, 5)
  const maxMatches = Math.max(...sortedHeroes.map(h => h.matches))

  return (
    <Card>
      <CardHeader>
        <CardTitle>Hero Performance</CardTitle>
        <CardDescription>Top 5 most played heroes.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {sortedHeroes.map((hero) => (
          <div key={hero.hero_name}>
            <div className="flex items-center gap-4 mb-2">
              <div className="w-10 flex-shrink-0">
                <img 
                  src={hero.hero_avatar} 
                  alt={hero.hero_name} 
                  className="w-full h-auto rounded-md"
                  style={{ aspectRatio: '280 / 380' }} 
                />
              </div>
              <div className="flex-1">
                <p className="font-semibold">{hero.hero_name}</p>
                <p className="text-sm text-muted-foreground">
                  {hero.win_rate.toFixed(0)}% Win Rate ({hero.matches} Matches)
                </p>
              </div>
              <p className="font-semibold">{hero.kda.toFixed(2)} KDA</p>
            </div>
            <Progress value={(hero.matches / maxMatches) * 100} className="h-2" />
          </div>
        ))}
      </CardContent>
    </Card>
  )
} 