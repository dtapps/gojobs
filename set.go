package gojobs

// SetDebug 设置调试模式
func (c *Client) SetDebug() {
	c.slog.status = true
}
