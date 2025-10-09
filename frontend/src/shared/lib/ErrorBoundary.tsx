import React from 'react'

interface ErrorBoundaryProps {
  children: React.ReactNode
  fallback?: (error: unknown, reset: () => void) => React.ReactNode
}

interface ErrorBoundaryState {
  hasError: boolean
  error: unknown
}

export class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = { hasError: false, error: null }

  static getDerivedStateFromError(error: unknown): Partial<ErrorBoundaryState> {
    return { hasError: true, error }
  }

  componentDidCatch(error: unknown, errorInfo: unknown) {
    // Intentionally minimal: hook logger/telemetry here later
    if (import.meta.env.DEV) {
      // eslint-disable-next-line no-console
      console.error('Unhandled error caught by ErrorBoundary:', error, errorInfo)
    }
  }

  private reset = () => {
    this.setState({ hasError: false, error: null })
  }

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.reset)
      }
      return null
    }
    return this.props.children
  }
}


