package installer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/the-maldridge/vInstaller/internal/config"
	"github.com/the-maldridge/vInstaller/internal/keys"
)

// Installer is a type to contain channels and other assets on the
// install process.
type Installer struct {
	Config *config.Config
	Output chan string
	Errors chan error
	Done   chan bool

	Meta *config.Meta

	target string
}

func mountSpecials() error {
	log.Println("Mounting special filesystems")
	return nil
}

func unmountSpecials() error {
	log.Println("Unmounting special filesystems")
	return nil
}

func (i *Installer) runCommand(cmdstr string) error {
	// This breaks quoted strings, and is an ugly hack.  But it
	// works.
	parts := strings.Split(cmdstr, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("could not get stderr pipe: %v", err)
		i.Errors <- err
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("could not get stdout pipe: %v", err)
		i.Errors <- err
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			msg := scanner.Text()
			i.Output <- msg
			log.Println(msg)
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			msg := scanner.Text()
			i.Output <- msg
			log.Println(msg)
		}
	}()
	log.Printf("$ %s", cmdstr)
	if err := cmd.Run(); err != nil {
		log.Printf("could not run cmd: %v", err)
		i.Errors <- err
		return err
	}
	if err != nil {
		i.Errors <- err
		log.Printf("could not wait for cmd: %v", err)
		return err
	}
	return nil
}

// Install attempts to put a system on disk
func (i *Installer) Install(target string) {
	defer i.closeChannels()
	i.target = target

	if i.Meta == nil {
		i.Meta = config.DefaultMeta()
	}

	mountSpecials()
	defer unmountSpecials()

	if err := i.verifyTargetDir(); err != nil {
		log.Fatal(err)
	}

	i.installBaseSystem()

	log.Println("System installed")

	i.Done <- true
}

func (i *Installer) closeChannels() {
	close(i.Output)
	close(i.Errors)
	close(i.Done)
}

func (i *Installer) verifyTargetDir() error {
	if stat, err := os.Stat(i.target); err != nil || !stat.IsDir() {
		log.Println(err)
		i.Errors <- err
		return err
	}
	return nil
}

func (i *Installer) installBaseSystem() error {
	return i.xbpsInstall([]string{"base-system"})
}

func (i *Installer) xbpsInstall(pkgs []string) error {
	baseDir := filepath.Join(i.target, "var/db/xbps/")

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		i.Output <- "Installing keys"
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			i.Errors <- err
			return err
		}

		if err := keys.RestoreAssets(baseDir, "keys"); err != nil {
			i.Errors <- err
			return err
		}
	}

	cmd := fmt.Sprintf("sudo xbps-install -y -S -i -R %s -M -r %s %s",
		i.Meta.Mirror,
		i.target,
		strings.Join(pkgs, " "),
	)

	// We drop the error here because xbps-install can return
	// non-zero in places where things should otherwise be fine.
	i.runCommand(cmd)
	return nil
}
