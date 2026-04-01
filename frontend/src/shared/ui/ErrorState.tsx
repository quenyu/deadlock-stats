import { AlertCircle } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { cn } from '@/shared/lib/utils'

interface ErrorStateProps {
  error: string | Error | null
  onRetry?: () => void
  title?: string
  className?: string
}

export function ErrorState({ 
  error, 
  onRetry, 
  title = 'Error',
  className = ''
}: ErrorStateProps) {
  if (!error) return null

  const errorMessage = error instanceof Error ? error.message : error

  return (
    <div className={cn("rounded-lg border border-destructive bg-destructive/10 p-4", className)}>
      <div className="flex items-start gap-3">
        <AlertCircle className="h-5 w-5 text-destructive mt-0.5" />
        <div className="flex-1">
          <h3 className="font-semibold text-destructive mb-1">{title}</h3>
          <p className="text-sm text-destructive/90 mb-3">{errorMessage}</p>
          {onRetry && (
            <Button 
              onClick={onRetry} 
              variant="outline" 
              size="sm"
              className="w-fit"
            >
              Try again
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}

