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
    label: 'Главная',
    icon: Home,
  },
  {
    label: 'Игроки',
    icon: Users,
    subItems: [
      {
        title: 'Поиск',
        href: routes.player.search,
        description: 'Найти игрока по Steam ID или никнейму',
        icon: Search,
      },
      ...(isAuthenticated && steamId
        ? [
            {
              title: 'Мой профиль',
              href: routes.player.profile(steamId),
              description: 'Посмотреть свою статистику и достижения',
              icon: User,
            },
          ]
        : []),
    ],
  },
  {
    label: 'Билды',
    icon: Swords,
    subItems: [
      {
        title: 'Все билды',
        href: routes.builds.list,
        description: 'Просмотреть все опубликованные билды',
        icon: Swords,
      },
      ...(isAuthenticated
        ? [
            {
              title: 'Создать билд',
              href: routes.builds.create,
              description: 'Опубликовать новый билд для сообщества',
              icon: PlusCircle,
            },
          ]
        : []),
    ],
  },
  {
    href: routes.analytics,
    label: 'Аналитика',
    icon: BarChart,
  },
  {
    href: routes.premium,
    label: 'Премиум',
    icon: Star,
  },
]

export function Navbar() {
  const { user } = useUserStore()
  const navigation = navConfig(!!user, user?.steam_id)

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
          <SheetTitle className="sr-only">Главное меню</SheetTitle>
          <SheetDescription className="sr-only">
            Навигационное меню для мобильных устройств
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
  