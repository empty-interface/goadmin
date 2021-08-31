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
	conn   *dbms.GormConnection
}

func NewServer(config ServerConfig) (*Server, error) {
	return &Server{
		config: config,
	}, nil
}
func (srv *Server) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc(routes.HomePath, routes.HandleHome)
	router.HandleFunc(routes.ConnectPath, routes.HandleConnect(srv.Connect))
	router.HandleFunc(routes.DBPath, routes.HandleDatabase)
	srv.router = router
}
func (srv *Server) Connect(driver, username, password, dbname string) error {
	driver = strings.ToLower(driver)
	config := dbms.Config{
		Host:     "localhost",
		Username: username,
		DBName:   dbname,
		Port:     "5432",
		Password: password,
	}
	conn, err := dbms.New(driver, config)
	if err != nil {
		return err
	}
	fmt.Printf("Connected to db using driver: %s\n", driver)
	srv.conn = conn
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
