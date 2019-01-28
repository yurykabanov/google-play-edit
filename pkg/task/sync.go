package task

import (
	"context"

	"github.com/yurykabanov/google-play-edit/pkg/play"
)

type ListingWithImages struct {
	Listing          *play.Listing
	ImageTypeSources ImageTypeSources
}

type Sync interface {
	Run(ctx context.Context, targetListingsWithImages []ListingWithImages) error
}

type sync struct {
	api         *play.Api
	accessToken *play.AccessToken
	packageName string
	editId      string

	upsert Upsert
}

func NewSync(
	api *play.Api,
	accessToken *play.AccessToken,
	packageName string,
	editId string,
) *sync {
	return &sync{
		api:         api,
		accessToken: accessToken,
		packageName: packageName,
		editId:      editId,

		upsert: NewUpsert(api, accessToken, packageName, editId),
	}
}

func (task *sync) Run(ctx context.Context, targetListingsWithImages []ListingWithImages, delete bool) error {
	if delete {
		originalListings, err := task.api.Listings.List(ctx, task.accessToken, task.packageName, task.editId)
		if err != nil {
			return err
		}

		targetLanguages := make(map[string]struct{})
		for _, l := range targetListingsWithImages {
			targetLanguages[l.Listing.Language] = struct{}{}
		}

		for _, origListing := range originalListings {
			if _, ok := targetLanguages[origListing.Language]; !ok {
				err := task.api.Listings.Delete(ctx, task.accessToken, task.packageName, task.editId, origListing.Language)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, lwi := range targetListingsWithImages {
		err := task.upsert.Run(ctx, lwi.Listing, lwi.ImageTypeSources)
		if err != nil {
			return err
		}
	}

	return nil
}
