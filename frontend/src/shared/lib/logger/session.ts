import { v4 as uuidv4 } from 'uuid'

export class SessionManager {
  private sessionId: string
  private userId: string | null = null

  constructor() {
    this.sessionId = this.generateSessionId()
  }

  private generateSessionId(): string {
    const stored = sessionStorage.getItem('logger-session-id')
    if (stored) {
      return stored
    }

    const newSessionId = uuidv4()
    sessionStorage.setItem('logger-session-id', newSessionId)
    return newSessionId
  }

  getSessionId(): string {
    return this.sessionId
  }

  setUserId(id: string | null): void {
    this.userId = id
  }

  getUserId(): string | null {
    return this.userId
  }

  reset(): void {
    this.sessionId = uuidv4()
    sessionStorage.setItem('logger-session-id', this.sessionId)
    this.userId = null
  }
}

