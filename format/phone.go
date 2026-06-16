package format

import (
	"regexp"
	"strings"
)

func NormalizePhone(countryCode, phone string) string {
	re := regexp.MustCompile(`[^0-9]`)
	phone = re.ReplaceAllString(phone, "")
	countryCode = re.ReplaceAllString(countryCode, "")

	if countryCode == "" || phone == "" {
		return ""
	}

	phone = strings.TrimPrefix(phone, countryCode)
	phone = strings.TrimPrefix(phone, "0")

	return countryCode + phone
}
