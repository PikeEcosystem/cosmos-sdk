package stakingplus

import (
	"github.com/PikeEcosystem/cosmos-sdk/codec"
	"github.com/PikeEcosystem/cosmos-sdk/codec/types"
	authzcodec "github.com/PikeEcosystem/cosmos-sdk/x/authz/codec"
	"github.com/PikeEcosystem/cosmos-sdk/x/foundation"
	fdncodec "github.com/PikeEcosystem/cosmos-sdk/x/foundation/codec"
	govcodec "github.com/PikeEcosystem/cosmos-sdk/x/gov/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/authz interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&CreateValidatorAuthorization{}, "lbm-sdk/CreateValidatorAuthorization", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*foundation.Authorization)(nil),
		&CreateValidatorAuthorization{},
	)
}

func init() {
	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(fdncodec.Amino)
}
