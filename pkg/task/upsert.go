package task

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"

	"github.com/yurykabanov/google-play-edit/pkg/play"
)

type Upsert interface {
	Run(ctx context.Context, listing *play.Listing, imageTypeSources ImageTypeSources) error
}

type upsert struct {
	api         *play.Api
	accessToken *play.AccessToken
	packageName string
	editId      string
}

func NewUpsert(
	api *play.Api,
	accessToken *play.AccessToken,
	packageName string,
	editId string,
) *upsert {
	return &upsert{
		api:         api,
		accessToken: accessToken,
		packageName: packageName,
		editId:      editId,
	}
}

type ImageTypeSources map[play.EditImageType]<-chan io.ReadSeeker

func (task *upsert) Run(ctx context.Context, listing *play.Listing, imageTypeSources ImageTypeSources) error {
	_, err := task.api.Listings.Update(ctx, task.accessToken, task.packageName, task.editId, listing)
	if err != nil {
		return err
	}

	for imageType, imageChan := range imageTypeSources {
		err = task.handleImageType(ctx, listing, imageType, imageChan)
		if err != nil {
			return err
		}
	}

	return nil
}

func (task *upsert) handleImageType(ctx context.Context, listing *play.Listing, imageType play.EditImageType, imagesChan <-chan io.ReadSeeker) error {
	images, err := task.api.Images.List(
		ctx,
		task.accessToken, task.packageName, task.editId,
		listing.Language, imageType,
	)
	if err != nil {
		return err
	}

	imagesCount := len(images)
	pos := 0

	for imageReader := range imagesChan {
		sha1sum, err := task.sha1sum(imageReader)
		if err != nil {
			return err
		}

		if pos < imagesCount {
			n, err := task.deleteImagesUntilSha1(ctx, images[pos:], sha1sum, listing, imageType)
			if err != nil {
				return err
			}
			pos += n

			continue
		}

		_, err = task.api.Images.Upload(
			ctx,
			task.accessToken, task.packageName, task.editId,
			listing.Language, imageType, imageReader,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (task *upsert) sha1sum(r io.ReadSeeker) (string, error) {
	hash := sha1.New()

	br := bufio.NewReader(r)
	_, err := br.WriteTo(hash)
	if err != nil {
		return "", err
	}

	_, err = r.Seek(0, 0)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (task *upsert) deleteImagesUntilSha1(
	ctx context.Context,
	images []play.Image,
	sha1sum string,
	listing *play.Listing,
	imageType play.EditImageType,
) (int, error) {
	pos := 0
	imagesCount := len(images)

	for {
		if pos >= imagesCount {
			break
		}

		image := &images[pos]

		if sha1sum == image.Sha1 {
			pos++
			break
		} else {
			err := task.api.Images.Delete(
				ctx,
				task.accessToken, task.packageName, task.editId,
				listing.Language, imageType, image.Id,
			)
			if err != nil {
				return 0, err
			}

			pos++
		}
	}

	return pos, nil
}
