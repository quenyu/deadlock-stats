import React from 'react'
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/shared/ui/pagination'
import { UserCard } from '@/shared/ui/UserCard'
import { User } from '@/entities/user'
import { PlayerSearchResult } from '@/shared/lib/validation'

type UserCardData = PlayerSearchResult | User

interface PaginatedResultsProps {
  results: UserCardData[]
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
  onPageChange: (page: number) => void
  onUserClick?: (user: UserCardData) => void
  showExtendedInfo?: boolean
  loading?: boolean
}

export const PaginatedResults: React.FC<PaginatedResultsProps> = ({
  results,
  totalCount,
  page,
  pageSize,
  totalPages,
  onPageChange,
  onUserClick,
  showExtendedInfo = false,
  loading = false
}) => {
  const generatePaginationItems = () => {
    const items = []
    const maxVisiblePages = 5
    const startPage = Math.max(1, page - Math.floor(maxVisiblePages / 2))
    const endPage = Math.min(totalPages, startPage + maxVisiblePages - 1)

    if (page > 1) {
      items.push(
        <PaginationItem key="prev">
          <PaginationPrevious 
            href="#" 
            onClick={(e) => {
              e.preventDefault()
              onPageChange(page - 1)
            }}
          />
        </PaginationItem>
      )
    }

    if (startPage > 1) {
      items.push(
        <PaginationItem key={1}>
          <PaginationLink 
            href="#" 
            isActive={page === 1}
            onClick={(e) => {
              e.preventDefault()
              onPageChange(1)
            }}
          >
            1
          </PaginationLink>
        </PaginationItem>
      )
      
      if (startPage > 2) {
        items.push(
          <PaginationItem key="ellipsis1">
            <PaginationEllipsis />
          </PaginationItem>
        )
      }
    }

    for (let i = startPage; i <= endPage; i++) {
      items.push(
        <PaginationItem key={i}>
          <PaginationLink 
            href="#" 
            isActive={page === i}
            onClick={(e) => {
              e.preventDefault()
              onPageChange(i)
            }}
          >
            {i}
          </PaginationLink>
        </PaginationItem>
      )
    }

    if (endPage < totalPages) {
      if (endPage < totalPages - 1) {
        items.push(
          <PaginationItem key="ellipsis2">
            <PaginationEllipsis />
          </PaginationItem>
        )
      }
      
      items.push(
        <PaginationItem key={totalPages}>
          <PaginationLink 
            href="#" 
            isActive={page === totalPages}
            onClick={(e) => {
              e.preventDefault()
              onPageChange(totalPages)
            }}
          >
            {totalPages}
          </PaginationLink>
        </PaginationItem>
      )
    }

    if (page < totalPages) {
      items.push(
        <PaginationItem key="next">
          <PaginationNext 
            href="#" 
            onClick={(e) => {
              e.preventDefault()
              onPageChange(page + 1)
            }}
          />
        </PaginationItem>
      )
    }

    return items
  }

  if (loading) {
    return (
      <div className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: pageSize }).map((_, i) => (
            <div key={i} className="h-32 bg-muted animate-pulse rounded-lg" />
          ))}
        </div>
      </div>
    )
  }

  if (results.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-muted-foreground">No results found</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="text-sm text-muted-foreground">
        Showing {((page - 1) * pageSize) + 1} to {Math.min(page * pageSize, totalCount)} of {totalCount} results
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {results.map((user) => {
          const key = 'steamId' in user ? user.steamId : user.steam_id
          return (
            <UserCard
              key={key}
              user={user}
              onClick={() => onUserClick?.(user)}
              showExtendedInfo={showExtendedInfo}
            />
          )
        })}
      </div>

      {totalPages > 1 && (
        <Pagination>
          <PaginationContent>
            {generatePaginationItems()}
          </PaginationContent>
        </Pagination>
      )}
    </div>
  )
} 