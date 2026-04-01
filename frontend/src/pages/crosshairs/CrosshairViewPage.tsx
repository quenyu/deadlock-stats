import { useParams } from 'react-router-dom'
import { useState, useRef } from 'react'
import { useCrosshair } from '@/shared/lib/react-query/hooks/useCrosshair'
import { CrosshairPreview } from '@/features/crosshairBuilder/ui/CrosshairPreview'
import { LikeButton } from '@/features/crosshairBuilder/ui/LikeButton'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Badge } from '@/shared/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { Button } from '@/shared/ui/button'
import { ArrowLeft, CalendarIcon, EyeIcon, Copy, Download } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { routes } from '@/shared/constants/routes'
import { createLogger } from '@/shared/lib/logger'
import previewImage from '@/shared/assets/images/preview_1.png'

const log = createLogger('CrosshairViewPage')

export const CrosshairViewPage = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { data: crosshair, isLoading, error } = useCrosshair(id!)
  const [crosshairPosition, setCrosshairPosition] = useState({ x: 0, y: 0 })
  const previewRef = useRef<HTMLDivElement>(null)

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!previewRef.current) return
    
    const rect = previewRef.current.getBoundingClientRect()
    const x = e.clientX - rect.left
    const y = e.clientY - rect.top
    
    setCrosshairPosition({ x, y })
  }

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2 mb-8"></div>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <div className="h-96 bg-gray-200 rounded"></div>
            <div className="space-y-4">
              <div className="h-6 bg-gray-200 rounded w-3/4"></div>
              <div className="h-4 bg-gray-200 rounded w-full"></div>
              <div className="h-4 bg-gray-200 rounded w-2/3"></div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !crosshair) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Crosshair not found</h1>
          <p className="text-muted-foreground mb-6">
            The crosshair you're looking for doesn't exist or has been removed.
          </p>
          <Button onClick={() => navigate(routes.crosshairs.list)}>
            <ArrowLeft size={16} className="mr-2" />
            Back to Gallery
          </Button>
        </div>
      </div>
    )
  }

  const settings = typeof crosshair.settings === 'string' 
    ? JSON.parse(crosshair.settings) 
    : crosshair.settings

  const hexToRgb = (hex: string) => {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
    return result ? {
      r: parseInt(result[1], 16),
      g: parseInt(result[2], 16),
      b: parseInt(result[3], 16)
    } : { r: 0, g: 0, b: 0 }
  }

  const handleCopySettings = () => {
    const rgb = hexToRgb(settings.color)
    
    const commands = [
      `citadel_crosshair_color_r ${rgb.r}`,
      `citadel_crosshair_color_g ${rgb.g}`,
      `citadel_crosshair_color_b ${rgb.b}`,
      `citadel_crosshair_pip_border ${settings.pipBorder ? 'true' : 'false'}`,
      `citadel_crosshair_pip_gap_static ${settings.pipGapStatic ? 'true' : 'false'}`,
      `citadel_crosshair_pip_opacity ${settings.pipOpacity}`,
      `citadel_crosshair_pip_width ${settings.thickness}`,
      `citadel_crosshair_pip_height ${settings.length}`,
      `citadel_crosshair_pip_gap ${settings.gap}`,
      `citadel_crosshair_dot_opacity ${settings.dot ? settings.opacity : 0}`,
      `citadel_crosshair_dot_outline_opacity ${settings.dotOutlineOpacity}`
    ]
    
    const config = commands.join('; ')
    navigator.clipboard.writeText(config)
    // TODO: Add toast notification
  }

  const handleDownloadSettings = () => {
    const settingsJson = JSON.stringify(settings, null, 2)
    const blob = new Blob([settingsJson], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${crosshair.title.replace(/[^a-z0-9]/gi, '_').toLowerCase()}_crosshair.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  log.info('Rendering crosshair view', { id: crosshair.id, title: crosshair.title })

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-6">
        <Button 
          variant="ghost" 
          onClick={() => navigate(routes.crosshairs.list)}
          className="mb-4"
        >
          <ArrowLeft size={16} className="mr-2" />
          Back to Gallery
        </Button>
      </div>

      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <h1 className="text-4xl font-bold mb-2">{crosshair.title}</h1>
            <p className="text-lg text-muted-foreground">
              {crosshair.description || 'No description provided'}
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Badge variant="secondary">
              <EyeIcon size={14} className="mr-1" />
              {crosshair.view_count} views
            </Badge>
            {!crosshair.is_public && (
              <Badge variant="outline">Private</Badge>
            )}
          </div>
        </div>
        
        <div className="flex items-center gap-3 text-sm text-muted-foreground">
          <Avatar className="w-8 h-8">
            <AvatarImage src={crosshair.author_avatar || `https://api.dicebear.com/7.x/initials/svg?seed=${crosshair.author_id}`} />
            <AvatarFallback>
              {(crosshair.author_name || crosshair.author_id).charAt(0).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <span className="font-medium">{crosshair.author_name || `User ${crosshair.author_id.substring(0, 8)}`}</span>
          <span>•</span>
          <div className="flex items-center gap-1">
            <CalendarIcon size={14} />
            <span>{new Date(crosshair.created_at).toLocaleDateString()}</span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Crosshair Preview - Large */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Preview</CardTitle>
            <CardDescription>
              See how your crosshair looks in-game
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div 
              ref={previewRef}
              onMouseMove={handleMouseMove}
              className="relative w-full aspect-[16/9] flex items-center justify-center rounded-2xl overflow-hidden shadow-2xl border-2 border-slate-200 dark:border-slate-700 bg-background/20 backdrop-blur-sm transition-all duration-300 hover:shadow-lg cursor-none"
            >
              <img 
                alt="Deadlock preview" 
                className="w-full h-full object-cover select-none" 
                src={previewImage}
                draggable={false}
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent pointer-events-none" />
              <div className="absolute top-4 right-4 px-3 py-1 bg-black/60 backdrop-blur-sm rounded-md text-xs font-medium text-white pointer-events-none">
                {settings.color}
              </div>
              <svg 
                width="120" 
                height="120" 
                viewBox="0 0 120 120" 
                className="absolute pointer-events-none drop-shadow-[0_0_12px_rgba(0,0,0,0.8)]"
                style={{ 
                  left: `${crosshairPosition.x}px`,
                  top: `${crosshairPosition.y}px`,
                  transform: 'translate(-50%, -50%)',
                  zIndex: 1
                }}
              >
                <CrosshairPreview {...settings} isInteractive={false} />
              </svg>
            </div>
            
            <div className="mt-6 flex gap-3">
              <Button className="flex-1" onClick={handleCopySettings}>
                <Copy size={16} className="mr-2" />
                Copy Config
              </Button>
              <Button variant="outline" onClick={handleDownloadSettings}>
                <Download size={16} className="mr-2" />
                Download JSON
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Settings Panel */}
        <Card>
          <CardHeader>
            <CardTitle>Settings</CardTitle>
            <CardDescription>
              Configuration details
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="space-y-3">
                <div className="flex items-center justify-between py-2 border-b border-slate-200 dark:border-slate-800">
                  <span className="text-sm font-medium text-muted-foreground">Color</span>
                  <div className="flex items-center gap-2">
                    <div 
                      className="w-6 h-6 rounded-full border-2 border-slate-300 dark:border-slate-600" 
                      style={{ backgroundColor: settings.color }}
                    />
                    <span className="text-sm font-mono">{settings.color}</span>
                  </div>
                </div>
                
                <div className="flex items-center justify-between py-2 border-b border-slate-200 dark:border-slate-800">
                  <span className="text-sm font-medium text-muted-foreground">Length</span>
                  <span className="text-sm font-semibold">{settings.length}px</span>
                </div>
                
                <div className="flex items-center justify-between py-2 border-b border-slate-200 dark:border-slate-800">
                  <span className="text-sm font-medium text-muted-foreground">Thickness</span>
                  <span className="text-sm font-semibold">{settings.thickness}px</span>
                </div>
                
                <div className="flex items-center justify-between py-2 border-b border-slate-200 dark:border-slate-800">
                  <span className="text-sm font-medium text-muted-foreground">Gap</span>
                  <span className="text-sm font-semibold">{settings.gap}px</span>
                </div>
                
                <div className="flex items-center justify-between py-2 border-b border-slate-200 dark:border-slate-800">
                  <span className="text-sm font-medium text-muted-foreground">Opacity</span>
                  <span className="text-sm font-semibold">{Math.round(settings.opacity * 100)}%</span>
                </div>
                
                <div className="flex items-center justify-between py-2 border-b border-slate-200 dark:border-slate-800">
                  <span className="text-sm font-medium text-muted-foreground">Dot</span>
                  <Badge variant={settings.dot ? "default" : "secondary"} className="text-xs">
                    {settings.dot ? 'Enabled' : 'Disabled'}
                  </Badge>
                </div>
                
                <div className="flex items-center justify-between py-2 border-b border-slate-200 dark:border-slate-800">
                  <span className="text-sm font-medium text-muted-foreground">Pip Border</span>
                  <Badge variant={settings.pipBorder ? "default" : "secondary"} className="text-xs">
                    {settings.pipBorder ? 'Yes' : 'No'}
                  </Badge>
                </div>
                
                <div className="flex items-center justify-between py-2">
                  <span className="text-sm font-medium text-muted-foreground">Static Gap</span>
                  <Badge variant={settings.pipGapStatic ? "default" : "secondary"} className="text-xs">
                    {settings.pipGapStatic ? 'Yes' : 'No'}
                  </Badge>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Like Button */}
      <div className="mt-8 flex justify-center">
        <LikeButton
          crosshairId={crosshair.id}
          liked={false}
          likesCount={crosshair.likes_count}
        />
      </div>
    </div>
  )
}
