package authware

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	application       *Application                  // application will contain the current application and information
	client            = newClient()                 // client is the HTTP request client to send HTTP requests
	baseAddress       = "https://api.authware.org/" // baseAddress is the address that is prepended to every request URL passed into doRequest
	cloudflareIssuer  = []byte("Cloudflare Inc")    // cloudflareIssuer is the bytes for a valid issuer on a Cloudflare certificate
	letsEncryptIssuer = []byte("Let's Encrypt")     // letsEncryptIssuer is the bytes for a valid issuer on a Let's Encrypt certificate
)

// newClient will create a new Go HTTP client to make requests to the Authware API
func newClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			TLSClientConfig: &tls.Config{
				VerifyConnection: func(state tls.ConnectionState) error {
					// Check if the first peer certificate issuers are matching our valid issuers
					if !bytes.Contains(state.PeerCertificates[0].RawIssuer, cloudflareIssuer) &&
						!bytes.Contains(state.PeerCertificates[0].RawIssuer, letsEncryptIssuer) {
						// Both are not matching, we return the tamperedCertificate error
						return TamperedCertificate
					}

					return nil
				},
			},
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
}

// addAuthenticationToken will add an authentication token to the backing store for the HTTP client
func addAuthenticationToken(token string) {
	application.AuthToken = token
}

func addApplication(app *Application) {
	application = app
}

// newRequest is a proxy method for the http.NewRequest method that simply adds the Authware specific headers automatically
func newRequest(method string, url string, body []byte) (req *http.Request, err error) {
	// Make a reader for the body
	bodyReader := bytes.NewReader(body)

	// Create a new HTTP request and prepend the base address to the URL
	rawReq, err := http.NewRequest(method, baseAddress+url, bodyReader)
	if err != nil {
		return
	}

	// Set authorization only if a valid authToken is set
	if application.AuthToken != "" {
		rawReq.Header.Set("Authorization", application.AuthToken)
	}

	// Set the users hardware ID in the header if the func is defined
	if application.HardwareIdentifierFunc != nil {
		rawReq.Header.Set("X-Authware-Hardware-ID", application.HardwareIdentifierFunc())
	}

	rawReq.Header.Set("X-Authware-App-Version", application.Version)
	rawReq.Header.Set("User-Agent", "Authware-Gopher/1.0.0")
	return rawReq, nil
}

// doRequest will take in the URL, request method and body, create a request and send it to the Authware API
func doRequest[T any](method string, url string, body any) (resp *T, err error) {
	// Serialize and create the body
	bodyBytes, jErr := json.Marshal(body)
	if jErr != nil {
		return nil, jErr
	}

	req, rErr := newRequest(method, url, bodyBytes)
	if rErr != nil {
		return nil, rErr
	}

	// Make the request
	rawResp, cErr := client.Do(req)
	if cErr != nil {
		return nil, cErr
	}

	// Defer closing the body
	defer rawResp.Body.Close()

	// Read all the bytes in the body
	respBytes, ioErr := io.ReadAll(rawResp.Body)
	if ioErr != nil {
		return nil, ioErr
	}

	// We check if it is an OK response code
	if rawResp.StatusCode != http.StatusOK &&
		rawResp.StatusCode != http.StatusCreated &&
		rawResp.StatusCode != http.StatusNoContent {

		// It wasn't, we need to deserialize the response to the default one and return the message
		// in the form of an error
		var defaultResp *DefaultResponse
		err = json.Unmarshal(respBytes, &defaultResp)
		if err != nil {
			return
		}

		// If there is errors in the slice then return the first one as the message
		if len(defaultResp.Errors) > 0 {
			return nil, fmt.Errorf(defaultResp.Errors[0])
		} else { // Otherwise just return the message on the response
			return nil, fmt.Errorf(defaultResp.Message)
		}
	} else {
		err = json.Unmarshal(respBytes, &resp)
		return
	}
}
