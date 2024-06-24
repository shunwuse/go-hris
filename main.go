package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/routes"
)

func main() {
	env := lib.NewEnv()

	r := gin.Default()

	// Routes
	routes := routes.NewRoutes()
	routes.Setup(r)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(":" + env.ServerPort)
}
