package main

import(
	"github.com/hashicorp/terraform/plugin"
	"terraform-provisioner-vcfile/vcfile"
)

func main()  {

	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: vcfile.Provisioner,
	})
}
