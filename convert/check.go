package convert

import (
	"fmt"
	"strings"
)

// check date field value
// like this:
//	    format: yyyy-MM-dd HH:mm:ss  1980-09-01 09:08:00
//      format: yyyy-MM-dd HH:mm:ss||[+-](\d+y\d+w\d+M\d+d\d+H\d+m\d+s)/[ywMdHms] 1980-07-01 09:08:00||+7w/d
//      format: now[+-](\d+y\d+w\d+M\d+d\d+H\d+m\d+s)/\d+[ywMdHms]  now+7y/d
func CheckDateValue(dateValue string) error {
	return nil
}

func CheckKeyword(keyword string) error {
	if len(strings.Split(keyword, " ")) == 1 {
		return nil
	} else {
		return fmt.Errorf("keyword:'%s' is invalid", keyword)
	}
}
