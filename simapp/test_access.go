package simapp

import (
	"testing"

	"github.com/PikeEcosystem/cosmos-sdk/baseapp"
	"github.com/PikeEcosystem/cosmos-sdk/client"

	"github.com/PikeEcosystem/cosmos-sdk/simapp/params"

	"github.com/PikeEcosystem/cosmos-sdk/codec"
	bankkeeper "github.com/PikeEcosystem/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/PikeEcosystem/cosmos-sdk/x/staking/keeper"
)

type TestSupport struct {
	t   testing.TB
	app *SimApp
}

func NewTestSupport(t testing.TB, app *SimApp) *TestSupport {
	return &TestSupport{t: t, app: app}
}

func (s TestSupport) AppCodec() codec.Codec {
	return s.app.appCodec
}

func (s TestSupport) StakingKeeper() stakingkeeper.Keeper {
	return s.app.StakingKeeper
}

func (s TestSupport) BankKeeper() bankkeeper.Keeper {
	return s.app.BankKeeper
}

func (s TestSupport) GetBaseApp() *baseapp.BaseApp {
	return s.app.BaseApp
}

func (s TestSupport) GetTxConfig() client.TxConfig {
	return params.MakeTestEncodingConfig().TxConfig
}
