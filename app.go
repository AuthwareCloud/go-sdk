package authware

import (
	"net/http"
	"time"
)

type Application struct {
	AuthToken                   string        // AuthToken is the current users authorization token
	Version                     string        // Version is the current local version of the application
	Id                          string        // Id is the identifier for your application as appearing on Authware
	Name                        *string       // Name is the name of the application as it appears on Authware
	DateCreated                 *time.Time    // DateCreated is the date and time the application was created
	IsHardwareIdCheckingEnabled *bool         // IsHardwareIdCheckingEnabled is a flag whether a users hardware ID is checked and enforced upon login and any request
	Apis                        *[]Api        // Apis is a list of APIs that are registered on the application
	UserCount                   *int          // UserCount is the amount of users registered on the application
	RequestCount                *int          // RequestCount is the amount of proxied API requests that have taken place on the application, this is not Authware API requests but proxied application API requests
	HardwareIdentifierFunc      func() string // HardwareIdentifierFunc is a function that is called to gather the users hardware identifier, by default this is not set and hardware IDs are ignored at the client level
	initialized                 bool          // initialized is a flag whether the application has been initialized properly or not
}

// InitializeApplication will gather the detailed information about the app and get it ready for usage with Authware, it is mandatory that you call this method before you use any Authware API call
func (a *Application) InitializeApplication() error {
	if a.initialized {
		return AppAlreadyInitialized
	}

	// Ensure that the app details are valid
	if a.Id == "" || a.Version == "" {
		return BadIdConfiguration
	}

	// set http backing app
	addApplication(a)

	// Make the request to get app info
	app, err := doRequest[backingApp](http.MethodPost, "app", initForm{Id: a.Id})
	if err != nil {
		return err
	}

	// Ensure that the HardwareIdentifierFunc is set if the app is set up to check hwids
	if app.IsHardwareIdCheckingEnabled && a.HardwareIdentifierFunc == nil {
		return BadHardwareIdConfiguration
	}

	a.Apis = &app.Apis
	a.Name = &app.Name
	a.DateCreated = &app.DateCreated
	a.IsHardwareIdCheckingEnabled = &app.IsHardwareIdCheckingEnabled
	a.UserCount = &app.UserCount
	a.RequestCount = &app.RequestCount
	a.initialized = true

	// reset after done
	addApplication(a)

	return nil
}

func (a *Application) Authenticate(username string, password string) error {
	// Ensure the app is initialized
	if !a.initialized {
		return AppNotInitialized
	}

	// Make the request to auth the user
	auth, err := doRequest[authResponse](http.MethodPost, "user/auth", loginForm{
		Id:       a.Id,
		Username: username,
		Password: password,
	})

	if err != nil {
		return err
	}

	addAuthenticationToken(auth.AuthToken)
	return nil
}

func (a *Application) Register(username string, password string, email string, token string) error {
	// Ensure the app is initialized
	if !a.initialized {
		return AppNotInitialized
	}

	// Make the request to auth the user
	_, err := doRequest[DefaultResponse](http.MethodPost, "user/register", registerForm{
		Id:           a.Id,
		Username:     username,
		Password:     password,
		EmailAddress: email,
		Token:        token,
	})

	if err != nil {
		return err
	}

	return nil
}
