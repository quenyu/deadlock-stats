export interface User {
  id: string
  steam_id: string
  nickname: string
  avatar_url: string
  profile_url: string
  created_at?: Date
  updated_at?: Date
  
  account_id?: number
  countrycode?: string
  last_updated?: number
  realname?: string
  
  is_deadlock_player: boolean
  deadlock_status_known: boolean
}

export interface SearchResult {
  results: User[]
  total_count: number
  page: number
  page_size: number
  total_pages: number
}

export interface SearchResponse {
  results: User[]
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
  searchTime?: number
  searchType?: string
  query?: string
  filters?: any
}