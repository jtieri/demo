package e2e

import (
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	coinztypes "github.com/jtieri/demo/x/coinz/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

var (
	adminAddress  = "cosmos1jtvahy53nvaftwq6r04dmnalv6re2av0xgp705"
	adminMnemonic = "today please woman quote finger problem patrol inflict usual crystal brick gesture myself permit clutch upset book clarify embrace mixture enact board awful almost"

	admin = coinztypes.Admin{
		Address: adminAddress,
	}

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
			Denom:          "mytoken",
			CoinType:       "118",
			GasPrices:      "0.0mytoken",
			GasAdjustment:  1.5,
			TrustingPeriod: "100h",
			PreGenesis:     nil,
			ModifyGenesis:  modifyGenesis(admin),
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
			Denom:          "mytoken",
			CoinType:       "118",
			GasPrices:      "0.0mytoken",
			GasAdjustment:  1.5,
			TrustingPeriod: "100h",
			PreGenesis:     nil,
			ModifyGenesis:  modifyGenesis(admin),
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
