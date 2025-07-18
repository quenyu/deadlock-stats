import { Menu } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  NavigationMenuContent,
} from '@/shared/ui/navigation-menu'
import { Sheet, SheetContent, SheetTrigger, SheetTitle, SheetDescription } from '@/shared/ui/sheet'
import * as React from 'react'
import { routes } from '@/shared/constants/routes'

const navLinks = [
  { href: routes.home, label: 'Home' },
  {
    label: 'Players',
    subItems: [
      { title: 'Search', href: routes.player.search, description: 'Search for players by Steam ID or nickname' },
      { title: 'Profile (example)', href: routes.player.profile('123'), description: 'Example player profile (replace id dynamically)' },
    ],
  },
  {
    label: 'Builds',
    subItems: [
      { title: 'All Builds', href: routes.builds.list, description: 'Browse all published builds' },
      { title: 'Create Build', href: routes.builds.create, description: 'Publish a new build' },
    ],
  },
  { href: routes.analytics, label: 'Analytics' },
  { href: routes.premium, label: 'Premium' },
]

export function Navbar() {
  return (
    <>
      {/* Desktop Menu */}
      <div className="hidden md:flex md:items-center md:space-x-4">
        <NavigationMenu viewport={false}>
          <NavigationMenuList>
            {navLinks.map((link) =>
              link.subItems ? (
                <NavigationMenuItem key={link.label}>
                  <NavigationMenuTrigger>{link.label}</NavigationMenuTrigger>
                  <NavigationMenuContent className="bg-[hsl(var(--popover))] text-[hsl(var(--popover-foreground))]">
                    <ul className="grid w-[400px] gap-3 p-4 md:w-[500px] md:grid-cols-2">
                      {link.subItems.map((sub) => (
                        <ListItem key={sub.title} title={sub.title} href={sub.href}>
                          {sub.description}
                        </ListItem>
                      ))}
                    </ul>
                  </NavigationMenuContent>
                </NavigationMenuItem>
              ) : (
                <NavigationMenuItem key={link.href}>
                  <NavigationMenuLink 
                    href={link.href}
                    className="group inline-flex h-9 w-max items-center justify-center rounded-md bg-background px-4 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus:outline-none disabled:pointer-events-none disabled:opacity-50 data-[active]:bg-accent/50 data-[state=open]:bg-accent/50"
                  >
                    {link.label}
                  </NavigationMenuLink>
                </NavigationMenuItem>
              )
            )}
          </NavigationMenuList>
        </NavigationMenu>
      </div>

      {/* Mobile Menu (Burger) */}
      <div className="md:hidden">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon">
              <Menu className="h-6 w-6" />
              <span className="sr-only">Toggle Menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-[240px] p-6 bg-[hsl(var(--popover))] z-[100]">
            <SheetTitle className="sr-only">Main Menu</SheetTitle>
            <SheetDescription className="sr-only">Navigation menu for mobile</SheetDescription>
            <a href="/" className="mb-6 flex items-center space-x-2">
              <span className="font-bold">Deadlock Stats</span>
            </a>
            <nav className="flex flex-col space-y-2">
              {navLinks.map((link) =>
                link.subItems ? (
                  <div key={link.label}>
                    <span className="px-3 py-2 text-base font-medium text-foreground">{link.label}</span>
                    <div className="mt-1 flex flex-col space-y-1">
                      {link.subItems.map((sub) => (
                        <a
                          key={sub.href}
                          href={sub.href}
                          className="ml-4 block rounded-md px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                        >
                          {sub.title}
                        </a>
                      ))}
                    </div>
                  </div>
                ) : (
                  <a
                    key={link.href}
                    href={link.href}
                    className="rounded-md px-3 py-2 text-base font-medium text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                  >
                    {link.label}
                  </a>
                )
              )}
            </nav>
          </SheetContent>
        </Sheet>
      </div>
    </>
  )
}

function ListItem({
  title,
  children,
  href,
  ...props
}: React.ComponentPropsWithoutRef<"li"> & { href: string }) {
  return (
    <li {...props}>
      <NavigationMenuLink asChild>
        <a 
          href={href}
          className="block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"
        >
          <div className="text-sm font-medium leading-none">{title}</div>
          <p className="text-muted-foreground line-clamp-2 text-sm leading-snug mt-1">
            {children}
          </p>
        </a>
      </NavigationMenuLink>
    </li>
  )
}
  