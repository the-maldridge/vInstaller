package frontend

import (
	"flag"
	"log"

	"github.com/the-maldridge/vInstaller/internal/config"
)

// The InstallerFrontend will fetch an installer config, and then
// confirm that the user is ready to proceed
type InstallerFrontend interface {
	GetInstallerConfig() (*config.Config, error)
	ConfirmInstallation() error
	ShowInstallationProgress(<-chan string, <-chan error, <-chan bool)
}

// Factory creates a new InstallerFrontend and returns it
type Factory func() (InstallerFrontend, error)

var (
	frontends map[string]Factory

	frontend = flag.String("frontend", "", "Frontend interface to use")
)

func init() {
	frontends = make(map[string]Factory)
}

// Register adds a new frontend
func Register(name string, factory Factory) {
	if _, ok := frontends[name]; ok {
		return
	}
	frontends[name] = factory
}

// New returns a frontend ready to use
func New() (InstallerFrontend, error) {
	if *frontend == "" && len(frontends) == 1 {
		log.Println("No frontend specified")
		log.Println("Initializing ", List()[0])
		return frontends[List()[0]]()
	}
	if f, ok := frontends[*frontend]; ok {
		log.Println("Using explicitely specified frontend")
		return f()
	}
	log.Println("Problem finding frontend")
	return nil, ErrUnknownFrontend
}

// List returns a list of all frontends
func List() []string {
	l := []string{}
	for f := range frontends {
		l = append(l, f)
	}
	return l
}
