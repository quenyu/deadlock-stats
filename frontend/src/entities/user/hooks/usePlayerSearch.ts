import { useState, useCallback } from 'react'
import { searchApi } from '../api/searchApi'
import type { SearchType, SearchFilters } from '../types/search'

export const usePlayerSearch = () => {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const searchPlayers = useCallback(async (
    query: string, 
    searchType: SearchType = 'all',
    page: number = 1,
    pageSize: number = 10
  ) => {
    if (query.length < 2) {
      return { 
        results: [], 
        total_count: 0, 
        page: 1, 
        page_size: 10, 
        total_pages: 0, 
        searchTime: 0 
      }
    }

    setLoading(true)
    setError(null)

    try {
      const response = await searchApi.searchPlayers(query, searchType, page, pageSize)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to search players'
      setError(errorMessage)
      return { 
        results: [], 
        total_count: 0, 
        page: 1, 
        page_size: 10, 
        total_pages: 0, 
        searchTime: 0 
      }
    } finally {
      setLoading(false)
    }
  }, [])

  const searchAutocomplete = useCallback(async (query: string, limit: number = 10) => {
    if (query.length < 2) {
      return { results: [], totalFound: 0, searchTime: 0 }
    }

    setLoading(true)
    setError(null)

    try {
      const response = await searchApi.searchAutocomplete(query, limit)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to search autocomplete'
      setError(errorMessage)
      return { results: [], totalFound: 0, searchTime: 0 }
    } finally {
      setLoading(false)
    }
  }, [])

  const searchWithFilters = useCallback(async (
    query: string,
    filters: SearchFilters,
    page: number = 1,
    pageSize: number = 20
  ) => {
    if (query.length < 2) {
      return { 
        results: [], 
        total_count: 0, 
        page: 1, 
        page_size: 20, 
        total_pages: 0, 
        searchTime: 0, 
        filters,
        query
      }
    }

    setLoading(true)
    setError(null)

    try {
      const response = await searchApi.searchWithFilters(query, filters, page, pageSize)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to search with filters'
      setError(errorMessage)
      return { 
        results: [], 
        total_count: 0, 
        page: 1, 
        page_size: 20, 
        total_pages: 0, 
        searchTime: 0, 
        filters,
        query
      }
    } finally {
      setLoading(false)
    }
  }, [])

  const getPopularPlayers = useCallback(async (page: number = 1, pageSize: number = 10) => {
    setLoading(true)
    setError(null)

    try {
      const response = await searchApi.getPopularPlayers(page, pageSize)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get popular players'
      setError(errorMessage)
      return { 
        results: [], 
        total_count: 0, 
        page: 1, 
        page_size: 10, 
        total_pages: 0, 
        searchTime: 0 
      }
    } finally {
      setLoading(false)
    }
  }, [])

  const getRecentlyActivePlayers = useCallback(async (page: number = 1, pageSize: number = 10) => {
    setLoading(true)
    setError(null)

    try {
      const response = await searchApi.getRecentlyActivePlayers(page, pageSize)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get recently active players'
      setError(errorMessage)
      return { 
        results: [], 
        total_count: 0, 
        page: 1, 
        page_size: 10, 
        total_pages: 0, 
        searchTime: 0 
      }
    } finally {
      setLoading(false)
    }
  }, [])

  const searchDebug = useCallback(async (
    query: string, 
    searchType: SearchType = 'all',
    page: number = 1,
    pageSize: number = 10
  ) => {
    if (query.length < 2) {
      return { 
        results: [], 
        total_count: 0, 
        page: 1, 
        page_size: 10, 
        total_pages: 0, 
        searchTime: 0,
        searchType,
        query,
        debugInfo: [] 
      }
    }

    setLoading(true)
    setError(null)

    try {
      const response = await searchApi.searchDebug(query, searchType, page, pageSize)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to search debug'
      setError(errorMessage)
      return { 
        results: [], 
        total_count: 0, 
        page: 1, 
        page_size: 10, 
        total_pages: 0, 
        searchTime: 0,
        searchType,
        query,
        debugInfo: [] 
      }
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    loading,
    error,
    searchPlayers,
    searchAutocomplete,
    searchWithFilters,
    getPopularPlayers,
    getRecentlyActivePlayers,
    searchDebug
  }
} 