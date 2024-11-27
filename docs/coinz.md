# coinz module

## Overview

The `coinz` module is a Go implementation of [Ethereum based Solidity contracts](https://basescan.org/address/0xb488b03971bf796792253c2d033f16482015f76c#code)
that define a mint and burn mechanism for a token. It was built using Ignite CLI so is powered by the Cosmo SDK.

The `coinz` module exposes methods for burning and minting a specific asset that is defined at genesis.

Minting new assets is a permissioned behavior that can only be performed by an Admin account, which must also be specified at genesis.

The coinz module can also be paused, which will disable the ability for tokens to be minted or burned.
Pausing or unpausing the module can only be done by the Admin account.

## Details

The `Admin` and `AssetMetadata` *MUST* be fully specified during genesis since the module does not
currently expose methods for initializing these values after the chain has started.

The `PauseState` must also be initialized at genesis, since we also do not expose methods for setting 
the state if it is not initialized.

### Types

---

#### Admin

`Admin` represents the authority address that is capable of performing permissioned behaviors
such as minting new tokens, updating the admin, and pausing the module.

The `Admin` *MUST* be initialized via the `genesis.json` file. Once initialized, the `Admin` can be updated to a new address via
`MsgUpdateAdmin` by the existing admin account.

```protobuf
message Admin {
  string address = 1; 
}
```

#### AssetMetadata

`AssetMetadata` represents the asset that the `coinz` module is able to mint and burn.
The `coinz` module will only be able to mint and burn a single asset. The asset amount
should be initialized to the initial total supply of the asset that will be minted at chain start.

The `AssetMetadata` *MUST* be initialized via the `genesis.json` file. Once initialized, there is no way to currently update the
`AssetMetadata`.

```protobuf
message AssetMetadata {
  cosmos.base.v1beta1.Coin asset = 1 [(gogoproto.nullable) = false];
}
```

### Msgs

---

#### MsgMint

`MsgMint` is used to mint new assets.  

If the from address does not equal the admin address the msg will fail to execute.
If the asset denom is not equal to the denom registered at genesis the msg will fail to execute.
Upon successful execution, new assets equal to the specified amount will be minted and transferred 
to the specified address.

```protobuf
message MsgMint {
  option (cosmos.msg.v1.signer) = "from";
  string                   from    = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string                   address = 2 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  cosmos.base.v1beta1.Coin amount  = 3 [(gogoproto.nullable)  = false                 ];
}
```

#### MsgBurn

`MsgBurn` is used to burn existing assets.

The msg can be executed by any account. If the asset denom is not equal to the denom registered at genesis the msg will fail to execute.
If the signers account does not have adequate funds to burn it will fail to execute. Upon successful
execution, the assets will be moved from the user account to the module account and burned, which will decrease the 
assets total supply.

```protobuf
message MsgBurn {
  option (cosmos.msg.v1.signer) = "from";
  string                   from   = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  cosmos.base.v1beta1.Coin amount = 2 [(gogoproto.nullable) = false];
}
```

#### MsgUpdateAdmin

`MsgUpdateAdmin` is used to update the `Admin` address.  

The msg will fail to execute if the from address does not equal the current `Admin` address.  
If successful, the `Admin` address will be updated to the specified address.

```protobuf
message MsgUpdateAdmin {
  option (cosmos.msg.v1.signer) = "from";
  string from    = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string address = 2 [(cosmos_proto.scalar) = "cosmos.AddressString"];
}
```

#### MsgUpdatePauseState

`MsgUpdatePauseState` is used to update the `PauseState`.
If the from address does not equal the `Admin` address the msg will fail to execute. 
If the specified state to update to is already the current state then a no-op occurs.

```protobuf
message MsgUpdatePauseState {
  option (cosmos.msg.v1.signer) = "from";
  string from   = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  bool   paused = 2;
}
```

### Queries

---

```protobuf
// QueryAdminAddressRequest is used to query the Admin address from state.
message QueryAdminAddressRequest {}

// QueryAdminAddressResponse is sent as a response to QueryAdminAddressRequest with the Admins address.
message QueryAdminAddressResponse {
  string address = 1;
}
```

```protobuf
// QueryGetPauseStateRequest is used to query the PauseState from state.
message QueryGetPauseStateRequest {}

// QueryGetPauseStateResponse is sent as a response to QueryGetPauseStateRequest with the PauseState value.
message QueryGetPauseStateResponse {
  bool paused = 1 [(gogoproto.jsontag) = "paused"];
}
```

## Notes

---

The coinz module should be able to have a more sophisticated authority schema, such that there could
be many different types of roles that are scoped to specific behaviors in the module.

e.g. `PauseAuthority`, `MintAuthority`, etc.

Not only that but the original Solidity implementation contains the concept of delayed actions, which are not
implemented in the coinz module.

Also, in the original Solidity contracts pausing the system should also pause token transfers. The `coinz` module
does not currently check if the module is paused for token transfers so they can still be performed even when the module
is paused. This could be implemented via an Ante handler that is executed in `CheckTx` when performing other
stateful validations of proposed txs.

The last missing feature from the Solidity contracts is the presence of a `BurnFrom` function
that would allow the burning of funds from an account that is different than the one that executed the tx.
This involves the use of allowances, which create the ability for one account to perform actions on behalf of another account.
This could be implemented in the `coinz` module via `authz`

These implementation details were outside the scope of the two-hour time boxed exercise and so the decision to use
a single `Admin` address for all of these roles was used instead and the other mentioned features were left out for the
initial implementation in order to focus on implementing the core features, writing tests, and ensuring sufficient
documentation existed.
