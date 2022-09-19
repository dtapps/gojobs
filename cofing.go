package gojobs

import (
	"context"
	"go.dtapp.net/goip"
	"os"
	"runtime"
)

func (c *Client) setConfig(ctx context.Context) {
	c.config.sdkVersion = Version
	c.config.systemOs = runtime.GOOS
	c.config.systemArch = runtime.GOARCH
	c.config.goVersion = runtime.Version()
	c.config.systemInsideIp = goip.GetInsideIp(ctx)
	hostname, _ := os.Hostname()
	c.config.systemHostName = hostname
}
