import {
  BookUser,
  LogOut,
  Settings,
  Swords,
  User as UserIcon,
} from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu'
import { useUserStore } from '@/entities/user'
import { AuthBySteamButton } from '@/features/AuthBySteam'
import { Button } from '@/shared/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/shared/ui/avatar'
import { useEffect } from 'react'
import { AppLink } from '@/shared/ui/AppLink/AppLink'
import { routes } from '@/shared/constants/routes'
import { createLogger } from '@/shared/lib/logger'

const log = createLogger('AuthWidget')

export function AuthWidget() {
  const { user, isLoading, error } = useUserStore()

  useEffect(() => {
    log.debug('Rendering with state', {
      user: user ? `${user.nickname} (${user.id})` : 'null',
      isLoading,
      error,
    })
  }, [user, isLoading, error])

  if (isLoading) {
    return (
      <div className="h-8 w-28 animate-pulse rounded-md bg-muted" />
    )
  }

  if (error) {
    log.error('Authentication error', error)
    return <AuthBySteamButton />
  }

  if (!user) {
    return <AuthBySteamButton />
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="lg"
          className="flex items-center gap-2 p-2 h-auto focus-visible:ring-0 focus-visible:ring-offset-0"
        >
          <Avatar className="h-8 w-8">
            <AvatarImage src={user.avatar_url} alt={user.nickname} />
            <AvatarFallback>
              {user.nickname.charAt(0).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <span className="text-sm font-medium hidden sm:inline-block">{user.nickname}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56" align="end" forceMount>
        <DropdownMenuLabel className="font-normal">
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">{user.nickname}</p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem asChild>
            <AppLink to={routes.player.profile(user.steam_id)}>
              <UserIcon className="mr-2 h-4 w-4" />
              <span>Profile</span>
            </AppLink>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <a
              href={user.profile_url}
              target="_blank"
              rel="noopener noreferrer"
              className="w-full flex items-center"
            >
              <BookUser className="mr-2 h-4 w-4" />
              <span>Steam Profile</span>
            </a>
          </DropdownMenuItem>
          <DropdownMenuItem disabled>
            <Swords className="mr-2 h-4 w-4" />
            <span>Matches</span>
          </DropdownMenuItem>
          <DropdownMenuItem disabled>
            <Settings className="mr-2 h-4 w-4" />
            <span>Settings</span>
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem onSelect={() => useUserStore.getState().logout()}>
          <LogOut className="mr-2 h-4 w-4" />
          <span>Logout</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
} 