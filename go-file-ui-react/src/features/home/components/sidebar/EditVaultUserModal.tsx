import { useState } from "react";
import { Controller, useFieldArray, useForm } from "react-hook-form";
import { PlusCircle, Trash2 } from "lucide-react";
import { ConfirmDialog } from "../../../../components/ConfirmDialog";
import { Button } from "../../../../components/ui/button";
import { Combobox } from "../../../../components/ui/combobox";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../../../../components/ui/dialog";
import { Field } from "../../../../components/ui/field";
import { Input } from "../../../../components/ui/input";
import { Label } from "../../../../components/ui/label";
import { useRemoveVaultUser } from "../../../vaults/api/removeVaultUser";
import { useUpdateVaultUser } from "../../../vaults/api/updateVaultUser";
import {
  getVaultUserRole,
  VaultUserRole,
  type VaultUser,
} from "../../../vaults/types";
import { useAddVaultUser } from "../../../vaults/api/addVaultUser";
import { toast } from "sonner";
import type { AxiosError } from "axios";

type PermissionField = {
  vaultUserId: number | null;
  path: string;
  role: VaultUserRole;
};

type FormValues = {
  permissions: PermissionField[];
};

type DeleteTarget = {
  vaultUserId: number;
  index: number;
  path: string;
};

const roleOptions = Object.values(VaultUserRole).map((v) => ({
  label: getVaultUserRole(Number(v) as VaultUserRole)!,
  value: Number(v) as VaultUserRole,
}));

export const EditVaultUserModal = ({
  open,
  onOpenChange,
  permissions,
  email,
  vaultId,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  permissions: VaultUser[];
  email: string;
  vaultId: number;
}) => {
  const [deleteTarget, setDeleteTarget] = useState<DeleteTarget | null>(null);
  const [openComboboxIndex, setOpenComboboxIndex] = useState<number | null>(
    null,
  );

  const { mutateAsync: updateVaultUser } = useUpdateVaultUser({ vaultId });
  const { mutate: removeVaultUser } = useRemoveVaultUser({ vaultId });
  const { mutate: addVaultUser } = useAddVaultUser({ vaultId });

  const form = useForm<FormValues>({
    defaultValues: {
      permissions: permissions.map((p) => ({
        vaultUserId: p.vaultUserId,
        path: p.path,
        role: p.role,
      })),
    },
  });

  const { fields, remove, append } = useFieldArray({
    control: form.control,
    name: "permissions",
  });

  const handleSave = async (values: FormValues) => {
    await Promise.all(
      values.permissions.map((p) => {
        if (p.vaultUserId != null) {
          return updateVaultUser({
            vaultUserId: p.vaultUserId,
            path: p.path,
            role: p.role,
          });
        } else {
          return addVaultUser({
            path: p.path,
            role: p.role,
            email,
          });
        }
      }),
    ).catch((error: AxiosError) => {
      toast.error(`${error.response?.data ??  "Something went wrong"}`)
    }).then(() => {
      onOpenChange(false);
      form.reset();
    });
  };

  const handleDeleteConfirm = () => {
    if (!deleteTarget) return;
    removeVaultUser(deleteTarget.vaultUserId, {
      onSuccess: () => remove(deleteTarget.index),
    });
    setDeleteTarget(null);
  };

  return (
    <>
      <Dialog
        open={open}
        onOpenChange={(o) => {
          onOpenChange(o);
          form.reset();
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="font-semibold">
              Edit permissions
            </DialogTitle>
            <DialogDescription className="text-sm text-muted-foreground">
              {email}
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={form.handleSubmit(handleSave)}>
            <div className="flex flex-col gap-3 overflow-y-auto max-h-[50vh] pr-1">
              {fields.map((field, index) => (
                <div
                  key={field.id}
                  className="group relative bg-muted/40 rounded-lg p-3"
                >
                  <input
                    type="hidden"
                    {...form.register(`permissions.${index}.vaultUserId`)}
                  />
                  <div className="flex gap-3 pr-8">
                    <Field className="flex-1">
                      <Label>Role</Label>
                      <Controller
                        name={`permissions.${index}.role`}
                        control={form.control}
                        render={({ field: f }) => (
                          <Combobox
                            items={roleOptions}
                            value={f.value}
                            onChange={(v) => f.onChange(v)}
                            open={openComboboxIndex === index}
                            onOpenChange={(o) =>
                              setOpenComboboxIndex(o ? index : null)
                            }
                          />
                        )}
                      />
                    </Field>
                    <Field className="flex-1">
                      <Label>Path</Label>
                      <Controller
                        name={`permissions.${index}.path`}
                        control={form.control}
                        render={({ field: f }) => <Input {...f} />}
                      />
                    </Field>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="absolute top-1 right-1 w-7 h-7 opacity-0 group-hover:opacity-100 transition-opacity text-muted-foreground hover:text-destructive"
                    onClick={() => {
                      if (field.vaultUserId != null) {
                        setDeleteTarget({
                          vaultUserId: field.vaultUserId,
                          index,
                          path: form.getValues(`permissions.${index}.path`),
                        });
                      } else {
                        remove(index);
                      }
                    }}
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </div>
              ))}
              <Button
                type="button"
                variant="outline"
                className="w-full border-dashed text-muted-foreground hover:text-foreground gap-2"
                onClick={() => append({ path: "/", role: 4, vaultUserId: null })}
              >
                <PlusCircle className="w-4 h-4" />
                Add permission
              </Button>
            </div>

            <DialogFooter className="mt-4">
              <DialogClose asChild>
                <Button type="button" variant="outline">
                  Cancel
                </Button>
              </DialogClose>
              <Button type="submit">Save</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={deleteTarget !== null}
        onOpenChange={(o) => !o && setDeleteTarget(null)}
        title="Remove permission?"
        description={`Access to "${deleteTarget?.path}" will be removed for ${email}.`}
        onConfirm={handleDeleteConfirm}
      />
    </>
  );
};
