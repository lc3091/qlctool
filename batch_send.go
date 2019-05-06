/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package commands

import (
	"github.com/abiosoft/ishell"
	"github.com/qlcchain/go-qlc/ledger"
	"github.com/spf13/cobra"
)

var blockIndex = 0
var l *ledger.Ledger

func batchSend() {
	var fromAccountP string
	var toAccountsP []string
	var tokenP string
	var amountP string
	var countP int

	if interactive {
		from := Flag{
			Name:  "from",
			Must:  true,
			Usage: "send account private key",
			Value: "",
		}
		to := Flag{
			Name:  "to",
			Must:  true,
			Usage: "receive accounts",
			Value: "",
		}
		token := Flag{
			Name:  "token",
			Must:  false,
			Usage: "token name for send action(defalut is QLC)",
			Value: "QLC",
		}
		amount := Flag{
			Name:  "amount",
			Must:  true,
			Usage: "send amount",
			Value: "",
		}
		count := Flag{
			Name:  "count",
			Must:  true,
			Usage: "send count",
			Value: 100,
		}
		c := &ishell.Cmd{
			Name: "batchsend",
			Help: "batch send transaction",
			Func: func(c *ishell.Context) {
				args := []Flag{from, to, token, amount, count}
				if HelpText(c, args) {
					return
				}
				if err := CheckArgs(c, args); err != nil {
					Warn(err)
					return
				}
				fromAccountP = StringVar(c.Args, from)
				toAccountsP = StringSliceVar(c.Args, to)
				tokenP := StringVar(c.Args, token)
				amountP := StringVar(c.Args, amount)
				countP,_ := IntVar(c.Args, count)

				var i int
				for i = 0; i < countP; i++ {
					for _, toAccount := range toAccountsP {
						if err := sendAction(fromAccountP, toAccount, tokenP, amountP); err != nil {
							Warn(err)
							return
						}
						blockIndex++
					}
				}

				Info("batch transaction done")
			},
		}
		shell.AddCmd(c)
	} else {
		var batchSendCmd = &cobra.Command{
			Use:   "batchsend",
			Short: "batch send transaction",
			Run: func(cmd *cobra.Command, args []string) {
				for _, toAccount := range toAccountsP {
					if err := sendAction(fromAccountP, toAccount, tokenP, amountP); err != nil {
						cmd.Println(err)
						return
					}
				}
				cmd.Println("batch transaction done")
			},
		}
		batchSendCmd.Flags().StringVar(&fromAccountP, "from", "", "send account private key")
		batchSendCmd.Flags().StringSliceVar(&toAccountsP, "to", []string{}, "receive accounts")
		batchSendCmd.Flags().StringVar(&tokenP, "token", "QLC", "token name for send action")
		batchSendCmd.Flags().StringVar(&amountP, "amount", "", "send amount")
		batchSendCmd.Flags().IntVar(&countP, "count", 100, "send count")
		rootCmd.AddCommand(batchSendCmd)
	}
}
