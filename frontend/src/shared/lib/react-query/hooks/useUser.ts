/**
 * User-related React Query hooks
 */

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { queryKeys, mutationKeys } from '../config'
import { api } from '@/shared/api/api'
import { userSchema } from '@/shared/lib/validation'
import type { User } from '@/shared/lib/validation'
import { validateApiResponse } from '@/shared/lib/validation'
import { createLogger } from '@/shared/lib/logger'

const log = createLogger('useUser')

/**
 * Fetch current user
 */
async function fetchCurrentUser(): Promise<User | null> {
  try {
    const response = await api.get('/users/me')
    
    if (response.status === 401) {
      return null
    }
    
    return validateApiResponse(userSchema, response.data, '/users/me')
  } catch (error) {
    log.error('Failed to fetch current user', { error })
    return null
  }
}

/**
 * Hook to get current user
 */
export function useCurrentUser() {
  return useQuery({
    queryKey: queryKeys.user.me(),
    queryFn: fetchCurrentUser,
    staleTime: 10 * 60 * 1000, // 10 minutes
    retry: 1,
  })
}

/**
 * Update user mutation
 */
export function useUpdateUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationKey: [mutationKeys.user.update],
    mutationFn: async (data: Partial<User>) => {
      const response = await api.put('/users/me', data)
      return validateApiResponse(userSchema, response.data, '/users/me')
    },
    onSuccess: (updatedUser) => {
      // Update cache
      queryClient.setQueryData(queryKeys.user.me(), updatedUser)
      log.info('User updated successfully')
    },
    onError: (error) => {
      log.error('Failed to update user', { error })
    },
  })
}

/**
 * Logout mutation
 */
export function useLogout() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async () => {
      await api.post('/auth/logout')
    },
    onSuccess: () => {
      // Clear all caches
      queryClient.clear()
      log.info('User logged out')
    },
    onError: (error) => {
      log.error('Logout failed', { error })
    },
  })
}

