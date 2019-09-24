package client

import "testing"

func TestNewEmptyClient(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Error(err)
	}
	context := &Context{Config: config}
	context.NewEmptyClient()
}
