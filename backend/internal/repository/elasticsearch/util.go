package elasticsearch

import "github.com/olivere/elastic/v7"

func GetValueWithHighlight(hit elastic.SearchHitHighlight, key string, defaultValue string) string {
	if value, ok := hit[key]; ok {
		return value[0]
	}
	return defaultValue
}
