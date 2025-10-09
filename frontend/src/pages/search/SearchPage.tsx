import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Input } from '@/shared/ui/input'
import { routes } from '@/shared/constants/routes'
import React from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Label } from '@/shared/ui/label'
import { PaginatedResults } from '@/shared/ui/PaginatedResults'
import { PageSizeSelector } from '@/shared/ui/PageSizeSelector'
import { usePlayerSearch } from '@/shared/lib/react-query/hooks'
import { SkeletonList } from '@/shared/ui/skeleton'
import { createLogger } from '@/shared/lib/logger'
import { PlayerSearchResult } from '@/shared/lib/validation'
import { User } from '@/entities/user'

type UserCardData = PlayerSearchResult | User

const log = createLogger('SearchPage')

export const SearchPage = () => {
  const [query, setQuery] = useState('')
  const [searchType, setSearchType] = useState('nickname')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const navigate = useNavigate()

  const shouldSearch = searchType === 'nickname' ? query.length >= 3 : query.length > 0
  
  const { data, isLoading, error } = usePlayerSearch(query, shouldSearch)
  
  const results = data?.players || []
  const totalCount = data?.meta.total_count || 0
  const totalPages = data?.meta.total_pages || 0

  const handleUserClick = (user: UserCardData) => {
    // Переход на профиль пользователя
    const steamId = 'steamId' in user ? user.steamId : ('steam_id' in user ? user.steam_id : '')
    navigate(routes.player.profile(steamId))
  }

  log.debug('Current search results', { count: results.length })

  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8">
      <div className="max-w-xl mx-auto">
        <h1 className="text-3xl font-bold text-center mb-6">Search for Players</h1>
        <div className="flex items-end gap-2">
          <div className="flex-grow">
            <Label htmlFor="search-input" className="mb-2 block">
              Player {searchType === 'nickname' ? 'Nickname' : 'SteamID'}
            </Label>
            <Input
              id="search-input"
              type="search"
              placeholder={searchType === 'nickname' ? 'Enter nickname...' : 'Enter SteamID...'}
              value={query}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setQuery(e.target.value)}
              className="h-12 text-lg"
            />
          </div>
          <div className="w-[140px]">
            <Label className="mb-2 block">Search by</Label>
            <Select value={searchType} onValueChange={(value) => {setQuery(''); setSearchType(value)}}>
              <SelectTrigger className="h-12 w-full text-base">
                <SelectValue placeholder="Search type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="nickname">Nickname</SelectItem>
                <SelectItem value="steamid">SteamID</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>

      <div className="mt-8">
        {error && <div className="text-center text-red-500 mb-4">{error.message}</div>}
        
        {totalCount > 0 && (
          <div className="flex justify-end mb-4">
            <PageSizeSelector 
              pageSize={pageSize} 
              onPageSizeChange={(newPageSize) => {
                setPageSize(newPageSize)
                setPage(1)
              }} 
            />
          </div>
        )}
        
        {isLoading ? (
          <SkeletonList count={pageSize} showAvatar avatarSize={48} lines={2} />
        ) : error ? (
          <div className="text-center py-8 text-destructive">
            Failed to search for players.
          </div>
        ) : results.length === 0 && shouldSearch ? (
          <div className="text-center py-8 text-muted-foreground">
            No players found.
          </div>
        ) : (
          <PaginatedResults
            results={results}
            totalCount={totalCount}
            page={page}
            pageSize={pageSize}
            totalPages={totalPages}
            onPageChange={setPage}
            onUserClick={handleUserClick}
            showExtendedInfo={true}
            loading={false}
          />
        )}
      </div>
    </div>
  )
}