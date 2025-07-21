import {
  Menu,
  Home,
  Users,
  Swords,
  BarChart,
  User,
  Search,
  PlusCircle,
  Star,
  Crosshair,
} from 'lucide-react'
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  NavigationMenuContent,
} from '@/shared/ui/navigation-menu'
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetTitle,
  SheetDescription,
} from '@/shared/ui/sheet'
import { Button } from '@/shared/ui/button'
import * as React from 'react'
import { routes } from '@/shared/constants/routes'
import { AppLink } from '@/shared/ui/AppLink/AppLink'
import { useUserStore } from '@/entities/user'

const navConfig = (isAuthenticated: boolean, steamId?: string) => [
  {
    href: routes.home,
    label: 'Home',
    icon: Home,
  },
  {
    label: 'Players',
    icon: Users,
    subItems: [
      {
        title: 'Search',
        href: routes.player.search,
        description: 'Find a player by Steam ID or nickname',
        icon: Search,
      },
      ...(isAuthenticated && steamId
        ? [
            {
              title: 'My Profile',
              href: routes.player.profile(steamId),
              description: 'View your stats and achievements',
              icon: User,
            },
          ]
        : []),
    ],
  },
  {
    label: 'Crosshairs',
    icon: Crosshair,
    subItems: [
      {
        title: 'All Crosshairs',
        href: routes.crosshairs.list,
        description: 'Browse all published crosshairs',
        icon: Crosshair,
      },
      ...(isAuthenticated
        ? [
            {
              title: 'Create Crosshair',
              href: routes.crosshairs.create,
              description: 'Publish a new crosshair for the community',
              icon: PlusCircle,
            },
          ]
        : []),
    ],
  },
  {
    label: 'Builds',
    icon: Swords,
    subItems: [
      {
        title: 'All Builds',
        href: routes.builds.list,
        description: 'Browse all published builds',
        icon: Swords,
      },
      ...(isAuthenticated
        ? [
            {
              title: 'Create Build',
              href: routes.builds.create,
              description: 'Publish a new build for the community',
              icon: PlusCircle,
            },
          ]
        : []),
    ],
  },
  {
    href: routes.analytics,
    label: 'Analytics',
    icon: BarChart,
  },
  {
    href: routes.premium,
    label: 'Premium',
    icon: Star,
  },
]

export function Navbar() {
  const { user, isLoading } = useUserStore()

  if (isLoading) {
    return (
      <>
        <div className="hidden h-10 w-full animate-pulse rounded-md bg-muted md:block" />
        <div className="h-10 w-10 animate-pulse rounded-md bg-muted md:hidden" />
      </>
    )
  }
  
  const navigation = navConfig(!!user, user ? user.steam_id : undefined)

  return (
    <>
      <DesktopNav navigation={navigation} />
      <MobileNav navigation={navigation} />
    </>
  )
}

function DesktopNav({ navigation }: { navigation: ReturnType<typeof navConfig> }) {
  return (
    <div className="hidden md:flex md:items-center md:space-x-4">
      <NavigationMenu>
        <NavigationMenuList>
          {navigation.map((link) =>
            link.subItems ? (
              <NavigationMenuItem key={link.label}>
                <NavigationMenuTrigger className="text-sm font-medium">
                  <link.icon className="mr-2 h-4 w-4" />
                  {link.label}
                </NavigationMenuTrigger>
                <NavigationMenuContent>
                  <ul className="grid w-[400px] gap-3 p-4 md:w-[500px] md:grid-cols-2">
                    {link.subItems.map((sub) => (
                      <ListItem
                        key={sub.title}
                        title={sub.title}
                        to={sub.href}
                        icon={sub.icon}
                      >
                        {sub.description}
                      </ListItem>
                    ))}
                  </ul>
                </NavigationMenuContent>
              </NavigationMenuItem>
            ) : (
              <NavigationMenuItem key={link.href}>
                <AppLink
                  to={link.href}
                  className="group inline-flex h-10 w-max items-center justify-center rounded-md bg-transparent px-4 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus:outline-none disabled:pointer-events-none disabled:opacity-50"
                >
                  <link.icon className="mr-2 h-4 w-4" />
                  {link.label}
                </AppLink>
              </NavigationMenuItem>
            )
          )}
        </NavigationMenuList>
      </NavigationMenu>
    </div>
  )
}

function MobileNav({ navigation }: { navigation: ReturnType<typeof navConfig> }) {
  return (
    <div className="md:hidden">
      <Sheet>
        <SheetTrigger asChild>
          <Button variant="ghost" size="icon">
            <Menu className="h-6 w-6" />
            <span className="sr-only">Toggle Menu</span>
          </Button>
        </SheetTrigger>
        <SheetContent side="left" className="w-[260px] p-4 z-[100]">
          <SheetTitle className="sr-only">Main Menu</SheetTitle>
          <SheetDescription className="sr-only">
            Navigation menu for mobile devices
          </SheetDescription>
          <AppLink to="/" className="mb-6 flex items-center space-x-2">
            <span className="text-lg font-bold">Deadlock Stats</span>
          </AppLink>
          <nav className="flex flex-col space-y-2">
            {navigation.map((link) =>
              link.subItems ? (
                <div key={link.label}>
                  <h4 className="px-3 py-2 text-sm font-semibold text-muted-foreground flex items-center">
                    <link.icon className="mr-2 h-4 w-4" />
                    {link.label}
                  </h4>
                  <div className="mt-1 flex flex-col space-y-1">
                    {link.subItems.map((sub) => (
                      <AppLink
                        key={sub.href}
                        to={sub.href}
                        className="ml-4 flex items-center gap-2 rounded-md px-3 py-2 text-base font-medium text-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                      >
                        <sub.icon className="h-4 w-4 text-muted-foreground" />
                        {sub.title}
                      </AppLink>
                    ))}
                  </div>
                </div>
              ) : (
                <AppLink
                  key={link.href}
                  to={link.href}
                  className="flex items-center gap-2 rounded-md px-3 py-2 text-base font-medium text-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                >
                  <link.icon className="h-4 w-4" />
                  {link.label}
                </AppLink>
              )
            )}
          </nav>
        </SheetContent>
      </Sheet>
    </div>
  )
}

const ListItem = React.forwardRef<
  HTMLAnchorElement,
  React.ComponentPropsWithoutRef<typeof AppLink> & { icon: React.ElementType }
>(({ className, title, children, icon: Icon, ...props }, ref) => {
  return (
    <li>
      <NavigationMenuLink asChild>
        <AppLink
          ref={ref}
          className={
            'block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground'
          }
          {...props}
        >
          <div className="flex items-center gap-2 text-sm font-medium leading-none">
            <Icon className="h-4 w-4" />
            {title}
          </div>
          <p className="line-clamp-2 text-sm leading-snug text-muted-foreground pl-6">
            {children}
          </p>
        </AppLink>
      </NavigationMenuLink>
    </li>
  )
})
ListItem.displayName = 'ListItem'
  