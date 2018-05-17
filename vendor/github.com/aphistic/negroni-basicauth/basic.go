package basicauth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/negroni"
)

func getCreds(req *http.Request) (string, string) {
	if req == nil {
		return "", ""
	}

	authVal := req.Header.Get("Authorization")
	authTokens := strings.Split(authVal, " ")
	if len(authTokens) < 2 || authTokens[0] != "Basic" {
		return "", ""
	}

	userpass, err := base64.StdEncoding.DecodeString(authTokens[1])
	if err != nil {
		return "", ""
	}

	userpassTokens := strings.Split(string(userpass), ":")
	if len(userpassTokens) < 2 {
		return "", ""
	}

	return userpassTokens[0], userpassTokens[1]
}

func unauthorized(res http.ResponseWriter, realm string) {
	res.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}

func BasicFunc(realm string, authfn func(string, string, *http.Request) bool) negroni.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		user, pass := getCreds(req)

		if len(user) == 0 || !authfn(user, pass, req) {
			unauthorized(res, realm)
			return
		}

		next(res, req)
	}
}

func BasicAuth(realm string, users map[string]string) negroni.HandlerFunc {
	return BasicFunc(realm, func(user, pass string, req *http.Request) bool {
		if userpass, ok := users[user]; ok {
			if pass == userpass {
				return true
			}
		}

		return false
	})
}
