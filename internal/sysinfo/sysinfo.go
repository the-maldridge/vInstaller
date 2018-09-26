package sysinfo

import (
	"fmt"
	"strings"
	"log"

	"github.com/jaypipes/ghw"
)

// System tells us everything about the machine itself, its used to
// drive other menus and information points in the installer.
type System struct {
	Mem *ghw.MemoryInfo
	Blk *ghw.BlockInfo
	Net *ghw.NetworkInfo
	CPU *ghw.CPUInfo
}

func (s *System) String() string {
	out := []string{"Your system appears to have the following characteristics:"}
	if s.CPU != nil {
		for i, c := range s.CPU.Processors {
			out = append(out, fmt.Sprintf("CPU %d: %s %s", i, c.Vendor, c.Model))
		}
	}
	if s.Mem != nil {
		out = append(out, fmt.Sprintf("Memory: %s", s.Mem.String()))
	}
	if s.Blk != nil {
		out = append(out, fmt.Sprintf("Disks:"))
		for i, d := range s.Blk.Disks {
			out = append(out, fmt.Sprintf("%d: %s", i, d))
			for _, p := range d.Partitions {
				out = append(out, fmt.Sprintf("    %s", p))
			}
		}
	}
	return strings.Join(out, "\n")
}

// DiscoverHardware fetches system info
func DiscoverHardware() *System {
	sys := new(System)
	var err error

	sys.Blk, err = ghw.Block()
	if err != nil {
		log.Printf("Error getting block storage info: %v", err)
	}

	sys.Net, err = ghw.Network()
	if err != nil {
		log.Printf("Error getting network info: %v", err)
	}

	sys.Mem, err = ghw.Memory()
	if err != nil {
		log.Printf("Error getting memory info: %v", err)
	}

	sys.CPU, err = ghw.CPU()
	if err != nil {
		log.Printf("Error getting CPU info: %v", err)
	}

	return sys
}

