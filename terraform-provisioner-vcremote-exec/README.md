```hcl-terraform

resource "null_resource" "foo" {

  provisioner "vcremote-exec" {
    vsphere_host = "stl0vcsa01.corp.ctf.local"
    vsphere_username = ""
    vsphere_password = ""
    vm_name = "stl0sqldvs02"
    guest_username = ""
    guest_password = ""
    //    commands = ["Get-Service","ipconfig","hostname"]
    // commands = ["hostname", "Install-WindowsFeature RSAT-AD-Powershell -IncludeAllSubFeature", "Install-WstsDomainController -ComputerName 'wng0dc02.corp.ctf.local'"]
    // commands = ["chef-client","ipconfig","hostname"]
    // commands = ["Write-Error 'This is error'"]
    // commands = ["Get-Service","Restart-Computer -Force","Get-Process"]
    script = <<EOT
                ipconfig
                hostname
                Get-Process
              EOT
  }

  triggers = {
    "always" = uuid()
 // "x" = timestamp()
  }
}

```
