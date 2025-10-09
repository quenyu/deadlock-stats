/**
 * List skeleton component
 */

import { Skeleton } from './Skeleton'
import { SkeletonText } from './SkeletonText'
import { cn } from '@/shared/lib/utils'

export interface SkeletonListProps {
  /** Number of items */
  count?: number
  
  /** Show avatar/icon */
  showAvatar?: boolean
  
  /** Avatar size */
  avatarSize?: number
  
  /** Number of text lines per item */
  lines?: number
  
  /** Show actions */
  showActions?: boolean
  
  /** Custom className */
  className?: string
}

/**
 * List skeleton with items
 */
export function SkeletonList({
  count = 5,
  showAvatar = true,
  avatarSize = 40,
  lines = 2,
  showActions = false,
  className,
}: SkeletonListProps) {
  return (
    <div className={cn('space-y-4', className)}>
      {Array.from({ length: count }).map((_, index) => (
        <div key={index} className="flex items-center gap-3">
          {showAvatar && (
            <Skeleton
              variant="circle"
              width={avatarSize}
              height={avatarSize}
              className="flex-shrink-0"
            />
          )}
          
          <div className="flex-1">
            <SkeletonText lines={lines} lastLineWidth={60} />
          </div>
          
          {showActions && (
            <Skeleton width={32} height={32} variant="circle" className="flex-shrink-0" />
          )}
        </div>
      ))}
    </div>
  )
}

