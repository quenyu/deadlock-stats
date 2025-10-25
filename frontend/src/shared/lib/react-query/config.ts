/**
 * React Query configuration
 */

import { QueryClient, DefaultOptions } from '@tanstack/react-query'

/**
 * Default query options for production
 */
export const defaultQueryOptions: DefaultOptions = {
  queries: {
    // Stale time - data calculating fresh in this time
    staleTime: 5 * 60 * 1000, // 5 minutes
    
    // Cache time - how long to keep unused data
    gcTime: 10 * 60 * 1000, // 10 minutes (was cacheTime)
    
    // Retry configuration
    retry: 3,
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
    
    // Refetch configuration
    refetchOnWindowFocus: true,
    refetchOnMount: true,
    refetchOnReconnect: true,
    
    // Error handling
    throwOnError: false,
  },
  mutations: {
    // Retry mutations
    retry: 1,
    retryDelay: 1000,
    
    // Error handling
    throwOnError: false,
  },
}

/**
 * Development query options
 */
export const developmentQueryOptions: DefaultOptions = {
  queries: {
    staleTime: 0, // Always refetch in development
    gcTime: 5 * 60 * 1000, // 5 minutes
    retry: 1, // Less retries in development
    retryDelay: 500,
    refetchOnWindowFocus: false, // Annoying in development
    refetchOnMount: true,
    refetchOnReconnect: true,
    throwOnError: false,
  },
  mutations: {
    retry: 0, // No retries in development
    throwOnError: false,
  },
}

/**
 * Create query client with configuration
 */
export function createQueryClient(isDevelopment = false): QueryClient {
  const options = isDevelopment ? developmentQueryOptions : defaultQueryOptions

  return new QueryClient({
    defaultOptions: options,
  })
}

/**
 * Query keys factory
 */
export const queryKeys = {
  // User queries
  user: {
    all: ['users'] as const,
    me: () => [...queryKeys.user.all, 'me'] as const,
    detail: (id: string) => [...queryKeys.user.all, 'detail', id] as const,
  },
  
  // Player queries
  player: {
    all: ['players'] as const,
    search: (query: string) => [...queryKeys.player.all, 'search', query] as const,
    profile: (steamId: string) => [...queryKeys.player.all, 'profile', steamId] as const,
    stats: (steamId: string) => [...queryKeys.player.all, 'stats', steamId] as const,
    matches: (steamId: string, filters?: Record<string, unknown>) => 
      [...queryKeys.player.all, 'matches', steamId, filters] as const,
  },
  
  // Crosshair queries
  crosshair: {
    all: ['crosshairs'] as const,
    published: () => [...queryKeys.crosshair.all, 'published'] as const,
    detail: (id: string) => [...queryKeys.crosshair.all, 'detail', id] as const,
    my: () => [...queryKeys.crosshair.all, 'my'] as const,
  },
  
  // Match queries
  match: {
    all: ['matches'] as const,
    detail: (id: string) => [...queryKeys.match.all, 'detail', id] as const,
    recent: (filters?: Record<string, unknown>) => 
      [...queryKeys.match.all, 'recent', filters] as const,
  },
} as const

/**
 * Mutation keys factory
 */
export const mutationKeys = {
  user: {
    update: 'user-update',
  },
  crosshair: {
    create: 'crosshair-create',
    update: 'crosshair-update',
    delete: 'crosshair-delete',
    like: 'crosshair-like',
    unlike: 'crosshair-unlike',
  },
} as const

