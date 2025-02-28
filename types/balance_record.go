package types

import (
	"math/big"
	"encoding/json"

	"github.com/lianxiangcloud/linkchain/libs/common"
)

const (
	AccountAddress uint32 = 0
	PrivateAddress uint32 = 1
	NoAddress      uint32 = 2
)

type Payload []byte

type BlockBalanceRecords struct {
	TxRecords []*TxBalanceRecords `json:"tx_records"`
	BlockHash common.Hash         `json:"block_hash"`
	BlockTime uint64              `json:"block_time"`
}

type TxBalanceRecords struct {
	Hash     common.Hash     `json:"hash"`
	Type     string          `json:"type"`
	Records  []BalanceRecord `json:"records"`
	Payloads []Payload       `json:"payloads"`
	Nonce    uint64          `json:"nonce"`
	GasLimit uint64          `json:"gas_limit"`
	GasPrice *big.Int        `json:"gas_price"`
	From     common.Address  `json:"from"`
	To       common.Address  `json:"to"`
	TokenId  common.Address  `json:"token_id"`
}

type BalanceRecord struct {
	From            common.Address `json:from`
	To              common.Address `json:to`
	FromAddressType uint32         `json:"from_address_type"`
	ToAddressType   uint32         `json:"to_address_type"`
	Type            string         `json:"type"`
	TokenID         common.Address `json:"token_id"`
	Amount          *big.Int       `json:"amount"`
	Hash            common.Hash    `json:"hash"`
}

type balanceRecordHash struct {
	From            common.Address `json:from`
	To              common.Address `json:to`
	FromAddressType uint32         `json:"from_address_type"`
	ToAddressType   uint32         `json:"to_address_type"`
	Type            string         `json:"type"`
	TokenID         common.Address `json:"token_id"`
	Amount          *big.Int       `json:"amount"`
	RandomNum       uint64         `json:"random_num"`
}

func NewTxBalanceRecords() *TxBalanceRecords {
	return &TxBalanceRecords{}
}

func GenBalanceRecord(from common.Address, to common.Address, fromAddressType uint32, toAddressType uint32, typeStr string, tokenId common.Address, amount *big.Int) BalanceRecord {
	finnalAmount := big.NewInt(0).Add(big.NewInt(0), amount)
	brh := balanceRecordHash{
		From:            from,
		To:              to,
		FromAddressType: fromAddressType,
		ToAddressType:   toAddressType,
		Type:            typeStr,
		TokenID:         tokenId,
		Amount:          finnalAmount,
		RandomNum:       common.RandUint64(),
	}
	hash := rlpHash(brh)

	return BalanceRecord{
		From:            brh.From,
		To:              brh.To,
		FromAddressType: brh.FromAddressType,
		ToAddressType:   brh.ToAddressType,
		Type:            brh.Type,
		TokenID:         brh.TokenID,
		Amount:          brh.Amount,
		Hash:            hash,
	}
}

func NewBlockBalanceRecords() *BlockBalanceRecords {
	return &BlockBalanceRecords{
		TxRecords: make([]*TxBalanceRecords, 0),
	}
}

func (b *BlockBalanceRecords) Json() []byte {
	jb, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return jb
}

func (b *BlockBalanceRecords) AddTxBalanceRecord(t *TxBalanceRecords)  {
	b.TxRecords = append(b.TxRecords, t)
}

func (b *BlockBalanceRecords) Clear() {
	b.TxRecords = make([]*TxBalanceRecords, 0)
}

func (b *BlockBalanceRecords) SetBlockHash(blockHash common.Hash) {
	b.BlockHash = blockHash
}

func (b *BlockBalanceRecords) SetBlockTime(blockTime uint64) {
	b.BlockTime = blockTime
}

func (t *TxBalanceRecords) SetOptions(hash common.Hash, typenane string, payloads []Payload, nonce uint64,
	gasLimit uint64, gasPrice *big.Int, from common.Address, to common.Address, tokenId common.Address) {
		t.Hash     = hash
		t.Type     = typenane
		t.Payloads = payloads
		t.Nonce    = nonce
		t.GasLimit = gasLimit
		t.GasPrice = gasPrice
		t.From     = from
		t.To       = to
		t.TokenId  = tokenId
}

func (t *TxBalanceRecords) AddBalanceRecord(br BalanceRecord) {
	t.Records = append(t.Records, br)
}

func (t *TxBalanceRecords) ClearBalanceRecord() {
	t.Records = make([]BalanceRecord, 0)
}