import { User } from ".."

// Matches backend: backend/internal/dto/search_filters.go SearchFilters
export interface SearchFilters {
  search_type?: 'all' | 'steamid' | 'nickname'
  min_rank?: number
  max_rank?: number
  min_matches?: number
  max_matches?: number
  min_win_rate?: number
  max_win_rate?: number
  min_kd_ratio?: number
  max_kd_ratio?: number
  sort_by?: 'rank' | 'matches' | 'win_rate' | 'kd_ratio' | 'nickname' | 'created_at' | 'updated_at'
  sort_order?: 'asc' | 'desc'
}

// Matches backend: backend/internal/dto/search_result.go SearchResult
export interface SearchResponse {
  query?: string
  searchType?: string
  searchTime?: number
  total_count: number
  page: number
  page_size: number
  total_pages: number
  results: User[]
}

export interface AutocompleteResponse {
  query: string
  limit: number
  searchTime: number
  totalFound: number
  results: User[]
}

export interface FilteredSearchResponse {
  query: string
  filters: SearchFilters
  searchTime?: number
  total_count: number
  page: number
  page_size: number
  total_pages: number
  results: User[]
}

export interface PopularPlayersResponse {
  searchTime?: number
  total_count: number
  page: number
  page_size: number
  total_pages: number
  results: User[]
}

export interface RecentlyActiveResponse {
  searchTime?: number
  total_count: number
  page: number
  page_size: number
  total_pages: number
  results: User[]
}

export interface SearchDebugInfo {
  steamID: string
  nickname: string
  avatarURL: string
  profileURL: string
  createdAt: string
  updatedAt: string
  isValid: boolean
}

export interface SearchDebugResponse {
  query: string
  searchType: string
  searchTime: number
  total_count: number
  page: number
  page_size: number
  total_pages: number
  results: User[]
  debugInfo: SearchDebugInfo[]
}

export type SearchType = 'nickname' | 'steamid' | 'all' 