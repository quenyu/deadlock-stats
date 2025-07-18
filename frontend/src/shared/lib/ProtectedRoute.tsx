import { Navigate } from 'react-router-dom'
import { useUserStore } from '@/entities/user'
import { routes } from '@/shared/constants/routes'
import { ReactElement } from 'react'

interface ProtectedRouteProps {
  children: ReactElement
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { user, isLoading } = useUserStore()

  if (isLoading) {
    return <div className="p-8 text-center text-muted-foreground">Loading...</div>
  }

  if (!user) {
    return <Navigate to={routes.home} replace />
  }

  return children
} 