package server

import (
	"fmt"
	"net/http"
	"strings"

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
}

func NewServer(config ServerConfig) (*Server, error) {
	return &Server{
		config: config,
	}, nil
}
func (srv *Server) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc(routes.HomePath, routes.HandleSession())
	router.HandleFunc(routes.ConnectPath, routes.HandleConnect(srv.Connect))
	router.HandleFunc(routes.DisconnectPath, routes.HandleDisconnect)
	srv.router = router
}

func (srv *Server) Connect(sess *routes.Session) error {
	driver := strings.ToLower(sess.Driver)
	cfg := dbms.NewConfig(sess.Username, sess.Password, sess.DBname)
	fmt.Printf("Connected to db using driver: %s\n", driver)
	conn, err := dbms.NewClient(sess.Driver, cfg)
	if err != nil {
		return err
	}
	sess.Conn = conn
	return nil
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
