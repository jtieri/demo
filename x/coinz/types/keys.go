package types

const (
	// ModuleName defines the module name
	ModuleName = "coinz"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_coinz"
)

var (
	ParamsKey = []byte("p_coinz")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
