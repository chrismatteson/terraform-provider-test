package main

import (
	"github.com/chrismatteson/terraform-provider-test/test"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: test.Provider})
}
