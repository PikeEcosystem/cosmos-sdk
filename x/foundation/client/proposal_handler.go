package client

import (
	"github.com/PikeEcosystem/cosmos-sdk/x/foundation/client/cli"
	govclient "github.com/PikeEcosystem/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewProposalCmdFoundationExec)
