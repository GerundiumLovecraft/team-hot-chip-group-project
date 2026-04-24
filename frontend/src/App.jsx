import { createBrowserRouter, RouterProvider, Outlet } from "react-router-dom";
import Navbar from "./components/Navbar/Navbar";
import { BrowseSpots } from "./pages/BrowseSpots/BrowseSpots";
import { LoginPage } from "./pages/Login/LoginPage";
import { SignupPage } from "./pages/Signup/SignupPage";
import "./App.css";

// import { SpotDetailsPage } from "./pages/SpotDetailsPage";
// import { NewSpotPage } from "./pages/NewSpotPage";
// import { Profile } from "./pages/Profile";

const Layout = () => {
  return (
    <>
      <Navbar />
      <Outlet />
    </>
  );
}

// What is this? Docs here: https://reactrouter.com/en/main/start/overview
const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      { path: "/", element: <BrowseSpots /> },
      { path: "/login", element: <LoginPage /> },
      { path: "/signup", element: <SignupPage /> },
   // Not yet created below pages
   // { path: "/spots/:id", element: <SpotDetailsPage /> },
   // { path: "/spots/new", element: <NewSpotPage /> },
   // { path: "/profile", element: <Profile /> },
    ],
  },
]);

const App = () => {
  return <RouterProvider router={router} />;
};

export default App;
