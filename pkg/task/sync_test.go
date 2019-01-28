package task

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/yurykabanov/google-play-edit/pkg/play"
)

type upsertMock struct {
	mock.Mock
}

func (mock *upsertMock) Run(ctx context.Context, listing *play.Listing, imageTypeSources ImageTypeSources) error {
	args := mock.Called(ctx, listing, imageTypeSources)

	return args.Error(0)
}

func TestSync_Run(t *testing.T) {
	listingApi := &mockListingApi{}

	imagesApi := &mockImagesApi{}

	api := &play.Api{
		Edits:    &mockEditApi{},
		Listings: listingApi,
		Images:   imagesApi,
	}

	upsert := &upsertMock{}

	task := NewSync(api, token, packageName, editId)
	task.upsert = upsert

	originalListings := []play.Listing{
		{Language: "aa-AA"},
		{Language: "cc-CC"},
	}

	targetListings := []ListingWithImages{
		{Listing: &play.Listing{Language: "aa-AA"}, ImageTypeSources: ImageTypeSources{}},
		{Listing: &play.Listing{Language: "bb-BB"}, ImageTypeSources: ImageTypeSources{}},
	}

	// It should list existing listings
	listingApi.On("List", ctx, token, packageName, editId).
		Return(originalListings, nil).
		Times(1)

	// It should delete language that is not on the list
	listingApi.On("Delete", ctx, token, packageName, editId, "cc-CC").
		Return(nil).
		Times(1)

	upsert.On("Run", ctx, targetListings[0].Listing, targetListings[0].ImageTypeSources).
		Return(nil).
		Times(1)
	upsert.On("Run", ctx, targetListings[1].Listing, targetListings[1].ImageTypeSources).
		Return(nil).
		Times(1)

	err := task.Run(ctx, targetListings, true)

	assert.Nil(t, err, "error should be nil")
}
