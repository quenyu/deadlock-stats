import { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useUserStore } from '@/entities/user'
import { routes } from '@/shared/constants/routes'

interface ProtectedRouteProps {
  children: ReactNode
  redirectTo?: string
}

export function ProtectedRoute({ 
  children, 
  redirectTo = routes.home 
}: ProtectedRouteProps) {
  const { user, isLoading } = useUserStore()
  
  if (isLoading) {
    return <div className="p-8 text-center text-muted-foreground">Downloading...</div>
  }
  
  if (!user) {
    return <Navigate to={redirectTo} replace />
  }

  return <>{children}</>
} 