package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kalitheniks/goadmin/dbms"
	"github.com/kalitheniks/goadmin/routes"
)

type ServerConfig struct {
	Port string
	Host string
}
type Server struct {
	config ServerConfig
	router *mux.Router
	conn   dbms.DBMS
}

func NewServer(config ServerConfig) (*Server, error) {
	return &Server{
		config: config,
	}, nil
}
func (srv *Server) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/", routes.HandleHome)
	srv.router = router
}
func (srv *Server) getAddr() string {
	return fmt.Sprintf("%s:%s", srv.config.Host, srv.config.Port)
}
func (srv *Server) ListenAndServe() error {
	srv.setupRoutes()
	addr := srv.getAddr()
	fmt.Printf("Server running on: %s\n", addr)
	return http.ListenAndServe(addr, srv.router)
}
