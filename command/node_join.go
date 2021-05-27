package command

import (
	"context"
	"fmt"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// NodeJoinCommand :
type NodeJoinCommand struct {
	UI cli.UI
	Command
}

func (c *NodeJoinCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())
	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	return flags
}

// Name :
func (c *NodeJoinCommand) Name() string {
	return "node join"
}

// Synopsis :
func (c *NodeJoinCommand) Synopsis() string {
	return "Join a network"
}

// Run :
func (c *NodeJoinCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.UI.Error("This command takes one argument")
		c.UI.Error(`For additional help, try 'drago node join --help'`)
		return 1
	}

	networkName := args[0]
	nodeID := ""
	networkID := ""

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	if nodeID, err = localAgentNodeID(api); err != nil {
		c.UI.Error(fmt.Sprintf("Error determining local node ID: %s", err))
		return 1
	}

	networks, err := api.Networks().List()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting networks: %s", err))
		return 1
	}

	for _, n := range networks {
		if n.Name == networkName {
			networkID = n.ID

			break
		}
	}

	if networkID == "" {
		c.UI.Error("Error: network not found")
		return 1
	}

	if err = api.Interfaces().Create(nodeID, networkID); err != nil {
		c.UI.Error(fmt.Sprintf("Error joining network: %s", err))
		return 1
	}

	c.UI.Output("Joined!")

	return 0
}

// Help :
func (c *NodeJoinCommand) Help() string {
	h := `
Usage: drago node join <network> [options]

  Have the local client node join an existing network.

  If ACLs are enabled, this option requires a token with the 'interface:write' capability.

General Options:
` + GlobalOptions()

	return strings.TrimSpace(h)
}
