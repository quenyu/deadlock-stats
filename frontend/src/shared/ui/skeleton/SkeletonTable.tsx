/**
 * Table skeleton component
 */

import { Skeleton } from './Skeleton'
import { cn } from '@/shared/lib/utils'

export interface SkeletonTableProps {
  /** Number of rows */
  rows?: number
  
  /** Number of columns */
  columns?: number
  
  /** Show header */
  showHeader?: boolean
  
  /** Column widths (percentage) */
  columnWidths?: number[]
  
  /** Custom className */
  className?: string
}

/**
 * Table skeleton with rows and columns
 */
export function SkeletonTable({
  rows = 5,
  columns = 4,
  showHeader = true,
  columnWidths,
  className,
}: SkeletonTableProps) {
  const defaultColumnWidths = Array(columns).fill(100 / columns)
  const widths = columnWidths || defaultColumnWidths

  return (
    <div className={cn('w-full', className)}>
      {showHeader && (
        <div className="flex gap-4 pb-4 border-b mb-4">
          {Array.from({ length: columns }).map((_, index) => (
            <Skeleton
              key={`header-${index}`}
              width={`${widths[index]}%`}
              height={20}
              className="flex-1"
            />
          ))}
        </div>
      )}
      
      <div className="space-y-3">
        {Array.from({ length: rows }).map((_, rowIndex) => (
          <div key={`row-${rowIndex}`} className="flex gap-4">
            {Array.from({ length: columns }).map((_, colIndex) => (
              <Skeleton
                key={`cell-${rowIndex}-${colIndex}`}
                width={`${widths[colIndex]}%`}
                height={16}
                className="flex-1"
              />
            ))}
          </div>
        ))}
      </div>
    </div>
  )
}

