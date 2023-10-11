package main

import(
	"github.com/hashicorp/terraform/plugin"
	"terraform-provisioner-vcchef/chef"
)

func main()  {

	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: chef.Provisioner,
	})
}
