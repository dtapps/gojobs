package gojobs

import (
	"context"
)

func (c *Client) Println(ctx context.Context, isPrint bool, v ...any) {
	if isPrint {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Info("", v...)
		}
	}
}
