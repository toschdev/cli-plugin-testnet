package cmd

import (
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/spf13/cobra"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/toschdev/ignite-testnet/network"
	"github.com/toschdev/ignite-testnet/network/networkchain"
)

// NewNetworkGenesisEditChangeParam creates a new command to send param change request.
func NewNetworkGenesisEditChangeParam() *cobra.Command {
	c := &cobra.Command{
		Use:   "modify-param [launch-id] [module-name] [param-name] [value (json, string, number)]",
		Short: "To request changes to a module parameter in the genesis file",
		RunE:  networkRequestChangeParamHandler,
		Args:  cobra.ExactArgs(4),
	}

	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkRequestChangeParamHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	module := args[1]
	param := args[2]
	value := []byte(args[3])

	n, err := nb.Network()
	if err != nil {
		return err
	}

	// fetch chain information
	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch))
	if err != nil {
		return err
	}

	// check validity of request
	err = c.CheckRequestChangeParam(
		cmd.Context(),
		module,
		param,
		value,
	)
	if err != nil {
		return err
	}

	// create the param change request
	paramChangeRequest := launchtypes.NewParamChange(
		launchID,
		module,
		param,
		value,
	)

	// simulate the param change request
	if err := verifyRequestsFromRequestContents(
		cmd.Context(),
		cacheStorage,
		nb,
		launchID,
		paramChangeRequest,
	); err != nil {
		return err
	}

	// send the request
	return n.SendRequest(cmd.Context(), launchID, paramChangeRequest)
}
