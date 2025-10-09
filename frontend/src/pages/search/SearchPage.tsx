import { useState, useEffect } from 'react'
import { Input } from '@/shared/ui/input'
import { api } from '@/shared/api/api'
import { type User } from '@/entities/user'
import { routes } from '@/shared/constants/routes'
import React from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Label } from '@/shared/ui/label'
import { PaginatedResults } from '@/shared/ui/PaginatedResults'
import { PageSizeSelector } from '@/shared/ui/PageSizeSelector'
import { createLogger } from '@/shared/lib/logger'

const log = createLogger('SearchPage')

export const SearchPage = () => {
  const [query, setQuery] = useState('')
  const [searchType, setSearchType] = useState('nickname')
  const [results, setResults] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [totalCount, setTotalCount] = useState(0)
  const [totalPages, setTotalPages] = useState(0)

  useEffect(() => {
    if (searchType === 'nickname' && query.length < 3) {
      setResults([])
      setError(null)
      setTotalCount(0)
      setTotalPages(0)
      return
    }
    if (query.length === 0) {
      setResults([])
      setError(null)
      setTotalCount(0)
      setTotalPages(0)
      return
    }

    const fetchResults = async () => {
      setLoading(true)
      setError(null)
      try {
        const response = await api.get(`/players/search?q=${query}&type=${searchType}&page=${page}&pageSize=${pageSize}`)

        log.debug('Search results received', { 
          resultCount: response.data.results?.length || response.data.length,
          totalCount: response.data.total_count 
        })

        setResults(response.data.results || response.data)
        setTotalCount(response.data.total_count || response.data.length)
        setTotalPages(response.data.total_pages || 1)
        
        if ((response.data.results || response.data).length === 0) {
          setError('No players found.')
        }
      } catch (err) {
        setError('Failed to search for players.')
        log.error('Search failed', { error: err })
      } finally {
        setLoading(false)
      }
    }

    const debounceTimeout = setTimeout(fetchResults, 300)
    return () => clearTimeout(debounceTimeout)
  }, [query, searchType, page, pageSize])

  useEffect(() => {
    setPage(1)
  }, [query, searchType])

  const handleUserClick = (user: User) => {
    // Переход на профиль пользователя
    window.location.href = routes.player.profile(user.steam_id)
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
        {error && <div className="text-center text-red-500 mb-4">{error}</div>}
        
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
        
        <PaginatedResults
          results={results}
          totalCount={totalCount}
          page={page}
          pageSize={pageSize}
          totalPages={totalPages}
          onPageChange={setPage}
          onUserClick={handleUserClick}
          showExtendedInfo={true}
          loading={loading}
        />
      </div>
    </div>
  )
}