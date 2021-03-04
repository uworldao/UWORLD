package types

import (
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/param"
)

var (
	ErrSignature        = errors.New("signature verification failed")
	ErrSigner           = errors.New("inconsistent signer")
	ErrNoSignature      = errors.New("no signature")
	ErrWrongSignature   = errors.New("wrong signature")
	ErrTxNonceRepeat    = errors.New("the nonce value is repeated, increase the nonce value")
	ErrCoinBase         = errors.New("wrong coin base reward")
	ErrAddress          = errors.New("wrong address")
	ErrTxType           = errors.New("unknown transaction type")
	ErrTxAmount         = errors.New("wrong amount")
	ErrTxHash           = errors.New("wrong transaction hash")
	ErrNonce            = errors.New("incorrect nonce")
	ErrTooBigNonce      = errors.New("too big nonce")
	ErrNotEnoughBalance = errors.New("balance is not enough")
	ErrNotEnoughFees    = fmt.Errorf(" %s is not enough", param.Token)
	ErrContract         = errors.New("lack of contract content")
	ErrContractAddr     = errors.New("wrong contract address")
	ErrTxHead           = errors.New("transaction head cant be nil")
	ErrTxBody           = errors.New("transaction body cant be nil")
)
