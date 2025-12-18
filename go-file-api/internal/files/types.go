package files

type FileResponse struct {
	Name string
	Key  string
}

type FileMetadata struct {
	Name        string `json:"name"`
	MimeType    string `json:"mimeType"`
	Size        int64  `json:"size"`
	Editable    bool   `json:"editable"`
	Previewable bool   `json:"previewable"`
}
