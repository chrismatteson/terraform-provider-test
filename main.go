package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/chrismatteson/terraform-provider-test/test"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: test.Provider})
}
