/**
 * Crosshair-related React Query hooks
 */

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { queryKeys, mutationKeys } from '../config'
import { api } from '@/shared/api/api'
import {
  crosshairSchema,
  publishedCrosshairsResponseSchema,
  createCrosshairSchema,
} from '@/shared/lib/validation'
import type {
  Crosshair,
  PublishedCrosshairsResponse,
  CreateCrosshair,
} from '@/shared/lib/validation'
import { validateApiResponse, safeParse } from '@/shared/lib/validation'
import { createLogger } from '@/shared/lib/logger'

const log = createLogger('useCrosshair')

/**
 * Get published crosshairs
 */
export function usePublishedCrosshairs(enabled = true) {
  return useQuery({
    queryKey: queryKeys.crosshair.published(),
    queryFn: async (): Promise<PublishedCrosshairsResponse> => {
      const response = await api.get('/crosshairs/published')
      return validateApiResponse(
        publishedCrosshairsResponseSchema,
        response.data,
        '/crosshairs/published'
      )
    },
    enabled,
    staleTime: 60 * 1000, // 1 minute
  })
}

/**
 * Get crosshair by ID
 */
export function useCrosshair(id: string, enabled = true) {
  return useQuery({
    queryKey: queryKeys.crosshair.detail(id),
    queryFn: async (): Promise<Crosshair> => {
      const response = await api.get(`/crosshairs/${id}`)
      return validateApiResponse(
        crosshairSchema,
        response.data,
        `/crosshairs/${id}`
      )
    },
    enabled: enabled && !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

/**
 * Create crosshair mutation
 */
export function useCreateCrosshair() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationKey: [mutationKeys.crosshair.create],
    mutationFn: async (data: CreateCrosshair) => {
      // Validate before sending
      const validationResult = safeParse(createCrosshairSchema, data)
      if (!validationResult.success) {
        throw new Error('Invalid crosshair data')
      }
      
      const response = await api.post('/crosshairs', validationResult.data)
      return validateApiResponse(crosshairSchema, response.data, '/crosshairs')
    },
    onSuccess: (newCrosshair) => {
      // Invalidate published crosshairs list
      queryClient.invalidateQueries({
        queryKey: queryKeys.crosshair.published(),
      })
      
      // Add to cache
      queryClient.setQueryData(
        queryKeys.crosshair.detail(newCrosshair.id),
        newCrosshair
      )
      
      log.info('Crosshair created successfully', { id: newCrosshair.id })
    },
    onError: (error) => {
      log.error('Failed to create crosshair', { error })
    },
  })
}

/**
 * Update crosshair mutation
 */
export function useUpdateCrosshair(id: string) {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationKey: [mutationKeys.crosshair.update, id],
    mutationFn: async (data: Partial<CreateCrosshair>) => {
      const response = await api.put(`/crosshairs/${id}`, data)
      return validateApiResponse(crosshairSchema, response.data, `/crosshairs/${id}`)
    },
    onSuccess: (updatedCrosshair) => {
      // Update cache
      queryClient.setQueryData(
        queryKeys.crosshair.detail(id),
        updatedCrosshair
      )
      
      // Invalidate lists
      queryClient.invalidateQueries({
        queryKey: queryKeys.crosshair.published(),
      })
      
      log.info('Crosshair updated successfully', { id })
    },
    onError: (error) => {
      log.error('Failed to update crosshair', { error, id })
    },
  })
}

/**
 * Delete crosshair mutation
 */
export function useDeleteCrosshair() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationKey: [mutationKeys.crosshair.delete],
    mutationFn: async (id: string) => {
      await api.delete(`/crosshairs/${id}`)
      return id
    },
    onSuccess: (id) => {
      // Remove from cache
      queryClient.removeQueries({
        queryKey: queryKeys.crosshair.detail(id),
      })
      
      // Invalidate lists
      queryClient.invalidateQueries({
        queryKey: queryKeys.crosshair.published(),
      })
      
      log.info('Crosshair deleted successfully', { id })
    },
    onError: (error, id) => {
      log.error('Failed to delete crosshair', { error, id })
    },
  })
}

/**
 * Like crosshair mutation
 */
export function useLikeCrosshair() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationKey: [mutationKeys.crosshair.like],
    mutationFn: async (id: string) => {
      const response = await api.post(`/crosshairs/${id}/like`)
      return validateApiResponse(crosshairSchema, response.data, `/crosshairs/${id}/like`)
    },
    onMutate: async (id) => {
      // Optimistic update
      await queryClient.cancelQueries({ queryKey: queryKeys.crosshair.detail(id) })
      
      const previousCrosshair = queryClient.getQueryData<Crosshair>(
        queryKeys.crosshair.detail(id)
      )
      
      if (previousCrosshair) {
        queryClient.setQueryData(queryKeys.crosshair.detail(id), {
          ...previousCrosshair,
          likes: previousCrosshair.likes + 1,
        })
      }
      
      return { previousCrosshair }
    },
    onError: (error, id, context) => {
      // Rollback on error
      if (context?.previousCrosshair) {
        queryClient.setQueryData(
          queryKeys.crosshair.detail(id),
          context.previousCrosshair
        )
      }
      log.error('Failed to like crosshair', { error, id })
    },
    onSuccess: (updatedCrosshair) => {
      log.info('Crosshair liked', { id: updatedCrosshair.id })
    },
  })
}

/**
 * Unlike crosshair mutation
 */
export function useUnlikeCrosshair() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationKey: [mutationKeys.crosshair.unlike],
    mutationFn: async (id: string) => {
      const response = await api.delete(`/crosshairs/${id}/like`)
      return validateApiResponse(crosshairSchema, response.data, `/crosshairs/${id}/like`)
    },
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: queryKeys.crosshair.detail(id) })
      
      const previousCrosshair = queryClient.getQueryData<Crosshair>(
        queryKeys.crosshair.detail(id)
      )
      
      if (previousCrosshair) {
        queryClient.setQueryData(queryKeys.crosshair.detail(id), {
          ...previousCrosshair,
          likes: Math.max(0, previousCrosshair.likes - 1),
        })
      }
      
      return { previousCrosshair }
    },
    onError: (error, id, context) => {
      if (context?.previousCrosshair) {
        queryClient.setQueryData(
          queryKeys.crosshair.detail(id),
          context.previousCrosshair
        )
      }
      log.error('Failed to unlike crosshair', { error, id })
    },
    onSuccess: (updatedCrosshair) => {
      log.info('Crosshair unliked', { id: updatedCrosshair.id })
    },
  })
}

