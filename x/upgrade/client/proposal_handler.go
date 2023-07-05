package client

import (
	govclient "github.com/PikeEcosystem/cosmos-sdk/x/gov/client"
	"github.com/PikeEcosystem/cosmos-sdk/x/upgrade/client/cli"
)

var (
	ProposalHandler       = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal)
	CancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal)
)
