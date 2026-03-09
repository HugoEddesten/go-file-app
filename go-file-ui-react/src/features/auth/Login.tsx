import z from "zod";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Card, CardTitle } from "../../components/ui/card";
import { FieldGroup } from "../../components/ui/field";
import { Form, FormDescription } from "../../components/ui/form";
import { Button } from "../../components/ui/button";
import { login } from "./api/login";
import { sendResetPasswordEmail } from "./api/sendResetPasswordEmail";
import { Link } from "react-router-dom";
import { useState } from "react";
import { FormField } from "../../components/form/FormField";

const loginSchema = z.object({
  email: z.string().min(1, { error: "Required" }),
  password: z.string().min(1, { error: "Required" }),
});

export const Login = () => {
  const [resetSent, setResetSent] = useState(false);

  const form = useForm<z.infer<typeof loginSchema>>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const handleSubmit = async (values: z.infer<typeof loginSchema>) => {
    try {
      await login(values);
    } catch (err: any) {
      const message = err?.response?.data ?? "Something went wrong";
      form.setError("root", { type: "server", message });
    }
  };

  const handleForgotPassword = async () => {
    const email = form.getValues("email");
    if (!email) {
      form.clearErrors();
      form.setError("email", { type: "manual", message: "Enter your email first" });
      return;
    }

    try {
      await sendResetPasswordEmail(email);
      setResetSent(true);
    } catch {
      form.setError("root", { type: "server", message: "Something went wrong" });
    }
  };

  return (
    <div className="flex w-full h-full items-center justify-center">
      <Card className="p-4 w-xl flex items-center justify-center">
        <CardTitle className="text-center text-2xl">
          Log in to you account
        </CardTitle>
        <Form {...form}>
          <form
            className="flex flex-col items-center gap-4 w-lg"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FieldGroup>
              <Controller
                control={form.control}
                name="email"
                render={(state) => <FormField label="Email" {...state} />}
              />
              <Controller
                control={form.control}
                name="password"
                render={(state) => (
                  <FormField
                    label="Password"
                    input={{ type: "password", autoComplete: "off" }}
                    {...state}
                  />
                )}
              />
            </FieldGroup>
            <FormDescription>
              Dont have an account?{" "}
              <Link to="/register" className="underline">
                Register here
              </Link>
            </FormDescription>
            <div className="grid grid-rows-2 w-full justify-items-center">
              {form.formState.errors.root && (
                <FormDescription className="text-destructive">
                  {form.formState.errors.root.message}
                </FormDescription>
              )}
              {resetSent && (
                <FormDescription className="text-center">
                  Password reset email sent. Check your inbox.
                </FormDescription>
              )}
              <div className="row-start-2 w-full justify-end flex gap-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={handleForgotPassword}
                >
                  Forgot password?
                </Button>
                <Button type="submit">Log in</Button>
              </div>
            </div>
          </form>
        </Form>
      </Card>
    </div>
  );
};
