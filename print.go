package gojobs

import (
	"context"
	"log"
)

func (c *Client) Println(ctx context.Context, isPrint bool, v ...any) {
	if isPrint {
		if c.slog.status {
			log.Println(v...)
		}
	}
}
