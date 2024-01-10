import AuthProvider from "@/components/providers/AuthProvider";
import { HashRouter, Route, Routes } from "react-router-dom";
import DashboardPage from "./pages/Dashboard";
import IndexPage from "./pages/IndexPage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";

function App() {
  return (
    <AuthProvider>
      <HashRouter>
        <Routes>
          <Route path="/" Component={IndexPage} />
          <Route path="/register" Component={RegisterPage} />
          <Route path="/login" Component={LoginPage} />
          <Route path="/dashboard" Component={DashboardPage} />
        </Routes>
      </HashRouter>
    </AuthProvider>
  );
}

export default App;
