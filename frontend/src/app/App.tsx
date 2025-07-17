import { RouterProvider } from 'react-router-dom'
import { router } from './providers/router'
import { ThemeProvider } from './providers/theme/ThemeProvider'

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <RouterProvider router={router} />
    </ThemeProvider>
  )
}

export default App 