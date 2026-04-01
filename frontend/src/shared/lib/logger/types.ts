import type { LogContext } from './config'

export interface ScopedLogger {
  debug: (message: string, context?: LogContext) => void
  log: (message: string, context?: LogContext) => void
  info: (message: string, context?: LogContext) => void
  warn: (message: string, context?: LogContext) => void
  error: (message: string, context?: LogContext) => void
  time: (label: string) => void
  timeEnd: (label: string) => void
  group: (label: string) => void
  groupEnd: () => void
  table: (data: unknown) => void
  createScope: (scopeName: string, scopeContext?: LogContext) => ScopedLogger
}

