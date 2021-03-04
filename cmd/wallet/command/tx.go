package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
	"github.com/uworldao/UWORLD/crypto/ecc/secp256k1"
	"github.com/uworldao/UWORLD/rpc"
	"github.com/uworldao/UWORLD/rpc/rpctypes"
	"github.com/uworldao/UWORLD/ut/transaction"
	"strconv"
	"time"
)

func init() {
	txCmds := []*cobra.Command{
		GetTransactionCmd,
		SendTransactionCmd,
	}
	RootCmd.AddCommand(txCmds...)
	RootSubCmdGroups["transaction"] = txCmds

}

var SendTransactionCmd = &cobra.Command{
	Use:     "SendTransaction {from} {to} {contract} {amount} {note} {password} {nonce}; Send a transaction;",
	Aliases: []string{"sendtransaction", "st", "ST"},
	Short:   "SendTransaction {from} {to} {contract} {amount} {note} {password} {nonce}; Send a transaction;",
	Example: `
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ UWD 10  "transaction note"
		OR
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ UWD 10  "transaction note" 123456
		OR
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ UWD 10  "transaction note" 123456 1
	`,
	Args: cobra.MinimumNArgs(5),
	Run:  SendTransaction,
}

func SendTransactionByAddrRpc(args []string) (string, error) {
	var passwd []byte
	var err error
	if len(args) > 5 {
		passwd = []byte(args[5])
	} else {
		return "", fmt.Errorf("read password failed! %s", err.Error())
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		return "", fmt.Errorf("wrong password")

	}

	tx, err := parseParams(args)
	if err != nil {
		return "", err

	}
	resp, err := GetAccountByRpc(tx.From().String())
	if err != nil {
		return "", err

	}
	if resp.Code != 0 {
		return "", fmt.Errorf(" err: code %d, message: %s", resp.Code, resp.Err)

	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		return "", fmt.Errorf(" err: ", err)

	}
	if tx.TxHead.Nonce == 0 {
		tx.TxHead.Nonce = account.Nonce + 1
	}
	if !signTx1(tx, privKey.Private) {
		return "", fmt.Errorf(" err: ", errors.New("signature failure"))
	}

	rs, err := sendTx1(tx)
	if err != nil {
		return "", fmt.Errorf(" err: ", err)
	} else if rs.Code != 0 {
		return "", fmt.Errorf(" err: code %d, message: %s", rs.Code, rs.Err)
	} else {
		return string(rs.Result), nil
	}
}

func SendTransaction(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 5 {
		passwd = []byte(args[5])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			log.Error(cmd.Use+" err: ", fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("wrong password"))
		return
	}

	tx, err := parseParams(args)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	resp, err := GetAccountByRpc(tx.From().String())
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code != 0 {
		log.Errorf(cmd.Use+" err: code %d, message: %s", resp.Code, resp.Err)
		return
	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if tx.TxHead.Nonce == 0 {
		tx.TxHead.Nonce = account.Nonce + 1
	}
	if !signTx(cmd, tx, privKey.Private) {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendTx(cmd, tx)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
	} else if rs.Code != 0 {
		log.Errorf(cmd.Use+" err: code %d, message: %s", rs.Code, rs.Err)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseParams(args []string) (*types.Transaction, error) {
	var err error
	var amount, nonce uint64
	from := hasharry.StringToAddress(args[0])
	to := hasharry.StringToAddress(args[1])
	contract := hasharry.StringToAddress(args[2])
	if fAmount, err := strconv.ParseFloat(args[3], 64); err != nil {
		return nil, errors.New("wrong amount")
	} else {
		if fAmount < 0 {
			return nil, errors.New("wrong amount")
		}
		if amount, err = types.NewAmount(fAmount); err != nil {
			return nil, errors.New("wrong amount")
		}
	}
	note := args[4]
	if len(args) > 6 {
		nonce, err = strconv.ParseUint(args[6], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	return transaction.NewTransaction(from.String(), to.String(), contract.String(), note, amount, nonce), nil
}

func signTx(cmd *cobra.Command, tx *types.Transaction, key string) bool {
	tx.SetHash()
	priv, err := secp256k1.ParseStringToPrivate(key)
	if err != nil {
		log.Error(cmd.Use+" err: ", errors.New("[key] wrong"))
		return false
	}
	if err := tx.SignTx(priv); err != nil {
		log.Error(cmd.Use+" err: ", errors.New("sign failed"))
		return false
	}
	return true
}
func signTx1(tx *types.Transaction, key string) bool {
	tx.SetHash()
	priv, err := secp256k1.ParseStringToPrivate(key)
	if err != nil {
		return false
	}
	if err := tx.SignTx(priv); err != nil {
		return false
	}
	return true
}
func sendTx(cmd *cobra.Command, tx *types.Transaction) (*rpc.Response, error) {
	rpcTx, err := types.TranslateTxToRpcTx(tx)
	if err != nil {
		return nil, err
	}
	rpcClient, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	jsonBytes, err := json.Marshal(rpcTx)
	if err != nil {
		return nil, err
	}
	re := &rpc.Bytes{Bytes: jsonBytes}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := rpcClient.Gc.SendTransaction(ctx, re)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
func sendTx1(tx *types.Transaction) (*rpc.Response, error) {
	rpcTx, err := types.TranslateTxToRpcTx(tx)
	if err != nil {
		return nil, err
	}
	rpcClient, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	jsonBytes, err := json.Marshal(rpcTx)
	if err != nil {
		return nil, err
	}
	re := &rpc.Bytes{Bytes: jsonBytes}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := rpcClient.Gc.SendTransaction(ctx, re)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

var GetTransactionCmd = &cobra.Command{
	Use:     "GetTransaction {txhash}; Get Transaction by hash;",
	Aliases: []string{"gettransaction", "gt", "GT"},
	Short:   "GetTransaction {txhash}; Get Transaction by hash;",
	Example: `
	GetTransaction 0xef7b92e552dca02c97c9d596d1bf69d0044d95dec4cee0e6a20153e62bce893b
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  GetTransaction,
}

func GetTransaction(cmd *cobra.Command, args []string) {
	resp, err := GetTransactionRpc(args[0])
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	}
	outputRespError(cmd.Use, resp)
}

func GetTransactionRpc(hashStr string) (*rpc.Response, error) {
	client, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	h := &rpc.Hash{Hash: hashStr}
	resp, err := client.Gc.GetTransaction(ctx, h)
	return resp, err
}

func SendTransactionRpc(tx string) (*rpc.Response, error) {

	rpcClient, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	re := &rpc.Bytes{Bytes: []byte(tx)}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := rpcClient.Gc.SendTransaction(ctx, re)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
