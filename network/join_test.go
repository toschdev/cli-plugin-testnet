package network

import (
	"context"
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/toschdev/ignite-testnet/network/networktypes"
	"github.com/toschdev/ignite-testnet/network/testutil"
)

const (
	TestDenom                     = "stake"
	TestAmountString              = "95000000"
	TestAmountInt                 = int64(95000000)
	TestAccountRequestID          = uint64(1)
	TestGenesisValidatorRequestID = uint64(2)
)

func TestJoin(t *testing.T) {
	t.Run("successfully get join request with custom public address", func(t *testing.T) {
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)

		expectedReqs := []launchtypes.RequestContent{
			launchtypes.NewGenesisValidator(
				testutil.LaunchID,
				addr,
				gentx.JSON(t),
				[]byte{},
				sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)),
				launchtypes.Peer{
					Id: testutil.NodeID,
					Connection: &launchtypes.Peer_TcpAddress{
						TcpAddress: testutil.TCPAddress,
					},
				},
			),
		}

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		reqs, err := network.GetJoinRequestContents(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
		)
		require.NoError(t, err)
		require.ElementsMatch(t, expectedReqs, reqs)
	})

	t.Run("successfully get join request with public address read from the gentx", func(t *testing.T) {
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)

		expectedReqs := []launchtypes.RequestContent{
			launchtypes.NewGenesisValidator(
				testutil.LaunchID,
				addr,
				gentx.JSON(t),
				[]byte{},
				sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)),
				launchtypes.Peer{
					Id: testutil.NodeID,
					Connection: &launchtypes.Peer_TcpAddress{
						TcpAddress: testutil.TCPAddress,
					},
				},
			),
		}

		reqs, err := network.GetJoinRequestContents(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
		)
		require.NoError(t, err)
		require.ElementsMatch(t, expectedReqs, reqs)
	})

	t.Run("successfully get join request with account request", func(t *testing.T) {
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)

		expectedReqs := []launchtypes.RequestContent{
			launchtypes.NewGenesisValidator(
				testutil.LaunchID,
				addr,
				gentx.JSON(t),
				[]byte{},
				sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)),
				launchtypes.Peer{
					Id: testutil.NodeID,
					Connection: &launchtypes.Peer_TcpAddress{
						TcpAddress: testutil.TCPAddress,
					},
				},
			),
			launchtypes.NewGenesisAccount(
				testutil.LaunchID,
				addr,
				sdk.NewCoins(sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt))),
			),
		}

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		reqs, err := network.GetJoinRequestContents(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithAccountRequest(sdk.NewCoins(sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)))),
			WithPublicAddress(testutil.TCPAddress),
		)
		require.NoError(t, err)
		require.ElementsMatch(t, expectedReqs, reqs)
	})

	t.Run("failed to get join request, failed to read node id", func(t *testing.T) {
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)
		expectedError := errors.New("failed to get node id")

		suite.ChainMock.
			On("NodeID", mock.Anything).
			Return("", expectedError).
			Once()

		_, err = network.GetJoinRequestContents(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
		)
		require.ErrorIs(t, err, expectedError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to get join request, failed to read gentx", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			gentxPath      = "invalid/path"
			suite, network = newSuite(account)
		)

		_, err := network.GetJoinRequestContents(context.Background(), suite.ChainMock, testutil.LaunchID, gentxPath)
		require.Error(t, err)
		suite.AssertAllMocks(t)
	})
}
