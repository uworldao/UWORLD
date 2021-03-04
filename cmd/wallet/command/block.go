package command

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/uworldao/UWORLD/rpc"
	"strconv"
	"time"
)

func init() {
	blockCmds := []*cobra.Command{GetBlockCmd}

	RootCmd.AddCommand(blockCmds...)
	RootSubCmdGroups["block"] = blockCmds
}

//GenerateCmd cpu mine block
var GetBlockCmd = &cobra.Command{
	Use:     "GetBlock {height/hash};",
	Short:   "GetBlock {height/hash}; Get block by height or hash;",
	Aliases: []string{"getblock", "gb", "GB"},
	Example: `
	GetBlock 1 
	GetBlock 0x4e32b712330c0d4ee45f06017390c5d1d3c26d0e6c7be4ea9a5036bdb6c72a07 
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  GetBlock,
}

func GetBlock(cmd *cobra.Command, args []string) {
	height, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		GetBlockByHash(cmd, args)
		return
	}
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := client.Gc.GetBlockByHeight(ctx, &rpc.Height{Height: height})
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

func GetBlockByHash(cmd *cobra.Command, args []string) {
	var err error
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := client.Gc.GetBlockByHash(ctx, &rpc.Hash{Hash: args[0]})
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
