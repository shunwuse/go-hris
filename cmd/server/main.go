package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/routes"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	env := lib.NewEnv()

	r := gin.Default()

	// Routes
	routes := routes.NewRoutes()
	routes.Setup(r)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(":" + env.ServerPort)
}
