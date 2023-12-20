import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { PageLayout } from "./components/layouts/PageLayout";

function App() {
  const goto = useNavigate();

  return (
    <PageLayout>
      <div className="text-blue-900 text-2xl font-bold col items-center space-y-8">
        <h1>Welcome to License Server & Manager</h1>
        <div className="row space-x-4">
          <Button
            className="px-10"
            onClick={() => {
              goto("/login");
            }}
          >
            Login
          </Button>
          <Button
            className="px-10"
            onClick={() => {
              goto("/register");
            }}
          >
            Register
          </Button>
        </div>
      </div>
    </PageLayout>
  );
}

export default App;
