import { useCrosshairStore } from '@/entities/crosshair/model/store'
import { Button } from '@/shared/ui/button'
import { Badge } from '@/shared/ui/badge'
import { HeartIcon } from 'lucide-react'

interface LikeButtonProps {
  crosshairId: string
  liked: boolean
  likesCount: number
}

export const LikeButton: React.FC<LikeButtonProps> = ({ crosshairId, liked, likesCount }) => {
  const like = useCrosshairStore(s => s.like)
  const unlike = useCrosshairStore(s => s.unlike)

  const handleClick = async () => {
    try {
      if (liked) {
        await unlike(crosshairId)
      } else {
        await like(crosshairId)
      }
    } catch (error) {
      console.error('Failed to toggle like:', error)
    }
  }

  return (
    <Button
      variant="outline"
      size="sm"
      onClick={handleClick}
      className={`gap-2 transition-colors ${
        liked 
          ? 'bg-red-50 border-red-200 text-red-600 hover:bg-red-100' 
          : 'hover:bg-gray-50'
      }`}
    >
      <HeartIcon 
        size={16} 
        className={liked ? 'fill-red-500 text-red-500' : 'text-gray-400'} 
      />
      <span>{likesCount}</span>
    </Button>
  )
} 