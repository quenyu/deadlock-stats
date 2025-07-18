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
  auth: {
    steam: '/api/v1/auth/steam/login',
    callback: '/api/v1/auth/steam/callback',
    logout: '/api/v1/auth/logout',
  },
  analytics: '/analytics',
  premium: '/premium',
} as const 