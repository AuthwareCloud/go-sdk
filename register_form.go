package authware

// registerForm is the underlying form for registering/creating new users in an application
type registerForm struct {
	Id           string `json:"app_id"`        // Id is the ID of the application
	Username     string `json:"username"`      // Username is the name to identify the new user by
	Password     string `json:"password"`      // Password is the text to allow the user to sign in to their account
	Token        string `json:"token"`         // Token is the license key/token used to grant the user time or a role on your application
	EmailAddress string `json:"email_address"` // EmailAddress is the contact address to message the user about notifications, password resets and other information
}
