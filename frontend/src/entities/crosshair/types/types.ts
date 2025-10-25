import { User } from '../../user'
import type { CrosshairSettings as CrosshairSettingsType } from '@/shared/lib/validation'

// Re-export from validation schema for consistency
export type CrosshairSettings = CrosshairSettingsType

export interface CrosshairPreset {
  name: string
  settings: CrosshairSettings
}

// Matches backend: backend/internal/domain/crosshair.go Crosshair
export interface Crosshair {
  id: string
  author_id: string
  author?: User
  title: string
  description: string
  settings: CrosshairSettings | string // Can be object or JSON string
  likes_count: number
  is_public: boolean
  view_count: number
  created_at: string
  updated_at: string
}

// Matches backend: backend/internal/domain/crosshair.go CrosshairLike
export interface CrosshairLike {
  id: string
  user_id: string
  crosshair_id: string
  created_at: string
}

// Extended type for UI with client-side fields
export interface CrosshairListItem extends Crosshair {
  is_liked: boolean // Client-side only
  author_name?: string // Denormalized from author
  author_avatar?: string // Denormalized from author
} 