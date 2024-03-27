package gojobs

import (
	"context"
	"log/slog"
)

func (c *Client) Println(ctx context.Context, isPrint bool, v ...any) {
	if isPrint {
		if c.slog.status {
			slog.InfoContext(ctx, "", v...)
		}
	}
}
