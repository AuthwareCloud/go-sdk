package authware

// loginForm is the underlying form for authentication requests
type loginForm struct {
	Id       string `json:"app_id"`   // Id is the ID of the application
	Username string `json:"username"` // Username is the username of the user
	Password string `json:"password"` // Password is the password corresponding to the user identified by their username
}
