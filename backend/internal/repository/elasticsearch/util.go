package elasticsearch

func GetValueWithHighlight(hit map[string]interface{}, key string, defaultValue string) string {
	if value, ok := hit[key]; ok {
		values := value.([]interface{})
		return values[0].(string)
	}
	return defaultValue
}
