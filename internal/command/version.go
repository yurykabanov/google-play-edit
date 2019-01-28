package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Build   = "unknown"
	Version = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Google Play Edit %s (build %s)", Version, Build)
	},
}
