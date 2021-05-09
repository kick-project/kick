package dfaults

// String returns default string if value is an empty string
func String(dfault string, value string) (result string) {
	if len(value) > 0 {
		result = value
	} else {
		result = dfault
	}
	return
}

// Interface returns a default value
func Interface(dfault interface{}, value interface{}) (result interface{}) {
	if value != nil {
		result = value
	} else {
		result = dfault
	}
	return
}
