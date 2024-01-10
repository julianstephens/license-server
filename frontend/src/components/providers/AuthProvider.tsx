import type { ChildrenProps } from "@/types/props";
import { createContext, useEffect, useContext, useState } from "react";
import { SignIn, SignOut, WhoAmI } from "@wails/backend/App";
import type { User } from "@auth0/auth0-react";

export const AuthStatus = {
  LOADING: "LOADING",
  SIGNED_IN: "SIGNED_IN",
  SIGNED_OUT: "SIGNED_OUT",
} as const;

export type IAuth = {
  authStatus: keyof typeof AuthStatus;
  signIn: VoidFunction;
  signOut: VoidFunction;
  profile: User;
};

const defaultState: IAuth = {
  authStatus: AuthStatus.LOADING,
  signIn: () => {},
  signOut: () => {},
  profile: {},
};

export const AuthContext = createContext(defaultState);

const AuthProvider = ({ children }: ChildrenProps) => {
  const [authStatus, setAuthStatus] = useState<keyof typeof AuthStatus>(
    AuthStatus.LOADING
  );

  const [profile, setProfile] = useState<Record<string, string>>({});

  useEffect(() => {
    WhoAmI()
      .then((userData) => {
        console.log("got user", userData);
        if (userData && Object.keys(userData).length > 0) {
          setAuthStatus(AuthStatus.SIGNED_IN);
          setProfile(userData);
        } else {
          setAuthStatus(AuthStatus.SIGNED_OUT);
          setProfile({});
        }
      })
      .catch((err) => {
        console.error("got err", err);
        setAuthStatus(AuthStatus.SIGNED_OUT);
        setProfile({});
      });
  }, [setAuthStatus, authStatus]);

  const signIn = () => {
    console.log("initiating sign in flow");
    SignIn()
      .then((res) => {
        console.log("did sign in", res);
        setAuthStatus(AuthStatus.SIGNED_IN);
      })
      .catch(console.error.bind(console));
  };

  const signOut = () => {
    SignOut()
      .then(() => {
        setAuthStatus(AuthStatus.SIGNED_OUT);
      })
      .catch(console.error.bind(console));
  };

  const state: IAuth = {
    authStatus,
    signIn,
    signOut,
    profile,
  };

  if (authStatus === AuthStatus.LOADING) {
    return null;
  }

  return <AuthContext.Provider value={state}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const { authStatus, signIn, signOut, profile } = useContext(AuthContext);

  return {
    authStatus,
    signIn,
    signOut,
    profile,
  };
};

export default AuthProvider;
