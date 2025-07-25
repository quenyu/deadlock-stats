import { useState, useEffect } from 'react'
import { Input } from '@/shared/ui/input'
import { api } from '@/shared/api/api'
import { type User } from '@/entities/user'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/shared/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { Button } from '@/shared/ui/button'
import { AppLink } from '@/shared/ui/AppLink/AppLink'
import { routes } from '@/shared/constants/routes'
import React from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Label } from '@/shared/ui/label'

export const SearchPage = () => {
  const [query, setQuery] = useState('')
  const [searchType, setSearchType] = useState('nickname')
  const [results, setResults] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (searchType === 'nickname' && query.length < 3) {
      setResults([])
      setError(null)
      return
    }
    if (query.length === 0) {
      setResults([])
      setError(null)
      return
    }

    const fetchResults = async () => {
      setLoading(true)
      setError(null)
      try {
        const response = await api.get<User[]>(`/players/search?q=${query}&type=${searchType}`)

        console.log("Data received from backend:", response.data);

        setResults(response.data)
        if (response.data.length === 0) {
          setError('No players found.')
        }
      } catch (err) {
        setError('Failed to search for players.')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    const debounceTimeout = setTimeout(fetchResults, 300)
    return () => clearTimeout(debounceTimeout)
  }, [query, searchType])

  console.log(results)

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
              <SelectTrigger className="h-12">
                <SelectValue />
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
        {loading && <div className="text-center">Searching...</div>}
        {error && <div className="text-center text-red-500">{error}</div>}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {results.map((user) => (
            <Card key={user.id}>
              <CardHeader className="items-center">
                <Avatar className="h-20 w-20">
                  <AvatarImage src={user.avatar_url} alt={user.nickname} />
                  <AvatarFallback>{user.nickname.charAt(0)}</AvatarFallback>
                </Avatar>
              </CardHeader>
              <CardContent className="text-center">
                <CardTitle>{user.nickname}</CardTitle>
              </CardContent>
              <CardFooter>
                <AppLink to={routes.player.profile(user.steam_id)} className="w-full">
                  <Button className="w-full">View Profile</Button>
                </AppLink>
              </CardFooter>
            </Card>
          ))}
        </div>
      </div>
    </div>
  )
}