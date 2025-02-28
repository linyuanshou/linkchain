package rpc

import (
	"errors"
	"fmt"
	"math"
	"runtime/debug"
	"time"

	"github.com/lianxiangcloud/linkchain/accounts"
	"github.com/lianxiangcloud/linkchain/accounts/keystore"
	"github.com/lianxiangcloud/linkchain/libs/common"
)

const (
	maxCue = 64
)

// PrivateAccountAPI provides an API to access accounts managed by this node.
// It offers methods to create, (un)lock en list accounts. Some methods accept
// passwords and are therefore considered private by default.
type PrivateAccountAPI struct {
	am        *accounts.Manager
	wallet    Wallet
	nonceLock *AddrLocker
	b         Backend
}

// NewPrivateAccountAPI create a new PrivateAccountAPI.
func NewPrivateAccountAPI(b Backend, nonceLock *AddrLocker) *PrivateAccountAPI {
	return &PrivateAccountAPI{
		am:        b.AccountManager(),
		wallet:    b.GetWallet(),
		nonceLock: nonceLock,
		b:         b,
	}
}

// ListAccounts will return a list of addresses for accounts this node manages.
func (s *PrivateAccountAPI) ListAccounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty
	for _, wallet := range s.am.Wallets() {
		for _, account := range wallet.Accounts() {
			addresses = append(addresses, account.Address)
		}
	}
	return addresses
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PrivateAccountAPI) NewAccount(password string, cue string) (common.Address, error) {
	if len(password) == 0 {
		return common.EmptyAddress, fmt.Errorf("password is empty")
	}
	if len(cue) > maxCue {
		return common.EmptyAddress, fmt.Errorf("cue is too long")
	}
	acc, err := fetchKeystore(s.am).NewAccount(password, cue)
	if err == nil {
		return acc.Address, nil
	}
	return common.EmptyAddress, err
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PrivateAccountAPI) GetCue(addr common.Address) (string, error) {
	return fetchKeystore(s.am).GetCue(accounts.Account{Address: addr})
}

// fetchKeystore retrives the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) *keystore.KeyStore {
	return am.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
func (s *PrivateAccountAPI) UnlockAccount(addr common.Address, password string, duration *uint64) (bool, error) {
	const max = uint64(time.Duration(math.MaxInt64) / time.Second)
	var d time.Duration
	if duration == nil {
		d = 300 * time.Second
	} else if *duration > max {
		return false, errors.New("unlock duration too large")
	} else {
		d = time.Duration(*duration) * time.Second
	}

	// get wallet curr eth address and lock it
	// currAddr, err := s.wallet.GetWalletEthAddress()
	// if err == nil {
	// 	err = fetchKeystore(s.am).Lock(*currAddr)
	// 	log.Info("UnlockAccount lock curr addr", "currAddr", currAddr, "err", err)
	// }

	err := fetchKeystore(s.am).TimedUnlock(accounts.Account{Address: addr}, password, d)
	debug.FreeOSMemory()

	if err == nil {
		for _, wallet := range s.am.Wallets() {
			for _, account := range wallet.Accounts() {
				if account.Address == addr {
					keypath := account.URL.Path
					err = s.wallet.OpenWallet(keypath, password)
					return err == nil, err
				}
			}
		}
	}

	return err == nil, err
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (s *PrivateAccountAPI) LockAccount(addr common.Address) bool {
	s.wallet.CloseWallet()
	return fetchKeystore(s.am).Lock(addr) == nil
}
