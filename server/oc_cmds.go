package server

// DONTCOVER

import (
	"fmt"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	cfg "github.com/PikeEcosystem/tendermint/config"
	osjson "github.com/PikeEcosystem/tendermint/libs/json"
	ostos "github.com/PikeEcosystem/tendermint/libs/os"
	"github.com/PikeEcosystem/tendermint/node"
	"github.com/PikeEcosystem/tendermint/p2p"
	pvm "github.com/PikeEcosystem/tendermint/privval"
	"github.com/PikeEcosystem/tendermint/types"
	ostversion "github.com/PikeEcosystem/tendermint/version"

	sdk "github.com/PikeEcosystem/cosmos-sdk/types"
)

// ShowNodeIDCmd - ported from Tendermint, dump node ID to stdout
func ShowNodeIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			nodeKey, err := p2p.LoadNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			fmt.Println(nodeKey.ID())
			return nil
		},
	}
}

// ShowValidatorCmd - ported from Tendermint, show this node's validator info
func ShowValidatorCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show-validator",
		Short: "Show this node's tendermint validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config
			return showValidator(cmd, cfg)
		},
	}

	return &cmd
}

func showValidator(cmd *cobra.Command, config *cfg.Config) error {
	var pv types.PrivValidator
	if config.PrivValidatorListenAddr != "" {
		chainID, err := loadChainID(config)
		if err != nil {
			return err
		}
		serverCtx := GetServerContextFromCmd(cmd)
		log := serverCtx.Logger
		pv, err = node.CreateAndStartPrivValidatorSocketClient(config.PrivValidatorListenAddr, chainID, log)
		if err != nil {
			return err
		}
	} else {
		keyFilePath := config.PrivValidatorKeyFile()
		if !ostos.FileExists(keyFilePath) {
			return fmt.Errorf("private validator file %s does not exist", keyFilePath)
		}
		pv = pvm.LoadFilePV(keyFilePath, config.PrivValidatorStateFile())
	}

	pubKey, err := pv.GetPubKey()
	if err != nil {
		return fmt.Errorf("can't get pubkey: %w", err)
	}

	bz, err := osjson.Marshal(pubKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private validator pubkey: %w", err)
	}

	fmt.Println(string(bz))
	return nil
}

func loadChainID(config *cfg.Config) (string, error) {
	stateDB, err := node.DefaultDBProvider(&node.DBContext{ID: "state", Config: config})
	if err != nil {
		return "", err
	}
	defer func() {
		var _ = stateDB.Close()
	}()
	genesisDocProvider := node.DefaultGenesisDocProviderFunc(config)
	_, genDoc, err := node.LoadStateFromDBOrGenesisDocProvider(stateDB, genesisDocProvider)
	if err != nil {
		return "", err
	}
	return genDoc.ChainID, nil
}

// ShowAddressCmd - show this node's validator address
func ShowAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address",
		Short: "Shows this node's tendermint validator consensus address",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			privValidator := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
			valConsAddr := (sdk.ConsAddress)(privValidator.GetAddress())
			fmt.Println(valConsAddr.String())
			return nil
		},
	}

	return cmd
}

// VersionCmd prints tendermint and ABCI version numbers.
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print tendermint libraries' version",
		Long: `Print protocols' and libraries' version numbers
against which this app has been compiled.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, err := yaml.Marshal(&struct {
				Tendermint      string
				ABCI          string
				BlockProtocol uint64
				P2PProtocol   uint64
			}{
				Tendermint:      ostversion.OCCoreSemVer,
				ABCI:          ostversion.ABCIVersion,
				BlockProtocol: ostversion.BlockProtocol,
				P2PProtocol:   ostversion.P2PProtocol,
			})
			if err != nil {
				return err
			}

			fmt.Println(string(bs))
			return nil
		},
	}
}
