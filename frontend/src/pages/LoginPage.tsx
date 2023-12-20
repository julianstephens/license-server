import { AuthForm } from "@/components/forms/AuthForm";
import { PageLayout } from "@/components/layouts/PageLayout";

const LoginPage = () => {
  return (
    <PageLayout>
      <AuthForm loginMode={true} />
    </PageLayout>
  );
};

export default LoginPage;
