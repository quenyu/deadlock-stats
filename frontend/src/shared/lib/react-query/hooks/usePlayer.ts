/**
 * Player-related React Query hooks
 */

import { useQuery, useInfiniteQuery, useQueryClient } from '@tanstack/react-query'
import { queryKeys } from '../config'
import { api } from '@/shared/api/api'
import {
  playerSearchResponseSchema,
  playerProfileSchema,
  recentMatchesResponseSchema,
} from '@/shared/lib/validation'
import type {
  PlayerSearchResponse,
  PlayerProfile,
  RecentMatchesResponse,
} from '@/shared/lib/validation'
import { validateApiResponse } from '@/shared/lib/validation'

/**
 * Search players
 */
export function usePlayerSearch(query: string, enabled = true) {
  return useQuery({
    queryKey: queryKeys.player.search(query),
    queryFn: async (): Promise<PlayerSearchResponse> => {
      const response = await api.get('/players/search', {
        params: { q: query },
      })
      return validateApiResponse(
        playerSearchResponseSchema,
        response.data,
        '/players/search'
      )
    },
    enabled: enabled && query.length > 0,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

/**
 * Get player profile
 */
export function usePlayerProfile(steamId: string, enabled = true) {
  return useQuery({
    queryKey: queryKeys.player.profile(steamId),
    queryFn: async (): Promise<PlayerProfile> => {
      const response = await api.get(`/players/${steamId}`)
      return validateApiResponse(
        playerProfileSchema,
        response.data,
        `/players/${steamId}`
      )
    },
    enabled: enabled && !!steamId,
    staleTime: 2 * 60 * 1000, // 2 minutes
    retry: 2,
  })
}

/**
 * Get player recent matches with infinite scroll
 */
export function usePlayerMatches(steamId: string, enabled = true) {
  return useInfiniteQuery({
    queryKey: queryKeys.player.matches(steamId),
    queryFn: async ({ pageParam = 1 }): Promise<RecentMatchesResponse> => {
      const response = await api.get(`/players/${steamId}/matches`, {
        params: {
          page: pageParam,
          limit: 10,
        },
      })
      return validateApiResponse(
        recentMatchesResponseSchema,
        response.data,
        `/players/${steamId}/matches`
      )
    },
    enabled: enabled && !!steamId,
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      const totalFetched = allPages.reduce((sum, page) => sum + page.matches.length, 0)
      return totalFetched < lastPage.total ? allPages.length + 1 : undefined
    },
    staleTime: 60 * 1000, // 1 minute
  })
}

/**
 * Prefetch player profile (for hover cards, etc)
 */
export function usePrefetchPlayerProfile() {
  const queryClient = useQueryClient()
  
  return (steamId: string) => {
    queryClient.prefetchQuery({
      queryKey: queryKeys.player.profile(steamId),
      queryFn: async () => {
        const response = await api.get(`/players/${steamId}`)
        return validateApiResponse(
          playerProfileSchema,
          response.data,
          `/players/${steamId}`
        )
      },
      staleTime: 2 * 60 * 1000,
    })
  }
}

