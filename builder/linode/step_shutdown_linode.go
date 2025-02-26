package linode

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/linode/linodego"
)

type stepShutdownLinode struct {
	client linodego.Client
}

func (s *stepShutdownLinode) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)
	instance := state.Get("instance").(*linodego.Instance)

	handleError := func(prefix string, err error) multistep.StepAction {
		return errorHelper(state, ui, prefix, err)
	}

	ui.Say("Shutting down Linode...")
	if err := s.client.ShutdownInstance(ctx, instance.ID); err != nil {
		return handleError("Error shutting down Linode", err)
	}

	_, err := s.client.WaitForInstanceStatus(ctx, instance.ID, linodego.InstanceOffline, int(c.StateTimeout.Seconds()))
	if err != nil {
		return handleError("Error waiting for Linode offline", err)
	}

	return multistep.ActionContinue
}

func (s *stepShutdownLinode) Cleanup(state multistep.StateBag) {}
