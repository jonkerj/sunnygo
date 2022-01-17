package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "sunnygo",
		Short: "tool to monitor your SMA inverter",
		Run:   root,
	}
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("sunnygo")
	flags := queryCmd.PersistentFlags()
	flags.String("url", "https://sunnyboy.local", "URL to Sunnboy")
	flags.String("right", "usr", "Permission level/username")
	flags.String("password", "", "Password")
	flags.String("device", "", "Device ID")
	viper.BindPFlags(flags)
}

func root(cmd *cobra.Command, args []string) {
	fmt.Println("the root command does nothing, use the subcommands")
}

func Execute() {
	rootCmd.Execute()
}
