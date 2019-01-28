package command

import (
	"context"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yurykabanov/google-play-edit/internal/pretty"
	"github.com/yurykabanov/google-play-edit/pkg/play"
)

var editListCmd = &cobra.Command{
	Use:   "list [id]",
	Short: "List edit with given ID",
	Long:  `Show detailed information about given edit.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := mustMakeHttpClient()
		token := mustAuthenticate(client)
		api := play.NewApi(play.WithApiHttpClient(client))

		packageName := viper.GetString("package-name")
		editId := args[0]

		edit, err := api.Edits.Get(context.Background(), token, packageName, editId)
		if err != nil {
			pretty.Errorf("Unable to query edit: %s", err.Error())
			os.Exit(1)
		}

		pretty.PrintEdit(edit)
		fmt.Println()

		listings, err := api.Listings.List(context.Background(), token, packageName, editId)
		if err != nil {
			pretty.Errorf("Unable to query listings: %s", err.Error())
			os.Exit(1)
		}

		for _, listing := range listings {
			pretty.PrintListing(&listing)
			fmt.Println()

			images, err := api.Images.List(context.Background(), token, packageName, editId, listing.Language, play.EditImagePhoneScreenshots)
			if err != nil {
				pretty.Errorf("Unable to query images: %s", err.Error())
				os.Exit(1)
			}

			fmt.Printf("%s:\n", aurora.Green("Phone screenshots"))
			for _, image := range images {
				pretty.PrintImage(&image)
			}

			fmt.Println()
			fmt.Println()
		}
	},
}
