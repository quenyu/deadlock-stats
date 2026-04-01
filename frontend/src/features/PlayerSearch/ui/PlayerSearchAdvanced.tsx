import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Filter, Users, Clock } from 'lucide-react'
import { Input } from '@/shared/ui/input'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { routes } from '@/shared/constants/routes'
import { usePlayerSearch, type SearchFilters, type SearchType } from '@/entities/user'
import { cn } from '@/shared/lib/utils'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { PaginatedResults } from '@/shared/ui/PaginatedResults'
import { PageSizeSelector } from '@/shared/ui/PageSizeSelector'

interface PlayerSearchAdvancedProps {
  className?: string
  showFilters?: boolean
  showPopular?: boolean
  showRecentlyActive?: boolean
}

export const PlayerSearchAdvanced = ({
  className,
  showPopular = true,
  showRecentlyActive = true
}: PlayerSearchAdvancedProps) => {
  const [query, setQuery] = useState('')
  const [searchType, setSearchType] = useState<SearchType>('all')
  const [filters, setFilters] = useState<SearchFilters>({
    search_type: 'all',
    sort_by: 'nickname',
    sort_order: 'asc'
  })
  const [showFilterPanel, setShowFilterPanel] = useState(false)
  const navigate = useNavigate()
  
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [totalCount, setTotalCount] = useState(0)

  const {
    loading: searchLoading,
    error,
    searchWithFilters,
    getPopularPlayers,
    getRecentlyActivePlayers
  } = usePlayerSearch()

  const [searchResults, setSearchResults] = useState<any[]>([])
  const [popularPlayers, setPopularPlayers] = useState<any[]>([])
  const [recentlyActivePlayers, setRecentlyActivePlayers] = useState<any[]>([])
  const [searchStats, setSearchStats] = useState<{ totalFound: number; searchTime: number }>({ totalFound: 0, searchTime: 0 })

  useEffect(() => {
    setFilters(prev => ({ ...prev, search_type: searchType }))
  }, [searchType])

  useEffect(() => {
    if (query.length < 2) {
      setSearchResults([])
      setSearchStats({ totalFound: 0, searchTime: 0 })
      setTotalCount(0)
      return
    }

    const fetchResults = async () => {
      const response = await searchWithFilters(query, filters, page, pageSize)
      setSearchResults(response.results)
      setTotalCount(response.total_count)
      setSearchStats({
        totalFound: response.total_count,
        searchTime: response.searchTime || 0
      })
    }

    const debounceTimeout = setTimeout(fetchResults, 300)
    return () => clearTimeout(debounceTimeout)
  }, [query, filters, page, pageSize, searchWithFilters])

  useEffect(() => {
    setPage(1)
  }, [query, filters])

  useEffect(() => {
    if (showPopular && query.length === 0) {
      const fetchPopular = async () => {
        const response = await getPopularPlayers(1, 10)
        setPopularPlayers(response.results)
      }
      fetchPopular()
    }
  }, [showPopular, query, getPopularPlayers])

  useEffect(() => {
    if (showRecentlyActive && query.length === 0) {
      const fetchRecentlyActive = async () => {
        const response = await getRecentlyActivePlayers(1, 10)
        setRecentlyActivePlayers(response.results)
      }
      fetchRecentlyActive()
    }
  }, [showRecentlyActive, query, getRecentlyActivePlayers])

  const handleUserClick = (user: any) => {
    navigate(routes.player.profile(user.steam_id))
  }

  const filteredResults = searchResults
  const filteredTotalCount = totalCount

  const renderFilterPanel = () => (
    <Card className="mb-4">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Filter className="h-5 w-5" />
          Search Filters
        </CardTitle>
      </CardHeader>
      <CardContent className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div>
          <label className="text-sm font-medium mb-2 block">Search Type</label>
          <Select value={searchType} onValueChange={v => setSearchType(v as SearchType)}>
            <SelectTrigger className="h-9 w-full text-base">
              <SelectValue placeholder="Search type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All</SelectItem>
              <SelectItem value="nickname">Nickname</SelectItem>
              <SelectItem value="steamid">SteamID</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div>
          <label className="text-sm font-medium mb-2 block">Sort By</label>
          <Select value={filters.sort_by} onValueChange={v => setFilters(f => ({ ...f, sort_by: v as any }))}>
            <SelectTrigger className="h-9 w-full text-base">
              <SelectValue placeholder="Sort by" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="nickname">Nickname</SelectItem>
              <SelectItem value="created_at">Created</SelectItem>
              <SelectItem value="updated_at">Updated</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div>
          <label className="text-sm font-medium mb-2 block">Sort Order</label>
          <Select value={filters.sort_order} onValueChange={v => setFilters(f => ({ ...f, sort_order: v as 'asc' | 'desc' }))}>
            <SelectTrigger className="h-9 w-full text-base">
              <SelectValue placeholder="Order" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="asc">Ascending</SelectItem>
              <SelectItem value="desc">Descending</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div>
          <label className="text-sm font-medium mb-2 block">Results Limit</label>
          <Select value={pageSize.toString()} onValueChange={v => setPageSize(Number(v))}>
            <SelectTrigger className="h-9 w-full text-base">
              <SelectValue placeholder="Limit" />
            </SelectTrigger>
            <SelectContent>
              {[10, 20, 50, 100].map(lim => (
                <SelectItem key={lim} value={lim.toString()}>{lim}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </CardContent>
    </Card>
  )

  return (
    <div className={cn("space-y-4", className)}>
      <div className="flex flex-col items-center gap-4">
        <div className="w-full max-w-2xl">
          <Input
            type="search"
            placeholder="Search players by nickname or Steam ID..."
            className="w-full text-lg"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            autoFocus
          />
        </div>
        <div className="flex gap-2">
          <Button
            variant={showFilterPanel ? "default" : "outline"}
            onClick={() => setShowFilterPanel(!showFilterPanel)}
            className="flex items-center gap-2"
          >
            <Filter className="h-4 w-4" />
            Filters
          </Button>
        </div>
      </div>

      {showFilterPanel && renderFilterPanel()}

      {error && (
        <div className="text-red-500 text-center p-4">
          {error}
        </div>
      )}

      {query.length > 0 ? (
        <div className="space-y-4">
          {searchStats.totalFound > 0 && (
            <div className="flex justify-between items-center">
              <div className="text-sm text-muted-foreground">
                Found {filteredTotalCount} players in {searchStats.searchTime}ms
              </div>
              <PageSizeSelector 
                pageSize={pageSize} 
                onPageSizeChange={(newPageSize) => {
                  setPageSize(newPageSize)
                  setPage(1)
                }} 
              />
            </div>
          )}
          <PaginatedResults
            results={filteredResults}
            totalCount={filteredTotalCount}
            page={1}
            pageSize={pageSize}
            totalPages={Math.ceil(filteredTotalCount / pageSize)}
            onPageChange={setPage}
            onUserClick={handleUserClick}
            showExtendedInfo={true}
            loading={searchLoading}
          />
        </div>
      ) : (
        <div className="space-y-6">
          {showPopular && popularPlayers.length > 0 && (
            <div>
              <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                <Users className="h-5 w-5" />
                Popular Players
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {popularPlayers.map((player, index) => (
                  <Card 
                    key={player.id || index} 
                    className="hover:shadow-md transition-shadow cursor-pointer"
                    onClick={() => handleUserClick(player)}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-center gap-4">
                        <Avatar className="h-12 w-12">
                          <AvatarImage src={player.avatar_url} alt={player.nickname} />
                          <AvatarFallback>{player.nickname.charAt(0).toUpperCase()}</AvatarFallback>
                        </Avatar>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium truncate">{player.nickname}</div>
                          <div className="text-sm text-muted-foreground truncate">
                            Steam ID: {player.steam_id}
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          )}

          {showRecentlyActive && recentlyActivePlayers.length > 0 && (
            <div>
              <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                <Clock className="h-5 w-5" />
                Recently Active
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {recentlyActivePlayers.map((player, index) => (
                  <Card 
                    key={player.id || index} 
                    className="hover:shadow-md transition-shadow cursor-pointer"
                    onClick={() => handleUserClick(player)}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-center gap-4">
                        <Avatar className="h-12 w-12">
                          <AvatarImage src={player.avatar_url} alt={player.nickname} />
                          <AvatarFallback>{player.nickname.charAt(0).toUpperCase()}</AvatarFallback>
                        </Avatar>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium truncate">{player.nickname}</div>
                          <div className="text-sm text-muted-foreground truncate">
                            Steam ID: {player.steam_id}
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
} 