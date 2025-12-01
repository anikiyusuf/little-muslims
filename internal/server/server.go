package server 

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yusufaniki/muslim_tech/internal/boostrap"
    "github.com/yusufaniki/muslim_tech/internal/router"
)


type Server struct {
	httpServer *http.Server
	app        *boostrap.Application
}


func NewServer(app *boostrap.Application) *Server{
	router := router.SetupRoutes(app)

	return &Server{
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%s", app.Config.Port),
			Handler: router,
			ReadTimeout: 10 * time.Second,		
			WriteTimeout: 10 * time.Second,
			IdleTimeout: time.Minute,

		},
		app: app,
	}
}

func (srv *Server) Start() error {
	srv.app.Logger.Infof("Server is starting on port ...")
	return srv.httpServer.ListenAndServe()
}

func (srv *Server) Shutdown(timeout time.Duration) error {
	srv.app.Logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return srv.httpServer.Shutdown(ctx)
}