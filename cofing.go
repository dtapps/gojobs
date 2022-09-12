package gojobs

import (
	"context"
	"go.dtapp.net/goarray"
	"go.dtapp.net/goip"
	"runtime"
)

func (c *Client) setConfig(ctx context.Context) {
	c.config.sdkVersion = Version
	c.config.systemOs = runtime.GOOS
	c.config.systemArch = runtime.GOARCH
	c.config.systemCpuQuantity = runtime.GOMAXPROCS(0)
	c.config.goVersion = runtime.Version()
	c.config.systemMacAddrS = goarray.TurnString(goip.GetMacAddr(ctx))
	c.config.systemInsideIp = goip.GetInsideIp(ctx)
}
