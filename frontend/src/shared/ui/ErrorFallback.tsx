import { Button } from '@/shared/ui/button'
import React from 'react'

interface ErrorFallbackProps {
  error: unknown
  onRetry?: () => void
}

export function ErrorFallback({ error, onRetry }: ErrorFallbackProps) {
  const message = (error instanceof Error && error.message) || 'Something went wrong.'
  return (
    <div className="flex min-h-[50vh] w-full flex-col items-center justify-center gap-4 p-6 text-center">
      <h2 className="text-2xl font-semibold">Unexpected error</h2>
      <p className="text-muted-foreground">{message}</p>
      {onRetry && (
        <Button onClick={onRetry} variant="default">
          Try again
        </Button>
      )}
    </div>
  )
}


