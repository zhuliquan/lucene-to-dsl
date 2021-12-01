package check

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// check date field value
// like this:
//	    format: yyyy-MM-dd HH:mm:ss  1980-09-01 09:08:00
//      format: yyyy-MM-dd HH:mm:ss||[+-](\d+y\d+w\d+M\d+d\d+H\d+m\d+s)/[ywMdHms] 1980-07-01 09:08:00||+7w/d
//      format: now[+-](\d+y\d+w\d+M\d+d\d+H\d+m\d+s)/\d+[ywMdHms]  now+7y/d
func CheckDateValue(dateValue string) error {
	time.ParseDuration()
	return nil
}

// check integer field value
func CheckIntegerValue(intValue string) error {
	if _, err := strconv.Atoi(intValue); err != nil {
		return fmt.Errorf("int_value: '%s' is invalid, err: %s", intValue, err.Error())
	} else {
		return nil
	}
}

// check float field value
func CheckFloatValue(floatValue string) error {
	if _, err := strconv.ParseFloat(floatValue, 64); err != nil {
		return fmt.Errorf("float_value: '%s' is invalid, err: %s", floatValue, err.Error())
	} else {
		return nil
	}
}

// es support ip ip_cidr query like this:
// {"term": {"ip_field": "172.168.1.0/24"}} or {"term": {"ip_field": "172.168.1.1"}}
// check ip field value
func CheckIpValue(ipValue string) error {
	if ip := net.ParseIP(ipValue); ip != nil {
		return nil
	} else if _, _, err := net.ParseCIDR(ipValue); err == nil {
		return nil
	} else {
		return fmt.Errorf("ipValue:'%s' is not valid ip / ip cidr", ipValue)
	}
}

func CheckKeyword(keyword string) error {
	if len(strings.Split(keyword, " ")) == 1 {
		return nil
	} else {
		return fmt.Errorf("keyword:'%s' is invalid", keyword)
	}
}

func CheckGeoValue(geoValue string) error {
	return nil
}
