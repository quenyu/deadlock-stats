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
import { AuthBySteamButton } from '@/features/AuthBySteam'
import { Button } from '@/shared/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/shared/ui/avatar'
import { AppLink } from '@/shared/ui/AppLink/AppLink'
import { routes } from '@/shared/constants/routes'
import { useCurrentUser, useLogout } from '@/shared/lib/react-query/hooks'
import { Skeleton } from '@/shared/ui/skeleton'

export function AuthWidget() {
  const { data: userFromSchema, isLoading } = useCurrentUser()
  const { mutate: logout } = useLogout()

  if (isLoading) {
    return <Skeleton width={112} height={32} variant="rounded" />
  }

  if (!userFromSchema) {
    return <AuthBySteamButton />
  }

  // Normalize for display
  const displayName = userFromSchema.nickname
  const avatarUrl = userFromSchema.avatar_url || ''
  const steamId = userFromSchema.steam_id
  const profileUrl = userFromSchema.profile_url || `https://steamcommunity.com/profiles/${steamId}`

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="lg"
          className="flex items-center gap-2 p-2 h-auto focus-visible:ring-0 focus-visible:ring-offset-0"
        >
          <Avatar className="h-8 w-8">
            <AvatarImage src={avatarUrl} alt={displayName} />
            <AvatarFallback>
              {displayName.charAt(0).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <span className="text-sm font-medium hidden sm:inline-block">{displayName}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56" align="end" forceMount>
        <DropdownMenuLabel className="font-normal">
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">{displayName}</p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem asChild>
            <AppLink to={routes.player.profile(steamId)}>
              <UserIcon className="mr-2 h-4 w-4" />
              <span>Profile</span>
            </AppLink>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <a
              href={profileUrl}
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
        <DropdownMenuItem onSelect={() => logout()}>
          <LogOut className="mr-2 h-4 w-4" />
          <span>Logout</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
} 