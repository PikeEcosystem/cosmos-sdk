package testutil

import (
	"fmt"
	"time"

	"github.com/PikeEcosystem/cosmos-sdk/client/flags"
	"github.com/PikeEcosystem/cosmos-sdk/testutil/rest"
	sdk "github.com/PikeEcosystem/cosmos-sdk/types"
	"github.com/PikeEcosystem/cosmos-sdk/x/authz"
	"github.com/PikeEcosystem/cosmos-sdk/x/authz/client/cli"
	banktypes "github.com/PikeEcosystem/cosmos-sdk/x/bank/types"
)

func (s *IntegrationTestSuite) TestQueryGrantGRPC() {
	val := s.network.Validators[0]
	grantee := s.grantee[1]
	grantsURL := val.APIAddress + "/cosmos/authz/v1beta1/grants?granter=%s&grantee=%s&msg_type_url=%s"
	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
	}{
		{
			"fail invalid granter address",
			fmt.Sprintf(grantsURL, "invalid_granter", grantee.String(), typeMsgSend),
			true,
			"decoding bech32 failed: invalid separator index -1",
		},
		{
			"fail invalid grantee address",
			fmt.Sprintf(grantsURL, val.Address.String(), "invalid_grantee", typeMsgSend),
			true,
			"decoding bech32 failed: invalid separator index -1",
		},
		{
			"fail with empty granter",
			fmt.Sprintf(grantsURL, "", grantee.String(), typeMsgSend),
			true,
			"empty address string is not allowed: invalid request",
		},
		{
			"fail with empty grantee",
			fmt.Sprintf(grantsURL, val.Address.String(), "", typeMsgSend),
			true,
			"empty address string is not allowed: invalid request",
		},
		{
			"fail invalid msg-type",
			fmt.Sprintf(grantsURL, val.Address.String(), grantee.String(), "invalidMsg"),
			true,
			"rpc error: code = NotFound desc = no authorization found for invalidMsg type: key not found",
		},
		{
			"valid query",
			fmt.Sprintf(grantsURL, val.Address.String(), grantee.String(), typeMsgSend),
			false,
			"",
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var g authz.QueryGrantsResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &g)
				require.NoError(err)
				require.Len(g.Grants, 1)
				g.Grants[0].UnpackInterfaces(val.ClientCtx.InterfaceRegistry)
				auth := g.Grants[0].GetAuthorization()
				require.Equal(auth.MsgTypeURL(), banktypes.SendAuthorization{}.MsgTypeURL())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryGrantsGRPC() {
	val := s.network.Validators[0]
	grantee := s.grantee[1]
	grantsURL := val.APIAddress + "/cosmos/authz/v1beta1/grants?granter=%s&grantee=%s"
	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errMsg    string
		preRun    func()
		postRun   func(*authz.QueryGrantsResponse)
	}{
		{
			"valid query: expect single grant",
			fmt.Sprintf(grantsURL, val.Address.String(), grantee.String()),
			false,
			"",
			func() {},
			func(g *authz.QueryGrantsResponse) {
				s.Require().Len(g.Grants, 1)
			},
		},
		{
			"valid query: expect two grants",
			fmt.Sprintf(grantsURL, val.Address.String(), grantee.String()),
			false,
			"",
			func() {
				_, err := ExecGrant(val, []string{
					grantee.String(),
					"generic",
					fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
					fmt.Sprintf("--%s=%s", cli.FlagMsgType, typeMsgVote),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
					fmt.Sprintf("--%s=%d", cli.FlagExpiration, time.Now().Add(time.Minute*time.Duration(120)).Unix()),
				})
				s.Require().NoError(err)
			},
			func(g *authz.QueryGrantsResponse) {
				s.Require().Len(g.Grants, 2)
			},
		},
		{
			"valid query: expect single grant with pagination",
			fmt.Sprintf(grantsURL+"&pagination.limit=1", val.Address.String(), grantee.String()),
			false,
			"",
			func() {},
			func(g *authz.QueryGrantsResponse) {
				s.Require().Len(g.Grants, 1)
			},
		},
		{
			"valid query: expect two grants with pagination",
			fmt.Sprintf(grantsURL+"&pagination.limit=2", val.Address.String(), grantee.String()),
			false,
			"",
			func() {},
			func(g *authz.QueryGrantsResponse) {
				s.Require().Len(g.Grants, 2)
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			tc.preRun()
			resp, _ := rest.GetRequest(tc.url)
			if tc.expectErr {
				s.Require().Contains(string(resp), tc.errMsg)
			} else {
				var authorizations authz.QueryGrantsResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &authorizations)
				s.Require().NoError(err)
				tc.postRun(&authorizations)
			}

		})
	}
}
