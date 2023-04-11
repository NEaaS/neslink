package main

import (
	"fmt"
	"syscall"

	"github.com/neaas/nescript"
	"github.com/neaas/neslink"
	"github.com/neaas/neslink/process"
	"golang.org/x/sys/unix"
)

func main() {
	script := nescript.NewScript("yes &").MustCompile()
	nesProc := *new(nescript.Process)

	if err := neslink.Do(neslink.NPNow(), neslink.NANewNs("example"), neslink.NAExecNescript(script, nil, &nesProc)); err != nil {
		panic(err)
	}

	tasks := make([]process.Task, 0)
	if err := neslink.Do(neslink.NPName("example"), process.NAGetProcesses(&tasks), neslink.NADeleteNamed("example")); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", tasks)

	unix.Kill(tasks[0].Process(), syscall.SIGKILL)
}
