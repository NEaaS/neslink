//go:build !linux
// +build !linux

package neslink

import (
	"fmt"
	"os"
)

// Inode provides the inode of the network namespace file.
func (n NsFd) Inode() (uint64, error) {
	return 0, fmt.Errorf("failed to obtain file info")
}

// Dev descirbes the device in which the network ns file resides.
func (n NsFd) Dev() (uint64, error) {
	return 0, fmt.Errorf("failed to obtain file info")
}

// close closes the file descriptor. This should be used to clean up any opned
// file descriptor.
func (n NsFd) close() error {
	fmt.Errorf("file descriptor can not be closed on non-linux builds")
}

// open opens the file for the network namespace and returns the file
// descriptor.
func (ns Namespace) open() (NsFd, error) {
	return NsFdNone, fmt.Errorf("netns file descriptor can not be opened on non-linux builds")
}

// set sets the current namespace to the one associated with the given file
// descriptor.
func (ns NsFd) set() error {
	fmt.Errorf("netns can not be set on non-linux builds")
}
