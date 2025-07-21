import { API_ENDPOINTS } from './api';

export const routes = {
  home: '/',
  player: {
    profile: (steamId = ':steamId') => `/player/${steamId}`,
    search: '/players/search',
  },
  builds: {
    list: '/builds',
    create: '/builds/create',
    view: (id = ':id') => `/builds/${id}`,
    edit: (id = ':id') => `/builds/${id}/edit`,
  },
  crosshairs: {
    list: '/crosshairs',
    create: '/crosshairs/create',
    view: (id = ':id') => `/crosshairs/${id}`,
    edit: (id = ':id') => `/crosshairs/${id}/edit`,
  },
  auth: {
    steam: API_ENDPOINTS.auth.steamLogin,
    logout: API_ENDPOINTS.auth.logout,
  },
  analytics: '/analytics',
  premium: '/premium',
} as const

export const isExternalRoute = (path: string): boolean => {
  return path.startsWith('http://') || path.startsWith('https://');
}; 