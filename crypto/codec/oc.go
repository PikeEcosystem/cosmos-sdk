package codec

import (
	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"

	picrypto "github.com/PikeEcosystem/tendermint/crypto"
	"github.com/PikeEcosystem/tendermint/crypto/encoding"

	"github.com/PikeEcosystem/cosmos-sdk/crypto/keys/ed25519"
	"github.com/PikeEcosystem/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/PikeEcosystem/cosmos-sdk/crypto/types"
	sdkerrors "github.com/PikeEcosystem/cosmos-sdk/types/errors"
)

// FrompiprotoPublicKey converts a OC's tmprotocrypto.PublicKey into our own PubKey.
func FrompiprotoPublicKey(protoPk tmprotocrypto.PublicKey) (cryptotypes.PubKey, error) {
	switch protoPk := protoPk.Sum.(type) {
	case *tmprotocrypto.PublicKey_Ed25519:
		return &ed25519.PubKey{
			Key: protoPk.Ed25519,
		}, nil
	case *tmprotocrypto.PublicKey_Secp256K1:
		return &secp256k1.PubKey{
			Key: protoPk.Secp256K1,
		}, nil
	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot convert %v from Tendermint public key", protoPk)
	}
}

// TopiprotoPublicKey converts our own PubKey to OC's tmprotocrypto.PublicKey.
func TopiprotoPublicKey(pk cryptotypes.PubKey) (tmprotocrypto.PublicKey, error) {
	switch pk := pk.(type) {
	case *ed25519.PubKey:
		return tmprotocrypto.PublicKey{
			Sum: &tmprotocrypto.PublicKey_Ed25519{
				Ed25519: pk.Key,
			},
		}, nil
	case *secp256k1.PubKey:
		return tmprotocrypto.PublicKey{
			Sum: &tmprotocrypto.PublicKey_Secp256K1{
				Secp256K1: pk.Key,
			},
		}, nil
	default:
		return tmprotocrypto.PublicKey{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot convert %v to Tendermint public key", pk)
	}
}

// FromOcPubKeyInterface converts OC's picrypto.PubKey to our own PubKey.
func FromOcPubKeyInterface(tmPk picrypto.PubKey) (cryptotypes.PubKey, error) {
	piprotoPk, err := encoding.PubKeyToProto(tmPk)
	if err != nil {
		return nil, err
	}

	return FrompiprotoPublicKey(piprotoPk)
}

// ToOcPubKeyInterface converts our own PubKey to OC's picrypto.PubKey.
func ToOcPubKeyInterface(pk cryptotypes.PubKey) (picrypto.PubKey, error) {
	piprotoPk, err := TopiprotoPublicKey(pk)
	if err != nil {
		return nil, err
	}

	return encoding.PubKeyFromProto(&piprotoPk)
}
