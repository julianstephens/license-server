import { PageLayout } from "@/components/layouts/PageLayout";
import { useAuth } from "@/components/providers/AuthProvider";

const DashboardPage = () => {
  const { profile, signOut } = useAuth();

  return (
    <PageLayout>
      <h1>Logged in as: {profile.given_name}</h1>
      <button className="button" onClick={() => signOut()}>
        Sign Out
      </button>
    </PageLayout>
  );
};
export default DashboardPage;
