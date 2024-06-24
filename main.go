package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/routes"
)

func main() {
	r := gin.Default()

	// Routes
	routes := routes.NewRoutes()
	routes.Setup(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
