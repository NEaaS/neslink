package main

import (
	"fmt"

	"github.com/neaas/go-openvswitch/ovs"
	"github.com/neaas/neslink"
	nlovs "github.com/neaas/neslink/ovs"
	"github.com/vishvananda/netlink"
)

func main() {
	// setup
	links := make([]netlink.Link, 0)
	ovsClient := ovs.New(ovs.Sudo())

	// 2. create new netns (example), add new bridge in it (egbr0), attach new
	// bridge to ovs, get all links in the new ns...
	if err := neslink.Do(
		neslink.NPNow(),
		neslink.NANewNs("example"),
		neslink.LANewBridge("egbr0"),
		nlovs.LAAttachBridge(ovsClient, neslink.LPName("egbr0")),
		neslink.NALinks(&links),
	); err != nil {
		panic(err)
	}

	// 3. get all ovs bridges
	ovsBrs, err := ovsClient.VSwitch.ListBridges()
	if err != nil {
		panic(err)
	}

	// 4. dump info
	fmt.Printf("1 | links found: %d\n", len(links))
	for _, l := range links {
		fmt.Printf("\t%d: %s - %s\n", l.Attrs().Index, l.Attrs().Name, l.Attrs().HardwareAddr)
	}
	fmt.Printf("1 | ovs brs found: %d\n", len(ovsBrs))
	for idx, l := range ovsBrs {
		fmt.Printf("\t%d: %s\n", idx, l)
	}

	// 5. in the new ns (example), detach the new bridge (egbr0) from ovs, then
	// delete the ns...
	if err := neslink.Do(
		neslink.NPName("example"),
		nlovs.LADetachBridge(ovsClient, neslink.LPName("egbr0")),
		neslink.NALinks(&links),
		neslink.NADeleteNamed("example"),
	); err != nil {
		panic(err)
	}

	// 6. get all ovs bridges
	ovsBrs, err = ovsClient.VSwitch.ListBridges()
	if err != nil {
		panic(err)
	}

	// 7. dump info again
	fmt.Printf("2 | links found: %d\n", len(links))
	for _, l := range links {
		fmt.Printf("\t%d: %s - %s\n", l.Attrs().Index, l.Attrs().Name, l.Attrs().HardwareAddr)
	}
	fmt.Printf("2 | ovs brs found: %d\n", len(ovsBrs))
	for idx, l := range ovsBrs {
		fmt.Printf("\t%d: %s\n", idx, l)
	}
}
