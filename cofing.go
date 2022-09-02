package gojobs

import (
	"context"
	"go.dtapp.net/goarray"
	"go.dtapp.net/goip"
	"runtime"
)

func (c *Client) setConfig(ctx context.Context) {
	c.config.runVersion = Version
	c.config.os = runtime.GOOS
	c.config.arch = runtime.GOARCH
	c.config.maxProCs = runtime.GOMAXPROCS(0)
	c.config.version = runtime.Version()
	c.config.macAddrS = goarray.TurnString(goip.GetMacAddr(ctx))
	c.config.insideIp = goip.GetInsideIp(ctx)
}
