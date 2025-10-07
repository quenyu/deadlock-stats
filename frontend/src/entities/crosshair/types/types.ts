export interface CrosshairSettings {
  color: string
  thickness: number
  length: number
  gap: number
  dot: boolean
  opacity: number
  pipOpacity: number
  dotOutlineOpacity: number
  hitMarkerDuration: number
  pipBorder: boolean
  pipGapStatic: boolean
}

export interface CrosshairPreset {
  name: string
  settings: CrosshairSettings
}

export interface CrosshairListItem {
  id: string
  title: string
  description: string
  settings: CrosshairSettings
  likes_count: number
  created_at: string
  author_id: string
  author_name?: string
  author_avatar?: string
  is_liked: boolean
  is_public: boolean
  view_count: number
}

export interface PublishedCrosshair {
  id: string
  settings: CrosshairSettings
  likes: number
  author_id: string
  createdAt: string
} 