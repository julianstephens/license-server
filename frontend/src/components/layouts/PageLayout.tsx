import React from "react";

export const PageLayout = ({ children }: { children: React.ReactNode }) => {
  return <div className="col full centered">{children}</div>;
};
