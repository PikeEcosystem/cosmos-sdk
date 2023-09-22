package teststaking

import (
	picrypto "github.com/PikeEcosystem/tendermint/crypto"
	pitypes "github.com/PikeEcosystem/tendermint/types"

	cryptocodec "github.com/PikeEcosystem/cosmos-sdk/crypto/codec"
	sdk "github.com/PikeEcosystem/cosmos-sdk/types"
	"github.com/PikeEcosystem/cosmos-sdk/x/staking/types"
)

// GetOcConsPubKey gets the validator's public key as an picrypto.PubKey.
func GetOcConsPubKey(v types.Validator) (picrypto.PubKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}

	return cryptocodec.ToOcPubKeyInterface(pk)
}

// ToOcValidator casts an SDK validator to a tendermint type Validator.
func ToOcValidator(v types.Validator, r sdk.Int) (*pitypes.Validator, error) {
	piPk, err := GetOcConsPubKey(v)
	if err != nil {
		return nil, err
	}

	return pitypes.NewValidator(piPk, v.ConsensusPower(r)), nil
}

// ToOcValidators casts all validators to the corresponding tendermint type.
func ToOcValidators(v types.Validators, r sdk.Int) ([]*pitypes.Validator, error) {
	validators := make([]*pitypes.Validator, len(v))
	var err error
	for i, val := range v {
		validators[i], err = ToOcValidator(val, r)
		if err != nil {
			return nil, err
		}
	}

	return validators, nil
}
