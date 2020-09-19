package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Share-Exp/apiutils/resterrors"
	"github.com/mercadolibre/golang-restclient/rest"
)

const (
	headerXPublic   = "X-Public"
	headerXClientID = "X-Client-ID"
	headerXCallerID = "X-Caller-ID"

	paramAccessToken = "access_token"
)

type accessToken struct {
	ID       string `json:"id"`
	UserID   int64  `json:"user_id"`
	ClientID int64  `json:"client_id"`
}

// IsPublic is a function to verify if the request is a public or private call.
// Receives a pointer of http.Request and returns a bool.
func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}

	return request.Header.Get(headerXPublic) == "true"
}

// GetCallerID is a function to recover from the header of the request the user ID.
func GetCallerID(request *http.Request) int64 {
	if request == nil {
		return 0
	}

	callerID, err := strconv.ParseInt(request.Header.Get(headerXCallerID), 10, 64)
	if err != nil {
		return 0
	}

	return callerID
}

// GetClientID is a function to recover from the header of the request the client ID.
func GetClientID(request *http.Request) int64 {
	if request == nil {
		return 0
	}

	clientID, err := strconv.ParseInt(request.Header.Get(headerXClientID), 10, 64)
	if err != nil {
		return 0
	}

	return clientID
}

// AuthenticateRequest is a function to validate the authentication in the request.
func AuthenticateRequest(request *http.Request, baseURL string) resterrors.RestErr {
	if request == nil {
		return nil
	}

	cleanRequest(request)

	accessTokenID := request.URL.Query().Get(paramAccessToken)

	if accessTokenID == "" {
		return nil
	}

	at, err := getAccessToken(accessTokenID, baseURL)
	if err != nil {
		if err.Status() == http.StatusNotFound {
			return nil
		}
		return err
	}

	request.Header.Add(headerXCallerID, fmt.Sprintf("%v", at.UserID))
	request.Header.Add(headerXClientID, fmt.Sprintf("%v", at.ClientID))

	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}

	request.Header.Del(headerXClientID)
	request.Header.Del(headerXCallerID)
}

func getAccessToken(accessTokenID string, baseURL string) (*accessToken, resterrors.RestErr) {
	oauthRestClient := rest.RequestBuilder{
		BaseURL: baseURL,
		Timeout: 200 * time.Millisecond,
	}

	response := oauthRestClient.Get(fmt.Sprintf("oauth/access_token/%s", accessTokenID))

	if response == nil || response.Response == nil {
		return nil, resterrors.NewInternalServerError("Invalid REST client response when trying to get access token.", errors.New("response error"))
	}

	if response.StatusCode > 299 {
		restErr, err := resterrors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, resterrors.NewInternalServerError("Invalid error interface when trying to get access token.", err)
		}
		return nil, restErr
	}

	var at accessToken
	if err := json.Unmarshal(response.Bytes(), &at); err != nil {
		return nil, resterrors.NewInternalServerError("Error when trying to unmarshall access token response.", errors.New("response error"))
	}

	return &at, nil
}
