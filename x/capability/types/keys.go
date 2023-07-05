package types

import (
	"fmt"

	sdk "github.com/PikeEcosystem/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "capability"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"
)

var (
	// KeyIndex defines the key that stores the current globally unique capability
	// index.
	KeyIndex = []byte("index")

	// KeyPrefixIndexCapability defines a key prefix that stores index to capability
	// name mappings.
	KeyPrefixIndexCapability = []byte("capability_index")

	// KeyMemInitialized defines the key that stores the initialized flag in the memory store
	KeyMemInitialized = []byte("mem_initialized")
)

// RevCapabilityKey returns a reverse lookup key for a given module and capability
// name.
func RevCapabilityKey(module, name string) []byte {
	return []byte(fmt.Sprintf("%s/rev/%s", module, name))
}

// FwdCapabilityKey returns a forward lookup key for a given module and capability
// reference.
func FwdCapabilityKey(module string, cap *Capability) []byte {
	// encode the key to a fixed length to avoid breaking consensus state machine
	// it's a hacky backport of https://github.com/cosmos/cosmos-sdk/pull/11737
	// the length 10 is picked so it's backward compatible on common architectures.
	key := fmt.Sprintf("%#010p", cap)
	if len(key) > 10 {
		key = key[len(key)-10:]
	}
	return []byte(fmt.Sprintf("%s/fwd/0x%s", module, key))
}

// IndexToKey returns bytes to be used as a key for a given capability index.
func IndexToKey(index uint64) []byte {
	return sdk.Uint64ToBigEndian(index)
}

// IndexFromKey returns an index from a call to IndexToKey for a given capability
// index.
func IndexFromKey(key []byte) uint64 {
	return sdk.BigEndianToUint64(key)
}
