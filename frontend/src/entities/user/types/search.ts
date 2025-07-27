import { User } from ".."

export interface SearchFilters {
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

export interface SearchResponse {
  query?: string
  searchType?: string
  searchTime?: number
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
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
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
  results: User[]
}

export interface PopularPlayersResponse {
  searchTime?: number
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
  results: User[]
}

export interface RecentlyActiveResponse {
  searchTime?: number
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
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
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
  results: User[]
  debugInfo: SearchDebugInfo[]
}

export type SearchType = 'nickname' | 'steamid' | 'all' 