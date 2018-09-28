package config

import (
	"fmt"
	"strings"
)

// Config represents the configuration of a system install.
type Config struct {
	TimeZone string
	Locale   string
	Keyboard string
	Hostname string

	RootPassword string

	Users []User

	GRUB struct {
		UseGraphical bool
		InstallTo    string
	}

	Filesystems []Filesystem
}

// User represents a system user
type User struct {
	Username string
	GECOS    string
	Password string
	Groups   []string
}

func (c Config) String() string {
	out := []string{"Your system configuration is as follows:"}

	out = append(out, fmt.Sprintf("Hostname: %s", c.Hostname))
	out = append(out, fmt.Sprintf("Keyboard: %s", c.Keyboard))
	out = append(out, fmt.Sprintf("Timezone: %s", c.TimeZone))
	out = append(out, fmt.Sprintf("Locale: %s", c.Locale))

	for i, u := range c.Users {
		out = append(out, fmt.Sprintf("User 100%d", i))
		out = append(out, fmt.Sprintf("  Username: %s", u.Username))
		out = append(out, fmt.Sprintf("  Name: %s", u.GECOS))
		out = append(out, fmt.Sprintf("  Groups: %s", strings.Join(u.Groups, ",")))
	}

	return strings.Join(out, "\n")
}

// Meta contains information about the install and what settings will
// affect its operation, such as mirrors.
type Meta struct {
	Mirror   string
	Services []string
}

// DefaultMeta returns the default metadata which should be safe to use
func DefaultMeta() *Meta {
	return &Meta{
		Mirror:   "http://mirrors.servercentral.com/voidlinux/current",
		Services: []string{"dhcpcd", "sshd"},
	}
}

// Filesystem represents a filesystem that is ready to go into the
// installer, and can be mapped onto /etc/fstab.
type Filesystem struct {
	FS      string
	MountTo string
	Type    string
	Options string
	Dump    int
	Pass    int
}
