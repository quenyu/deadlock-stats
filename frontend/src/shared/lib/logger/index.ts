/**
 * Logger - Modular logging system
 * 
 * @example
 * ```ts
 * import { logger, createLogger } from '@/shared/lib/logger'
 * 
 * // Global logger
 * logger.info('Application started')
 * 
 * // Scoped logger
 * const log = createLogger('MyComponent')
 * log.debug('Component mounted')
 * log.error('Failed to load data', { error })
 * ```
 */

import { Logger } from './logger'

export type { LoggerConfig, LogLevel, LogContext, LogEntry } from './config'
export type { ScopedLogger } from './types'

export { Logger } from './logger'
export { LogFormatter } from './formatter'
export { ConsoleTransport } from './transport'
export { SessionManager } from './session'

export { 
  defaultConfig, 
  developmentConfig, 
  productionConfig,
  shouldLog,
  LOG_LEVELS,
} from './config'

export const logger = new Logger()

/**
 * Create scoped logger
 * 
 * @param scopeName - Scope name (e.g., component name)
 * @param context - Additional context to include in all logs
 * @returns Scoped logger instance
 * 
 * @example
 * ```ts
 * const log = createLogger('UserProfile', { userId: '123' })
 * log.info('Profile loaded') // [UserProfile] Profile loaded { userId: '123' }
 * ```
 */
export const createLogger = (scopeName: string, context?: Record<string, unknown>) => {
  return logger.createScope(scopeName, context)
}

