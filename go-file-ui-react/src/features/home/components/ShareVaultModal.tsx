import z from "zod";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
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
import { useParams } from "react-router-dom";

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

  const { vaultId: vaultIdParam } = useParams<{ vaultId: string }>();
  const vaultId = Number(vaultIdParam);

  const form = useForm<z.infer<typeof shareVaultFormSchema>>({
    resolver: zodResolver(shareVaultFormSchema),
    defaultValues: {
      email: "",
      path: defaultPath,
      role: defaultRole,
    },
  });

  const handleShare = async (values: z.infer<typeof shareVaultFormSchema>) => {
    await api.post(`vault/assign-user/${vaultId}`, values);
    form.reset();
    onOpenChange(false);
  };

  return (
    <Dialog
      onOpenChange={(open) => {
        form.reset();
        onOpenChange(open);
      }}
      open={open}
    >
      {children && (
        <DialogTrigger asChild className="">
          {children}
        </DialogTrigger>
      )}

      <DialogContent>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleShare)}
            className="flex flex-col gap-4"
          >
            <DialogHeader>
              <DialogTitle>Share folder</DialogTitle>
              <DialogDescription>
                Give a user access to this folders and its sub-folders
              </DialogDescription>
            </DialogHeader>
            <FieldGroup>
              <Controller
                name="email"
                control={form.control}
                render={({ field }) => (
                  <Field>
                    <Label>Email</Label>
                    <Input autoComplete="off" {...field} value={field.value} />
                  </Field>
                )}
              />
              <Controller
                name="role"
                control={form.control}
                render={({ field }) => (
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
                render={({ field }) => (
                  <Field>
                    <Label>Share from path</Label>
                    <Input autoComplete="off" {...field} value={field.value} />
                  </Field>
                )}
              />
            </FieldGroup>
            <DialogFooter>
              <DialogClose asChild>
                <Button variant={"outline"}>Close</Button>
              </DialogClose>
              <Button>Share</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};
