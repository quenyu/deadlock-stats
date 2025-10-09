import { createBrowserRouter } from 'react-router-dom'
import { routes } from '@/shared/constants/routes'
import { HomePage } from '@/pages/home'
import { NotFoundPage } from '@/pages/not-found'
import { ErrorPage } from '@/pages/error'
import { PlayerProfilePage } from '@/pages/player-profile/PlayerProfilePage'
import { Layout } from '@/widgets/layout'
import { ProtectedRoute } from '@/shared/lib/ProtectedRoute'
import { PlayerSearchPage } from '@/pages/search/PlayerSearchPage'
import { CrosshairBuilder } from '@/pages/crosshairs/CrosshairBuilder'
import { CrosshairsPage } from '@/pages/crosshairs/CrosshairsPage'
import { CrosshairViewPage } from '@/pages/crosshairs/CrosshairViewPage'

// Temporary stub components until features are implemented
const ComingSoon = () => <div className="p-8 text-center text-muted-foreground">Coming soon...</div>

const rootPath = routes.home;

export const router = createBrowserRouter([
  {
    path: rootPath,
    element: <Layout />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: routes.home,
        element: <HomePage />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.player.profile(),
        element: <PlayerProfilePage />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.player.search,
        element: <PlayerSearchPage />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.builds.list,
        element: <ComingSoon />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.builds.create,
        element: (
          <ProtectedRoute>
            <ComingSoon />
          </ProtectedRoute>
        ),
        errorElement: <ErrorPage />,
      },
      {
        path: routes.builds.view(),
        element: <ComingSoon />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.builds.edit(),
        element: (
          <ProtectedRoute>
            <ComingSoon />
          </ProtectedRoute>
        ),
        errorElement: <ErrorPage />,
      },
      {
        path: routes.crosshairs.list,
        element: <CrosshairsPage />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.crosshairs.create,
        element: <CrosshairBuilder />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.crosshairs.view(),
        element: <CrosshairViewPage />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.analytics,
        element: <ComingSoon />,
        errorElement: <ErrorPage />,
      },
      {
        path: routes.premium,
        element: (
          <ProtectedRoute>
            <ComingSoon />
          </ProtectedRoute>
        ),
        errorElement: <ErrorPage />,
      },
      {
        path: '*',
        element: <NotFoundPage />,
      },
    ],
  },
])