package utils

import (
	"fmt"
	"regexp"
)

type watchAnime struct {
	URL     string `json:"url"`
	Quality string `json:"quality"`
}

func ScrapeSources(htmlData string) []watchAnime {
	var sources []watchAnime
	sourcesDataMatch := regexp.MustCompile(`sources:\s*\[\s*(.*?)\s*\],?\s*poster:`).FindStringSubmatch(htmlData)

	if len(sourcesDataMatch) > 1 {
		sourcesData := sourcesDataMatch[1]
		sourceRegex := regexp.MustCompile(`{\s*src:\s*'(.*?)'.*?type:\s*'(.*?)'.*?size:\s*(\d+),?\s*}`)
		matches := sourceRegex.FindAllStringSubmatch(sourcesData, -1)

		for _, match := range matches {
			url := match[1]
			size := match[3]
			quality := fmt.Sprintf("%sp", size)
			sources = append(sources, watchAnime{URL: url, Quality: quality})
		}
	}
	return sources
}
