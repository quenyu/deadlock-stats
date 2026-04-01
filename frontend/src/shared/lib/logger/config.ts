export type LogLevel = 'debug' | 'log' | 'info' | 'warn' | 'error'

export interface LoggerConfig {
  enabled: boolean
  
  level: LogLevel
  
  includeTimestamp: boolean
  
  includeSessionId: boolean
  
  includeUserId: boolean
  
  /** Callback for external logging services (e.g., Sentry) */
  onLog?: (entry: LogEntry) => void
}

export interface LogEntry {
  timestamp: string
  level: LogLevel
  scope?: string
  message: string
  context?: LogContext
  sessionId?: string
  userId?: string
}

export type LogContext = Record<string, unknown>

/**
 * Default configuration for production
 */
export const defaultConfig = (): LoggerConfig => ({
  enabled: import.meta.env.DEV,
  level: import.meta.env.DEV ? 'debug' : 'info',
  includeTimestamp: true,
  includeSessionId: true,
  includeUserId: true,
})

/**
 * Development configuration
 */
export const developmentConfig = (): LoggerConfig => ({
  enabled: true,
  level: 'debug',
  includeTimestamp: true,
  includeSessionId: true,
  includeUserId: true,
})

/**
 * Production configuration
 */
export const productionConfig = (): LoggerConfig => ({
  enabled: false,
  level: 'error',
  includeTimestamp: true,
  includeSessionId: true,
  includeUserId: true,
})

export const LOG_LEVELS: Record<LogLevel, number> = {
  debug: 0,
  log: 1,
  info: 2,
  warn: 3,
  error: 4,
}

export function shouldLog(currentLevel: LogLevel, configLevel: LogLevel, enabled: boolean): boolean {
  if (currentLevel === 'error') {
    return true
  }
  
  if (!enabled) {
    return false
  }
  
  return LOG_LEVELS[currentLevel] >= LOG_LEVELS[configLevel]
}

