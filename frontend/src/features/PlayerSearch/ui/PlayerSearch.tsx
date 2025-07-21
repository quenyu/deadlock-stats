import { useState, useEffect, useRef } from 'react'
import { Search } from 'lucide-react'
import { Input } from '@/shared/ui/input'
import { api } from '@/shared/api/api'
import { type User } from '@/entities/user'
import { AppLink } from '@/shared/ui/AppLink/AppLink'
import { routes } from '@/shared/constants/routes'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import React from 'react'

export const PlayerSearch = () => {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [isOpen, setIsOpen] = useState(false)
  const searchRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handler = (event: MouseEvent) => {
      if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  useEffect(() => {
    if (query.length < 3) {
      setResults([])
      return
    }

    const fetchResults = async () => {
      setLoading(true)
      try {
        const response = await api.get<User[]>(`/players/search?q=${query}`)
        setResults(response.data)
      } catch (error) {
        console.error('Failed to search players:', error)
      } finally {
        setLoading(false)
      }
    }

    const debounceTimeout = setTimeout(fetchResults, 300)
    return () => clearTimeout(debounceTimeout)
  }, [query])

  return (
    <div className="relative w-full md:w-64" ref={searchRef}>
      <div className="relative">
        <Input
          type="search"
          placeholder="Search players..."
          className="pl-10"
          value={query}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setQuery(e.target.value)}
          onFocus={() => setIsOpen(true)}
        />
        <div className="absolute inset-y-0 left-0 flex items-center pl-3">
          <Search className="h-5 w-5 text-muted-foreground" />
        </div>
      </div>
      {isOpen && (query.length > 0) && (
        <div className="absolute mt-1 w-full rounded-md border bg-popover shadow-lg z-10">
          {loading && <div className="p-4 text-sm text-muted-foreground">Loading...</div>}
          {!loading && results.length === 0 && query.length >= 3 && (
            <div className="p-4 text-sm text-muted-foreground">No results found.</div>
          )}
          {results.length > 0 && (
            <ul className="max-h-60 overflow-auto">
              {results.map((user) => (
                <li key={user.id}>
                  <AppLink
                    to={routes.player.profile(user.steam_id)}
                    className="flex items-center gap-3 p-3 hover:bg-accent"
                    onClick={() => {
                      setQuery('')
                      setIsOpen(false)
                    }}
                  >
                    <Avatar className="h-8 w-8">
                      <AvatarImage src={user.avatar_url} alt={user.nickname} />
                      <AvatarFallback>{user.nickname.charAt(0)}</AvatarFallback>
                    </Avatar>
                    <span className="font-medium">{user.nickname}</span>
                  </AppLink>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </div>
  )
} 