package authware

// DefaultResponse is the default response format returned by the Authware API, this is globally used across many endpoints and internally used to parse errors
type DefaultResponse struct {
	Code    int      `json:"code"`    // Code is the status code returned by the API, this is unrelated to the HTTP status code and contains more specific response data
	Message string   `json:"message"` // Message is the message returned by the API and normally contains human-readable information about the response
	Errors  []string `json:"errors"`  // Errors is the list of human-readable errors caused by the request, typically used for validation issues with user inputs
}
