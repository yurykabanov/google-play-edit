package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const appCommand = "google-play-edit"

var rootCmd = &cobra.Command{
	Use:   appCommand,
	Short: "Google Play Edit is a tool performing application editing and bulk screenshots uploading",
	Long:  `Google Play Edit handles Android application's listings and screenshots updates.'`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(editListCmd)
	rootCmd.AddCommand(editInsertCmd)
	rootCmd.AddCommand(editCommitCmd)

	rootCmd.PersistentFlags().String("account", "", "Google Service Account JSON file path")
	rootCmd.PersistentFlags().String("token", "", "Access Token for Google API")

	rootCmd.PersistentFlags().Bool("print-token", false, "Print Access Token for later usage")

	rootCmd.PersistentFlags().String("proxy", "", "HTTP Proxy")
	rootCmd.PersistentFlags().Bool("proxy-insecure", false, "Skip TLS verification")

	rootCmd.PersistentFlags().String("package-name", "", "Application Package Name")

	viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account"))
	viper.BindPFlag("print-token", rootCmd.PersistentFlags().Lookup("print-token"))
	viper.BindPFlag("proxy", rootCmd.PersistentFlags().Lookup("proxy"))
	viper.BindPFlag("proxy-insecure", rootCmd.PersistentFlags().Lookup("proxy-insecure"))
	viper.BindPFlag("package-name", rootCmd.PersistentFlags().Lookup("package-name"))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
