package authware

import "fmt"

var (
	TamperedCertificate        = fmt.Errorf("server certificate validation failed, tampering with https certificates may have occurred")                                                                                                                                                    // TamperedCertificate is thrown when a TLS certificate cannot be verified, it is only typically returned when the library is out-of-date (which would mean that it would be returned 100% of the time), or a client is tampering with HTTPS responses with an HTTP debugger. Using an HTTP debugger whilst using an Authware application is a security risk as they could potentially bypass authentication entirely by spoofing responses
	BadIdConfiguration         = fmt.Errorf("invalid application configuration, ensure you set the ID and version before calling InitializeApplication")                                                                                                                                    // BadConfiguration is thrown when calling a function with bad data such as a missing ID or version string
	BadHardwareIdConfiguration = fmt.Errorf("invalid application configuration, ensure that you set the HardwareIdentifierFunc to a valid hardware ID fetching function. if you do not want to validate hardware IDs then you need to disable the functionality on the authware dashboard") // BadHardwareIdConfiguration is thrown when the HardwareIdentifierFunc is not set when the application is configured to check user hardware IDs
	AppNotInitialized          = fmt.Errorf("the application must be initialized before calling this function, you can initialize it by calling InitializeApplication")
	AppAlreadyInitialized      = fmt.Errorf("application already initialized") // AppAlreadyInitialized is thrown when InitializeApplication is called when the app is already initialized
)
