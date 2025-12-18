import type {
  ControllerFieldState,
  ControllerRenderProps,
  FieldValues,
  Path,
  UseFormStateReturn,
} from "react-hook-form";
import { Field, FieldError, FieldLabel } from "../ui/field";
import { Input } from "../ui/input";

type FormFieldProps<
  TFieldValues extends FieldValues,
  TName extends Path<TFieldValues>
> = {
  label?: string;
  input?: React.InputHTMLAttributes<HTMLInputElement>;

  field: ControllerRenderProps<TFieldValues, TName>;
  fieldState: ControllerFieldState;
  formState: UseFormStateReturn<TFieldValues>;
};

export function FormField<
  TFieldValues extends FieldValues,
  TName extends Path<TFieldValues>
>({
  field,
  fieldState,
  label,
}: FormFieldProps<TFieldValues, TName>) {
  // const { field, fieldState, formState } = state;

  return (
    <Field data-invalid={fieldState.invalid}>
      {label && <FieldLabel htmlFor={field.name}>{label}</FieldLabel>}
      <Input
        id={field.name}
        type={field.name}
        aria-invalid={fieldState.invalid}
        {...field}
      />
      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
    </Field>
  );
}
