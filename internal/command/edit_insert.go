package command

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yurykabanov/google-play-edit/internal/loader"
	"github.com/yurykabanov/google-play-edit/internal/pretty"
	"github.com/yurykabanov/google-play-edit/pkg/play"
	"github.com/yurykabanov/google-play-edit/pkg/task"
)

var editInsertCmd = &cobra.Command{
	Use:   "insert [listings-file]",
	Short: "Insert edit",
	Long:  `Create new edit and update listings and screenshots.
This command will not delete any existing listings.

Listings file could be CSV, YAML or JSON file.`,
	Args:  cobra.ExactArgs(1),

	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("phone-screenshots", cmd.Flags().Lookup("phone-screenshots"))
	},

	Run: func(cmd *cobra.Command, args []string) {
		client := mustMakeHttpClient()
		token := mustAuthenticate(client)
		api := play.NewApi(play.WithApiHttpClient(client))

		packageName := viper.GetString("package-name")

		edit, err := api.Edits.Insert(context.Background(), token, packageName)
		if err != nil {
			pretty.Errorf("Unable to insert new edit: %s", err.Error())
			os.Exit(1)
		}
		editId := edit.Id

		pretty.PrintEdit(edit)
		fmt.Println()

		listings, err := loader.LoadListingsFromFile(args[0])
		if err != nil {
			pretty.Errorf("Unable to read new listings from file: %s", err.Error())
			os.Exit(1)
		}

		upsert := task.NewUpsert(api, token, packageName, editId)

		for _, listing := range listings {
			pretty.PrintListing(&listing)
			fmt.Println()

			images, err := loader.FindImagesForLang(viper.GetString("phone-screenshots"), listing.Language)
			if err != nil {
				pretty.Errorf("Unable to find images for lang %s", listing.Language)
				os.Exit(1)
			}

			ch := make(chan io.ReadSeeker)

			go func(ch chan<- io.ReadSeeker) {
				for _, image := range images {
					f, err := os.Open(image)
					if err != nil {
						pretty.Errorf("Unable to find images for lang %s", listing.Language)
						os.Exit(1)
					}
					ch <- f
				}
				close(ch)
			}(ch)

			err = upsert.Run(context.Background(), &listing, task.ImageTypeSources{
				play.EditImagePhoneScreenshots: ch,
			})
			if err != nil {
				pretty.Errorf("Unable to update/create listing and sync screenshots: %s", err.Error())
				os.Exit(1)
			}

			fmt.Println()
		}
	},
}

func init() {
	editInsertCmd.Flags().String("phone-screenshots", "", "Directory with phone screenshots (should be in subdirectories for each language)")
}
