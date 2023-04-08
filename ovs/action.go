package ovs

import (
	"errors"
	"fmt"

	"github.com/neaas/go-openvswitch/ovs"
	"github.com/willfantom/neslink"
)

// LAAttachBridge produces an action that when run will add an existing bridge
// to openvswitch or create a new bridge.
func LAAttachBridge(client *ovs.Client, provider neslink.LinkProvider) neslink.LinkAction {
	return neslink.LAGeneric(
		"attach-ovs-bridge",
		func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(neslink.ErrNoLink, err)
			} else {
				if l.Type() != "bridge" {
					return fmt.Errorf("provided interface is not a bridge")
				}
				if err := client.VSwitch.AddBridge(l.Attrs().Name); err != nil {
					return fmt.Errorf("ovs bridge could not be attached: %w", err)
				}
			}
			return nil
		},
	)
}

// LADetachBridge produces an action that when run will detach a bridge from
// openvswitch.
func LADetachBridge(client *ovs.Client, provider neslink.LinkProvider) neslink.LinkAction {
	return neslink.LAGeneric(
		"del-ovs-bridge",
		func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(neslink.ErrNoLink, err)
			} else {
				if err := client.VSwitch.DeleteBridge(l.Attrs().Name); err != nil {
					return fmt.Errorf("ovs bridge could not be detached: %w", err)
				}
			}
			return nil
		},
	)
}
