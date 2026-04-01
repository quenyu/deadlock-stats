const isDev = import.meta.env.DEV

export const logger = {
    //eslint-
  log: (...args: any[]) => {
    if (isDev) {
      console.log('[LOG]', ...args)
    }
  },
  
  error: (...args: any[]) => {
    if (isDev) {
      console.error('[ERROR]', ...args)
    }
    // В production можно отправлять в Sentry или другую систему мониторинга
  },
  
  warn: (...args: any[]) => {
    if (isDev) {
      console.warn('[WARN]', ...args)
    }
  },
  
  info: (...args: any[]) => {
    if (isDev) {
      console.info('[INFO]', ...args)
    }
  },
  
  debug: (...args: any[]) => {
    if (isDev) {
      console.debug('[DEBUG]', ...args)
    }
  }
}

export default logger