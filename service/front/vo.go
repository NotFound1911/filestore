package main

type SignupReq struct {
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}
type LoginReq struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type DownloadReq struct {
	FileName string `json:"file_name"`
	FileSha1 string `json:"file_sha1"`
}
