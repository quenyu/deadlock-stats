import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useUserStore } from '@/entities/user'

export function SteamCallbackPage() {
  const fetchUser = useUserStore((s) => s.fetchUser)
  const navigate = useNavigate()

  useEffect(() => {
    async function handle() {
      await fetchUser()
      navigate('/', { replace: true })
    }
    handle()
  }, [fetchUser, navigate])

  return <div className="p-8 text-center text-muted-foreground">Authenticating...</div>
} 