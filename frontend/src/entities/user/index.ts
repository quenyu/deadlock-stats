import useUserStore from './model/store'
import { fetchCurrentUser } from './api/fetchCurrentUser'
import type { User } from './types/types'

export { useUserStore, fetchCurrentUser }
export type { User }