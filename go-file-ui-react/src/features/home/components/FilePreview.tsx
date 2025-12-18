import { Card } from "../../../components/ui/card";
import { Field, FieldGroup } from "../../../components/ui/field";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner";
import { getRenderer, putPrefixOnByteCount } from "../../../helpers/renderer";
import { useFileDetails } from "../api/getFile";

export const FilePreview = ({ fileKey }: { fileKey: string }) => {
  const fileUrl = `${
    import.meta.env.VITE_API_URL
  }files/download/${encodeURIComponent(fileKey)}`;

  const { data: meta, isLoading } = useFileDetails({ path: fileKey });
  console.log(meta)
  const renderPreview = () => {
    if (!meta) return;

    switch (getRenderer(meta.mimeType)) {
      case "image":
        return <img src={fileUrl} className="object-scale-down max-h-full rounded-md"/>;
      case "pdf":
        return <iframe src={fileUrl} className="h-full rounded-md" />;
      case "text":
        // return <TextPreview path={fileUrl} />;
      case "video":
        return <iframe src={fileUrl} className="h-full rounded-md"/>
      default:
        return <NoPreview />;
    }
  };

  return (
    <Card className="md:col-span-2 p-2 flex">
      {isLoading ? (
        <MaximizedSpinner />
      ) : (
        <div className="flex flex-col gap-4 h-full">
          <FieldGroup className="gap-2">
            <Field>
              <Label>Name</Label>
              <Input value={meta?.name}/>
            </Field>
            <Field>
              <Label>Type</Label>
              <Input readOnly value={meta?.mimeType}/>
            </Field>
            <Field>
              <Label>Size</Label>
              <Input readOnly value={putPrefixOnByteCount(meta?.size ?? 0)?.stringResult}/>
            </Field>
          </FieldGroup>
          {renderPreview()}

        </div>

      )}
    </Card>
  );
};

const NoPreview = () => {
  return (
    <h3 className="text-center">Preview not availible</h3>
  )
}