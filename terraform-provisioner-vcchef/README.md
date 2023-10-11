```hcl-terraform

resource "null_resource" "vm" {

  provisioner "vcchef" {

    environment     = "_default"
    client_options  = ["chef_license 'accept'"]
    # run_list        = ["mc_wsts_os_standard::_uninstall"]
    run_list        = ["mc_wsts_os_tweaks"]
    node_name       = "stl0dc19"
    server_url      = "https://172.18.132.229/organizations/mastercard/"
    recreate_client = true           ### make this false if you are using inside VM
    user_name       = "chefadmin"
    user_key        = file("./chefadmin.pem")
    skip_install    = true
    # If you have a self signed cert on your chef server change this to :verify_none
    ssl_verify_mode  = ":verify_none"

    vsphere_host     = "stl0vcsa01.corp.ctf.local"
    vsphere_username = ""
    vsphere_password = ""
    vm_name          = "stl0dc19"
    guest_username   = ""
    guest_password   = ""
  }

  triggers = {
    build_number = uuid()
  }


}


```
