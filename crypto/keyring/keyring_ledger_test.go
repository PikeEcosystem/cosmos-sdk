//go:build ledger || test_ledger_mock
// +build ledger test_ledger_mock

package keyring

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PikeEcosystem/cosmos-sdk/crypto/hd"
	sdk "github.com/PikeEcosystem/cosmos-sdk/types"
)

func TestInMemoryCreateLedger(t *testing.T) {
	kb := NewInMemory()

	ledger, err := kb.SaveLedgerKey("some_account", hd.Secp256k1, "link", 438, 3, 1)
	if err != nil {
		require.Error(t, err)
		require.Equal(t, "ledger nano S: support for ledger devices is not available in this executable", err.Error())
		require.Nil(t, ledger)
		t.Skip("ledger nano S: support for ledger devices is not available in this executable")
		return
	}

	// The mock is available, check that the address is correct
	pubKey := ledger.GetPubKey()
	expectedPkStr := "PubKeySecp256k1{038F17B38DF1EFC0714D1D3BA0AC1388C32E8B38AD87FD769BAC0B4A11DCE0EBE1}"
	require.Equal(t, expectedPkStr, pubKey.String())

	// Check that restoring the key gets the same results
	restoredKey, err := kb.Key("some_account")
	require.NoError(t, err)
	require.NotNil(t, restoredKey)
	require.Equal(t, "some_account", restoredKey.GetName())
	require.Equal(t, TypeLedger, restoredKey.GetType())
	pubKey = restoredKey.GetPubKey()
	require.Equal(t, expectedPkStr, pubKey.String())

	path, err := restoredKey.GetPath()
	require.NoError(t, err)
	require.Equal(t, "m/44'/438'/3'/0/1", path.String())
}

// TestSignVerify does some detailed checks on how we sign and validate
// signatures
func TestSignVerifyKeyRingWithLedger(t *testing.T) {
	dir := t.TempDir()

	kb, err := New("keybasename", "test", dir, nil)
	require.NoError(t, err)

	i1, err := kb.SaveLedgerKey("key", hd.Secp256k1, "link", 438, 0, 0)
	if err != nil {
		require.Equal(t, "ledger nano S: support for ledger devices is not available in this executable", err.Error())
		t.Skip("ledger nano S: support for ledger devices is not available in this executable")
		return
	}
	require.Equal(t, "key", i1.GetName())

	d1 := []byte("my first message")
	s1, pub1, err := kb.Sign("key", d1)
	require.NoError(t, err)

	s2, pub2, err := SignWithLedger(i1, d1)
	require.NoError(t, err)

	require.True(t, pub1.Equals(pub2))
	require.True(t, bytes.Equal(s1, s2))

	require.Equal(t, i1.GetPubKey(), pub1)
	require.Equal(t, i1.GetPubKey(), pub2)
	require.True(t, pub1.VerifySignature(d1, s1))
	require.True(t, i1.GetPubKey().VerifySignature(d1, s1))
	require.True(t, bytes.Equal(s1, s2))

	localInfo, _, err := kb.NewMnemonic("test", English, sdk.FullFundraiserPath, DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(t, err)
	_, _, err = SignWithLedger(localInfo, d1)
	require.Error(t, err)
	require.Equal(t, "not a ledger object", err.Error())
}

func TestAltKeyring_SaveLedgerKey(t *testing.T) {
	dir := t.TempDir()

	keyring, err := New(t.Name(), BackendTest, dir, nil)
	require.NoError(t, err)

	// Test unsupported Algo
	_, err = keyring.SaveLedgerKey("key", notSupportedAlgo{}, "link", 438, 0, 0)
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrUnsupportedSigningAlgo.Error())

	ledger, err := keyring.SaveLedgerKey("some_account", hd.Secp256k1, "link", 438, 3, 1)
	if err != nil {
		require.Equal(t, "ledger nano S: support for ledger devices is not available in this executable", err.Error())
		t.Skip("ledger nano S: support for ledger devices is not available in this executable")
		return
	}

	// The mock is available, check that the address is correct
	require.Equal(t, "some_account", ledger.GetName())
	pubKey := ledger.GetPubKey()
	expectedPkStr := "PubKeySecp256k1{038F17B38DF1EFC0714D1D3BA0AC1388C32E8B38AD87FD769BAC0B4A11DCE0EBE1}"
	require.Equal(t, expectedPkStr, pubKey.String())

	// Check that restoring the key gets the same results
	restoredKey, err := keyring.Key("some_account")
	require.NoError(t, err)
	require.NotNil(t, restoredKey)
	require.Equal(t, "some_account", restoredKey.GetName())
	require.Equal(t, TypeLedger, restoredKey.GetType())
	pubKey = restoredKey.GetPubKey()
	require.Equal(t, expectedPkStr, pubKey.String())

	path, err := restoredKey.GetPath()
	require.NoError(t, err)
	require.Equal(t, "m/44'/438'/3'/0/1", path.String())
}
