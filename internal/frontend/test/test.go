package test

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/the-maldridge/vInstaller/internal/config"
	"github.com/the-maldridge/vInstaller/internal/frontend"
	"github.com/the-maldridge/vInstaller/internal/sysinfo"
)

// Frontend binds the loader and allows some prompts to ensure things
// look good before proceeding
type Frontend struct {
	sysinfo *sysinfo.System
	config  *config.Config
}

// New returns a ready to use PromptFrontend
func New() (frontend.InstallerFrontend, error) {
	return new(Frontend), nil
}

func init() {
	frontend.Register("auto", New)
}

// GetInstallerConfig fetches config from somewhere else to provide to
// the installer.
func (f *Frontend) GetInstallerConfig() (*config.Config, error) {
	f.config = &config.Config{
		TimeZone: "America/Los_Angeles",
		Locale:   "en_US.UTF-8",
		Keyboard: "us",
		Hostname: "test",

		RootPassword: "toor",

		Users: []config.User{
			config.User{
				Username: "void",
				GECOS:    "Void User",
				Password: "void",
				Groups:   []string{"wheel"},
			},
		},
	}
	return f.config, nil
}

// ConfirmInstallation may prompt for confirmation, or it may observe
// that there is an automatic override in place which precludes the
// need to proceed.
func (f *Frontend) ConfirmInstallation() error {
	fmt.Println(f.config)

	fmt.Printf("Do you wish to proceed with installation? (yes/no) ")
	reader := bufio.NewReader(os.Stdin)
	proceed, _ := reader.ReadString('\n')
	if strings.TrimSpace(proceed) == "yes" {
		return nil
	}
	return frontend.ErrInstallationAborted
}

// ShowInstallationProgress shows the output of the installation.
func (f *Frontend) ShowInstallationProgress(output <-chan string, errors <-chan error, done <-chan bool) {
	poll := true
	for poll {
		select {
		case o := <-output:
			fmt.Println(o)
		case e := <-errors:
			fmt.Println(e)
		case <-done:
			poll = false
		}
	}
}
