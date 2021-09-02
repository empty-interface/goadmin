package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/empty-interface/goadmin/dbms"
	"github.com/empty-interface/goadmin/routes"
	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Port string
	Host string
}
type Server struct {
	config     ServerConfig
	router     *mux.Router
	httpServer *http.Server
}

func NewServer(config ServerConfig) (*Server, error) {
	return &Server{
		config: config,
	}, nil
}
func (srv *Server) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc(routes.DisconnectPath, routes.HandleDisconnect)
	router.Handle(routes.HomePath, routes.HandleSession(srv.Connect))
	router.Handle(routes.ConnectPath, routes.HandleConnect(srv.Connect))
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
func (srv *Server) Close(ctx context.Context) {
	fmt.Println("\nClosing")
	done := make(chan bool, 1)
	go func() {
		srv.httpServer.Shutdown(ctx)
		manager := routes.GetGlobalSessionManager()
		if manager != nil {
			manager.Close()
		}
		done <- true
	}()
	select {
	case <-ctx.Done():
		fmt.Println("Timeout ,", ctx.Err().Error())
	case <-done:
		fmt.Println("\nDone closing server")
	}
}
func (srv *Server) ListenAndServe() {
	srv.setupRoutes()
	addr := srv.getAddr()
	srv.httpServer = &http.Server{
		Addr: addr, Handler: srv.router,
	}
	fmt.Printf("Server running on: %s\n", addr)
	go func() {
		if err := srv.httpServer.ListenAndServe(); err == http.ErrServerClosed {
			fmt.Println("Closed gracefully")
		} else {
			fmt.Println("Server closed err:", err.Error())
		}
	}()
}
