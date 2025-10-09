/**
 * Enterprise-grade logging utility for Deadlock Stats
 * 
 * Features:
 * - Environment-aware (dev/production)
 * - Multiple log levels (debug, log, info, warn, error)
 * - Structured logging with context
 * - Ready for external monitoring integration (Sentry, LogRocket, etc.)
 * - Performance tracking
 * - TypeScript type-safe
 * 
 * @example
 * ```typescript
 * import { logger } from '@/shared/lib/logger'
 * 
 * logger.info('User logged in', { userId: '123' })
 * logger.error('API call failed', { endpoint: '/api/users', error })
 * logger.debug('Component rendered', { props })
 * ```
 */

export type LogLevel = 'debug' | 'log' | 'info' | 'warn' | 'error'

export type LogContext = Record<string, unknown>

export interface LogEntry {
  level: LogLevel
  message: string
  context?: LogContext
  timestamp: number
}

export interface LoggerConfig {
  enabled: boolean
  minLevel: LogLevel
  enableTimestamps: boolean
  enableStackTrace: boolean
  onLog?: (entry: LogEntry) => void
}

const LOG_LEVELS: Record<LogLevel, number> = {
  debug: 0,
  log: 1,
  info: 2,
  warn: 3,
  error: 4,
}

class Logger {
  private config: LoggerConfig
  private sessionId: string

  constructor(config?: Partial<LoggerConfig>) {
    const isDev = import.meta.env.DEV
    const isProd = import.meta.env.PROD

    this.config = {
      enabled: isDev,
      minLevel: isDev ? 'debug' : 'warn',
      enableTimestamps: true,
      enableStackTrace: isProd,
      ...config,
    }

    this.sessionId = this.generateSessionId()

    if (isDev) {
      this.info('Logger initialized', {
        mode: isDev ? 'development' : 'production',
        minLevel: this.config.minLevel,
        sessionId: this.sessionId,
      })
    }
  }

  /**
   * Debug level logging - verbose diagnostic information
   * Only shows in development
   */
  debug(message: string, context?: LogContext): void {
    this.log('debug', message, context)
  }

  /**
   * General logging - standard application flow
   */
  log(message: string, context?: LogContext): void
  log(level: LogLevel, message: string, context?: LogContext): void
  log(
    messageOrLevel: string | LogLevel,
    messageOrContext?: string | LogContext,
    context?: LogContext
  ): void {
    let level: LogLevel
    let message: string
    let finalContext: LogContext | undefined

    if (this.isLogLevel(messageOrLevel)) {
      level = messageOrLevel
      message = messageOrContext as string
      finalContext = context
    } else {
      level = 'log'
      message = messageOrLevel
      finalContext = messageOrContext as LogContext | undefined
    }

    this.writeLog(level, message, finalContext)
  }

  /**
   * Info level - important informational messages
   */
  info(message: string, context?: LogContext): void {
    this.log('info', message, context)
  }

  /**
   * Warning level - potentially harmful situations
   */
  warn(message: string, context?: LogContext): void {
    this.log('warn', message, context)
  }

  /**
   * Error level - error events that might still allow the app to continue
   */
  error(message: string, error?: Error | unknown, context?: LogContext): void {
    const errorContext: LogContext = {
      ...context,
      ...(error instanceof Error && {
        errorName: error.name,
        errorMessage: error.message,
        errorStack: this.config.enableStackTrace ? error.stack : undefined,
      }),
      ...(error && !(error instanceof Error) && { error }),
    }

    this.log('error', message, errorContext)

    // Send to error tracking service in production
    if (import.meta.env.PROD && this.config.onLog) {
      this.sendToExternalService('error', message, errorContext)
    }
  }

  /**
   * Performance timing utility
   */
  time(label: string): void {
    if (this.shouldLog('debug')) {
      // eslint-disable-next-line no-console
      console.time(label)
    }
  }

  /**
   * End performance timing
   */
  timeEnd(label: string): void {
    if (this.shouldLog('debug')) {
      // eslint-disable-next-line no-console
      console.timeEnd(label)
    }
  }

  /**
   * Group related logs together
   */
  group(label: string): void {
    if (this.shouldLog('debug') && this.config.enabled) {
      // eslint-disable-next-line no-console
      console.group(label)
    }
  }

  /**
   * End log group
   */
  groupEnd(): void {
    if (this.shouldLog('debug') && this.config.enabled) {
      // eslint-disable-next-line no-console
      console.groupEnd()
    }
  }

  /**
   * Table display for structured data
   */
  table(data: unknown): void {
    if (this.shouldLog('debug') && this.config.enabled) {
      // eslint-disable-next-line no-console
      console.table(data)
    }
  }

  /**
   * Create a scoped logger with predefined context
   */
  createScope(scopeName: string, scopeContext?: LogContext): ScopedLogger {
    return new ScopedLogger(this, scopeName, scopeContext)
  }

  private writeLog(level: LogLevel, message: string, context?: LogContext): void {
    if (!this.shouldLog(level)) {
      return
    }

    const entry: LogEntry = {
      level,
      message,
      context: {
        ...context,
        sessionId: this.sessionId,
        timestamp: this.config.enableTimestamps ? Date.now() : undefined,
      },
      timestamp: Date.now(),
    }

    // Call external log handler if provided
    if (this.config.onLog) {
      this.config.onLog(entry)
    }

    // Console output
    if (this.config.enabled) {
      this.consoleOutput(entry)
    }
  }

  private consoleOutput(entry: LogEntry): void {
    const { level, message, context } = entry
    const prefix = this.config.enableTimestamps
      ? `[${new Date(entry.timestamp).toISOString()}]`
      : ''

    const formattedMessage = prefix ? `${prefix} ${message}` : message

    switch (level) {
      case 'debug':
        // eslint-disable-next-line no-console
        console.debug(formattedMessage, context || '')
        break
      case 'log':
        // eslint-disable-next-line no-console
        console.log(formattedMessage, context || '')
        break
      case 'info':
        // eslint-disable-next-line no-console
        console.info(formattedMessage, context || '')
        break
      case 'warn':
        // eslint-disable-next-line no-console
        console.warn(formattedMessage, context || '')
        break
      case 'error':
        // eslint-disable-next-line no-console
        console.error(formattedMessage, context || '')
        break
    }
  }

  private shouldLog(level: LogLevel): boolean {
    if (!this.config.enabled && level !== 'error') {
      return false
    }
    return LOG_LEVELS[level] >= LOG_LEVELS[this.config.minLevel]
  }

  private isLogLevel(value: unknown): value is LogLevel {
    return typeof value === 'string' && value in LOG_LEVELS
  }

  private generateSessionId(): string {
    return `${Date.now()}-${Math.random().toString(36).substring(2, 9)}`
  }

  private sendToExternalService(
    level: LogLevel,
    message: string,
    context: LogContext
  ): void {
    // Hook for Sentry, LogRocket, or other monitoring services
    // Example:
    // if (window.Sentry) {
    //   Sentry.captureMessage(message, { level, extra: context })
    // }
    if (this.config.onLog) {
      this.config.onLog({ level, message, context, timestamp: Date.now() })
    }
  }

  /**
   * Update logger configuration at runtime
   */
  configure(config: Partial<LoggerConfig>): void {
    this.config = { ...this.config, ...config }
  }

  /**
   * Get current configuration
   */
  getConfig(): Readonly<LoggerConfig> {
    return { ...this.config }
  }
}

/**
 * Scoped logger for component/module-specific logging
 */
class ScopedLogger {
  constructor(
    private parent: Logger,
    private scope: string,
    private scopeContext?: LogContext
  ) {}

  private addScope(message: string): string {
    return `[${this.scope}] ${message}`
  }

  private mergeContext(context?: LogContext): LogContext {
    return { ...this.scopeContext, ...context }
  }

  debug(message: string, context?: LogContext): void {
    this.parent.debug(this.addScope(message), this.mergeContext(context))
  }

  log(message: string, context?: LogContext): void {
    this.parent.log(this.addScope(message), this.mergeContext(context))
  }

  info(message: string, context?: LogContext): void {
    this.parent.info(this.addScope(message), this.mergeContext(context))
  }

  warn(message: string, context?: LogContext): void {
    this.parent.warn(this.addScope(message), this.mergeContext(context))
  }

  error(message: string, error?: Error | unknown, context?: LogContext): void {
    this.parent.error(this.addScope(message), error, this.mergeContext(context))
  }
}

// Singleton instance
export const logger = new Logger()

// For testing or custom instances
export { Logger }

/**
 * Create a logger for a specific component or module
 * 
 * @example
 * ```typescript
 * const log = createLogger('AuthWidget', { component: 'auth' })
 * log.info('User authenticated')
 * ```
 */
export function createLogger(scope: string, context?: LogContext): ScopedLogger {
  return logger.createScope(scope, context)
}

/**
 * Configure global logger
 * 
 * @example
 * ```typescript
 * configureLogger({
 *   enabled: true,
 *   minLevel: 'info',
 *   onLog: (entry) => {
 *     // Send to Sentry
 *   }
 * })
 * ```
 */
export function configureLogger(config: Partial<LoggerConfig>): void {
  logger.configure(config)
}

