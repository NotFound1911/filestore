package v1

type ResumeReq struct {
	FileName string `json:"file_name"`
	FileSha1 string `json:"file_sha1"`
}
