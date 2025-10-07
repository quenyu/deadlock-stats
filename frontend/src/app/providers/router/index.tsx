import { createBrowserRouter } from 'react-router-dom'
import { routes } from '@/shared/constants/routes'
import { HomePage } from '@/pages/home'
import { NotFoundPage } from '@/pages/not-found'
import { PlayerProfilePage } from '@/pages/player-profile/PlayerProfilePage'
import { Layout } from '@/widgets/layout'
import { ProtectedRoute } from '@/shared/lib/ProtectedRoute'
import { PlayerSearchPage } from '@/pages/search/PlayerSearchPage'
import { CrosshairBuilder } from '@/pages/crosshairs/CrosshairBuilder'
import { CrosshairsPage } from '@/pages/crosshairs/CrosshairsPage'

// Temporary stub components until features are implemented
const ComingSoon = () => <div className="p-8 text-center text-muted-foreground">Coming soon...</div>

const rootPath = routes.home;

export const router = createBrowserRouter([
  {
    path: rootPath,
    element: <Layout />,
    errorElement: <NotFoundPage />,
    children: [
      {
        path: routes.home,
        element: <HomePage />,
      },
      {
        path: routes.player.profile(),
        element: <PlayerProfilePage />,
      },
      {
        path: routes.player.search,
        element: <PlayerSearchPage />,
      },
      {
        path: routes.builds.list,
        element: <ComingSoon />,
      },
      {
        path: routes.builds.create,
        element: (
          <ProtectedRoute>
            <ComingSoon />
          </ProtectedRoute>
        ),
      },
      {
        path: routes.builds.view(),
        element: <ComingSoon />,
      },
      {
        path: routes.builds.edit(),
        element: (
          <ProtectedRoute>
            <ComingSoon />
          </ProtectedRoute>
        ),
      },
      {
        path: routes.crosshairs.list,
        element: <CrosshairsPage />,
      },
      {
        path: routes.crosshairs.create,
        element: <CrosshairBuilder />,
      },
      {
        path: routes.analytics,
        element: <ComingSoon />,
      },
      {
        path: routes.premium,
        element: (
          <ProtectedRoute>
            <ComingSoon />
          </ProtectedRoute>
        ),
      },
    ],
  },
])