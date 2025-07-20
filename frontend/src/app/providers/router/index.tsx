import { createBrowserRouter } from 'react-router-dom'
import { Layout } from '@/widgets/layout'
import { HomePage } from '@/pages/home'
import { NotFoundPage } from '@/pages/not-found'
import { routes } from '@/shared/constants/routes'
import { ProtectedRoute } from '@/shared/lib/ProtectedRoute'

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
        index: true,
        element: <HomePage />,
      },
      {
        path: routes.player.profile(),
        element: <ComingSoon />,
      },
      {
        path: routes.player.search,
        element: <ComingSoon />,
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