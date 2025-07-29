import { create } from 'zustand'
import { CrosshairSettings } from '../types/types'
import { PRESETS } from '../lib/presets'

const STORAGE_KEY = 'crosshair-settings'

interface CrosshairStore {
  settings: CrosshairSettings
  setSettings: (settings: Partial<CrosshairSettings>) => void
  loadPreset: (name: string) => void
  reset: () => void
  publish: () => Promise<void>
  like: () => Promise<void>
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
  publish: async () => {
    // TODO: do publish
    alert('Publish crosshair (stub)')
  },
  like: async () => {
    // TODO: do like
    alert('Like crosshair (stub)')
  },
})) 