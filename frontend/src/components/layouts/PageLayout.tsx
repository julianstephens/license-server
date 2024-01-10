import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { AuthStatus, useAuth } from "../providers/AuthProvider";

export const PageLayout = ({ children }: { children: React.ReactNode }) => {
  const { authStatus } = useAuth();
  const goto = useNavigate();

  useEffect(() => {
    if (authStatus === AuthStatus.SIGNED_OUT) {
      goto("/");
    }
  }, [authStatus]);
  return <div className="col full centered">{children}</div>;
};
