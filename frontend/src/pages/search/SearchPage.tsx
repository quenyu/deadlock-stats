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

export const SearchPage = () => {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (query.length < 3) {
      setResults([])
      setError(null)
      return
    }

    const fetchResults = async () => {
      setLoading(true)
      setError(null)
      try {
        const response = await api.get<User[]>(`/players/search?q=${query}`)
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
  }, [query])

  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8">
      <div className="max-w-xl mx-auto">
        <h1 className="text-3xl font-bold text-center mb-6">Search for Players</h1>
        <Input
          type="search"
          placeholder="Enter player nickname..."
          value={query}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setQuery(e.target.value)}
          className="h-12 text-lg"
        />
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