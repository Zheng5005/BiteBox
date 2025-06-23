import Navbar from "./components/Navbar";
import Login from "./pages/Login";
import MainPage from "./pages/Main";
import SignUp from "./pages/SignUp";

const App: React.FC = () => {
  return <>
    <Navbar />
    <MainPage />
    <Login />
    <SignUp />
  </>
};

export default App;
