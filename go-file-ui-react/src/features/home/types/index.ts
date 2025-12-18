export type FileRenderer =
  | "image"
  | "pdf"
  | "text"
  | "video"
  | "unsupported";

export type FileMetadata = {
  name: string;
  mimeType: string;
  size: number;
  editable: boolean;
  previewable: boolean;
  renderer?: FileRenderer;
};