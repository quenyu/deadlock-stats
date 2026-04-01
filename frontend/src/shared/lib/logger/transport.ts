import type { LogEntry } from './config'
import { LogFormatter } from './formatter'

export class ConsoleTransport {
  static output(entry: LogEntry): void {
    const { level, context } = entry
    const formattedMessage = LogFormatter.formatMessage(entry)

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

  static time(label: string): void {
    // eslint-disable-next-line no-console
    console.time(label)
  }

  static timeEnd(label: string): void {
    // eslint-disable-next-line no-console
    console.timeEnd(label)
  }

  static group(label: string): void {
    // eslint-disable-next-line no-console
    console.group(label)
  }

  static groupEnd(): void {
    // eslint-disable-next-line no-console
    console.groupEnd()
  }

  static table(data: unknown): void {
    // eslint-disable-next-line no-console
    console.table(data)
  }
}

