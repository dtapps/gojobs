package gojobs

// GetCurrentIp 获取当前ip
func (c *Client) GetCurrentIp() string {
	return c.config.systemOutsideIp
}

// GetSubscribeAddress 获取订阅地址
func (c *Client) GetSubscribeAddress() string {
	return c.cache.cornKeyPrefix + "_" + c.cache.cornKeyCustom
}
