package coinz

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/jtieri/demo/api/demo/coinz"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "AdminAddress",
					Use:            "admin-address",
					Short:          "Query admin-address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},

				{
					RpcMethod: "PauseState",
					Use:       "show-pause-state",
					Short:     "Query the modules pause state",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "UpdateAdmin",
					Use:            "update-admin [address]",
					Short:          "Send a update-admin tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "Mint",
					Use:            "mint [address] [amount]",
					Short:          "Send a mint tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "amount"}},
				},

				{
					RpcMethod:      "UpdatePauseState",
					Use:            "update-pause-state [paused]",
					Short:          "Update pause-state",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "paused"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
