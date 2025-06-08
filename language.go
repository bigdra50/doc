package main

import (
	"fmt"
	"os"
)

// supportedLanguages maps language codes to language names
var supportedLanguages = map[string]string{
	"ja": "Japanese",
	"en": "English",
	"ko": "Korean",
	"zh": "Chinese",
	"ru": "Russian",
	"es": "Spanish",
	"fr": "French",
	"de": "German",
	"it": "Italian",
	"pt": "Portuguese",
	"nl": "Dutch",
	"sv": "Swedish",
	"no": "Norwegian",
	"da": "Danish",
	"fi": "Finnish",
	"pl": "Polish",
	"cs": "Czech",
	"hu": "Hungarian",
	"ro": "Romanian",
	"bg": "Bulgarian",
	"hr": "Croatian",
	"sk": "Slovak",
	"sl": "Slovenian",
	"et": "Estonian",
	"lv": "Latvian",
	"lt": "Lithuanian",
	"mt": "Maltese",
	"el": "Greek",
	"tr": "Turkish",
	"ar": "Arabic",
	"he": "Hebrew",
	"hi": "Hindi",
	"th": "Thai",
	"vi": "Vietnamese",
	"id": "Indonesian",
	"ms": "Malay",
	"tl": "Filipino",
	"sw": "Swahili",
	"am": "Amharic",
}

// validateLanguageCode validates a language code against the default supported languages
func validateLanguageCode(code string) error {
	if _, exists := supportedLanguages[code]; !exists {
		return fmt.Errorf("unsupported language code: %s", code)
	}
	return nil
}

// validateLanguageCodeWithMap validates a language code against a provided map
func validateLanguageCodeWithMap(code string, supportedLangs map[string]string) error {
	if _, exists := supportedLangs[code]; !exists {
		return fmt.Errorf("unsupported language code: %s", code)
	}
	return nil
}

// showSupportedLanguages displays all supported language codes
func showSupportedLanguages() {
	fmt.Fprintf(os.Stderr, "Supported language codes:\n")

	// Sort for consistent output
	codes := make([]string, 0, len(supportedLanguages))
	for code := range supportedLanguages {
		codes = append(codes, code)
	}

	// Simple sort
	for i := 0; i < len(codes); i++ {
		for j := i + 1; j < len(codes); j++ {
			if codes[i] > codes[j] {
				codes[i], codes[j] = codes[j], codes[i]
			}
		}
	}

	for _, code := range codes {
		fmt.Fprintf(os.Stderr, "  %s - %s\n", code, supportedLanguages[code])
	}
}

// getSimilarLanguageCodes finds language codes that start with the input string
func getSimilarLanguageCodes(input string) []string {
	var similar []string
	for code := range supportedLanguages {
		if len(code) >= len(input) && code[:len(input)] == input {
			similar = append(similar, code)
		}
	}

	// Simple sort
	for i := 0; i < len(similar); i++ {
		for j := i + 1; j < len(similar); j++ {
			if similar[i] > similar[j] {
				similar[i], similar[j] = similar[j], similar[i]
			}
		}
	}

	return similar
}

// getSimilarLanguageCodesWithMap finds similar language codes in a provided map
func getSimilarLanguageCodesWithMap(input string, supportedLangs map[string]string) []string {
	var similar []string
	for code := range supportedLangs {
		if len(code) >= len(input) && code[:len(input)] == input {
			similar = append(similar, code)
		}
	}

	// Simple sort
	for i := 0; i < len(similar); i++ {
		for j := i + 1; j < len(similar); j++ {
			if similar[i] > similar[j] {
				similar[i], similar[j] = similar[j], similar[i]
			}
		}
	}

	return similar
}
