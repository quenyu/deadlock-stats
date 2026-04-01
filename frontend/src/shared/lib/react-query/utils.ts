/**
 * React Query utility functions
 */

import { QueryClient } from '@tanstack/react-query'
import { queryKeys } from './config'

/**
 * Invalidate all user-related queries
 */
export function invalidateUserQueries(queryClient: QueryClient) {
  return queryClient.invalidateQueries({
    queryKey: queryKeys.user.all,
  })
}

/**
 * Invalidate all player-related queries
 */
export function invalidatePlayerQueries(queryClient: QueryClient) {
  return queryClient.invalidateQueries({
    queryKey: queryKeys.player.all,
  })
}

/**
 * Invalidate all crosshair-related queries
 */
export function invalidateCrosshairQueries(queryClient: QueryClient) {
  return queryClient.invalidateQueries({
    queryKey: queryKeys.crosshair.all,
  })
}

/**
 * Clear all caches
 */
export function clearAllCaches(queryClient: QueryClient) {
  return queryClient.clear()
}

/**
 * Prefetch queries for a list of items
 */
export async function prefetchList<T>(
  queryClient: QueryClient,
  items: T[],
  getQueryKey: (item: T) => readonly unknown[],
  fetchFn: (item: T) => Promise<unknown>
) {
  await Promise.all(
    items.map((item) =>
      queryClient.prefetchQuery({
        queryKey: getQueryKey(item),
        queryFn: () => fetchFn(item),
      })
    )
  )
}

/**
 * Optimistic update helper
 */
export function optimisticUpdate<T>(
  queryClient: QueryClient,
  queryKey: readonly unknown[],
  updateFn: (old: T | undefined) => T
) {
  queryClient.setQueryData<T>(queryKey, updateFn)
}

/**
 * Cancel and rollback helper
 */
export async function cancelAndRollback<T>(
  queryClient: QueryClient,
  queryKey: readonly unknown[]
): Promise<T | undefined> {
  await queryClient.cancelQueries({ queryKey })
  return queryClient.getQueryData<T>(queryKey)
}

