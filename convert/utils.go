package convert

import (
	"strings"

	"github.com/araddon/dateparse"
)

func ToUpper(x string) (string, error) {
	return strings.ToUpper(x), nil
}

func ToLower(x string) (string, error) {
	return strings.ToLower(x), nil
}

func ParseDate(x string) (string, error) {
	if t, err := dateparse.ParseAny(x); err != nil {
		return "", err
	} else {
		return t.Format("2006-01-02 15:04:05"), nil
	}

}
