package play

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	editListingsApiBaseUrl   = apiBaseUrl + "/edits/%s/listings"
	editListingsApiDelete    = editListingsApiBaseUrl + "/%s"
	editListingsApiDeleteAll = editListingsApiBaseUrl
	editListingsApiGet       = editListingsApiBaseUrl + "/%s"
	editListingsApiList      = editListingsApiBaseUrl
	editListingsApiPatch     = editListingsApiBaseUrl + "/%s"
	editListingsApiUpdate    = editListingsApiBaseUrl + "/%s"
)

type editListingsApi struct {
	client *http.Client
}

type Listing struct {
	Language         string `json:"language"`
	Title            string `json:"title"`
	FullDescription  string `json:"fullDescription"`
	ShortDescription string `json:"shortDescription"`
	Video            string `json:"video"`
}

type ListingList struct {
	Kind     string    `json:"kind"`
	Listings []Listing `json:"listings"`
}

func decodeListingResponse(decoder *json.Decoder) (*Listing, error) {
	var listing Listing

	err := decoder.Decode(&listing)
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

func decodeListingListResponse(decoder *json.Decoder) (*ListingList, error) {
	var listingList ListingList

	err := decoder.Decode(&listingList)
	if err != nil {
		return nil, err
	}
	return &listingList, nil
}

func (api *editListingsApi) Delete(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
	lang string,
) error {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(editListingsApiDelete, packageName, editId, lang), nil)
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

func (api *editListingsApi) DeleteAll(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
) error {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(editListingsApiDeleteAll, packageName, editId), nil)
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

func (api *editListingsApi) Get(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
	lang string,
) (*Listing, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(editListingsApiGet, packageName, editId, lang), nil)
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

	return decodeListingResponse(dec)
}

func (api *editListingsApi) List(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
) ([]Listing, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(editListingsApiList, packageName, editId), nil)
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

	list, err := decodeListingListResponse(dec)
	if err != nil {
		return nil, err
	}

	return list.Listings, nil
}

func (api *editListingsApi) Update(
	ctx context.Context,
	token *AccessToken,
	packageName string,
	editId string,
	listing *Listing,
) (*Listing, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	err := enc.Encode(listing)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(editListingsApiUpdate, packageName, editId, listing.Language), &buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Content-Type", "application/json")

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

	return decodeListingResponse(dec)
}
