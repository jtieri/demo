package e2e

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	coinztypes "github.com/jtieri/demo/x/coinz/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"go.uber.org/zap/zaptest"
)

var (
	defaultAdminAddr = "cosmos1jtvahy53nvaftwq6r04dmnalv6re2av0xgp705"
	defaultAdminSeed = "today please woman quote finger problem patrol inflict usual crystal brick gesture myself permit clutch upset book clarify embrace mixture enact board awful almost"

	defaultAdmin = coinztypes.Admin{
		Address: defaultAdminAddr,
	}

	coinzDenom    = "myToken"
	initialSupply = 0

	defaultAssetMetadata = coinztypes.AssetMetadata{Asset: sdk.NewCoin(coinzDenom, math.NewInt(int64(initialSupply)))}

	numVals  = 1
	numNodes = 0

	image = ibc.DockerImage{
		Repository: "demod",
		Version:    "local",
		UIDGID:     "1025:1025",
	}

	chainASpec = interchaintest.ChainSpec{
		ChainConfig: ibc.ChainConfig{
			Type:           "cosmos",
			Name:           "demo",
			ChainID:        "demo-1",
			Images:         []ibc.DockerImage{image},
			Bin:            "demod",
			Bech32Prefix:   "cosmos",
			Denom:          "stake",
			CoinType:       "118",
			GasPrices:      "0.0stake",
			GasAdjustment:  1.5,
			TrustingPeriod: "100h",
			PreGenesis:     nil,
			ModifyGenesis:  modifyGenesis(defaultAdmin),
			EncodingConfig: encodingConfig(),
		},
		NumValidators: &numVals,
		NumFullNodes:  &numNodes,
	}

	chainBSpec = interchaintest.ChainSpec{
		ChainConfig: ibc.ChainConfig{
			Type:           "cosmos",
			Name:           "demo",
			ChainID:        "demo-2",
			Images:         []ibc.DockerImage{image},
			Bin:            "demod",
			Bech32Prefix:   "cosmos",
			Denom:          "stake",
			CoinType:       "118",
			GasPrices:      "0.0stake",
			GasAdjustment:  1.5,
			TrustingPeriod: "100h",
			PreGenesis:     nil,
			ModifyGenesis:  modifyGenesis(defaultAdmin),
			EncodingConfig: encodingConfig(),
		},
		NumValidators: &numVals,
		NumFullNodes:  &numNodes,
	}
)

// modifyGenesis is used in the ChainConfig to make changes to the genesis file before starting the chain.
func modifyGenesis(admin coinztypes.Admin) func(cfg ibc.ChainConfig, bz []byte) ([]byte, error) {
	return func(cfg ibc.ChainConfig, bz []byte) ([]byte, error) {
		genesis := []cosmos.GenesisKV{
			cosmos.NewGenesisKV("app_state.coinz.admin", admin),
			cosmos.NewGenesisKV("app_state.coinz.asset", defaultAssetMetadata),
		}

		return cosmos.ModifyGenesis(genesis)(cfg, bz)
	}
}

// encodingConfig is used in the ChainConfig to wire up the appropriate encodings needed to decode
// msgs related to the coinz module from within interchaintest.
// This is useful for debugging output generated from tests run via interchaintest..
func encodingConfig() *testutil.TestEncodingConfig {
	c := cosmos.DefaultEncoding()

	coinztypes.RegisterInterfaces(c.InterfaceRegistry)

	return &c
}

// startDemoChain is a helper method used for configuring an interchaintest environment.
// It handles all the configuration for spinning up any number of chains given a slice of ChainSpecs.
func startDemoChain(
	ctx context.Context,
	t *testing.T,
	reporter *testreporter.RelayerExecReporter,
	chainSpecs []*interchaintest.ChainSpec,
) ([]ibc.Chain, error) {
	t.Helper()

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), chainSpecs)

	chains, err := cf.Chains(t.Name())
	if err != nil {
		return nil, err
	}

	ic := interchaintest.NewInterchain()

	for _, chain := range chains {
		ic.AddChain(chain)
	}

	client, network := interchaintest.DockerSetup(t)

	err = ic.Build(ctx, reporter, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
	})
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		err := ic.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	return chains, nil
}
