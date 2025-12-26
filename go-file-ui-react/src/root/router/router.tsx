import { createBrowserRouter, Outlet, RouterProvider } from "react-router-dom"
import { PublicRoute } from "./components/PublicRoute"
import { ProtectedRoute } from "./components/ProtectedRoute"
import { Login } from "../../features/auth/Login"
import { Home } from "../../features/home/Home"
import { Register } from "../../features/auth/Register"
import { Vaults } from "../../features/vaults/Vaults"
import { DefaultLayout } from "../layout/defaultLayout"


export const router = createBrowserRouter([
  {
    element: (
      <DefaultLayout>
        <Outlet />
      </DefaultLayout>
    ),
    children: [
      {
        element: (
          <PublicRoute />
        ),
        children: [
          { path: "/login", element: <Login /> },
          { path: "/register", element: <Register /> }
        ]
      },
      {
        element: <ProtectedRoute />,
        children: [
          { path: "/", element: <Vaults /> },
          { path: "/profile", element: <div>profile</div> },
        ]
      },
      {
        element: <ProtectedRoute requireVaultId={true} />,
        children: [
          { path: "/vault", element: <Home />},
        ]
      }
    ]
  }
])

export const AppRouter  = () => {
  return <RouterProvider router={router}/>
}