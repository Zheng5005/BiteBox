import { Route, Routes } from "react-router";
import Navbar from "./components/Navbar";
import Login from "./pages/Login";
import MainPage from "./pages/Main";
import SignUp from "./pages/SignUp";
import RecipeDetails from "./pages/RecipeDetails";
import { AuthProvider } from "./context/AuthContext";
import PostRecipe from "./pages/PostRecipe";

const App: React.FC = () => {
  return <>
    <AuthProvider>
      <Navbar />
      <Routes >
        {/* Global */}
        <Route path="/" element={<MainPage />} />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/details/:id" element={<RecipeDetails />} />
        <Route path="/post" element={<PostRecipe />} />

        {/* Auth Users */}
      </Routes>
    </AuthProvider>
  </>
};

export default App;
