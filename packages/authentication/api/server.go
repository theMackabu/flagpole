package api

import (
	"log"
	"net/http"
	"strings"

	"flagpole_auth/api/users/transport"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/gildas/go-logger"
)

type HttpServer struct {
	Router   *mux.Router
	handlers transport.Handlers
}

func (server *HttpServer) CreateRouter() {
	server.Router = mux.NewRouter().StrictSlash(true)
	server.SetRoutes()
}

func (server *HttpServer) StartServer(address string, origins string) {
	log.Println("server started:", address)
	json := logger.Create("authentication")

	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(origins, ","),
		AllowCredentials: true,
	})

	log.Fatal(http.ListenAndServe(address, c.Handler(json.HttpHandler()(transport.Logger(server.Router)))))
}

func (server *HttpServer) SetRoutes() {
	server.Router.HandleFunc("/_/login", transport.SetJSONResponse(server.handlers.Login)).Methods("POST")
	server.Router.HandleFunc("/_/signup", transport.SetJSONResponse(transport.ValidateRequestBody(server.handlers.Signup))).Methods("POST")
	server.Router.HandleFunc("/_/refresh", transport.SetJSONResponse(transport.IsAuthorized(server.handlers.Refresh))).Methods("POST")
}
