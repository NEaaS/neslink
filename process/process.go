package process

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/neaas/neslink"
	"golang.org/x/sys/unix"
)

// Task represents an OS task.
type Task struct {
	processID  int
	taskID     int
	netnsDev   uint64
	netnsInode uint64
}

const (
	processRoot string = "/proc"
)

var (
	procRegex *regexp.Regexp = regexp.MustCompile("^(?:(?:/proc/)([0-9]*)(?:/task/)(.*)(?:/ns/net))$")
)

// Process returns the process ID of the task.
func (t Task) Process() int {
	return t.processID
}

// Task returns the task ID.
func (t Task) Task() int {
	return t.taskID
}

// NsDev returns the network ns Dev ID for the task.
func (t Task) NsDev() uint64 {
	return t.netnsDev
}

// NsInode returns the network ns Inode ID for the task.
func (t Task) NsInode() uint64 {
	return t.netnsInode
}

// NAGetProcesses creates an action that when run will populate the given Task
// slice with the tasks that reside in the given the executing network
// namespace.
func NAGetProcesses(tasks *[]Task) neslink.NsAction {
	return neslink.NAGeneric(
		"get-ns-processes",
		func() error {
			currNs, err := neslink.NPNow().Provide()
			if err != nil {
				return fmt.Errorf("failed to get the active ns for the process list")
			}
			var currNsStat unix.Stat_t
			if err := unix.Stat(currNs.String(), &currNsStat); err != nil {
				return fmt.Errorf("failed to get the active ns id for the process list")
			}
			currNsDev := currNsStat.Dev
			currNsIno := currNsStat.Ino
			nsTasks := make([]Task, 0)
			if err := filepath.WalkDir(processRoot, func(s string, d fs.DirEntry, e error) error {
				match := procRegex.FindStringSubmatch(s)
				if match == nil {
					return nil
				}
				if len(match) != 3 {
					return nil
				}
				var stat unix.Stat_t
				if err := unix.Stat(match[0], &stat); err != nil {
					return nil
				}
				if stat.Dev == currNsDev && stat.Ino == currNsIno {
					procID, _ := strconv.Atoi(match[1])
					taskID, _ := strconv.Atoi(match[2])
					nsTasks = append(nsTasks, Task{procID, taskID, stat.Dev, stat.Ino})
				}
				return nil
			}); err != nil {
				return fmt.Errorf("failed to get processes via path: %w", err)
			}
			*tasks = nsTasks
			return nil
		},
	)
}
