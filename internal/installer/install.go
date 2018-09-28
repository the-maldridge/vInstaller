package installer

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/the-maldridge/vInstaller/internal/config"
	"github.com/the-maldridge/vInstaller/internal/keys"

	"github.com/mattn/go-shellwords"
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
	args, err := shellwords.Parse(cmdstr)
	if err != nil {
		log.Printf("could not get stderr pipe: %v", err)
		i.Errors <- err
		return err
	}
	cmd := exec.Command(args[0], args[1:]...)
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
	var err error
	i.target, err = filepath.Abs(target)
	if err != nil {
		log.Println("Couldn't determine target path")
	}

	if i.Meta == nil {
		i.Meta = config.DefaultMeta()
	}

	mountSpecials()
	defer unmountSpecials()

	if err := i.verifyTargetDir(); err != nil {
		log.Fatal(err)
	}

	if err := i.installBaseSystem(); err != nil {
		log.Println(err)
		return
	}
	if err := i.configureHostname(); err != nil {
		log.Println(err)
		return
	}
	if err := i.configureRCconf(); err != nil {
		log.Println(err)
		return
	}
	if err := i.configureLocaleconf(); err != nil {
		log.Println(err)
		return
	}
	if err := i.configureFStab(); err != nil {
		log.Println(err)
		return
	}
	if err := i.addUsers(); err != nil {
		log.Println(err)
		return
	}

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

	cmd := fmt.Sprintf("xbps-install -y -S -i -R %s -M -r %s %s",
		i.Meta.Mirror,
		i.target,
		strings.Join(pkgs, " "),
	)

	// We drop the error here because xbps-install can return
	// non-zero in places where things should otherwise be fine.
	i.runCommand(cmd)
	return nil
}

func (i *Installer) configureHostname() error {
	// Write the hosts file out
	i.Output <- "Configuring network names"
	i.Output <- "  Configuring /etc/hosts"
	log.Println("Configuring /etc/hosts")
	t, err := fetchTemplate("hosts")
	if err != nil {
		i.Errors <- err
		return err
	}
	f, err := os.OpenFile(filepath.Join(i.target, "etc/hosts"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		i.Errors <- err
		return err
	}
	if err := t.Execute(f, i.Config.Hostname); err != nil {
		i.Errors <- err
		return err
	}
	if err := f.Close(); err != nil {
		i.Errors <- err
		return err
	}
	i.Output <- "    /etc/hosts has been configured"

	// Set the hostname
	i.Output <- "  Configuring /etc/hostname"
	log.Println("Configuring /etc/hostname")
	hostname := []byte(strings.Split(i.Config.Hostname, ".")[0])
	if err := ioutil.WriteFile(filepath.Join(i.target, "etc/hostname"), hostname, 0644); err != nil {
		i.Errors <- err
		return err
	}
	i.Output <- "    /etc/hostname has been configured"

	return nil
}

func (i *Installer) configureRCconf() error {
	// Write the hosts file out
	i.Output <- "Configuring /etc/rc.conf"
	log.Println("Configuring /etc/rc.conf")
	t, err := fetchTemplate("rc.conf")
	if err != nil {
		i.Errors <- err
		return err
	}
	f, err := os.OpenFile(filepath.Join(i.target, "etc/rc.conf"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		i.Errors <- err
		return err
	}
	data := struct {
		TimeZone string
		Keyboard string
	}{
		TimeZone: i.Config.TimeZone,
		Keyboard: i.Config.Keyboard,
	}

	if err := t.Execute(f, data); err != nil {
		i.Errors <- err
		return err
	}
	if err := f.Close(); err != nil {
		i.Errors <- err
		return err
	}
	i.Output <- "  /etc/rc.conf has been configured"
	return nil
}

func (i *Installer) configureLocaleconf() error {
	// Write the hosts file out
	i.Output <- "Configuring /etc/locale.conf"
	log.Println("Configuring /etc/locale.conf")
	t, err := fetchTemplate("locale.conf")
	if err != nil {
		i.Errors <- err
		return err
	}
	f, err := os.OpenFile(filepath.Join(i.target, "etc/locale.conf"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		i.Errors <- err
		return err
	}
	if err := t.Execute(f, i.Config.Locale); err != nil {
		i.Errors <- err
		return err
	}
	if err := f.Close(); err != nil {
		i.Errors <- err
		return err
	}
	i.Output <- "  /etc/locale.conf has been configured"
	return nil
}

func (i *Installer) configureFStab() error {
	// Write the hosts file out
	i.Output <- "Configuring /etc/fstab"
	log.Println("Configuring /etc/fstab")
	t, err := fetchTemplate("fstab")
	if err != nil {
		i.Errors <- err
		return err
	}
	f, err := os.OpenFile(filepath.Join(i.target, "etc/fstab"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		i.Errors <- err
		return err
	}
	if err := t.Execute(f, i.Config.Filesystems); err != nil {
		i.Errors <- err
		return err
	}
	if err := f.Close(); err != nil {
		i.Errors <- err
		return err
	}
	i.Output <- "  /etc/locale.conf has been configured"
	return nil
}

func (i *Installer) enableServices() error {
	i.Output <- "Enabling Services"
	serviceDir := filepath.Join(i.target, "etc/runit/runsvdir/default/")
	for _, s := range i.Meta.Services {
		i.Output <- fmt.Sprintf("  %s", s)
		if err := os.Symlink(filepath.Join(serviceDir, s), filepath.Join("/etc/sv/", s)); err != nil {
			i.Errors <- err
			return err
		}
	}
	return nil
}

func (i *Installer) addUsers() error {
	i.Output <- "Adding user account(s)"
	log.Println("Adding user accounts")

	for _, u := range i.Config.Users {
		cmd := fmt.Sprintf("chroot %s useradd -m -U -G %s -c '%s' %s",
			i.target,
			strings.Join(u.Groups, ","),
			u.GECOS,
			u.Username,
		)
		i.runCommand(cmd)

		cmd = fmt.Sprintf("sh -c 'echo %s:%s | chroot %s chpasswd -c SHA512'",
			u.Username,
			u.Password,
			i.target,
		)
		i.runCommand(cmd)
	}

	i.Output <- "  User accounts added"
	return nil
}

func fetchTemplate(name string) (*template.Template, error) {
	templateString, err := Asset("templates/" + name)
	if err != nil {
		log.Fatal(err)
	}

	return template.New(name).Parse(string(templateString))
}
