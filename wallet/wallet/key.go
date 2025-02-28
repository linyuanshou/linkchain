package wallet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/lianxiangcloud/linkchain/accounts/keystore"
	"github.com/lianxiangcloud/linkchain/libs/crypto"
	lktypes "github.com/lianxiangcloud/linkchain/libs/cryptonote/types"
	"github.com/lianxiangcloud/linkchain/libs/cryptonote/xcrypto"
	"github.com/lianxiangcloud/linkchain/libs/log"
	"github.com/lianxiangcloud/linkchain/wallet/types"
)

var (
	ErrStrToAddressInvalid  = errors.New("str to address invalid")
	ErrStrToAddressCheckSum = errors.New("str to address check sum fail")
)

//WordsToKey Converts seed words to bytes (secret key)
func WordsToKey(words string) (lktypes.SecretKey, error) {
	return xcrypto.WordsToBytes(words)
}

//KeyToWords Converts bytes (secret key) to seed words
func KeyToWords(key lktypes.SecretKey) (string, error) {
	return xcrypto.BytesToWords(key, "English")
}

//GenerateKeys return spend secret key and spend public key
func GenerateKeys(key lktypes.SecretKey) (lktypes.SecretKey, lktypes.PublicKey) {
	return xcrypto.GenerateKeys(key)
}

//WordsToAccount recovery utxo account from the words
func WordsToAccount(words string) (*AccountBase, error) {
	rk, err := WordsToKey(words)
	if err != nil {
		return nil, err
	}
	log.Debug("WordsToAccount", "rk", fmt.Sprintf("%x", rk))
	return RecoveryKeyToAccount(rk)
}

//RecoveryKeyToAccount recovery utxo account from the key
func RecoveryKeyToAccount(rk lktypes.SecretKey) (*AccountBase, error) {
	spendSK, spendPK := GenerateKeys(rk)
	hash := crypto.Keccak256(spendSK[:])
	var rk2 lktypes.SecretKey
	copy(rk2[:], hash)
	viewSK, viewPK := GenerateKeys(rk2)
	acc := lktypes.AccountKey{
		Addr: lktypes.AccountAddress{
			SpendPublicKey: spendPK,
			ViewPublicKey:  viewPK,
		},
		SpendSKey: spendSK,
		ViewSKey:  viewSK,
		SubIdx:    uint64(0),
	}

	address := AddressToStr(&acc, uint64(0))
	acc.Address = address

	ab := AccountBase{
		KeyIndex:          make(map[lktypes.PublicKey]uint64),
		CreationTimestamp: time.Now().Unix(),
		MainSecKey:        rk,
	}
	ab.Keys = append(ab.Keys, &acc)
	ab.KeyIndex[acc.Addr.SpendPublicKey] = 0
	ab.CurrIdx = uint64(0)
	return &ab, nil
}

//GetSubaddr return a subaddr
func GetSubaddr(key *lktypes.AccountKey, index uint64) string {
	//TODO put spendPK into AccountKey.KeyIndex map
	return AddressToStr(key, index)
}

//AddressToStr return prefix + spend public key + view public key + checksum
func AddressToStr(key *lktypes.AccountKey, index uint64) string {
	//main address
	if index == 0 {
		return addressToStr(uint64(types.GetConfig().CRYPTONOTE_PUBLIC_ADDRESS_BASE58_PREFIX), key.Addr)
	}
	//subaddress
	addr := xcrypto.GetSubaddress(key, uint32(index))
	return addressToStr(uint64(types.GetConfig().CRYPTONOTE_PUBLIC_SUBADDRESS_BASE58_PREFIX), addr)
}

func addressToStr(prefix uint64, addr lktypes.AccountAddress) string {
	addrLen := types.GetConfig().CRYPTONOTE_PREFIX_LENGTH + 2*types.GetConfig().CRYPTONOTE_ADDRESS_LENGTH + types.GetConfig().CRYPTONOTE_CHECKSUM_LENGTH
	idx := 0
	buff := make([]byte, addrLen)
	binary.PutUvarint(buff, uint64(prefix))
	idx += types.GetConfig().CRYPTONOTE_PREFIX_LENGTH
	copy(buff[idx:], addr.SpendPublicKey[:])
	idx += types.GetConfig().CRYPTONOTE_ADDRESS_LENGTH
	copy(buff[idx:], addr.ViewPublicKey[:])
	idx += types.GetConfig().CRYPTONOTE_ADDRESS_LENGTH
	hash := xcrypto.FastHash(buff[:idx])
	checksum := hash[:types.GetConfig().CRYPTONOTE_CHECKSUM_LENGTH]
	copy(buff[idx:], checksum)
	idx += types.GetConfig().CRYPTONOTE_CHECKSUM_LENGTH
	str := xcrypto.Base58Encode(buff)
	return str
}

//StrToAddress parse address str and return utxo address
func StrToAddress(str string) (*lktypes.AccountAddress, error) {
	addrLen := types.GetConfig().CRYPTONOTE_PREFIX_LENGTH + 2*types.GetConfig().CRYPTONOTE_ADDRESS_LENGTH + types.GetConfig().CRYPTONOTE_CHECKSUM_LENGTH
	data := xcrypto.Base58Decode(str)
	if len(data) != addrLen {
		return nil, ErrStrToAddressInvalid
	}
	checksum := data[addrLen-types.GetConfig().CRYPTONOTE_CHECKSUM_LENGTH:]
	data = data[:addrLen-types.GetConfig().CRYPTONOTE_CHECKSUM_LENGTH]
	hash := xcrypto.FastHash(data)
	expectsum := hash[:types.GetConfig().CRYPTONOTE_CHECKSUM_LENGTH]
	if !bytes.Equal(checksum, expectsum) {
		return nil, ErrStrToAddressCheckSum
	}
	prefix, n := binary.Uvarint(data)
	if n != types.GetConfig().CRYPTONOTE_PREFIX_LENGTH {
		return nil, ErrStrToAddressInvalid
	}
	if prefix != uint64(types.GetConfig().CRYPTONOTE_PUBLIC_ADDRESS_BASE58_PREFIX) &&
		prefix != uint64(types.GetConfig().CRYPTONOTE_PUBLIC_SUBADDRESS_BASE58_PREFIX) {
		return nil, ErrStrToAddressInvalid
	}
	data = data[types.GetConfig().CRYPTONOTE_PREFIX_LENGTH:]
	var addr lktypes.AccountAddress
	copy(addr.SpendPublicKey[:], data[:types.GetConfig().CRYPTONOTE_ADDRESS_LENGTH])
	copy(addr.ViewPublicKey[:], data[types.GetConfig().CRYPTONOTE_ADDRESS_LENGTH:])
	return &addr, nil
}

//KeyFromAccount return secret key from the account keystore file. we use this key as the recovery key of utxo
func KeyFromAccount(keyjson []byte, passwd string) (lktypes.SecretKey, error) {
	accKey, err := keystore.DecryptKey(keyjson, passwd)
	if err != nil {
		return lktypes.SecretKey{}, err
	}
	sk := crypto.FromECDSA(accKey.PrivateKey)
	var key lktypes.SecretKey
	copy(key[:], sk[:])
	log.Debug("KeyFromAccount", "key", fmt.Sprintf("%x", key))
	return key, nil
}

func IsSubaddress(str string) bool {
	var prefix uint64
	// var blob []byte
	xcrypto.Base58DecodeAddr(&prefix, str)
	log.Debug("IsSubaddress", "addr", str, "prefix", prefix)
	if types.GetConfig().CRYPTONOTE_PUBLIC_SUBADDRESS_BASE58_PREFIX-int(prefix) == 0 {
		log.Debug("IsSubaddress", "Is Subaddress", str)
		return true
	}
	log.Debug("IsSubaddress", "Not Subaddress", str)
	return false
}
