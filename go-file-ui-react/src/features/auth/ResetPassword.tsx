import z from "zod";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Card, CardTitle } from "../../components/ui/card";
import { FieldGroup } from "../../components/ui/field";
import { Form, FormDescription } from "../../components/ui/form";
import { Button } from "../../components/ui/button";
import { resetPassword } from "./api/resetPassword";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useEffect } from "react";
import { FormField } from "../../components/form/FormField";

const resetPasswordSchema = z
  .object({
    password: z
      .string()
      .min(8, { error: "Password needs to be at least 8 characters long" }),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

export const ResetPassword = () => {
  const { token } = useParams<{ token: string }>();
  const navigate = useNavigate();

  const form = useForm<z.infer<typeof resetPasswordSchema>>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      password: "",
      confirmPassword: "",
    },
  });

  useEffect(() => {
    const subscription = form.watch(() => {
      if (form.formState.errors.root) {
        form.clearErrors("root");
      }
    });
    return () => subscription.unsubscribe();
  }, [form]);

  const handleSubmit = async (values: z.infer<typeof resetPasswordSchema>) => {
    try {
      await resetPassword(token!, values.password);
      navigate("/login");
    } catch (err: any) {
      const message = err?.response?.data ?? "Something went wrong";
      form.setError("root", { type: "server", message });
    }
  };

  return (
    <div className="flex w-full h-full items-center justify-center">
      <Card className="p-4 w-xl flex justify-center items-center">
        <CardTitle className="text-center text-2xl">
          Reset your password
        </CardTitle>
        <Form {...form}>
          <form
            className="flex flex-col items-center gap-4 w-lg"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FieldGroup>
              <Controller
                control={form.control}
                name="password"
                render={(state) => (
                  <FormField
                    label="New password"
                    input={{ type: "password", autoComplete: "off" }}
                    {...state}
                  />
                )}
              />
              <Controller
                control={form.control}
                name="confirmPassword"
                render={(state) => (
                  <FormField
                    label="Confirm password"
                    input={{ type: "password", autoComplete: "off" }}
                    {...state}
                  />
                )}
              />
            </FieldGroup>
            <FormDescription>
              Remembered it?{" "}
              <Link to="/login" className="underline">
                Log in here
              </Link>
            </FormDescription>
            <div className="grid grid-rows-2 justify-items-center">
              {form.formState.errors.root && (
                <FormDescription className="text-destructive">
                  {form.formState.errors.root.message}
                </FormDescription>
              )}
              <Button className="row-start-2 w-fit" type="submit">
                Reset password
              </Button>
            </div>
          </form>
        </Form>
      </Card>
    </div>
  );
};
