/**
 * Skeleton components - Modular loading skeletons
 * 
 * @example
 * ```tsx
 * import { Skeleton, SkeletonCard, SkeletonProfile } from '@/shared/ui/skeleton'
 * 
 * // Basic skeleton
 * <Skeleton width={200} height={20} />
 * 
 * // Card skeleton
 * <SkeletonCard showImage lines={3} showActions />
 * 
 * // Profile skeleton
 * <SkeletonProfile variant="horizontal" showStats />
 * 
 * // List skeleton
 * <SkeletonList count={5} showAvatar />
 * ```
 */

// Base components
export { Skeleton } from './Skeleton'
export type { SkeletonProps } from './Skeleton'

// Text skeleton
export { SkeletonText } from './SkeletonText'
export type { SkeletonTextProps } from './SkeletonText'

// Avatar skeleton
export { SkeletonAvatar } from './SkeletonAvatar'
export type { SkeletonAvatarProps } from './SkeletonAvatar'

// Composite components
export { SkeletonCard } from './SkeletonCard'
export type { SkeletonCardProps } from './SkeletonCard'

export { SkeletonList } from './SkeletonList'
export type { SkeletonListProps } from './SkeletonList'

export { SkeletonTable } from './SkeletonTable'
export type { SkeletonTableProps } from './SkeletonTable'

export { SkeletonProfile } from './SkeletonProfile'
export type { SkeletonProfileProps } from './SkeletonProfile'

