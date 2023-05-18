package transport

import (
	"encoding/json"
	"io"
	"net/http"

	jwt "flagpole_auth/api/auth"
	"flagpole_auth/api/users/data"
	"flagpole_auth/api/users/models"
	"github.com/Edmartt/go-password-hasher/hasher"
	"github.com/bitly/go-simplejson"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

var log = logger.Create("authentication", &logger.StdoutStream{Unbuffered: true})

type Handlers struct {
	user        *models.User
	logResponse LoginResponse
	sigResponse SignupResponse
	wrapper     jwt.JWTWrapper
}

func init() {
	data.RepoAccessInterface = data.UserRepository{}
}

func (h *Handlers) Login(w http.ResponseWriter, request *http.Request) {
	reqBody, requestError := io.ReadAll(request.Body)
	json_body := simplejson.New()

	if requestError != nil {
		log.Errorf(requestError.Error())
	}

	json.Unmarshal(reqBody, &h.user)
	searchedUser := data.RepoAccessInterface.Find(h.user.Username)

	if searchedUser.Username == h.user.Username {
		if hasher.CheckHash(searchedUser.Password, h.user.Password) {
			newToken, err := h.wrapper.GenerateJWT(h.user.Username, 5)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.user = nil

				return
			}

			h.logResponse.Token = newToken
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&h.logResponse)
			h.user = nil

			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		json_body.Set("error", "invalid username or password")
		json.NewEncoder(w).Encode(json_body)
		h.user = nil

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	json_body.Set("error", "invalid username or password")
	json.NewEncoder(w).Encode(json_body)
	h.user = nil
}

func (h *Handlers) Refresh(w http.ResponseWriter, request *http.Request) {
	reqBody, requestError := io.ReadAll(request.Body)
	json_body := simplejson.New()

	if requestError != nil {
		log.Errorf(requestError.Error())
	}

	json.Unmarshal(reqBody, &h.user)
	searchedUser := data.RepoAccessInterface.Find(h.user.Username)

	if searchedUser.Username != "" {
		newToken, err := h.wrapper.GenerateJWT(h.user.Username, 5)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.user = nil

			return
		}

		h.logResponse.Token = newToken
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&h.logResponse)
		h.user = nil

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	json_body.Set("error", "invalid username or token")
	json.NewEncoder(w).Encode(json_body)
	h.user = nil
}

func (h *Handlers) Signup(w http.ResponseWriter, request *http.Request) {
	requestError := json.NewDecoder(request.Body).Decode(&h.user)
	h.user.Id = uuid.New().String()

	hashedPassword := hasher.ConvertToHash(h.user.Password)
	h.user.Password = hashedPassword

	if requestError != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.user = nil

		return
	}

	_, createErr := data.RepoAccessInterface.Create(h.user)

	if createErr != nil {
		h.sigResponse.Created = false
		h.sigResponse.Message = "user already exists"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(h.sigResponse)
		h.user = nil
	} else {
		h.sigResponse.Created = true
		h.sigResponse.Message = "created new user"

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(h.sigResponse)
		h.user = nil
	}
}
