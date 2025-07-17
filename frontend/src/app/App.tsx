import { useEffect } from 'react'
import { RouterProvider } from 'react-router-dom'
import { ThemeProvider } from './providers/theme'
import { router } from './providers/router'
import useUserStore from '@/entities/user/model/store'

export function App() {
  const fetchUser = useUserStore((state) => state.fetchUser)

  useEffect(() => {
    fetchUser()
  }, [fetchUser])

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <RouterProvider router={router} />
    </ThemeProvider>
  )
} 