package pretty

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"

	"github.com/yurykabanov/google-play-edit/pkg/play"
)

func Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	lines := strings.Split(message, "\n")

	maxLen := 0
	for _, line := range lines {
		if l := len(line); l > maxLen {
			maxLen = l
		}
	}

	fmt.Println(aurora.Black(strings.Repeat(" ", maxLen+4)).BgRed())
	for _, line := range lines {
		if line != "" {
			fmt.Println(aurora.Black(fmt.Sprintf("  %-"+strconv.Itoa(maxLen)+"s  ", line)).BgRed())
		}
	}
	fmt.Println(aurora.Black(strings.Repeat(" ", maxLen+4)).BgRed())
}

func PrintEdit(edit *play.Edit) {
	expiry, _ := strconv.ParseInt(edit.ExpiryTimeSeconds, 10, 64)

	fmt.Println(aurora.Red("Edit").Bold())
	fmt.Printf("%s: %s\n", aurora.Green("ID").Bold(), edit.Id)
	fmt.Printf("%s: %s\n", aurora.Green("Expires at").Bold(), time.Unix(expiry, 0).Format(time.RFC3339))
}

func PrintListing(listing *play.Listing) {
	fmt.Println(aurora.Red("Listing").Bold())
	fmt.Printf("%s: %s\n", aurora.Green("Language"), aurora.Red(listing.Language))
	fmt.Printf("%s: %s\n", aurora.Green("Title"), listing.Title)
	fmt.Printf("%s: %s\n", aurora.Green("Short Description"), listing.ShortDescription)
	fmt.Printf("%s:\n%s\n", aurora.Green("Full Description"), aurora.Gray(listing.FullDescription))
}

func PrintImage(image *play.Image) {
	fmt.Printf("  - %s: %-20s [%-40s] %s\n", aurora.Green("Image"), image.Id, image.Sha1, aurora.Gray(image.Url))
}
