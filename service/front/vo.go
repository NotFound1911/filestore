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
	Sha1 string `json:"sha1"  form:"sha1"`
	Name string `json:"name" form:"name"`
}
