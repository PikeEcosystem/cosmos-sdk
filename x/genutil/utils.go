package genutil

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	cfg "github.com/PikeEcosystem/tendermint/config"
	osted25519 "github.com/PikeEcosystem/tendermint/crypto/ed25519"
	ostos "github.com/PikeEcosystem/tendermint/libs/os"
	"github.com/PikeEcosystem/tendermint/p2p"
	"github.com/PikeEcosystem/tendermint/privval"
	pitypes "github.com/PikeEcosystem/tendermint/types"
	"github.com/cosmos/go-bip39"

	cryptocodec "github.com/PikeEcosystem/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/PikeEcosystem/cosmos-sdk/crypto/types"
)

// ExportGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
func ExportGenesisFile(genDoc *pitypes.GenesisDoc, genFile string) error {
	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genFile)
}

// ExportGenesisFileWithTime creates and writes the genesis configuration to disk.
// An error is returned if building or writing the configuration to file fails.
func ExportGenesisFileWithTime(
	genFile, chainID string, validators []pitypes.GenesisValidator,
	appState json.RawMessage, genTime time.Time,
) error {
	genDoc := pitypes.GenesisDoc{
		GenesisTime: genTime,
		ChainID:     chainID,
		Validators:  validators,
		AppState:    appState,
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genFile)
}

// InitializeNodeValidatorFiles creates private validator and p2p configuration files.
func InitializeNodeValidatorFiles(config *cfg.Config) (nodeID string, valPubKey cryptotypes.PubKey, err error) {
	return InitializeNodeValidatorFilesFromMnemonic(config, "")
}

// InitializeNodeValidatorFiles creates private validator and p2p configuration files using the given mnemonic.
// If no valid mnemonic is given, a random one will be used instead.
func InitializeNodeValidatorFilesFromMnemonic(config *cfg.Config, mnemonic string) (nodeID string, valPubKey cryptotypes.PubKey, err error) {
	if len(mnemonic) > 0 && !bip39.IsMnemonicValid(mnemonic) {
		return "", nil, fmt.Errorf("invalid mnemonic")
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return "", nil, err
	}

	nodeID = string(nodeKey.ID())

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := ostos.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return "", nil, err
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := ostos.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return "", nil, err
	}

	var filePV *privval.FilePV
	if len(mnemonic) == 0 {
		filePV = privval.LoadOrGenFilePV(pvKeyFile, pvStateFile)
	} else {
		privKey := osted25519.GenPrivKeyFromSecret([]byte(mnemonic))
		filePV = privval.NewFilePV(privKey, pvKeyFile, pvStateFile)
	}

	tmValPubKey, err := filePV.GetPubKey()
	if err != nil {
		return "", nil, err
	}

	valPubKey, err = cryptocodec.FromOcPubKeyInterface(tmValPubKey)
	if err != nil {
		return "", nil, err
	}

	return nodeID, valPubKey, nil
}
