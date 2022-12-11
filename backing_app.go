package authware

import "time"

// backingApp is the backend app type that is used to deserialize the app details to
type backingApp struct {
	Name                        string    `json:"name"`
	Id                          string    `json:"id"`
	Version                     string    `json:"version"`
	DateCreated                 time.Time `json:"date_created"`
	IsHardwareIdCheckingEnabled bool      `json:"is_hwid_checking_enabled"`
	Apis                        []Api     `json:"apis,omitempty"`
	UserCount                   int       `json:"user_count,omitempty"`
	RequestCount                int       `json:"request_count,omitempty"`
}
