package main

import (
	"flag"
	"log"

	"github.com/the-maldridge/vInstaller/internal/frontend"
	_ "github.com/the-maldridge/vInstaller/internal/frontend/prompt"
)

func main() {
	flag.Parse()
	log.Println("Welcome to the installer!")

	f, err := frontend.New()
	if err != nil {
		log.Fatal("Bad frontend: ", err)
	}

	f.GetInstallerConfig()
}
