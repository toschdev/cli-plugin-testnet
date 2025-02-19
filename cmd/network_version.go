package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"

	"github.com/toschdev/ignite-testnet/network/networktypes"
)

// NewNetworkVersion creates a new version command to get the version of the app
// The version of the app to use to interact with a chain might be specified by the coordinator
func NewNetworkVersion() *cobra.Command {
	c := &cobra.Command{
		Use:   "version",
		Short: "Version of the app",
		Long: `The version of the app to use to interact with a chain might be specified by the coordinator.
`,
		RunE: networkVersion,
		Args: cobra.NoArgs,
	}
	return c
}

func networkVersion(_ *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	return session.Printf("%s\n", networktypes.Version)
}
