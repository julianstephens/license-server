import { createContext } from "react";

export type IAuthContext = {
  authenticated: boolean;
  setAuthenticated: (isAuthenticated: boolean) => void;
};

export const AuthContext = createContext<IAuthContext>({
  authenticated: false,
  setAuthenticated: () => {},
});
