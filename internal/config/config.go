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
