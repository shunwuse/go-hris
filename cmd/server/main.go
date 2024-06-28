package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/routes"
)

// @title           Swagger HRIS API
// @version         1.0
// @description     This is a sample server HRIS API.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

func main() {
	env := lib.NewEnv()

	r := gin.Default()

	// Routes
	routes := routes.NewRoutes()
	routes.Setup(r)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(":" + env.ServerPort)
}
