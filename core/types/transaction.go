package types

import (
	"encoding/json"
	"fmt"
	"github.com/uworldao/UWORLD/common/encode/rlp"
	hash2 "github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/crypto/ecc/secp256k1"
	"github.com/uworldao/UWORLD/crypto/hash"
	"github.com/uworldao/UWORLD/param"
	"github.com/uworldao/UWORLD/ut"
	"strconv"
)

const (
	NormalTransaction TransactionType = iota
	ContractTransaction
	LoginCandidate
	/*LogoutCandidate
	VoteToCandidate*/
)
const MaxNote = 256

const CoinBase = "coinbase"

type TransactionType uint8

type TransactionHead struct {
	TxHash     hash2.Hash
	TxType     TransactionType
	From       hash2.Address
	Nonce      uint64
	Fees       uint64
	Time       uint64
	Note       string
	SignScript *SignScript
}

type Transaction struct {
	TxHead *TransactionHead
	TxBody ITransactionBody
}

func (t *Transaction) IsCoinBase() bool {
	return t.TxHead.From.IsEqual(hash2.StringToAddress(CoinBase))
}

func (t *Transaction) Size() uint64 {
	bytes, _ := t.EncodeToBytes()
	return uint64(len(bytes))
}

func (t *Transaction) VerifyTx() error {
	if err := t.verifyHead(); err != nil {
		return err
	}

	if err := t.verifyBody(); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) verifyHead() error {
	if t.TxHead == nil {
		return ErrTxHead
	}

	if err := t.verifyTxType(); err != nil {
		return err
	}

	if err := t.verifyTxHash(); err != nil {
		return err
	}

	if err := t.verifyTxFrom(); err != nil {
		return err
	}

	if err := t.verifyTxNote(); err != nil {
		return err
	}

	if err := t.verifyTxFees(); err != nil {
		return err
	}

	if err := t.verifyTxSinger(); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) verifyBody() error {
	if t.TxBody == nil {
		return ErrTxBody
	}

	if err := t.verifyAmount(); err != nil {
		return err
	}

	if err := t.TxBody.VerifyBody(t.TxHead.From); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) VerifyCoinBaseTx(height, sumFees uint64) error {
	if err := t.verifyTxSize(); err != nil {
		return err
	}

	if err := t.verifyCoinBaseAmount(height, sumFees); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) verifyTxFees() error {
	var fees uint64
	switch t.TxHead.TxType {
	case NormalTransaction:
		fees = param.Fees
	case ContractTransaction:
		fees = param.TokenConsumption
	}
	if t.TxHead.Fees != fees {
		return fmt.Errorf("transaction costs %d fees", fees)
	}
	return nil
}

func (t *Transaction) verifyTxSinger() error {
	if !Verify(t.TxHead.TxHash, t.TxHead.SignScript) {
		return ErrSignature
	}

	if !VerifySigner(param.Net, t.TxHead.From, t.TxHead.SignScript.PubKey) {
		return ErrSigner
	}
	return nil
}

func (t *Transaction) verifyTxSize() error {
	// TODO change maxsize
	switch t.TxHead.TxType {
	case NormalTransaction:
		fallthrough
	case ContractTransaction:
		return nil
		/*case LogoutCandidate:
			fallthrough
		case LoginCandidate:
			fallthrough
		case VoteToCandidate:
			if t.Size() > MaxNoDataTxSize {
				return ErrTxSize
			}*/
	}
	return nil
}

func (t *Transaction) verifyCoinBaseAmount(height, amount uint64) error {
	nTx := t.TxBody.(*NormalTransactionBody)
	sumAmount := CalCoinBase(height, param.CoinHeight) + amount
	if sumAmount != nTx.Amount {
		return ErrCoinBase
	}
	return nil
}

func (t *Transaction) verifyAmount() error {
	nTx, ok := t.TxBody.(*NormalTransactionBody)
	if ok && nTx.Amount < param.MinAllowedAmount {
		return fmt.Errorf("the minimum amount of the transaction must not be less than %d", param.MinAllowedAmount)
	}
	return nil
}

func (t *Transaction) verifyTxFrom() error {
	if !ut.CheckUWDAddress(param.Net, t.From().String()) {
		return ErrAddress
	}
	return nil
}

func (t *Transaction) verifyTxType() error {
	switch t.TxHead.TxType {
	case NormalTransaction:
		return nil
	case ContractTransaction:
		return nil
		/*case VoteToCandidate:
			return nil
		case LoginCandidate:
			return nil
		case LogoutCandidate:
			return nil*/
	}
	return ErrTxType
}

func (t *Transaction) verifyTxHash() error {
	newTx := t.copy()
	newTx.SetHash()
	if newTx.Hash().IsEqual(t.Hash()) {
		return nil
	}
	return ErrTxHash
}

func (t *Transaction) verifyTxNote() error {
	if len(t.TxHead.Note) > MaxNote {
		return fmt.Errorf("the length of the transaction note must not be greater than %d", MaxNote)
	}
	return nil
}

func (t *Transaction) EncodeToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(t)
}

func (t *Transaction) SignTx(key *secp256k1.PrivateKey) error {
	var err error
	if t.TxHead.SignScript, err = Sign(key, t.TxHead.TxHash); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) SetHash() error {
	t.TxHead.TxHash = hash2.Hash{}
	t.TxHead.SignScript = &SignScript{}
	rpcTx, err := TranslateTxToRpcTx(t)
	if err != nil {
		return err
	}
	mBytes, err := json.Marshal(rpcTx)
	if err != nil {
		return err
	}
	t.TxHead.TxHash = hash.Hash(mBytes)
	return nil
}

func (t *Transaction) copy() *Transaction {
	header := &TransactionHead{
		TxHash:     t.TxHead.TxHash,
		TxType:     t.TxHead.TxType,
		From:       t.TxHead.From,
		Nonce:      t.TxHead.Nonce,
		Fees:       t.TxHead.Fees,
		Time:       t.TxHead.Time,
		Note:       t.TxHead.Note,
		SignScript: t.TxHead.SignScript,
	}
	return &Transaction{
		TxHead: header,
		TxBody: t.TxBody,
	}
}

func (t *Transaction) NonceKey() string {
	return t.TxHead.From.String() + "_" + strconv.FormatUint(t.TxHead.Nonce, 10)
}

func (t *Transaction) Hash() hash2.Hash {
	return t.TxHead.TxHash
}

func (t *Transaction) From() hash2.Address {
	return t.TxHead.From
}

func (t *Transaction) GetFees() uint64 {
	return t.TxHead.Fees
}

func (t *Transaction) GetNonce() uint64 {
	return t.TxHead.Nonce
}

func (t *Transaction) GetTime() uint64 {
	return t.TxHead.Time
}

func (t *Transaction) GetNote() string {
	return t.TxHead.Note
}

func (t *Transaction) GetTxType() TransactionType {
	return t.TxHead.TxType
}

func (t *Transaction) GetSignScript() *SignScript {
	return t.TxHead.SignScript
}

func (t *Transaction) GetTxHead() *TransactionHead {
	return t.TxHead
}

func (t *Transaction) GetTxBody() ITransactionBody {
	return t.TxBody
}

func (t *Transaction) TranslateToRlpTransaction() *RlpTransaction {
	rlpTx := &RlpTransaction{}
	rlpTx.TxHead = t.TxHead
	rlpTx.TxBody, _ = rlp.EncodeToBytes(t.TxBody)
	return rlpTx
}

type TxLocation struct {
	TxRoot  hash2.Hash
	TxIndex uint32
	Height  uint64
}

func (t *TxLocation) GetHeight() uint64 {
	return t.Height
}

func CalCoinBase(height, startHeight uint64) uint64 {
	if height < startHeight {
		return 0
	}
	height = height - startHeight + 1
	if height%(60*60*24/param.BlockInterval) != 0 {
		return 0
	}
	var mouth uint64
	if height%(60*60*24*30/param.BlockInterval) == 0 {
		mouth = height / (60 * 60 * 24 * 30 / param.BlockInterval)
	} else {
		mouth = height/(60*60*24*30/param.BlockInterval) + 1
	}

	coins, ok := param.DayCoin[mouth]
	if !ok {
		return 0
	}
	count, _ := NewAmount(coins)
	return count
}
