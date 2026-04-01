import type { LoggerConfig, LogLevel, LogContext } from './config'
import { defaultConfig, shouldLog } from './config'
import type { ScopedLogger } from './types'
import { LogFormatter } from './formatter'
import { ConsoleTransport } from './transport'
import { SessionManager } from './session'

export class Logger {
  private config: LoggerConfig
  private session: SessionManager

  constructor(config?: Partial<LoggerConfig>) {
    this.config = { ...defaultConfig(), ...config }
    this.session = new SessionManager()
    this.initialize()
  }

  private initialize(): void {
    if (import.meta.env.PROD) {
      window.onerror = (message, source, lineno, colno, error) => {
        this.error('Unhandled window error', { 
          message, 
          source, 
          lineno, 
          colno, 
          error: error?.message 
        })
      }

      window.onunhandledrejection = (event) => {
        this.error('Unhandled promise rejection', { 
          reason: event.reason 
        })
      }
    }
  }

  setUserId(id: string | null): void {
    this.session.setUserId(id)
  }

  getUserId(): string | null {
    return this.session.getUserId()
  }

  private shouldLog(level: LogLevel): boolean {
    return shouldLog(level, this.config.level, this.config.enabled)
  }

  private formatMessage(
    level: LogLevel,
    message: string,
    context?: LogContext,
    scope?: string
  ) {
    const sessionId = this.config.includeSessionId 
      ? this.session.getSessionId() 
      : undefined
      
    const userId = this.config.includeUserId 
      ? this.session.getUserId() || undefined
      : undefined

    return LogFormatter.formatEntry(
      level,
      message,
      context,
      scope,
      sessionId,
      userId,
      this.config.includeTimestamp
    )
  }

  private output(entry: ReturnType<typeof this.formatMessage>): void {
    if (!this.shouldLog(entry.level)) {
      return
    }

    ConsoleTransport.output(entry)

    if (this.config.onLog) {
      this.config.onLog(entry)
    }
  }

  debug(message: string, context?: LogContext): void {
    this.output(this.formatMessage('debug', message, context))
  }

  log(message: string, context?: LogContext): void {
    this.output(this.formatMessage('log', message, context))
  }

  info(message: string, context?: LogContext): void {
    this.output(this.formatMessage('info', message, context))
  }

  warn(message: string, context?: LogContext): void {
    this.output(this.formatMessage('warn', message, context))
  }

  error(message: string, context?: LogContext): void {
    this.output(this.formatMessage('error', message, context))
  }

  time(label: string): void {
    if (this.shouldLog('debug')) {
      ConsoleTransport.time(label)
    }
  }

  timeEnd(label: string): void {
    if (this.shouldLog('debug')) {
      ConsoleTransport.timeEnd(label)
    }
  }

  group(label: string): void {
    if (this.shouldLog('debug') && this.config.enabled) {
      ConsoleTransport.group(label)
    }
  }

  groupEnd(): void {
    if (this.shouldLog('debug') && this.config.enabled) {
      ConsoleTransport.groupEnd()
    }
  }

  table(data: unknown): void {
    if (this.shouldLog('debug') && this.config.enabled) {
      ConsoleTransport.table(data)
    }
  }

  createScope(scopeName: string, scopeContext?: LogContext): ScopedLogger {
    return {
      debug: (message, context) => 
        this.output(this.formatMessage('debug', message, { ...scopeContext, ...context }, scopeName)),
      
      log: (message, context) => 
        this.output(this.formatMessage('log', message, { ...scopeContext, ...context }, scopeName)),
      
      info: (message, context) => 
        this.output(this.formatMessage('info', message, { ...scopeContext, ...context }, scopeName)),
      
      warn: (message, context) => 
        this.output(this.formatMessage('warn', message, { ...scopeContext, ...context }, scopeName)),
      
      error: (message, context) => 
        this.output(this.formatMessage('error', message, { ...scopeContext, ...context }, scopeName)),
      
      time: (label) => this.time(`${scopeName}:${label}`),
      
      timeEnd: (label) => this.timeEnd(`${scopeName}:${label}`),
      
      group: (label) => this.group(`${scopeName}:${label}`),
      
      groupEnd: () => this.groupEnd(),
      
      table: (data) => this.table(data),
      
      createScope: (subScopeName, subScopeContext) => 
        this.createScope(`${scopeName}:${subScopeName}`, { ...scopeContext, ...subScopeContext }),
    }
  }
}

