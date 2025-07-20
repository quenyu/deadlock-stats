import { Outlet } from 'react-router-dom'
import { AuthWidget } from '@/widgets/auth/AuthWidget'
import { Navbar } from '@/widgets/Navbar'
import { ThemeToggle } from '@/features/theme-toggle/ThemeToggle'

export function Layout() {
  return (
    <div className="relative flex min-h-screen flex-col bg-background">
      <header className="sticky top-0 z-50 w-full border-b-0 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-6">
            <a href="/" className="hidden md:flex items-center space-x-2">
              <span className="text-lg font-bold">Deadlock Stats</span>
            </a>
            <Navbar />
          </div>
          <div className="flex items-center space-x-4">
            <ThemeToggle />
            <AuthWidget />
          </div>
        </div>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  )
} 