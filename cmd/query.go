package cmd

import (
	"github.com/jonkerj/sunnygo/pkg/tree"
	"github.com/jonkerj/sunnygo/pkg/webconnect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	queryCmd = &cobra.Command{
		Use:   "query",
		Short: "Query equipment for available fields",
		Run:   query,
	}
)

func init() {
	rootCmd.AddCommand(queryCmd)
}

func query(cmd *cobra.Command, args []string) {
	w, err := webconnect.NewWebconnect(viper.GetString("url"))
	if err != nil {
		panic(err)
	}

	if err := w.Login(viper.GetString("right"), viper.GetString("password")); err != nil {
		panic(err)
	}
	defer w.Logout()

	meta, err := w.DownloadMeta()
	if err != nil {
		panic(err)
	}

	vr, err := w.DownloadValues()
	if err != nil {
		panic(err)
	}

	root, err := webconnect.NodifyAllValues("0199-xxxxx1A5", meta, vr)
	if err != nil {
		panic(err)
	}

	tree.PrintTree(root, 0)
}
