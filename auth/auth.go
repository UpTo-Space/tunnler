package main

import (
	"net/http"
	"strings"

	"github.com/UpTo-Space/tunnler/common"
)

func (as *authServer) isAuthenticated(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(common.AuthHeaderKey)

		if authHeader == "" {
			w.Write([]byte("Auth missing"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authHeader)

		if len(fields) != 2 {
			w.Write([]byte("Invalid or missing Bearer Token."))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authType := fields[0]
		if strings.ToLower(authType) != common.AuthHeaderBearerType {
			w.Write([]byte("Authorization type not supported"))
			w.WriteHeader(http.StatusUnauthorized)
		}

		token := fields[1]
		_, err := as.tokenMaker.VerifyToken(token)
		if err != nil {
			w.Write([]byte("Invalid token"))
			w.WriteHeader(http.StatusUnauthorized)
		}

		h(w, r)
	})
}
