package main

import (
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/middlewares"
	"github.com/shunwuse/go-hris/routes"
)

type Server struct {
	env               lib.Env
	router            lib.RequestHandler
	routes            routes.Routes
	commonMiddlewares middlewares.CommonMiddlewares
	logger            lib.Logger
	database          lib.Database
}

func NewServer(
	env lib.Env,
	router lib.RequestHandler,
	routes routes.Routes,
	commonMiddlewares middlewares.CommonMiddlewares,
	logger lib.Logger,
	database lib.Database,
) Server {
	return Server{
		env:               env,
		router:            router,
		routes:            routes,
		commonMiddlewares: commonMiddlewares,
		logger:            logger,
		database:          database,
	}
}

func (server *Server) Run() {
	server.logger.Info("Starting to run server...")

	// setup common middlewares
	server.commonMiddlewares.Setup(server.router.Gin)

	// setup routes
	server.routes.Setup(server.router.Gin)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	server.router.Gin.Run(":" + server.env.ServerPort)
}
