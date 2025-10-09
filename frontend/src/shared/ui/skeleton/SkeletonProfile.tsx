/**
 * Profile/Player card skeleton component
 */

import { Skeleton } from './Skeleton'
import { SkeletonAvatar } from './SkeletonAvatar'
import { SkeletonText } from './SkeletonText'
import { cn } from '@/shared/lib/utils'

export interface SkeletonProfileProps {
  /** Layout variant */
  variant?: 'vertical' | 'horizontal'
  
  /** Avatar size */
  avatarSize?: 'sm' | 'md' | 'lg' | 'xl'
  
  /** Show stats */
  showStats?: boolean
  
  /** Number of stat items */
  statCount?: number
  
  /** Custom className */
  className?: string
}

/**
 * Profile skeleton for player/user cards
 */
export function SkeletonProfile({
  variant = 'horizontal',
  avatarSize = 'lg',
  showStats = true,
  statCount = 4,
  className,
}: SkeletonProfileProps) {
  const isVertical = variant === 'vertical'

  return (
    <div
      className={cn(
        'rounded-lg border bg-card p-6',
        isVertical ? 'flex flex-col items-center text-center space-y-4' : 'space-y-4',
        className
      )}
    >
      {/* Avatar and Name */}
      <div
        className={cn(
          'flex gap-4',
          isVertical ? 'flex-col items-center' : 'items-center'
        )}
      >
        <SkeletonAvatar size={avatarSize} />
        
        <div className={cn('flex-1', isVertical ? 'w-full' : '')}>
          <Skeleton height={24} width={isVertical ? '60%' : 200} className="mb-2" />
          <Skeleton height={16} width={isVertical ? '40%' : 150} />
        </div>
      </div>

      {/* Stats */}
      {showStats && (
        <div className={cn('grid gap-4', `grid-cols-${Math.min(statCount, 4)}`)}>
          {Array.from({ length: statCount }).map((_, index) => (
            <div key={index} className="space-y-1">
              <Skeleton height={12} width="60%" />
              <Skeleton height={20} width="80%" />
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

