import z from "zod";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTrigger,
} from "../../../components/ui/dialog";
import { useState, type ReactNode } from "react";
import { getVaultUserRole, VaultUserRole } from "../../vaults/types";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form } from "../../../components/ui/form";
import { Field, FieldGroup } from "../../../components/ui/field";
import { Input } from "../../../components/ui/input";
import { Combobox } from "../../../components/ui/combobox";
import { Label } from "../../../components/ui/label";
import { Button } from "../../../components/ui/button";
import { api } from "../../../lib/api";
import { useOutletContext } from "react-router-dom";

const shareVaultFormSchema = z.object({
  email: z.email(),
  role: z.enum(VaultUserRole),
  path: z.string(),
});

export const ShareVaultModal = ({
  children,
  defaultRole = VaultUserRole.VIEWER,
  defaultPath = "/",
  open,
  onOpenChange,
}: {
  children?: ReactNode;
  defaultRole?: VaultUserRole;
  defaultPath?: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) => {
  const [comboboxOpen, setComboboxOpen] = useState(false);
  const roleValues = Object.values(VaultUserRole).map((v) => {
    return {
      label: getVaultUserRole(Number(v) as VaultUserRole),
      value: Number(v),
    };
  });

  const vaultId = useOutletContext() as number;

  const form = useForm<z.infer<typeof shareVaultFormSchema>>({
    resolver: zodResolver(shareVaultFormSchema),
    defaultValues: {
      email: '',
      path: defaultPath,
      role: defaultRole,
    },
  });

  const handleShare = async (values: z.infer<typeof shareVaultFormSchema>) => {
    await api.post(`vault/assign-user/${vaultId}`, values)

    onOpenChange(false);
  }

  return (
    <Dialog onOpenChange={onOpenChange} open={open}>
      {children && (
        <DialogTrigger asChild className="">
          {children}
        </DialogTrigger>
      )}

      <DialogContent>
        <DialogHeader>
          <h3>Share vault</h3>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleShare)}>
            <FieldGroup>
              <Controller
                name="email"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field>
                    <Label>Email</Label>
                    <Input autoComplete="off" {...field} value={field.value} />
                  </Field>
                )}
              />
              <Controller
                name="role"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field>
                    <Label>Role</Label>
                    <Combobox
                      items={roleValues}
                      onChange={(v) => field.onChange(v)}
                      value={field.value}
                      open={comboboxOpen}
                      onOpenChange={setComboboxOpen}
                      />
                  </Field>
                )}
              />
              <Controller
                name="path"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field>
                    <Label>Share from path</Label>
                    <Input autoComplete="off" {...field} value={field.value} />
                  </Field>
                )}
              />
            </FieldGroup>
            <DialogFooter>
              <Button>Share</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};
