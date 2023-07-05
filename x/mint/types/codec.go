package types

import (
	"github.com/PikeEcosystem/cosmos-sdk/codec"
	cryptocodec "github.com/PikeEcosystem/cosmos-sdk/crypto/codec"
	sdk "github.com/PikeEcosystem/cosmos-sdk/types"
)

var amino = codec.NewLegacyAmino()

func init() {
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
