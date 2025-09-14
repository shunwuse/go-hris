package main

// @title           Swagger HRIS API
// @version         1.0
// @description     This is a sample server HRIS API.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

func main() {
	server := InitializeServer()

	server.Run()
}
