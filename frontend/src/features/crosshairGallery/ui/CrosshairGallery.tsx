import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useCrosshairStore } from '@/entities/crosshair/model/store'
import { LikeButton } from '@/features/crosshairBuilder/ui/LikeButton'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/shared/ui/card'
import { Badge } from '@/shared/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { EyeIcon, CalendarIcon } from 'lucide-react'
import { CrosshairPreview } from '@/features/crosshairBuilder/ui/CrosshairPreview'
import { createLogger } from '@/shared/lib/logger'
import { routes } from '@/shared/constants/routes'

const log = createLogger('CrosshairGallery')

export const CrosshairGallery = () => {
  const published = useCrosshairStore(s => s.published)
  const loading = useCrosshairStore(s => s.loading)
  const loadPublished = useCrosshairStore(s => s.loadPublished)
  const navigate = useNavigate()

  useEffect(() => {
    loadPublished()
  }, [loadPublished])

  const handleCrosshairClick = (crosshairId: string) => {
    navigate(routes.crosshairs.view(crosshairId))
  }

  log.debug('Gallery state', { publishedCount: published.length, loading })
  
  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {Array.from({ length: 6 }).map((_, i) => (
          <Card key={i} className="animate-pulse">
            <CardHeader>
              <div className="h-4 bg-gray-200 rounded w-3/4"></div>
              <div className="h-3 bg-gray-200 rounded w-1/2"></div>
            </CardHeader>
            <CardContent>
              <div className="h-32 bg-gray-200 rounded"></div>
            </CardContent>
            <CardFooter>
              <div className="h-8 bg-gray-200 rounded w-20"></div>
            </CardFooter>
          </Card>
        ))}
      </div>
    )
  }

  if (published.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-muted-foreground text-lg mb-2">
          No published crosshairs yet
        </div>
        <div className="text-sm text-muted-foreground">
          Be the first to share your crosshair!
        </div>
      </div>
    )
  }

  log.info('Rendering crosshair gallery', { count: published.length })

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {published.map((crosshair) => (
        <Card 
          key={crosshair.id} 
          className="hover:shadow-lg transition-shadow cursor-pointer"
          onClick={() => handleCrosshairClick(crosshair.id)}
        >
          <CardHeader>
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <CardTitle className="text-lg mb-1">{crosshair.title}</CardTitle>
                <CardDescription className="line-clamp-2">
                  {crosshair.description || 'No description'}
                </CardDescription>
              </div>
              <div className="flex items-center gap-1">
                <Badge variant="secondary" className="text-xs">
                  <EyeIcon size={12} className="mr-1" />
                  {crosshair.view_count}
                </Badge>
                {!crosshair.is_public && (
                  <Badge variant="outline" className="text-xs">
                    Private
                  </Badge>
                )}
              </div>
            </div>
          </CardHeader>

          <CardContent>
            <div className="relative bg-gray-900 rounded-lg p-4 mb-4">
              <div className="w-full h-32 flex items-center justify-center">
                <svg width="80" height="80" viewBox="0 0 120 120" className="drop-shadow-lg">
                  <CrosshairPreview 
                    {...(typeof crosshair.settings === 'string' ? JSON.parse(crosshair.settings) : crosshair.settings)} 
                    isInteractive={false} 
                  />
                </svg>
              </div>
            </div>

            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Avatar className="w-6 h-6">
                <AvatarImage src={crosshair.author_avatar || `https://api.dicebear.com/7.x/initials/svg?seed=${crosshair.author_id}`} />
                <AvatarFallback className="text-xs">
                  {(crosshair.author_name || crosshair.author_id).charAt(0).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <span>{crosshair.author_name || crosshair.author_id}</span>
              <span>•</span>
              <div className="flex items-center gap-1">
                <CalendarIcon size={12} />
                <span>{new Date(crosshair.created_at).toLocaleDateString()}</span>
              </div>
            </div>
          </CardContent>

          <CardFooter>
            <div onClick={(e) => e.stopPropagation()}>
              <LikeButton
                crosshairId={crosshair.id}
                liked={crosshair.is_liked}
                likesCount={crosshair.likes_count}
              />
            </div>
          </CardFooter>
        </Card>
      ))}
    </div>
  )
} 