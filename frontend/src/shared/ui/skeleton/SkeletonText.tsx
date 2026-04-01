/**
 * Text skeleton component
 */

import { Skeleton, SkeletonProps } from './Skeleton'
import { cn } from '@/shared/lib/utils'

export interface SkeletonTextProps extends Omit<SkeletonProps, 'height' | 'variant'> {
  /** Number of lines */
  lines?: number
  
  /** Line height */
  lineHeight?: number
  
  /** Gap between lines */
  gap?: number
  
  /** Last line width (percentage or 'auto') */
  lastLineWidth?: number | 'auto'
}

/**
 * Text skeleton with multiple lines
 */
export function SkeletonText({
  lines = 1,
  lineHeight = 16,
  gap = 8,
  lastLineWidth = 'auto',
  className,
  ...props
}: SkeletonTextProps) {
  return (
    <div className={cn('flex flex-col', className)} style={{ gap }}>
      {Array.from({ length: lines }).map((_, index) => {
        const isLastLine = index === lines - 1
        const width = isLastLine && lastLineWidth !== 'auto' 
          ? `${lastLineWidth}%` 
          : '100%'

        return (
          <Skeleton
            key={index}
            height={lineHeight}
            width={width}
            {...props}
          />
        )
      })}
    </div>
  )
}

