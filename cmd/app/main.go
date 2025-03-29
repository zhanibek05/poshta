package main

import (
	"flag"
	"poshta/internal/app"
)

func main() {
	configFile := flag.String("config", "./configs/.env", "Path to configuration file")
	flag.Parse()

	app.Run(*configFile)
}
