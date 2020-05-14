package dfaults

// String returns dfault string if value is an empty string
func String(dfault string, value string) (result string) {
	if len(value) > 0 {
		result = value
	} else {
		result = dfault
	}
	return
}
