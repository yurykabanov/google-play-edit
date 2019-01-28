package command

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yurykabanov/google-play-edit/internal/pretty"
	"github.com/yurykabanov/google-play-edit/pkg/play"
)

var editCommitCmd = &cobra.Command{
	Use:   "commit [id]",
	Short: "Commit edit with given ID",
	Long:  `Commits all changes made to the edit.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := mustMakeHttpClient()
		token := mustAuthenticate(client)
		api := play.NewApi(play.WithApiHttpClient(client))

		packageName := viper.GetString("package-name")
		editId := args[0]

		_, err := api.Edits.Commit(context.Background(), token, packageName, editId)
		if err != nil {
			pretty.Errorf("Unable to commit edit: %s", err.Error())
			os.Exit(1)
		}
	},
}
