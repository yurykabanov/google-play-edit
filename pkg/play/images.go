package play

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	editImagesApiBaseUrl   = apiBaseUrl + "/edits/%s/listings/%s/%s"
	editImagesApiDelete    = editImagesApiBaseUrl + "/%s"
	editImagesApiDeleteAll = editImagesApiBaseUrl
	editImagesApiList      = editImagesApiBaseUrl
	editImagesApiUpload    = "https://www.googleapis.com/upload/androidpublisher/v3/applications/%s/edits/%s/listings/%s/%s"
)

type EditImageType string

const (
	EditImageFeatureGraphic       EditImageType = "featureGraphic"
	EditImageIcon                 EditImageType = "icon"
	EditImagePhoneScreenshots     EditImageType = "phoneScreenshots"
	EditImagePromoGraphic         EditImageType = "promoGraphic"
	EditImageSevenInchScreenshots EditImageType = "sevenInchScreenshots"
	EditImageTenInchScreenshots   EditImageType = "tenInchScreenshots"
	EditImageTvBanner             EditImageType = "tvBanner"
	EditImageTvScreenshots        EditImageType = "tvScreenshots"
	EditImageWearScreenshots      EditImageType = "wearScreenshots"
)

var acceptedMimeTypes = []string{
	"image/jpeg",
	"image/png",
}

type InvalidImage struct {
	MimeType string
}

func (err InvalidImage) Error() string {
	return fmt.Sprintf("unacceptable image mime type '%s'", err.MimeType)
}

type Image struct {
	Id   string `json:"id"`
	Url  string `json:"url"`
	Sha1 string `json:"sha1"`
}

type DeletedImages struct {
	Deleted []Image `json:"deleted"`
}

type ImageList struct {
	Images []Image `json:"images"`
}

func decodeImageResponse(decoder *json.Decoder) (*Image, error) {
	var image struct{ Image Image `json:"image"` }

	err := decoder.Decode(&image)
	if err != nil {
		return nil, err
	}
	return &image.Image, nil
}

func decodeDeletedImagesResponse(decoder *json.Decoder) (*DeletedImages, error) {
	var deletedImages DeletedImages

	err := decoder.Decode(&deletedImages)
	if err != nil {
		return nil, err
	}
	return &deletedImages, nil
}

func decodeImageListResponse(decoder *json.Decoder) (*ImageList, error) {
	var imageList ImageList

	err := decoder.Decode(&imageList)
	if err != nil {
		return nil, err
	}
	return &imageList, nil
}

type editImagesApi struct {
	client *http.Client
}

func (api *editImagesApi) Delete(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
	lang string,
	imageType EditImageType,
	imageId string,
) error {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(editImagesApiDelete, packageName, editId, lang, imageType, imageId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return decodeApiErrorResponse(dec)
	}

	return nil
}

func (api *editImagesApi) DeleteAll(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
	lang string,
	imageType EditImageType,
) ([]Image, error) {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(editImagesApiDeleteAll, packageName, editId, lang, imageType), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, decodeApiErrorResponse(dec)
	}

	list, err := decodeDeletedImagesResponse(dec)
	if err != nil {
		return nil, err
	}

	return list.Deleted, nil
}

func (api *editImagesApi) List(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
	lang string,
	imageType EditImageType,
) ([]Image, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(editImagesApiList, packageName, editId, lang, imageType), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, decodeApiErrorResponse(dec)
	}

	list, err := decodeImageListResponse(dec)
	if err != nil {
		return nil, err
	}

	return list.Images, nil
}

func (api *editImagesApi) Upload(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
	lang string,
	imageType EditImageType,
	imageReader io.ReadSeeker,
) (*Image, error) {
	mimeType, err := detectMimeType(imageReader)
	if err != nil {
		return nil, err
	}

	if !api.ensureValidMimeType(mimeType) {
		return nil, InvalidImage{MimeType: mimeType}
	}

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf(editImagesApiUpload, packageName, editId, lang, imageType), imageReader)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Content-Type", mimeType)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, decodeApiErrorResponse(dec)
	}

	return decodeImageResponse(dec)
}

func detectMimeType(r io.ReadSeeker) (string, error) {
	buf := make([]byte, 512)

	n, err := r.Read(buf)

	if err != nil && err != io.EOF {
		return "", err
	}

	_, err = r.Seek(0, 0)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf[:n]), nil
}

func (api *editImagesApi) ensureValidMimeType(mimeType string) bool {
	for _, accepted := range acceptedMimeTypes {
		if accepted == mimeType {
			return true
		}
	}

	return false
}
