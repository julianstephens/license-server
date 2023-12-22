import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { PageLayout } from "./components/layouts/PageLayout";
import { Auth } from "../wailsjs/go/backend/App";

function App() {
  const goto = useNavigate();

  return (
    <PageLayout>
      <div className="col items-center space-y-8">
        <h1>Welcome to License Server & Manager</h1>
        <div className="row space-x-4">
          <Button
            className="px-10"
            onClick={() => {
              Auth().catch(console.error.bind(console));
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
