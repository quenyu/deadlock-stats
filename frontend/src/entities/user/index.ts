import useUserStore from './model/store'
import { fetchCurrentUser } from './api/fetchCurrentUser'
import { searchApi } from './api/searchApi'
import { usePlayerSearch } from './hooks/usePlayerSearch'
import type { User } from './types/types'
import type {
  SearchFilters,
  SearchResponse,
  AutocompleteResponse,
  FilteredSearchResponse,
  PopularPlayersResponse,
  RecentlyActiveResponse,
  SearchDebugResponse,
  SearchType
} from './types/search'

export { useUserStore, fetchCurrentUser, searchApi, usePlayerSearch }
export type { 
  User,
  SearchFilters,
  SearchResponse,
  AutocompleteResponse,
  FilteredSearchResponse,
  PopularPlayersResponse,
  RecentlyActiveResponse,
  SearchDebugResponse,
  SearchType
}