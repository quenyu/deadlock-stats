/**
 * Avatar skeleton component
 */

import { Skeleton, SkeletonProps } from './Skeleton'
import { cn } from '@/shared/lib/utils'

export interface SkeletonAvatarProps extends Omit<SkeletonProps, 'variant'> {
  /** Avatar size */
  size?: 'sm' | 'md' | 'lg' | 'xl' | number
}

const sizeMap = {
  sm: 32,
  md: 40,
  lg: 56,
  xl: 80,
}

/**
 * Circle avatar skeleton
 */
export function SkeletonAvatar({
  size = 'md',
  className,
  ...props
}: SkeletonAvatarProps) {
  const avatarSize = typeof size === 'number' ? size : sizeMap[size]

  return (
    <Skeleton
      variant="circle"
      width={avatarSize}
      height={avatarSize}
      className={className}
      {...props}
    />
  )
}

