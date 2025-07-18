import { Button } from '@/shared/ui/button'
import iconSteam from '@/shared/assets/icons/icon-steam.svg'
import { AppLink } from '@/shared/ui/AppLink/AppLink'
import { routes } from '@/shared/constants/routes'

export function AuthBySteamButton() {
  return (
    <Button
      asChild
      variant="secondary"
      size="lg"
      className="group flex items-center gap-2 px-4 transition-transform duration-200 hover:scale-105 focus:scale-105 focus:outline-none"
      aria-label="Sign in with Steam"
    >
      <AppLink to={routes.auth.steam} className="flex items-center gap-2">
        <img src={iconSteam} alt="Steam logo" className="h-5 w-5" />
        <span className="hidden sm:inline-block whitespace-nowrap">Sign&nbsp;in with Steam</span>
      </AppLink>
    </Button>
  )
}