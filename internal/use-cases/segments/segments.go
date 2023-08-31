package usecases_segments

import (
	"strings"
)

func FormatSegmnetName(segmentName string) string {
	formated := strings.ToUpper(segmentName)
	formated = strings.Replace(formated, " ", "_", -1)
	return formated
}
