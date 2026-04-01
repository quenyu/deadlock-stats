/**
 * Card skeleton component
 */

import { Skeleton } from './Skeleton'
import { SkeletonText } from './SkeletonText'
import { cn } from '@/shared/lib/utils'

export interface SkeletonCardProps {
  /** Show image/avatar */
  showImage?: boolean
  
  /** Image position */
  imagePosition?: 'top' | 'left'
  
  /** Image size */
  imageSize?: number
  
  /** Number of text lines */
  lines?: number
  
  /** Show actions area */
  showActions?: boolean
  
  /** Custom className */
  className?: string
}

/**
 * Card skeleton with image and text
 */
export function SkeletonCard({
  showImage = true,
  imagePosition = 'top',
  imageSize = 200,
  lines = 3,
  showActions = false,
  className,
}: SkeletonCardProps) {
  const isHorizontal = imagePosition === 'left'

  return (
    <div
      className={cn(
        'rounded-lg border bg-card p-4',
        isHorizontal ? 'flex gap-4' : 'space-y-3',
        className
      )}
    >
      {showImage && (
        <Skeleton
          variant={imagePosition === 'top' ? 'rounded' : 'circle'}
          width={isHorizontal ? imageSize / 2 : '100%'}
          height={isHorizontal ? imageSize / 2 : imageSize}
          className={isHorizontal ? 'flex-shrink-0' : ''}
        />
      )}
      
      <div className={cn('flex-1', isHorizontal ? '' : 'space-y-3')}>
        <SkeletonText lines={lines} lastLineWidth={70} />
        
        {showActions && (
          <div className="flex gap-2 mt-4">
            <Skeleton width={80} height={32} variant="rounded" />
            <Skeleton width={80} height={32} variant="rounded" />
          </div>
        )}
      </div>
    </div>
  )
}

