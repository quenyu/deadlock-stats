import { Outlet } from 'react-router-dom'

export function Layout() {
  return (
    <div className="relative flex min-h-screen flex-col">
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  )
} 