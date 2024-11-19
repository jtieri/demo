package keeper

import (
	"github.com/jtieri/demo/x/coinz/types"
)

var _ types.QueryServer = Keeper{}
