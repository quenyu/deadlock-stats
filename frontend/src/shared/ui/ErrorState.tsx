import { AlertCircle } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Alert, AlertDescription, AlertTitle } from '@/shared/ui/alert'

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
    <Alert variant="destructive" className={className}>
      <AlertCircle className="h-4 w-4" />
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription className="flex flex-col gap-2">
        <p>{errorMessage}</p>
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
      </AlertDescription>
    </Alert>
  )
}

