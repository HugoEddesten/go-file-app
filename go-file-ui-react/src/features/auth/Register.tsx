import z from "zod";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Card, CardTitle } from "../../components/ui/card";
import { FieldGroup } from "../../components/ui/field";
import { Form, FormDescription } from "../../components/ui/form";
import { Button } from "../../components/ui/button";
import { register } from "./api/register";
import { getInviteInfo, type InviteInfo } from "./api/getInviteInfo";
import { Link, useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import { FormField } from "../../components/form/FormField";

const registerSchema = z.object({
  email: z.string().min(5, { error: "Invalid email" }),
  password: z
    .string()
    .min(8, { error: "Password needs to be at least 8 characters long" }),
});

export const Register = () => {
  const { token } = useParams<{ token?: string }>();
  const [inviteInfo, setInviteInfo] = useState<InviteInfo | null>(null);
  const [inviteError, setInviteError] = useState<string | null>(null);

  const form = useForm<z.infer<typeof registerSchema>>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  useEffect(() => {
    if (!token) return;

    getInviteInfo(token)
      .then((info) => {
        setInviteInfo(info);
        form.setValue("email", info.email);
      })
      .catch(() => {
        setInviteError("This invite link is invalid or has expired.");
      });
  }, [token]);

  useEffect(() => {
    const subscription = form.watch(() => {
      if (form.formState.errors.root) {
        form.clearErrors("root");
      }
    });

    return () => subscription.unsubscribe();
  }, [form]);

  const handleSubmit = async (values: z.infer<typeof registerSchema>) => {
    try {
      await register(values);
    } catch (err: any) {
      const message = err?.response?.data ?? "Something went wrong";

      form.setError("root", {
        type: "server",
        message,
      });
    }
  };

  return (
    <div className="flex w-full h-full items-center justify-center px-4">
      <Card className="p-4 w-full max-w-xl flex justify-center items-center">
        <CardTitle className="text-center text-2xl">
          Register an account
        </CardTitle>

        <Form {...form}>
          {inviteInfo && (
            <FormDescription className="text-center">
              You've been invited to <strong>{inviteInfo.vaultName}</strong>
            </FormDescription>
          )}

          {inviteError && (
            <FormDescription className="text-destructive text-center">
              {inviteError}
            </FormDescription>
          )}
          <form
            className="flex flex-col items-center gap-4 w-full"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FieldGroup>
              <Controller
                control={form.control}
                name="email"
                render={(state) => (
                  <FormField
                    label="Email"
                    {...state}
                    field={{
                      ...state.field,
                      disabled: !!inviteInfo,
                    }}
                  />
                )}
              />
              <Controller
                control={form.control}
                name="password"
                render={(state) => <FormField label="Password" {...state} />}
              />
            </FieldGroup>
            <FormDescription>
              Already have an account?{" "}
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
                Register
              </Button>
            </div>
          </form>
        </Form>
      </Card>
    </div>
  );
};
