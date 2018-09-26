package prompt

import (
	"fmt"
	"bufio"
	"os"
	"strings"

	"github.com/the-maldridge/vInstaller/internal/frontend"
	"github.com/the-maldridge/vInstaller/internal/sysinfo"
	"github.com/the-maldridge/vInstaller/internal/config"
)

// Frontend is a simple frontend that just asks questions in a
// terminal and builds an installer configuration from that.
type Frontend struct {
	sysinfo *sysinfo.System
	config *config.Config
}

// New returns a ready to use PromptFrontend
func New() (frontend.InstallerFrontend, error) {
	return new(Frontend), nil
}

func init() {
	frontend.Register("prompt", New)
}

func prompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text
}

// GetInstallerConfig prompts the user for configuration values.
func (f *Frontend) GetInstallerConfig() error {
	fmt.Println("Welcome to the Void Linux Installer")
	fmt.Println("")
	fmt.Println("Please wait while the installer inspects your system...")
	fmt.Println("")
	f.sysinfo = sysinfo.DiscoverHardware()

	fmt.Println(f.sysinfo)
	fmt.Println("")

	f.config = new(config.Config)
	
	f.promptTimeZone()
	f.promptLocale()
	f.promptGRUB()
	f.promptKeyboard()
	f.promptRootPassword()
	f.promptUsers()

	fmt.Println(f.config)
	return nil
}

// ConfirmInstallation confirms that the user is ready to proceed with
// potentially destructive actions.
func (f *Frontend) ConfirmInstallation() error {
	return nil
}

// ShowInstallationProgress shows the output of the installation.
func (f *Frontend) ShowInstallationProgress() {

}

func (f *Frontend) promptTimeZone() {
	f.config.TimeZone = prompt("Enter your timezone: ")
}

func (f *Frontend) promptLocale() {
	f.config.Locale = prompt("Please enter your GLibC Locale: ")
}

func (f *Frontend) promptGRUB() {
	useGrub := strings.TrimSpace(prompt("Use GRUB? (Y/n): "))
	if useGrub == "" || strings.Contains(strings.ToLower(useGrub), "y") {
		graphical := strings.TrimSpace(prompt("Use graphical GRUB? (Y/n): "))
		if graphical == "" || strings.Contains(strings.ToLower(graphical), "y") {
			f.config.GRUB.UseGraphical = true
		}
		target := strings.TrimSpace(prompt(fmt.Sprintf("Install GRUB to: (/dev/%s)", f.sysinfo.Blk.Disks[0].Name)))
		if target == "" {
			f.config.GRUB.InstallTo = "/dev/"+ f.sysinfo.Blk.Disks[0].Name
			return
		}
		f.config.GRUB.InstallTo = target
	}
}

func (f *Frontend) promptKeyboard() {
	f.config.Keyboard = prompt("Please enter your keyboard layout")
}

func (f *Frontend) promptHostname() {
	f.config.Hostname = prompt("System Hostname: ")
}

func (f *Frontend) promptRootPassword() {
	f.config.RootPassword = prompt("Root Password: ")
}

func (f *Frontend) promptUsers() {
	addUsers := strings.TrimSpace(prompt("Do you wish to add a user? (Y/n) "))
	if addUsers == "" || strings.Contains(addUsers, "y") {
		u := config.User{
			Username: strings.TrimSpace(prompt("Username: ")),
			GECOS: strings.TrimSpace(prompt("Name for the user: ")),
			Password: strings.TrimSpace(prompt("Password: ")),
		}
		groups := prompt("Additional groups (comma seperated): ")
		u.Groups = strings.Split(groups, ",")
		f.config.Users = append(f.config.Users, u)
	}
}
