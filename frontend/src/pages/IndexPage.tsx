import { PageLayout } from "@/components/layouts/PageLayout";
import { Button } from "@/components/ui/button";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { AuthStatus, useAuth } from "@/components/providers/AuthProvider";

const IndexPage = () => {
  const goto = useNavigate();
  const { authStatus, signIn } = useAuth();

  useEffect(() => {
    console.log("auth status", authStatus);
    if (authStatus === AuthStatus.SIGNED_IN) {
      goto("/dashboard");
    }
  }, [authStatus]);

  return (
    <PageLayout>
      <div className="col items-center space-y-8">
        <h1>Welcome to License Server & Manager</h1>
        <div className="row space-x-4">
          <Button
            className="px-10"
            onClick={() => {
              signIn();
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
};

export default IndexPage;
