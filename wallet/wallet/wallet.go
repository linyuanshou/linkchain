package wallet

import (
	"math/big"
	"sync"
	"time"

	"github.com/lianxiangcloud/linkchain/accounts"
	"github.com/lianxiangcloud/linkchain/libs/common"
	lkctypes "github.com/lianxiangcloud/linkchain/libs/cryptonote/types"
	dbm "github.com/lianxiangcloud/linkchain/libs/db"
	"github.com/lianxiangcloud/linkchain/libs/hexutil"
	"github.com/lianxiangcloud/linkchain/libs/log"
	"github.com/lianxiangcloud/linkchain/libs/ser"
	tctypes "github.com/lianxiangcloud/linkchain/types"
	cfg "github.com/lianxiangcloud/linkchain/wallet/config"
	"github.com/lianxiangcloud/linkchain/wallet/types"
	wtypes "github.com/lianxiangcloud/linkchain/wallet/types"
)

const (
	defaultUTXOGas                = 0x7a120
	defaultRefreshUTXOGasInterval = 60 * time.Second
)

func init() {
	tctypes.RegisterUTXOTxData()
	ser.RegisterInterface((*tctypes.Input)(nil), nil)
}

// Wallet user wallet
type Wallet struct {
	// cmn.BaseService

	Logger   log.Logger
	lock     sync.Mutex
	walletDB dbm.DB
	utxoGas  *big.Int

	// config
	config *cfg.Config

	addrMap     map[common.Address]*LinkAccount
	currAccount *LinkAccount // latest unlock account

	accManager *accounts.Manager
}

// NewWallet returns a new, ready to go.
func NewWallet(config *cfg.Config,
	logger log.Logger, db dbm.DB, accManager *accounts.Manager) (*Wallet, error) {

	wallet := &Wallet{
		config:     config,
		walletDB:   db,
		accManager: accManager,
		addrMap:    make(map[common.Address]*LinkAccount),
	}
	wallet.utxoGas = new(big.Int).Mul(new(big.Int).SetUint64(defaultUTXOGas), new(big.Int).SetInt64(tctypes.ParGasPrice))

	// wallet.BaseService = *cmn.NewBaseService(logger, "Wallet", wallet)
	wallet.Logger = logger

	return wallet, nil
}

// OpenWallet ,open wallet with password
func (w *Wallet) OpenWallet(keystoreFile string, password string) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	la, err := NewLinkAccount(w.walletDB, w.Logger, keystoreFile, password)
	if err != nil {
		w.Logger.Error("OpenWallet NewLinkAccount fail", "err", err)
		return err
	}
	addr := la.getEthAddress()

	w.Logger.Info("OpenWallet", "address", addr)

	laOld, ok := w.addrMap[addr]
	if ok {
		w.currAccount = laOld

		return nil
	}

	w.addrMap[addr] = la
	w.currAccount = la

	// default start account refresh
	la.OnStart()

	return nil
}

// IsWalletClosed return true if currAccount is nil
func (w *Wallet) IsWalletClosed() bool {
	return w.currAccount == nil
}

// OnStart starts the Wallet. It implements cmn.Service.
func (w *Wallet) Start() error {
	w.Logger.Info("starting Wallet OnStart")
	go w.refreshUTXOGas()
	return nil
}

// OnStop stops the Wallet. It implements cmn.Service.
func (w *Wallet) Stop() {

	for _, account := range w.addrMap {
		w.Logger.Info("OnStop", "addr", account.getEthAddress())

		account.OnStop()
	}

	w.Logger.Info("Stopping Wallet")
}

func (w *Wallet) refreshUTXOGas() {
	w.Logger.Debug("refreshUTXOGas")
	refresh := time.NewTicker(defaultRefreshUTXOGasInterval)
	defer refresh.Stop()

	for {
		select {
		case <-refresh.C:
			utxoGas, err := w.getUTXOGas()
			if err != nil {
				w.Logger.Error("refreshUTXOGas", "err", err)
				continue
			}
			newUtxoGas := new(big.Int).Mul(new(big.Int).SetUint64(utxoGas), new(big.Int).SetInt64(tctypes.ParGasPrice))
			w.utxoGas.Set(newUtxoGas)
			w.Logger.Info("refreshUTXOGas set utxoGas", "utxoGas", w.utxoGas.String())
		}
	}
}

// GetBalance rpc get balance
func (w *Wallet) GetBalance(index uint64, token *common.Address) (*big.Int, error) {
	if w.IsWalletClosed() {
		return nil, wtypes.ErrWalletNotOpen
	}
	return w.currAccount.GetBalance(index, token), nil
}

// GetAddress rpc get address
func (w *Wallet) GetAddress(index uint64) (string, error) {
	if w.IsWalletClosed() {
		return "", wtypes.ErrWalletNotOpen
	}
	return w.currAccount.GetAddress(index)
}

// GetHeight rpc get height
func (w *Wallet) GetHeight() (localHeight uint64, remoteHeight uint64) {
	if w.IsWalletClosed() {
		return 0, 0
	}
	return w.currAccount.GetHeight()
}

// CreateSubAccount return new sub address and sub index
func (w *Wallet) CreateSubAccount(maxSub uint64) error {
	if w.IsWalletClosed() {
		return wtypes.ErrWalletNotOpen
	}
	return w.currAccount.CreateSubAccount(maxSub)
}

// AutoRefreshBlockchain set autoRefresh
func (w *Wallet) AutoRefreshBlockchain(autoRefresh bool) error {
	if w.IsWalletClosed() {
		return wtypes.ErrWalletNotOpen
	}
	return w.currAccount.AutoRefreshBlockchain(autoRefresh)
}

// GetAccountInfo return eth_account and utxo_accounts
func (w *Wallet) GetAccountInfo(tokenID *common.Address) (*types.GetAccountInfoResult, error) {
	if w.IsWalletClosed() {
		return nil, wtypes.ErrWalletNotOpen
	}
	return w.currAccount.GetAccountInfo(tokenID)
}

// RescanBlockchain ,reset wallet block and transfer info
func (w *Wallet) RescanBlockchain() error {
	if w.IsWalletClosed() {
		return wtypes.ErrWalletNotOpen
	}
	return w.currAccount.RescanBlockchain()
}

// GetWalletEthAddress ,return wallet eth address
func (w *Wallet) GetWalletEthAddress() (*common.Address, error) {
	if w.IsWalletClosed() {
		return nil, wtypes.ErrWalletNotOpen
	}
	addr := w.currAccount.getEthAddress()
	return &addr, nil
}

// CloseWallet ,close wallet
func (w *Wallet) CloseWallet() error {
	w.currAccount.OnStop()

	addr := w.currAccount.getEthAddress()
	delete(w.addrMap, addr)
	if len(w.addrMap) == 0 {
		w.currAccount = nil
	}
	for _, v := range w.addrMap {
		w.currAccount = v
		w.Logger.Info("CloseWallet reset currAccount", "currAccount", w.currAccount.getEthAddress())
		break
	}

	return nil
}

// getGOutIndex return curr idx
func (w *Wallet) getGOutIndex(token common.Address) uint64 {
	if w.IsWalletClosed() {
		return 0
	}
	return w.currAccount.GetGOutIndex(token)
}

// Status return wallet status
func (w *Wallet) Status() *types.StatusResult {
	if w.IsWalletClosed() {
		return nil
	}
	return w.currAccount.Status()
}

// GetTxKey return transaction's tx secKey
func (w *Wallet) GetTxKey(hash *common.Hash) (*lkctypes.Key, error) {
	if w.IsWalletClosed() {
		return nil, wtypes.ErrWalletNotOpen
	}
	return w.currAccount.GetTxKey(hash)
}

// GetMaxOutput return tokenID max output idx
func (w *Wallet) GetMaxOutput(tokenID common.Address) (*hexutil.Uint64, error) {
	if w.IsWalletClosed() {
		return nil, wtypes.ErrWalletNotOpen
	}
	return w.currAccount.GetMaxOutput(tokenID)
}

// GetUTXOTx return UTXOTransaction
func (w *Wallet) GetUTXOTx(hash common.Hash) (*tctypes.UTXOTransaction, error) {
	if w.IsWalletClosed() {
		return nil, wtypes.ErrWalletNotOpen
	}
	return w.currAccount.GetUTXOTx(hash)
}

// SelectAddress return
func (w *Wallet) SelectAddress(addr common.Address) error {
	la, ok := w.addrMap[addr]
	if !ok {
		return wtypes.ErrWalletNotOpen
	}
	w.currAccount = la

	w.Logger.Info("SelectAddress", "address", la.getEthAddress())

	return nil
}

// SetRefreshBlockInterval return
func (w *Wallet) SetRefreshBlockInterval(interval time.Duration) error {
	if w.IsWalletClosed() {
		return wtypes.ErrWalletNotOpen
	}
	w.currAccount.SetRefreshBlockInterval(interval)
	return nil
}
