package internal

import "github.com/shirou/gopsutil/v4/host"

// HostID returns the unique host ID provided by the OS.
func HostID() (string, error) {
	return host.HostID()
}
