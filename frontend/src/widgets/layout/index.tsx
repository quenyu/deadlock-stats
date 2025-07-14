import * as React from 'react'
import { Outlet } from 'react-router-dom'

export function Layout() {
  return (
    <div className="min-h-screen bg-background">
      <main className="container mx-auto py-4">
        <Outlet />
      </main>
    </div>
  )
} 