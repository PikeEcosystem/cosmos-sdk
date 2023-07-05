package testutil

import (
	"fmt"

	"github.com/PikeEcosystem/tendermint/libs/cli"

	"github.com/PikeEcosystem/cosmos-sdk/client"
	"github.com/PikeEcosystem/cosmos-sdk/testutil"
	clitestutil "github.com/PikeEcosystem/cosmos-sdk/testutil/cli"
	bankcli "github.com/PikeEcosystem/cosmos-sdk/x/bank/client/cli"
)

func MsgSendExec(clientCtx client.Context, from, to, amount fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{from.String(), to.String(), amount.String()}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, bankcli.NewSendTxCmd(), args)
}

func QueryBalancesExec(clientCtx client.Context, address fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{address.String(), fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, bankcli.GetBalancesCmd(), args)
}
