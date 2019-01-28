package loader

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/yurykabanov/google-play-edit/pkg/play"
)

var (
	ErrInsufficientColumns = errors.New(fmt.Sprintf("insufficient columns, csv must contain at least four following columns: Language, Title, ShortDescription, FullDescription"))
)

func LoadListingsFromFile(path string) ([]play.Listing, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var listings []play.Listing

	ext := filepath.Ext(path)

	if ext == ".csv" {
		rdr := csv.NewReader(f)
		for {
			row, err := rdr.Read()
			if err != nil {
				if err == io.EOF {
					return listings, nil
				}
				return nil, err
			}

			listing := play.Listing{}

			if len(row) < 4 {
				return nil, ErrInsufficientColumns
			}

			listing.Language = row[0]
			listing.Title = row[1]
			listing.ShortDescription = row[2]
			listing.FullDescription = row[3]

			if len(row) == 5 {
				listing.Video = row[4]
			}

			listings = append(listings, listing)
		}
	}

	switch ext {
	case ".yaml":
	case ".yml":
		dec := yaml.NewDecoder(f)
		err = dec.Decode(&listings)
	case ".json":
		dec := json.NewDecoder(f)
		err = dec.Decode(&listings)
	default:
		return nil, errors.New(fmt.Sprintf("unknown format: %s", ext))
	}
	if err != nil {
		return nil, err
	}

	return listings, nil
}

func FindImagesForLang(path, lang string) ([]string, error) {
	join := filepath.Join(path, lang, "*")
	fmt.Println(join)
	files, err := filepath.Glob(join)

	return files, err
}
