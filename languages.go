package linguist

// Convenience function that returns the color associated
// with the language, in HTML Hex notation (e.g. "#123ABC")
// from the languages.yaml provided by github.com/github/linguist
//
// returns empty string if there is no associated color for the language
func GetColor(language string) string {
	if c, ok := colors[language]; ok {
		return c
	}
	return ""
}
