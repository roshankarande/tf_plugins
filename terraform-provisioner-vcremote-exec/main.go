package main

import(
	"github.com/hashicorp/terraform/plugin"
	remoteExec "terraform-provisioner-vcremote-exec/vcremote-exec"
)

func main()  {

	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: remoteExec.Provisioner,
	})
}
