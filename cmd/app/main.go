package main

import (
	"flag"
	"poshta/internal/app"
)


// @title           Poshta API
// @version         1.0
// @description     Secure messanger API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.


func main() {
	configFile := flag.String("config", ".env", "Path to configuration file")
	flag.Parse()
	
	app.Run(*configFile)
}
