package i18nUtils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	LocalesMap Locales
)

func init() {
	var err error
	LocalesMap, err = readLocales("utils/i18nUtils/locales")
	if err != nil {
		panic(err)
	}

}

func readLocales(dirname string) (Locales, error) {
	locales := make(Locales)

	files, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	referenceKeys := make(map[string]bool)
	var referenceLocale string

	for _, file := range files {

		if filepath.Ext(file.Name()) == ".json" {
			locale := strings.Split(file.Name(), ".")[0]

			content, err := os.ReadFile(filepath.Join(dirname, file.Name()))
			if err != nil {
				return nil, err
			}

			var translations Locale
			err = json.Unmarshal(content, &translations)
			if err != nil {
				return nil, err
			}

			locales[locale] = translations

			// check for missing keys
			if referenceLocale == "" {
				referenceLocale = locale
				for key := range translations {
					referenceKeys[key] = true
				}
			} else {
				for key := range translations {
					if !referenceKeys[key] {
						return nil, fmt.Errorf("locale %s has extra key: %s", locale, key)
					}
				}
				for refKey := range referenceKeys {
					if _, exists := translations[refKey]; !exists {
						return nil, fmt.Errorf("locale %s is missing key: %s", referenceLocale, refKey)
					}
				}
			}
		}
	}

	return locales, nil
}
