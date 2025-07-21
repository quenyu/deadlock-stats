export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const API_ENDPOINTS = {
  auth: {
    steamLogin: `${API_BASE_URL}/auth/steam/login`,
    logout: `${API_BASE_URL}/auth/logout`,
  },
  player: {
    profile: (steamId: string) => `${API_BASE_URL}/players/${steamId}`,
  },
} 