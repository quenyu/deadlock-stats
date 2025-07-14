import * as React from 'react'
import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <div className="flex h-[80vh] flex-col items-center justify-center">
      <h1 className="text-4xl font-bold">404 - Page Not Found</h1>
      <p className="mt-4 text-muted-foreground">
        The page you are looking for does not exist.
      </p>
      <Link to="/" className="mt-8 underline">
        Return to Home
      </Link>
    </div>
  )
} 