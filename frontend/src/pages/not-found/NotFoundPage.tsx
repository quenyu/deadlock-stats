import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <div className="container flex h-screen w-screen flex-col items-center justify-center">
      <div className="mx-auto flex max-w-[420px] flex-col items-center justify-center text-center">
        <h1 className="text-4xl font-bold">404</h1>
        <p className="mb-4 text-muted-foreground">
          Page not found. Check the address or go back.
        </p>
        <Link
          to="/"
          className="rounded bg-primary px-4 py-2 font-medium text-primary-foreground hover:bg-primary/90"
        >
          Go to Homepage
        </Link>
      </div>
    </div>
  )
} 