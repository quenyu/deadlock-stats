import React from 'react'
import ReactDOM from 'react-dom/client'
import { App } from './app/App'
import { ErrorBoundary } from '@/shared/lib/ErrorBoundary'
import { ErrorFallback } from '@/shared/ui/ErrorFallback'
import { QueryProvider } from '@/shared/lib/react-query'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryProvider>
      <ErrorBoundary fallback={(error, reset) => <ErrorFallback error={error} onRetry={reset} />}>
        <App />
      </ErrorBoundary>
    </QueryProvider>
  </React.StrictMode>,
)
