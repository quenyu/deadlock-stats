import { api } from '@/shared/api/api'
import type {
  SearchResponse,
  AutocompleteResponse,
  FilteredSearchResponse,
  PopularPlayersResponse,
  RecentlyActiveResponse,
  SearchDebugResponse,
  SearchFilters,
  SearchType
} from '../types/search'

export const searchApi = {
  searchPlayers: async (
    query: string, 
    searchType: SearchType = 'all',
    page: number = 1,
    pageSize: number = 10
  ): Promise<SearchResponse> => {
    const response = await api.get<SearchResponse>(`/players/search`, {
      params: { query, searchType, page, pageSize }
    })
    return response.data
  },

  searchAutocomplete: async (query: string, limit: number = 10): Promise<AutocompleteResponse> => {
    const response = await api.get<AutocompleteResponse>(`/players/search/autocomplete`, {
      params: { query, limit }
    })
    return response.data
  },

  searchWithFilters: async (
    query: string,
    filters: SearchFilters,
    page: number = 1,
    pageSize: number = 20
  ): Promise<FilteredSearchResponse> => {
    const response = await api.get<FilteredSearchResponse>(`/players/search/filters`, {
      params: { query, page, pageSize, ...filters }
    })
    return response.data
  },

  getPopularPlayers: async (page: number = 1, pageSize: number = 10): Promise<PopularPlayersResponse> => {
    const response = await api.get<PopularPlayersResponse>(`/players/popular`, {
      params: { page, pageSize }
    })
    return response.data
  },

  getRecentlyActivePlayers: async (page: number = 1, pageSize: number = 10): Promise<RecentlyActiveResponse> => {
    const response = await api.get<RecentlyActiveResponse>(`/players/recently-active`, {
      params: { page, pageSize }
    })
    return response.data
  },

  searchDebug: async (
    query: string, 
    searchType: SearchType = 'all',
    page: number = 1,
    pageSize: number = 10
  ): Promise<SearchDebugResponse> => {
    const response = await api.get<SearchDebugResponse>(`/players/search/debug`, {
      params: { query, searchType, page, pageSize }
    })
    return response.data
  }
} 