package main

import (
	"flag"
	"log"

	"github.com/the-maldridge/vInstaller/internal/frontend"
	_ "github.com/the-maldridge/vInstaller/internal/frontend/prompt"
	_ "github.com/the-maldridge/vInstaller/internal/frontend/test"

	"github.com/the-maldridge/vInstaller/internal/installer"
)

var (
	targetDir = flag.String("target", "/target", "Mountpoint for the target filesystem")
)

func main() {
	flag.Parse()
	log.Println("Welcome to the installer!")

	f, err := frontend.New()
	if err != nil {
		log.Fatal("Bad frontend: ", err)
	}

	cfg, err := f.GetInstallerConfig()
	if err != nil {
		log.Println(err)
	}

	if err := f.ConfirmInstallation(); err != nil {
		log.Fatal(err)
	}

	output := make(chan string, 50)
	errors := make(chan error, 10)
	done := make(chan bool)

	installer := &installer.Installer{
		Config: cfg,
		Output: output,
		Errors: errors,
		Done:   done,
	}

	go installer.Install(*targetDir)

	f.ShowInstallationProgress(output, errors, done)
}
