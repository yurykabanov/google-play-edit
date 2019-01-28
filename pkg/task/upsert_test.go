package task

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yurykabanov/google-play-edit/pkg/play"
)

var (
	token       = &play.AccessToken{AccessToken: "access_token", ExpiresIn: 3600, TokenType: "Bearer"}
	packageName = "com.example.project"
	editId      = "some_edit_id"

	sha1hashes = map[string]string{
		"aaa": "7e240de74fb1ed08fa08d38063f6a6a91462a815",
		"bbb": "5cb138284d431abd6a053a56625ec088bfb88912",
		"ccc": "f36b4825e5db2cf7dd2d2593b3f5c24c0311d8b2",
		"ddd": "9c969ddf454079e3d439973bbab63ea6233e4087",
		"eee": "637a81ed8e8217bb01c15c67c39b43b0ab4e20f1",
	}

	images = []io.ReadSeeker{
		bytes.NewReader([]byte("aaa")),
		bytes.NewReader([]byte("ccc")),
		bytes.NewReader([]byte("eee")),
	}

	ctx, _ = context.WithTimeout(context.Background(), 5 * time.Second)
)

func TestUpsert_Run(t *testing.T) {
	listingApi := &mockListingApi{}

	imagesApi := &mockImagesApi{}

	api := &play.Api{
		Edits:    &mockEditApi{},
		Listings: listingApi,
		Images:   imagesApi,
	}

	task := NewUpsert(api, token, packageName, editId)

	listing := &play.Listing{Language: "zz-ZZ"}
	ch := make(chan io.ReadSeeker, len(images))

	for _, image := range images {
		ch <- image
	}
	close(ch)

	// It should update edit's listing
	listingApi.On("Update", ctx, token, packageName, editId, listing).Return(&play.Listing{}, nil).
		Times(1)

	// It will eventually query list of images (to compare them with given images)
	imagesApi.On("List", ctx, token, packageName, editId, listing.Language, play.EditImagePhoneScreenshots).
		Return([]play.Image{
			{Id: "id_aaa", Sha1: sha1hashes["aaa"]},
			{Id: "id_bbb", Sha1: sha1hashes["bbb"]},
			{Id: "id_ccc", Sha1: sha1hashes["ccc"]},
			{Id: "id_ddd", Sha1: sha1hashes["ddd"]},
		}, nil).
		Times(1)

	// Given following lists:
	// - original: ["aaa", "bbb", "ccc", "ddd"       ]
	// - target:   ["aaa",        "ccc",        "eee"]
	//
	// Task should
	// - delete "bbb"
	// - delete "ccc"
	// - upload "eee"
	// and should not touch
	// - "aaa"
	// - "ccc"
	// as they are present in the same order in both original and target lists

	imagesApi.On("Delete", ctx, token, packageName, editId, listing.Language, play.EditImagePhoneScreenshots, "id_bbb").
		Return(nil).
		Times(1)
	imagesApi.On("Delete", ctx, token, packageName, editId, listing.Language, play.EditImagePhoneScreenshots, "id_ddd").
		Return(nil).
		Times(1)
	imagesApi.On("Upload", ctx, token, packageName, editId, listing.Language, play.EditImagePhoneScreenshots, images[2]).
		Return(&play.Image{}, nil).
		Times(1)

	err := task.Run(ctx, listing, ImageTypeSources{
		play.EditImagePhoneScreenshots: ch,
	})

	assert.Nil(t, err, "error should be nil")
}
