package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/MinoKim82/go-kis-cli/pkg/domestic"
	"github.com/spf13/cobra"
)

var domesticCmd = &cobra.Command{
	Use:     "domestic",
	Aliases: []string{"dom"},
	Short:   "Domestic stock commands",
	Long:    `Commands for querying and trading domestic stocks.`,
}

var quoteCmd = &cobra.Command{
	Use:   "quote [stock_code]",
	Short: "Get the current quote for a domestic stock",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		code := args[0]
		resp, err := domestic.GetQuote(code)
		if err != nil {
			fmt.Printf("Error fetching quote: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Code\tPrice\tChange\tChange %\tVolume\t")

		out := resp.Output
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", code, out.StckPrpr, out.PrdyVrss, out.PrdyCtrt+"%", out.AcmlVol)

		w.Flush()
	},
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get the current account balance",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := domestic.GetBalance()
		if err != nil {
			fmt.Printf("Error fetching balance: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Code\tName\tQty\tPrice\tProfit/Loss\tProfit%\t")

		for _, item := range resp.Output1 {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s%%\t\n", item.Pdno, item.PrdtName, item.HldgQty, item.Prpr, item.EvluPflsAmt, item.EvluPflsRt)
		}
		w.Flush()

		fmt.Println("--------------------------------")
		if len(resp.Output2) > 0 {
			fmt.Printf("Total Evaluation Amount: %s\n", resp.Output2[0].TotEvluAmt)
		}
	},
}

func init() {
	rootCmd.AddCommand(domesticCmd)
	domesticCmd.AddCommand(quoteCmd)
	domesticCmd.AddCommand(balanceCmd)
}
