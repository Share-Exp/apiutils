package apiutils

import (
	"net/http"
	"strings"

	"github.com/Share-Exp/apiutils/resterrors"
)

// GetAuthorizationToken receives the request and extracts the token.
func GetAuthorizationToken(request *http.Request) (string, resterrors.RestErr) {
	if request == nil {
		return "", resterrors.NewBadRequestError("Request can't be nil.")
	}

	token := strings.TrimSpace(request.Header["Authorization"][0])
	token = strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))

	return token, nil
}
