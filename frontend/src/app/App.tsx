import { useEffect } from 'react'
import { RouterProvider } from 'react-router-dom'
import { Toaster } from 'sonner'
import { ThemeProvider } from './providers/theme'
import { router } from './providers/router'
import { useUserStore } from '@/entities/user'

export function App() {
  const fetchUser = useUserStore((state) => state.fetchUser)

  useEffect(() => {
    fetchUser()
  }, [fetchUser])

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <RouterProvider router={router} />
      <Toaster 
        position="top-right" 
        richColors 
        closeButton
        duration={4000}
      />
    </ThemeProvider>
  )
} 