import React from 'react'
import { Card, CardContent, CardHeader } from '@/shared/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { Badge } from '@/shared/ui/badge'
import { User } from '@/entities/user'

interface UserCardProps {
  user: User
  onClick?: () => void
  showExtendedInfo?: boolean
}

export const UserCard: React.FC<UserCardProps> = ({ 
  user, 
  onClick, 
  showExtendedInfo = false 
}) => {
  const formatDate = (date?: Date) => {
    if (!date) return 'N/A'
    return new Date(date).toLocaleDateString()
  }

  const formatLastUpdated = (timestamp?: number) => {
    if (!timestamp) return 'N/A'
    return new Date(timestamp * 1000).toLocaleDateString()
  }

  return (
    <Card 
      className={`min-w-[220px] max-w-full cursor-pointer transition-all hover:shadow-md ${onClick ? 'hover:scale-[1.02]' : ''}`}
      onClick={onClick}
    >
      <CardHeader className="pb-3">
        <div className="flex items-center space-x-3">
          <Avatar className="h-12 w-12">
            <AvatarImage src={user.avatar_url} alt={user.nickname} />
            <AvatarFallback>
              {user.nickname?.charAt(0).toUpperCase() || 'U'}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 min-w-0">
            <h3 className="font-semibold text-lg truncate">{user.nickname}</h3>
            {user.realname && (
              <p className="text-sm text-muted-foreground truncate">
                {user.realname}
              </p>
            )}
          </div>
          {user.countrycode && (
            <Badge variant="secondary" className="text-xs">
              {user.countrycode}
            </Badge>
          )}
        </div>
      </CardHeader>
      
      <CardContent className="pt-0">
        <div className="space-y-2">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Steam ID:</span>
            <span className="font-mono text-xs break-all max-w-[120px] overflow-hidden inline-block align-bottom">{user.steam_id}</span>
          </div>
          
          {showExtendedInfo && (
            <>
              {user.account_id && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Account ID:</span>
                  <span>{user.account_id}</span>
                </div>
              )}
              
              {user.created_at && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Created:</span>
                  <span>{formatDate(user.created_at)}</span>
                </div>
              )}
              
              {user.last_updated && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Last Updated:</span>
                  <span>{formatLastUpdated(user.last_updated)}</span>
                </div>
              )}
            </>
          )}
          
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
        </div>
      </CardContent>
    </Card>
  )
} 