package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/MinoKim82/go-kis-cli/pkg/overseas"
	"github.com/spf13/cobra"
)

var overseasCmd = &cobra.Command{
	Use:     "overseas",
	Aliases: []string{"ovs"},
	Short:   "Overseas stock commands",
	Long:    `Commands for querying and trading overseas (US, etc.) stocks.`,
}

var ovsQuoteCmd = &cobra.Command{
	Use:   "quote [exchange] [symbol]",
	Short: "Get the current quote for an overseas stock",
	Long:  `Exchange codes: NAS (Nasdaq), NYS (NYSE), AMS (Amex), etc. Example: kis ovs quote NAS AAPL`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		exchange, symbol := args[0], args[1]
		resp, err := overseas.GetQuote(exchange, symbol)
		if err != nil {
			fmt.Printf("Error fetching overseas quote: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Symbol\tPrice\tChange\tChange %\tVolume\t")

		out := resp.Output
		fmt.Fprintf(w, "%s\t%s\t%s\t%s%%\t%s\t\n", symbol, out.Last, out.Diff, out.Rate, out.Tvol)

		w.Flush()
	},
}

var ovsBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get the current overseas account balance",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := overseas.GetBalance()
		if err != nil {
			fmt.Printf("Error fetching overseas balance: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Code\tName\tQty\tPrice\tProfit/Loss($)\tProfit%\t")

		for _, item := range resp.Output1 {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s%%\t\n", item.OvrsPdno, item.OvrsItemName, item.OvrsCblcQty, item.NowPric2, item.EvluPflsAmt, item.EvluPflsRt)
		}
		w.Flush()

		fmt.Println("--------------------------------")
		fmt.Printf("Total Evaluation Amount (USD): %s\n", resp.Output2.FrcrEvluTamt)
	},
}

var ovsBuyCmd = &cobra.Command{
	Use:   "buy [exchange] [symbol] [qty] [price]",
	Short: "Buy an overseas stock",
	Long:  `Buy an overseas (US, etc.) stock. Example: kis ovs buy NAS AAPL 10 150`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		exchange, symbol, qty, price := args[0], args[1], args[2], args[3]
		resp, err := overseas.BuyOrder(exchange, symbol, qty, price)
		if err != nil {
			fmt.Printf("Overseas Buy Order Failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Overseas Buy Order Placed Successfully!\n")
		fmt.Printf("Order No: %s | Time: %s\n", resp.Output.ODNO, resp.Output.ORD_TMD)
	},
}

var ovsSellCmd = &cobra.Command{
	Use:   "sell [exchange] [symbol] [qty] [price]",
	Short: "Sell an overseas stock",
	Long:  `Sell an overseas (US, etc.) stock. Example: kis ovs sell NAS AAPL 10 155`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		exchange, symbol, qty, price := args[0], args[1], args[2], args[3]
		resp, err := overseas.SellOrder(exchange, symbol, qty, price)
		if err != nil {
			fmt.Printf("Overseas Sell Order Failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Overseas Sell Order Placed Successfully!\n")
		fmt.Printf("Order No: %s | Time: %s\n", resp.Output.ODNO, resp.Output.ORD_TMD)
	},
}

func init() {
	rootCmd.AddCommand(overseasCmd)
	overseasCmd.AddCommand(ovsQuoteCmd)
	overseasCmd.AddCommand(ovsBalanceCmd)
	overseasCmd.AddCommand(ovsBuyCmd)
	overseasCmd.AddCommand(ovsSellCmd)
}
