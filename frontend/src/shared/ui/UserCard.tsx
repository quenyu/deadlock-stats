import React from 'react'
import { Card, CardContent, CardHeader } from '@/shared/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { Badge } from '@/shared/ui/badge'
import { User } from '@/entities/user'
import { PlayerSearchResult } from '@/shared/lib/validation'

type UserCardData = User | PlayerSearchResult

interface UserCardProps {
  user: UserCardData
  onClick?: () => void
  showExtendedInfo?: boolean
}

// Type guard to check if user is PlayerSearchResult
function isPlayerSearchResult(user: UserCardData): user is PlayerSearchResult {
  return 'steamId' in user && !('steam_id' in user)
}

export const UserCard: React.FC<UserCardProps> = ({ 
  user, 
  onClick, 
  showExtendedInfo = false 
}) => {
  const isPlayer = isPlayerSearchResult(user)
  
  // Normalize data for both types
  const nickname = isPlayer ? user.personaName : user.nickname
  const avatarUrl = isPlayer ? user.avatar : user.avatar_url
  const steamId = isPlayer ? user.steamId : user.steam_id
  const totalMatches = isPlayer ? user.totalMatches : undefined
  const winrate = isPlayer ? user.winrate : undefined
  const rank = isPlayer ? user.rank : undefined
  const lastMatchDate = isPlayer ? user.lastMatchDate : undefined
  
  // Only for User type
  const realname = !isPlayer ? user.realname : undefined
  const countrycode = !isPlayer ? user.countrycode : undefined
  const lastUpdated = !isPlayer ? user.last_updated : undefined
  
  const formatDate = (date?: string) => {
    if (!date) return 'N/A'
    return new Date(date).toLocaleDateString()
  }

  const formatLastUpdated = (timestamp?: number | string) => {
    if (!timestamp) return 'N/A'
    const date = typeof timestamp === 'number' ? new Date(timestamp * 1000) : new Date(timestamp)
    return date.toLocaleDateString()
  }

  return (
    <Card 
      className={`min-w-[220px] max-w-full cursor-pointer transition-all hover:shadow-md ${onClick ? 'hover:scale-[1.02]' : ''}`}
      onClick={onClick}
    >
      <CardHeader className="pb-3">
        <div className="flex items-center space-x-3">
          <Avatar className="h-12 w-12">
            <AvatarImage src={avatarUrl} alt={nickname} />
            <AvatarFallback>
              {nickname?.charAt(0).toUpperCase() || 'U'}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 min-w-0">
            <h3 className="font-semibold text-lg truncate">{nickname}</h3>
            {realname && (
              <p className="text-sm text-muted-foreground truncate">
                {realname}
              </p>
            )}
          </div>
          {countrycode && (
            <Badge variant="secondary" className="text-xs">
              {countrycode}
            </Badge>
          )}
        </div>
      </CardHeader>
      
      <CardContent className="pt-0">
        <div className="space-y-2">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Steam ID:</span>
            <span className="font-mono text-xs break-all max-w-[120px] overflow-hidden inline-block align-bottom">{steamId}</span>
          </div>
          
          {showExtendedInfo && (
            <>
              {!isPlayer && user.account_id && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Account ID:</span>
                  <span>{user.account_id}</span>
                </div>
              )}
              
              {!isPlayer && user.created_at && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Created:</span>
                  <span>{formatDate(user.created_at)}</span>
                </div>
              )}
              
              {lastUpdated && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Last Updated:</span>
                  <span>{formatLastUpdated(lastUpdated)}</span>
                </div>
              )}
              
              {totalMatches !== undefined && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Matches:</span>
                  <span className="font-semibold">{totalMatches}</span>
                </div>
              )}
              
              {winrate !== undefined && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Winrate:</span>
                  <span className="font-semibold text-green-600">{winrate.toFixed(1)}%</span>
                </div>
              )}
              
              {rank && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Rank:</span>
                  <Badge variant="secondary">{rank}</Badge>
                </div>
              )}
              
              {lastMatchDate && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Last Match:</span>
                  <span>{formatLastUpdated(lastMatchDate)}</span>
                </div>
              )}
            </>
          )}
          
          {!isPlayer && (
            <>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Source:</span>
                <Badge variant={user.id ? "default" : "outline"} className="text-xs">
                  {user.id ? "Local" : "Steam"}
                </Badge>
              </div>
              
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Status:</span>
                <Badge>Deadlock Player</Badge>
              </div>
            </>
          )}
        </div>
      </CardContent>
    </Card>
  )
} 