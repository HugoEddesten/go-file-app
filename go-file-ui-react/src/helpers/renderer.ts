type Renderer = "image" | "pdf" | "text" | "video" | "unsupported";

export function getRenderer(mime: string): Renderer {
  if (mime.startsWith("image/")) return "image";
  if (mime === "application/pdf") return "pdf";
  if (mime.startsWith("text/")) return "text";
  if (mime.startsWith("video/")) return "video";
  return "unsupported";
}

export function putPrefixOnByteCount(bytes: number) {
  let iterationCount = 0;
  let result = bytes;

  while (iterationCount < 5) {
    const transformation = bytes / Math.pow(1000, iterationCount);
    if (transformation < 99) {
      result = transformation;
      break;
    }
    iterationCount++;
  }

  switch (iterationCount) {
    case 0:
      return { unit: "B", amount: result, stringResult: `${result.toFixed(2)} B` };
    case 1:
      return { unit: "KB", amount: result, stringResult: `${result.toFixed(2)} KB` };
    case 2:
      return { unit: "MB", amount: result, stringResult: `${result.toFixed(2)} MB` };
    case 3:
      return { unit: "GB", amount: result, stringResult: `${result.toFixed(2)} GB` };
    case 4:
      return { unit: "TB", amount: result, stringResult: `${result.toFixed(2)} TB` };
  }
}
