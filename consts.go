package gojobs

import (
	"fmt"
	"go.dtapp.net/goip"
)

var ip string

func configIp() {
	ip = goip.GetOutsideIp()
}

const prefix = "cron_%s:"

func prefixSprintf() string {
	return fmt.Sprintf(prefix, ip)
}
