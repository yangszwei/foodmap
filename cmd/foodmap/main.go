package main

import (
	"flag"
	"foodmap/internal/registry"
	"log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", ".env", "Configuration path")
	if err := registry.Start(configPath); err != nil {
		log.Fatalln(err)
	}
}
