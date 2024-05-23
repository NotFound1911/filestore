package v1

type ResumeReq struct {
	FileName string `json:"file_name"`
	FileSha1 string `json:"file_sha1"`
}

type InitMultiUploadFileReq struct {
	FileName string `json:"file_name"`
	FileSha1 string `json:"file_sha1"`
	FileSize int64  `json:"file_size"`
}

type MultiUploadFileMergeReq struct {
	UploadId int64  `json:"upload_id"`
	FileName string `json:"file_name"`
	FileSha1 string `json:"file_sha1"`
	FileSize int64  `json:"file_size"`
}
