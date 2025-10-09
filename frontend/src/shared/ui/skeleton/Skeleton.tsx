/**
 * Base Skeleton component
 */

import { cn } from '@/shared/lib/utils'

export interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {
  /** Width of skeleton */
  width?: string | number
  
  /** Height of skeleton */
  height?: string | number
  
  /** Border radius variant */
  variant?: 'default' | 'rounded' | 'circle'
  
  /** Animation style */
  animation?: 'pulse' | 'wave' | 'none'
}

/**
 * Base Skeleton component with customizable props
 */
export function Skeleton({
  className,
  width,
  height,
  variant = 'default',
  animation = 'pulse',
  style,
  ...props
}: SkeletonProps) {
  const variantStyles = {
    default: 'rounded',
    rounded: 'rounded-lg',
    circle: 'rounded-full',
  }

  const animationStyles = {
    pulse: 'animate-pulse',
    wave: 'animate-shimmer',
    none: '',
  }

  return (
    <div
      className={cn(
        'bg-muted',
        variantStyles[variant],
        animationStyles[animation],
        className
      )}
      style={{
        width: typeof width === 'number' ? `${width}px` : width,
        height: typeof height === 'number' ? `${height}px` : height,
        ...style,
      }}
      {...props}
    />
  )
}

