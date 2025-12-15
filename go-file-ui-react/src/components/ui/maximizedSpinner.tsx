import { Spinner } from "./spinner"

export const MaximizedSpinner = () => {
  return (
    <div className="flex justify-center items-center w-full h-full">
      <Spinner className="w-16 h-16"/>
    </div>
  );
}