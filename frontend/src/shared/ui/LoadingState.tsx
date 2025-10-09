import { Loader2 } from 'lucide-react'
import { cn } from '@/shared/lib/utils'

interface LoadingStateProps {
  message?: string
  className?: string
  size?: 'sm' | 'md' | 'lg'
}

const sizeClasses = {
  sm: 'h-4 w-4',
  md: 'h-8 w-8',
  lg: 'h-12 w-12',
}

export function LoadingState({ 
  message = 'Loading...', 
  className = '',
  size = 'md'
}: LoadingStateProps) {
  return (
    <div className={cn('flex flex-col items-center justify-center gap-3 p-8', className)}>
      <Loader2 className={cn('animate-spin text-muted-foreground', sizeClasses[size])} />
      {message && (
        <p className="text-sm text-muted-foreground">{message}</p>
      )}
    </div>
  )
}

