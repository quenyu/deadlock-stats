import * as React from 'react'

export function HomePage() {
  return (
    <div className="flex flex-col items-center justify-center">
      <h1 className="text-4xl font-bold">Welcome to Deadlock Stats</h1>
      <p className="mt-4 text-muted-foreground">
        Monitor and analyze deadlocks in your applications
      </p>
    </div>
  )
} 