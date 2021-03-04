package command

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/uworldao/UWORLD/crypto/mnemonic"
)

func init() {
	mnemonicCmds := []*cobra.Command{
		EntropyCmd,
		MnemonicCmd,
		MnemonicToEcCmd,
	}

	RootCmd.AddCommand(mnemonicCmds...)
	RootSubCmdGroups["mnemonic"] = mnemonicCmds
}

var EntropyCmd = &cobra.Command{
	Use:     "Entropy ;Generate a cryptographically secure pseudorandom entropy(seed);",
	Aliases: []string{"entropy", "E", "e"},
	Short:   "Entropy ;Generate a cryptographically secure pseudorandom entropy(seed);",
	Example: `
	Entropy 
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  Entropy,
}

func Entropy(cmd *cobra.Command, args []string) {
	if entropy, err := mnemonic.Entropy(); err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	} else {
		fmt.Println()
		fmt.Println(entropy)
	}
}

var MnemonicCmd = &cobra.Command{
	Use:     "Mnemonic {entropy};Create a mnemonic world-list (BIP39) from an entropy;",
	Aliases: []string{"mnemonic", "M", "m"},
	Short:   "Mnemonic {entropy};Create a mnemonic world-list (BIP39) from an entropy;",
	Example: `
	Mnemonic dad1f695098e409da517aa09d91bb163ea749c3f9ee564cb75e223a78f460a1e
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  Mnemonic,
}

func Mnemonic(cmd *cobra.Command, args []string) {
	if mnemonic, err := mnemonic.Mnemonic(args[0]); err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	} else {
		fmt.Println()
		fmt.Println(mnemonic)
	}
}

var MnemonicToEcCmd = &cobra.Command{
	Use:     "MnemonicToEc {Mnemonic};Create a new EC private key from a mnemonic;",
	Aliases: []string{"mnemonictoec", "MTE", "mte"},
	Short:   "MnemonicToEc {Mnemonic};Create a new EC private key from a mnemonic;",
	Example: `
	Mnemonic 'suspect moral pipe basic tomato excite nephew vocal antenna silver unable sick point evoke wrist syrup gospel forum joy elder jump perfect chronic select'
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  MnemonicToEc,
}

func MnemonicToEc(cmd *cobra.Command, args []string) {
	if ec, err := mnemonic.MnemonicToEc(args[0]); err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	} else {
		fmt.Println()
		fmt.Println(ec.String())
	}
}
