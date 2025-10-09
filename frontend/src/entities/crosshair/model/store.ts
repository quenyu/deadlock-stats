import { create } from 'zustand'
import { CrosshairSettings, CrosshairListItem } from '../types/types'
import { PRESETS } from '../lib/presets'
import { publishCrosshair } from '../api/publishCrosshair'
import { likeCrosshair, unlikeCrosshair } from '../api/likeCrosshair'
import { fetchPublishedCrosshairs } from '../api/fetchPublishedCrosshairs'
import { logger } from '@/shared/lib/logger'

const STORAGE_KEY = 'crosshair-settings'

interface CrosshairStore {
  settings: CrosshairSettings
  published: CrosshairListItem[]
  loading: boolean
  setSettings: (settings: Partial<CrosshairSettings>) => void
  loadPreset: (name: string) => void
  reset: () => void
  loadPublished: () => Promise<void>
  publish: (title: string, description: string, isPublic?: boolean) => Promise<void>
  like: (id: string) => Promise<void>
  unlike: (id: string) => Promise<void>
}

const defaultSettings = PRESETS.default

export const useCrosshairStore = create<CrosshairStore>((set, get) => ({
  settings: (() => {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      try {
        return { ...defaultSettings, ...JSON.parse(saved) }
      } catch {
        return defaultSettings
      }
    }
    return defaultSettings
  })(),
  published: [],
  loading: false,
  setSettings: (patch) => set(state => {
    const newSettings = { ...state.settings, ...patch }
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newSettings))
    return { settings: newSettings }
  }),
  loadPreset: (name) => set(() => {
    const preset = PRESETS[name] || defaultSettings
    localStorage.setItem(STORAGE_KEY, JSON.stringify(preset))
    return { settings: preset }
  }),
  reset: () => set(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(defaultSettings))
    return { settings: defaultSettings }
  }),
  loadPublished: async () => {
    set({ loading: true })
    try {
      const data = await fetchPublishedCrosshairs()
      set({ published: data.crosshairs })
    } catch (error) {
      logger.error('Failed to load published crosshairs', error)
    } finally {
      set({ loading: false })
    }
  },
  publish: async (title, description, isPublic = true) => {
    const { settings } = get()
    try {
      await publishCrosshair({ 
        title, 
        description, 
        settings, 
        is_public: isPublic 
      })
      await get().loadPublished()
    } catch (error) {
      logger.error('Failed to publish crosshair', error)
      throw error
    }
  },
  like: async (id) => {
    try {
      await likeCrosshair(id)
      await get().loadPublished()
    } catch (error) {
      logger.error('Failed to like crosshair', error)
      throw error
    }
  },
  unlike: async (id) => {
    try {
      await unlikeCrosshair(id)
      await get().loadPublished()
    } catch (error) {
      logger.error('Failed to unlike crosshair', error)
      throw error
    }
  },
})) 