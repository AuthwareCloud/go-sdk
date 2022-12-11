package authware

// authResponse is a successful authentication response containing an authorization token for a user
type authResponse struct {
	AuthToken string `json:"auth_token"` // AuthToken is the authorization token to grant a user privileges on authenticated endpoints
}
