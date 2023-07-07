package cmd

import (
	"fmt"
	"os"
	"strconv"

	"steamcalc/market/buff"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var rootCmd = &cobra.Command{
	Use:   "steamcalc [MIN] [MAX]",
	Short: "Crawl third-party marketplace items and calculate resale discounts.",

	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		min, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		max, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		slog.Info("Crawling buff...", slog.Float64("min", min), slog.Float64("max", max))
		buff.Crawl(min, max)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
