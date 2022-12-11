package authware

// Api represents an API attached to an application on Authware
type Api struct {
	Id   string `json:"id"`   // Id is the identifier used to execute the API and identify it on Authware
	Name string `json:"name"` // Name is the friendly name used to identify the API in a human-readable way
}
