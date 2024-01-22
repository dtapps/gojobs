package gojobs

func newClient() (*Client, error) {
	client, err := NewClient(&ClientConfig{})
	return client, err
}
