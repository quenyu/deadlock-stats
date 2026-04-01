
import type { LogEntry, LogLevel, LogContext } from './config'

export class LogFormatter {
  static formatEntry(
    level: LogLevel,
    message: string,
    context?: LogContext,
    scope?: string,
    sessionId?: string,
    userId?: string,
    includeTimestamp = true
  ): LogEntry {
    const entry: LogEntry = {
      timestamp: includeTimestamp ? new Date().toISOString() : '',
      level,
      message,
    }

    if (scope) entry.scope = scope
    if (context) entry.context = context
    if (sessionId) entry.sessionId = sessionId
    if (userId) entry.userId = userId

    return entry
  }

  static formatPrefix(scope?: string, timestamp?: string): string {
    const parts: string[] = []
    
    if (timestamp) {
      parts.push(`[${new Date(timestamp).toISOString()}]`)
    }
    
    if (scope) {
      parts.push(`[${scope}]`)
    }
    
    return parts.join(' ')
  }

  static formatMessage(entry: LogEntry): string {
    const prefix = this.formatPrefix(entry.scope, entry.timestamp)
    return prefix ? `${prefix} ${entry.message}` : entry.message
  }

  static serializeContext(context?: LogContext): string {
    if (!context || Object.keys(context).length === 0) {
      return ''
    }
    
    try {
      return JSON.stringify(context, null, 2)
    } catch (error) {
      return String(context)
    }
  }
}

