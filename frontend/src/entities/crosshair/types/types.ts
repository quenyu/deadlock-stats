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

export interface PublishedCrosshair {
  id: string
  settings: CrosshairSettings
  likes: number
  author: string
  createdAt: string
} 