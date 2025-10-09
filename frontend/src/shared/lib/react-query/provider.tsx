/**
 * React Query Provider component
 */

import { ReactNode } from 'react'
import { QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { createQueryClient } from './config'

const queryClient = createQueryClient(import.meta.env.DEV)

interface QueryProviderProps {
  children: ReactNode
}

/**
 * React Query Provider
 * Wrap your app with this component
 */
export function QueryProvider({ children }: QueryProviderProps) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {import.meta.env.DEV && (
        <ReactQueryDevtools
          initialIsOpen={false}
          position="bottom-right"
        />
      )}
    </QueryClientProvider>
  )
}

