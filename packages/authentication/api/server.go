package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"flagpole_auth/api/users/transport"
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

func (server *HttpServer) StartServer(address string, origins string) {
	fmt.Println("server started: ", address)

	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(origins, ","),
		AllowCredentials: true,
	})

	log.Fatal(http.ListenAndServe(address, c.Handler(transport.Logger(server.Router))))
}

func (server *HttpServer) SetRoutes() {
	server.Router.HandleFunc("/api/login", transport.SetJSONResponse(server.handlers.Login)).Methods("POST")
	server.Router.HandleFunc("/api/signup", transport.SetJSONResponse(transport.ValidateRequestBody(server.handlers.Signup))).Methods("POST")
	server.Router.HandleFunc("/api/refresh", transport.SetJSONResponse(transport.IsAuthorized(server.handlers.Refresh))).Methods("GET")
}
