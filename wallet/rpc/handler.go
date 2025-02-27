package rpc

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lianxiangcloud/linkchain/libs/common"
	lkctypes "github.com/lianxiangcloud/linkchain/libs/cryptonote/types"
	"github.com/lianxiangcloud/linkchain/libs/cryptonote/xcrypto"
	"github.com/lianxiangcloud/linkchain/libs/hexutil"
	"github.com/lianxiangcloud/linkchain/libs/log"
	"github.com/lianxiangcloud/linkchain/libs/ser"
	"github.com/lianxiangcloud/linkchain/types"
	wtypes "github.com/lianxiangcloud/linkchain/wallet/types"
	"github.com/lianxiangcloud/linkchain/wallet/wallet"
)

// PublicTransactionPoolAPI exposes methods for the RPC interface
type PublicTransactionPoolAPI struct {
	b         Backend
	nonceLock *AddrLocker
	wallet    Wallet
}

// NewPublicTransactionPoolAPI creates a new RPC service with methods specific for the transaction pool.
func NewPublicTransactionPoolAPI(b Backend, nonceLock *AddrLocker) *PublicTransactionPoolAPI {
	return &PublicTransactionPoolAPI{b, nonceLock, b.GetWallet()}
}

func (s *PublicTransactionPoolAPI) signUTXOTransaction(ctx context.Context, args wtypes.SendUTXOTxArgs) (*wtypes.SignUTXOTransactionResult, error) {
	args.SetDefaults()

	log.Debug("signTx", "input", args)
	//utxo not support token, next version will support
	if *args.TokenID != common.EmptyAddress {
		return nil, wtypes.ErrUTXONotSupportToken
	}
	destsCnt := len(args.Dests)
	// tosCnt := len(args.Tos)
	if destsCnt == 0 {
		return nil, fmt.Errorf("need more dests")
	}

	dests := make([]types.DestEntry, 0)
	hasOneAccountOutput := false
	utxoDestsCnt := 0
	for i := 0; i < destsCnt; i++ {
		toAddress := args.Dests[i].Addr
		if len(toAddress) == 95 {
			if utxoDestsCnt >= wtypes.UTXO_DESTS_MAX_NUM {
				return nil, wtypes.ErrUTXODestsOverLimit
			}
			// utxo address
			addr, err := wallet.StrToAddress(args.Dests[i].Addr)
			if err != nil {
				return nil, fmt.Errorf("error dests addr:%s", toAddress)
			}

			var remark [32]byte
			copy(remark[:], args.Dests[i].Remark[:])
			log.Debug("signUTXOTransaction", "Remark", args.Dests[i].Remark, "len", len(args.Dests[i].Remark), "remark", remark)

			dests = append(dests, &types.UTXODestEntry{Addr: *addr, Amount: args.Dests[i].Amount.ToInt(), IsSubaddress: wallet.IsSubaddress(args.Dests[i].Addr), Remark: remark})
			utxoDestsCnt++
		} else {
			if !common.IsHexAddress(toAddress) {
				return nil, fmt.Errorf("error dests addr:%s", toAddress)
			}
			if hasOneAccountOutput {
				// can not sign more than one account output
				return nil, fmt.Errorf("account output too more")
			}
			addr := common.HexToAddress(toAddress)
			dests = append(dests, &types.AccountDestEntry{To: addr, Amount: args.Dests[i].Amount.ToInt(), Data: args.Dests[i].Data})
			hasOneAccountOutput = true
		}

	}

	txs, err := s.wallet.CreateUTXOTransaction(args.From, uint64(*args.Nonce), args.SubAddrs, dests, *args.TokenID, args.From, nil)
	if err != nil {
		return nil, err
	}

	var signedtxs []wtypes.SignUTXORet
	for _, tx := range txs {
		bz, err := ser.EncodeToBytes(tx)
		if err != nil {
			return nil, err
		}

		keys := tx.GetInputKeyImages()
		for i := 0; i < len(keys); i++ {
			log.Debug("signUTXOTransaction", "args.SubAddrs", args.SubAddrs, "keyimage", keys[i])
		}
		gas := hexutil.Uint64(tx.Gas())
		signedtxs = append(signedtxs, wtypes.SignUTXORet{Raw: fmt.Sprintf("0x%s", hex.EncodeToString(bz)), Hash: tx.Hash(), Gas: gas})
	}
	return &wtypes.SignUTXOTransactionResult{Txs: signedtxs}, nil
}

// SignUTXOTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (s *PublicTransactionPoolAPI) SignUTXOTransaction(ctx context.Context, args wtypes.SendUTXOTxArgs) (*wtypes.SignUTXOTransactionResult, error) {
	return s.signUTXOTransaction(ctx, args)
}

// SendUTXOTransaction send utxo tx
func (s *PublicTransactionPoolAPI) SendUTXOTransaction(ctx context.Context, args wtypes.SendUTXOTxArgs) (*wtypes.SendUTXOTransactionResult, error) {
	signRet, err := s.signUTXOTransaction(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(signRet.Txs) > 1 {
		return nil, fmt.Errorf("Transaction would be too large.  try transfer_split")
	}

	ret := s.wallet.Transfer([]string{signRet.Txs[0].Raw})
	ret[0].Gas = signRet.Txs[0].Gas

	return &wtypes.SendUTXOTransactionResult{Txs: ret}, nil
}

// SendUTXOTransactionSplit send utxo tx
func (s *PublicTransactionPoolAPI) SendUTXOTransactionSplit(ctx context.Context, args wtypes.SendUTXOTxArgs) (*wtypes.SendUTXOTransactionResult, error) {
	signRet, err := s.signUTXOTransaction(ctx, args)
	if err != nil {
		return nil, err
	}
	var signedRaw []string
	for i := 0; i < len(signRet.Txs); i++ {
		signedRaw = append(signedRaw, signRet.Txs[i].Raw)
	}
	ret := s.wallet.Transfer(signedRaw)
	for index := 0; index < len(signRet.Txs); index++ {
		ret[index].Gas = signRet.Txs[index].Gas
	}
	return &wtypes.SendUTXOTransactionResult{Txs: ret}, nil
}

// BlockHeight get block height
func (s *PublicTransactionPoolAPI) BlockHeight(ctx context.Context) (*wtypes.BlockHeightResult, error) {
	localHeight, remoteHeight := s.wallet.GetHeight()

	return &wtypes.BlockHeightResult{LocalHeight: hexutil.Uint64(localHeight), RemoteHeight: hexutil.Uint64(remoteHeight)}, nil
}

// Balance get account Balance
func (s *PublicTransactionPoolAPI) Balance(ctx context.Context, args wtypes.BalanceArgs) (*wtypes.BalanceResult, error) {
	if args.TokenID == nil {
		args.TokenID = &common.EmptyAddress
	}

	address, err := s.wallet.GetAddress(uint64(args.AccountIndex))
	if err != nil {
		return nil, err
	}
	balance, err := s.wallet.GetBalance(uint64(args.AccountIndex), args.TokenID)
	if err != nil {
		return nil, err
	}

	return &wtypes.BalanceResult{Balance: (*hexutil.Big)(balance), Address: address, TokenID: args.TokenID}, err
}

// CreateSubAccount create sub account to max sub index
func (s *PublicTransactionPoolAPI) CreateSubAccount(ctx context.Context, maxSub hexutil.Uint64) (bool, error) {
	err := s.wallet.CreateSubAccount(uint64(maxSub))
	if err != nil {
		return false, err
	}
	return true, nil
}

// Balance get account Balance
func (s *PublicTransactionPoolAPI) AutoRefreshBlockchain(ctx context.Context, autoRefresh bool) (bool, error) {
	err := s.wallet.AutoRefreshBlockchain(autoRefresh)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetAccountInfo get all sub accounts Balance
func (s *PublicTransactionPoolAPI) GetAccountInfo(ctx context.Context, tokenID *common.Address) (*wtypes.GetAccountInfoResult, error) {
	if tokenID == nil {
		tokenID = &common.EmptyAddress
	}
	ret, err := s.wallet.GetAccountInfo(tokenID)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// RescanBlockchain reset wallet block and transfer info
func (s *PublicTransactionPoolAPI) RescanBlockchain(ctx context.Context) (bool, error) {
	err := s.wallet.RescanBlockchain()
	if err != nil {
		return false, err
	}
	return true, nil
}

// Status return wallet status
func (s *PublicTransactionPoolAPI) Status(ctx context.Context) (*wtypes.StatusResult, error) {
	status := s.wallet.Status()
	return status, nil
}

// GetTxKey return tx key
func (s *PublicTransactionPoolAPI) GetTxKey(ctx context.Context, hash common.Hash) (*lkctypes.Key, error) {
	return s.wallet.GetTxKey(&hash)
}

// CheckTxKey Check a transaction in the blockchain with its secret key.
// func (s *PublicTransactionPoolAPI) CheckTxKey(ctx context.Context, args wtypes.CheckTxKeyArgs) (*wtypes.CheckTxKeyResult, error) {
// 	block, amount, err := s.wallet.CheckTxKey(&args.TxHash, &args.TxKey, args.DestAddr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &wtypes.CheckTxKeyResult{BlockID: *block, Amount: amount}, err
// }

// GetMaxOutput return max output
func (s *PublicTransactionPoolAPI) GetMaxOutput(ctx context.Context, tokenID common.Address) (*hexutil.Uint64, error) {
	return s.wallet.GetMaxOutput(tokenID)
}

// GetProofKey return proof key
func (s *PublicTransactionPoolAPI) GetProofKey(ctx context.Context, args wtypes.ProofKeyArgs) (*wtypes.ProofKeyRet, error) {
	tx, err := s.wallet.GetUTXOTx(args.Hash)
	if err != nil {
		return nil, err
	}
	addr, err := wallet.StrToAddress(args.Addr)
	if err != nil {
		return nil, err
	}
	txKey, err := s.wallet.GetTxKey(&args.Hash)
	if err != nil {
		return nil, err
	}
	derivationKey, err := xcrypto.GenerateKeyDerivation(addr.ViewPublicKey, lkctypes.SecretKey(*txKey))
	if err != nil {
		return nil, err
	}
	outIdx := 0
	for _, output := range tx.Outputs {
		if utxoOutput, ok := output.(*types.UTXOOutput); ok {
			otAddr, err := xcrypto.DerivePublicKey(derivationKey, outIdx, addr.SpendPublicKey)
			if err != nil {
				return nil, err
			}
			if bytes.Equal(otAddr[:], utxoOutput.OTAddr[:]) {
				return &wtypes.ProofKeyRet{
					ProofKey: fmt.Sprintf("%x", derivationKey[:]),
				}, nil
			}
			outIdx++
		}
	}
	return nil, wtypes.ErrNoTransInTx
}

// CheckProofKey verify proof key
func (s *PublicTransactionPoolAPI) CheckProofKey(ctx context.Context, args wtypes.VerifyProofKeyArgs) (*wtypes.VerifyProofKeyRet, error) {
	tx, err := s.wallet.GetUTXOTx(args.Hash)
	if err != nil {
		return nil, err
	}
	addr, err := wallet.StrToAddress(args.Addr)
	if err != nil {
		return nil, err
	}
	k, err := hex.DecodeString(args.Key)
	if err != nil || len(k) != lkctypes.COMMONLEN {
		return nil, wtypes.ErrArgsInvalid
	}
	var key lkctypes.Key
	copy(key[:], k[:])
	ret := &wtypes.VerifyProofKeyRet{
		Records: make([]*wtypes.VerifyProofKey, 0),
	}
	outIdx := 0
	for _, output := range tx.Outputs {
		if utxoOutput, ok := output.(*types.UTXOOutput); ok {
			otAddr, err := xcrypto.DerivePublicKey(lkctypes.KeyDerivation(key), outIdx, addr.SpendPublicKey)
			if err != nil {
				return nil, err
			}
			if bytes.Equal(otAddr[:], utxoOutput.OTAddr[:]) {
				ecdh := &lkctypes.EcdhTuple{
					Amount: tx.RCTSig.RctSigBase.EcdhInfo[outIdx].Amount,
				}
				scalar, err := xcrypto.DerivationToScalar(lkctypes.KeyDerivation(key), outIdx)
				if err != nil {
					return nil, err
				}
				ok := xcrypto.EcdhDecode(ecdh, lkctypes.Key(scalar), false)
				if !ok {
					return nil, err
				}
				ret.Records = append(ret.Records, &wtypes.VerifyProofKey{
					Hash:   args.Hash,
					Addr:   args.Addr,
					Amount: (*hexutil.Big)(types.Hash2BigInt(ecdh.Amount)),
				})
			}
			outIdx++
		}
	}
	return ret, nil
}

// SelectAddress set wallet curr account
func (s *PublicTransactionPoolAPI) SelectAddress(ctx context.Context, addr common.Address) (bool, error) {
	err := s.wallet.SelectAddress(addr)
	return err == nil, err
}

// SetRefreshBlockInterval set wallet curr account
func (s *PublicTransactionPoolAPI) SetRefreshBlockInterval(ctx context.Context, interval time.Duration) (bool, error) {
	if interval <= time.Duration(0) {
		return false, fmt.Errorf("interval must be greater than 0")
	}
	sec := interval * time.Second
	err := s.wallet.SetRefreshBlockInterval(sec)
	return err == nil, err
}
