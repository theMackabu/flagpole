package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	jwt "flagpole_auth/api/auth"
	"github.com/bitly/go-simplejson"
	"flagpole_auth/api/users/models"
	"github.com/gildas/go-logger"
)

func SetJSONResponse(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		handler.ServeHTTP(w, r)
	}
}

func ValidateRequestBody(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rBody, rError := io.ReadAll(r.Body)

		user := models.User{}

		if rError != nil {
			log.Println(rError.Error())
		}

		json.Unmarshal(rBody, &user)

		r.Body = ioutil.NopCloser(bytes.NewBuffer(rBody))

		if len(user.Username) < 5 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("The username must be 5 characters min.")
			return
		}

		if len(user.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("The password must be 8 characters min.")
			return
		}

		handler(w, r)
	}
}

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := simplejson.New()
		sentToken := r.Header.Get("Authorization")

		if sentToken == "" {
			data.Set("missing_headers", "Authorization")
			
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(data)
			return
		}

		tokenFromRequest := strings.Split(sentToken, "Bearer")

		if len(tokenFromRequest) == 2 {
			sentToken = strings.TrimSpace(tokenFromRequest[1])
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Incorrect Format")
			return
		}

		jwtWrapper := jwt.JWTWrapper{}
		claims, err := jwtWrapper.ValidateToken(sentToken)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			data.Set("error", err.Error())
			json.NewEncoder(w).Encode(data)
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Attribute)
		handler(w, r.WithContext(ctx))
	}
}

func Logger(next http.Handler) http.Handler {
	 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Must(logger.FromContext(r.Context()))
		next.ServeHTTP(w, r)
	 })
}