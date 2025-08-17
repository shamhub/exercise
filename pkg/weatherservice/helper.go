package weatherservice

import (
	"fmt"
	"regexp"
	"time"

	"github.com/shamhub/exercise/pkg/errorlib"
)

func findGridCoordinates(endpoint string) (string, error) {

	re := regexp.MustCompile(`(\d+),(\d+)`)

	// Find the matches in the string
	matches := re.FindStringSubmatch(endpoint)

	if len(matches) < 3 {
		return "", errorlib.NewResponseError(500, "could not find valid coordinates from forecast grid data endpoint")
	}

	return matches[0], nil
}

func getCurrentDate() string {

	currentTime := time.Now()

	// To get only the date components
	year := currentTime.Year()
	month := currentTime.Month()
	day := currentTime.Day()

	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
