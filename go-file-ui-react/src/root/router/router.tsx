import { createBrowserRouter, RouterProvider } from "react-router-dom"
import { PublicRoute } from "./components/PublicRoute"
import { ProtectedRoute } from "./components/ProtectedRoute"
import { Login } from "../../features/auth/Login"
import { Home } from "../../features/home/Home"
import { Register } from "../../features/auth/Register"


export const router = createBrowserRouter([
  {
    element: <PublicRoute />,
    children: [
      { path: "/login", element: <Login /> },
      { path: "/register", element: <Register /> }
    ]
  },
  {
    element: <ProtectedRoute />,
    children: [
      { path: "/", element: <Home /> },
      { path: "/profile", element: <div>profile</div> }
    ]
  }
])

export const AppRouter  = () => {
  return <RouterProvider router={router}/>
}