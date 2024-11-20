package e2e

import (
	"context"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
)

// TestChainStart asserts that it is possible to start the demo chain via interchaintest.
func TestChainStart(t *testing.T) {
	var (
		ctx  = context.Background()
		rep  = testreporter.NewNopReporter()
		eRep = rep.RelayerExecReporter(t)
	)

	chains, err := startDemoChain(ctx, t, eRep, []*interchaintest.ChainSpec{&chainASpec})
	require.NoError(t, err)
	require.Len(t, chains, 1)

	return
}
