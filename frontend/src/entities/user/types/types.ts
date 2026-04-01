// Matches backend: backend/internal/domain/user.go User
// Extended with fields from backend/internal/dto/user_search_result.go UserSearchResult
// Import and re-export from validation schema for type safety
import type { User as UserFromValidation } from '@/shared/lib/validation'

export type User = UserFromValidation

// Matches backend: backend/internal/dto/search_result.go SearchResult
export interface SearchResult {
  results: User[]
  total_count: number
  page: number
  page_size: number
  total_pages: number
}