package command

import (
	"context"
	"github.com/jhdriver/UWORLD/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	nodeCmds := []*cobra.Command{
		GetLastHeightCmd,
		GetTxPoolTxs,
		GetPeersCmd,
		NodeInfoCmd,
	}
	RootCmd.AddCommand(nodeCmds...)
	RootSubCmdGroups["node"] = nodeCmds
}

//GenerateCmd cpu mine block
var GetTxPoolTxs = &cobra.Command{
	Use:     "GetTxPool",
	Short:   "GetTxPool; Get transactions in the transaction pool;",
	Aliases: []string{"gettxpool", "gtp", "GTP"},
	Example: `
	GetTxPool 
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  GetTxPool,
}

func GetTxPool(cmd *cobra.Command, args []string) {

	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := client.Gc.GetPoolTxs(ctx, &rpc.Null{})
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

var GetPeersCmd = &cobra.Command{
	Use:     "GetPeers",
	Short:   "GetPeers; Get transactions in the transaction pool;",
	Aliases: []string{"getpeers", "gp", "GP"},
	Example: `
	Peers 
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  Peers,
}

func Peers(cmd *cobra.Command, args []string) {
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	resp, err := client.Gc.Peers(ctx, &rpc.Null{})
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

var GetLastHeightCmd = &cobra.Command{
	Use:     "GetLastHeight",
	Short:   "GetLastHeight; Get last height of node;",
	Aliases: []string{"getlastheight", "glh", "GLP"},
	Example: `
	GetLastHeight 
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  GetLastHeight,
}

func GetLastHeight(cmd *cobra.Command, args []string) {
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	resp, err := client.Gc.GetLastHeight(ctx, &rpc.Null{})
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

var NodeInfoCmd = &cobra.Command{
	Use:     "NodeInfo ;Gets the current node information",
	Short:   "NodeInfo ;Gets the current node information;",
	Aliases: []string{"nodeinfo", "NI", "ni"},
	Example: `
	NodeInfo
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  NodeInfo,
}

func NodeInfo(cmd *cobra.Command, args []string) {
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	resp, err := client.Gc.NodeInfo(ctx, &rpc.Null{})
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
