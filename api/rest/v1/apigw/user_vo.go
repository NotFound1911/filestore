package apigw

type SignupReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
type LoginJWTReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
