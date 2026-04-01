/**
 * React Query - Modular data fetching and caching
 * 
 * @example
 * ```tsx
 * import { QueryProvider } from '@/shared/lib/react-query'
 * import { useCurrentUser, usePlayerSearch } from '@/shared/lib/react-query/hooks'
 * 
 * // Wrap app with provider
 * <QueryProvider>
 *   <App />
 * </QueryProvider>
 * 
 * // Use hooks
 * const { data: user, isLoading } = useCurrentUser()
 * const { data: players } = usePlayerSearch('query')
 * ```
 */

// Export provider
export { QueryProvider } from './provider'

// Export configuration
export { createQueryClient, queryKeys, mutationKeys } from './config'

// Export all hooks
export * from './hooks'

// Export utilities
export * from './utils'

// Re-export commonly used React Query exports
export {
  useQuery,
  useMutation,
  useQueryClient,
  useInfiniteQuery,
  useSuspenseQuery,
  type UseQueryOptions,
  type UseMutationOptions,
  type UseInfiniteQueryOptions,
} from '@tanstack/react-query'

