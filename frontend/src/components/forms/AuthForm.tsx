import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "../ui/button";
import { useState } from "react";
import { model } from "@wails/models";
import { Loader } from "../Loader";

export const AuthForm = ({ loginMode }: { loginMode?: boolean }) => {
  const [isLoading, setLoading] = useState(false);
  const [title] = useState(loginMode ? "Sign In" : "Create an Account");
  const [submit] = useState(loginMode ? "Login" : "Register");
  const form = useForm<model.AuthRequest>();

  const onSubmit = (data: model.AuthRequest) => {
    setLoading(true);
    console.log(data);
  };

  return (
    <>
      <h1>{title}</h1>
      {isLoading ? (
        <Loader />
      ) : (
        <Form {...form}>
          <form className="space-y-8" onSubmit={form.handleSubmit(onSubmit)}>
            {!loginMode && (
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="John Doe" {...field} />
                    </FormControl>
                    <FormDescription>
                      This is your public display name.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input
                      type="email"
                      placeholder="john.doe@email.com"
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>This is your email address.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" {...field} />
                  </FormControl>
                  <FormDescription>This is your password.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button className="block mx-auto mt-8" type="submit">
              {submit}
            </Button>
          </form>
        </Form>
      )}
    </>
  );
};
