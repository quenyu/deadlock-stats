import { useUserStore } from '@/entities/user'
import { AuthBySteamButton } from '@/features/AuthBySteam'
import { Button } from '@/shared/ui/button'

export function AuthWidget() {
  const { user, isLoading } = useUserStore()

  if (isLoading) {
    return <div className="text-sm text-muted-foreground">Loading...</div>
  }

  if (!user) {
    return <AuthBySteamButton />
  }

  return (
    <div className="flex items-center gap-3">
      <img
        src={user.avatar_url}
        alt={user.nickname}
        className="h-8 w-8 rounded-full"
      />
      <span className="text-sm font-medium">{user.nickname}</span>
      <Button variant="ghost" size="sm" onClick={() => useUserStore.getState().logout()}>
        Logout
      </Button>
    </div>
  )
} 