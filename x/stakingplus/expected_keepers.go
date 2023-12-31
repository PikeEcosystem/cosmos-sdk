package stakingplus

import (
	sdk "github.com/PikeEcosystem/cosmos-sdk/types"
)

// FoundationKeeper defines the expected foundation keeper
type FoundationKeeper interface {
	Accept(ctx sdk.Context, grantee sdk.AccAddress, msg sdk.Msg) error
}
