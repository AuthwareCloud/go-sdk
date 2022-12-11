package authware

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	application       *Application                                                     // application will contain the current application and information
	client            = newClient()                                                    // client is the HTTP request client to send HTTP requests
	baseAddress       = "https://api.authware.org/"                                    // baseAddress is the address that is prepended to every request URL passed into doRequest
	cloudflareIssuer  = []byte("CN=Cloudflare Inc ECC CA-3, O=Cloudflare, Inc., C=US") // cloudflareIssuer is the bytes for a valid issuer on a Cloudflare certificate
	letsEncryptIssuer = []byte(", O=Let's Encrypt, C=US")                              // letsEncryptIssuer is the bytes for a valid issuer on a Let's Encrypt certificate
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

					// Now we validate the DNS names (domains) of the certificates
					validDomain := false

					// Enumerate over every DNS domain in the first peer certificate
					for _, domain := range state.PeerCertificates[0].PermittedDNSDomains {
						// If the current enumerated domain contains 'authware.org'
						if strings.Contains(domain, "authware.org") {
							// We set the valid domain to true and break out of the loop
							validDomain = true
							break
						}
					}

					// Now we check whether the domain validity check above passed
					if !validDomain {
						// If it didn't, return the tamperedCertificate error
						return TamperedCertificate
					}

					return nil
				},
			},
			TLSHandshakeTimeout:   10,
			ResponseHeaderTimeout: 10,
			ExpectContinueTimeout: 10,
		},
		Timeout: 10 * time.Second,
	}
}

// addAuthenticationToken will add an authentication token to the backing store for the HTTP client
func addAuthenticationToken(token string) {
	application.AuthToken = token
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
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := newRequest(method, url, bodyBytes)
	if err != nil {
		return
	}

	// Make the request
	rawResp, err := client.Do(req)
	if err != nil {
		return
	}

	// Defer closing the body
	defer rawResp.Body.Close()

	// Read all the bytes in the body
	respBytes, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return
	}

	// We check if it is an OK response code
	if rawResp.StatusCode != http.StatusOK &&
		rawResp.StatusCode != http.StatusCreated &&
		rawResp.StatusCode != http.StatusNoContent {

		// It wasn't, we need to deserialize the response to the default one and return the message
		// in the form of an error
		var defaultResp *DefaultResponse
		err := json.Unmarshal(respBytes, &defaultResp)
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
