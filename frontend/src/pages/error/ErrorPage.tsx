import { useRouteError, isRouteErrorResponse, Link } from 'react-router-dom'
import { AlertCircle, Home, RefreshCw } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

export function ErrorPage() {
  const error = useRouteError()

  let errorMessage = 'An unexpected error occurred'
  let errorStatus = 500
  let errorStatusText = 'Internal Error'

  if (isRouteErrorResponse(error)) {
    errorStatus = error.status
    errorStatusText = error.statusText
    errorMessage = error.data?.message || error.statusText
  } else if (error instanceof Error) {
    errorMessage = error.message
  }

  const handleReload = () => {
    window.location.reload()
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
            <AlertCircle className="h-6 w-6 text-destructive" />
          </div>
          <CardTitle className="text-2xl">
            {errorStatus === 404 ? 'Page Not Found' : 'Oops! Something went wrong'}
          </CardTitle>
          <CardDescription>
            {errorStatus !== 404 && (
              <span className="font-mono text-sm">
                Error {errorStatus}: {errorStatusText}
              </span>
            )}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-center text-muted-foreground">{errorMessage}</p>
          <div className="flex flex-col gap-2 sm:flex-row">
            <Button asChild variant="default" className="flex-1">
              <Link to="/">
                <Home className="mr-2 h-4 w-4" />
                Go Home
              </Link>
            </Button>
            {errorStatus !== 404 && (
              <Button onClick={handleReload} variant="outline" className="flex-1">
                <RefreshCw className="mr-2 h-4 w-4" />
                Reload Page
              </Button>
            )}
          </div>
          {import.meta.env.DEV && error instanceof Error && (
            <details className="mt-4 rounded-lg bg-muted p-4">
              <summary className="cursor-pointer text-sm font-medium">
                Error Details (Development Only)
              </summary>
              <pre className="mt-2 overflow-auto text-xs">
                {error.stack}
              </pre>
            </details>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

