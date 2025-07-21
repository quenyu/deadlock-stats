import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { Card, CardContent, CardHeader } from '@/shared/ui/card'
import type { PlayerProfile } from '@/entities/player'
import { RankBadge } from '@/shared/ui/RankBadge/RankBadge'
import { Crown, Swords, Shield } from 'lucide-react'

export const PlayerInfoCard = ({ profile }: { profile: PlayerProfile }) => {
  return (
    <Card>
      <CardHeader className="flex flex-col sm:flex-row items-center gap-6 p-6">
        <Avatar className="h-24 w-24 border-4 border-primary/20">
          <AvatarImage src={profile.avatar_url} alt={profile.nickname} />
          <AvatarFallback className="text-3xl">{profile.nickname.charAt(0)}</AvatarFallback>
        </Avatar>
        <div className="flex-1 text-center sm:text-left">
          <h1 className="text-3xl font-bold">{profile.nickname}</h1>
          <RankBadge rank={profile.rank_name} subrank={profile.sub_rank} rankImage={profile.rank_image} />
        </div>
      </CardHeader>
      <CardContent className="grid grid-cols-1 sm:grid-cols-3 gap-4 text-center p-6 bg-muted/50">
        <div className="flex flex-col items-center">
          <Crown className="w-6 h-6 text-yellow-500 mb-2" />
          <p className="text-2xl font-bold">{profile.win_rate.toFixed(1)}%</p>
          <p className="text-sm text-muted-foreground">Win Rate</p>
        </div>
        <div className="flex flex-col items-center">
          <Swords className="w-6 h-6 text-red-500 mb-2" />
          <p className="text-2xl font-bold">{profile.kd_ratio.toFixed(2)}</p>
          <p className="text-sm text-muted-foreground">KDA Ratio</p>
        </div>
        <div className="flex flex-col items-center">
          <Shield className="w-6 h-6 text-blue-500 mb-2" />
          <p className="text-2xl font-bold">{profile.total_matches}</p>
          <p className="text-sm text-muted-foreground">Total Matches</p>
        </div>
      </CardContent>
    </Card>
  )
} 