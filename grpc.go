package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/goip"
)

var ip string

func configIp() {
	ip = goip.GetOutsideIp(context.Background())
}

const prefix = "cron:"

const prefixIp = "cron_%s:"

func prefixSprintf(str string) string {
	return fmt.Sprintf(prefixIp, str)
}
