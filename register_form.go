package authware

type registerForm struct {
	Id           string `json:"app_id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Token        string `json:"token"`
	EmailAddress string `json:"email_address"`
}
