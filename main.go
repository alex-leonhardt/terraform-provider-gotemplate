package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/rapid7/terraform-provider-gotemplate/gotemplate"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: gotemplate.Provider})
}
