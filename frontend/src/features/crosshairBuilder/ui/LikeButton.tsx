import { useCrosshairStore } from '@/entities/crosshair/model/store'
import Component from '@/shared/ui/components/comp-117'

export const LikeButton = () => {
  const like = useCrosshairStore(s => s.like)
  return (
    <Component onClick={like} />
  )
} 