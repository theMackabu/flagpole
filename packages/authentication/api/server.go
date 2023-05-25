package api

import (
	"net/http"
	"strings"

	"flagpole/packages/authentication/api/users/transport"
	"github.com/gildas/go-logger"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type HttpServer struct {
	Router   *mux.Router
	handlers transport.Handlers
}

func (server *HttpServer) CreateRouter() {
	server.Router = mux.NewRouter().StrictSlash(true)
	server.SetRoutes()
}

func (server *HttpServer) StartServer(address string, origins string, log *logger.Logger) {
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(origins, ","),
		AllowCredentials: true,
	})

	log.Record("address", address).Infof("server started")
	log.Record("server", http.ListenAndServe(address, c.Handler(log.HttpHandler()(transport.Logger(server.Router))))).Fatalf("Failed to start server...")
}

func (server *HttpServer) SetRoutes() {
	server.Router.HandleFunc("/_/login", transport.SetJSONResponse(server.handlers.Login)).Methods("POST")
	server.Router.HandleFunc("/_/signup", transport.SetJSONResponse(transport.ValidateRequestBody(server.handlers.Signup))).Methods("POST")
	server.Router.HandleFunc("/_/refresh", transport.SetJSONResponse(transport.IsAuthorized(server.handlers.Refresh))).Methods("POST")
}
