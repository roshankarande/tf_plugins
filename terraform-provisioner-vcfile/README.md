
```hcl-terraform
resource "null_resource" "foo" {

  provisioner "vcfile" {
    vsphere_host = "stl0vcsa01.corp.ctf.local"
    vsphere_username = ""
    vsphere_password = ""
    vm_name = "stl0sqldvs02"
    guest_username = ""
    guest_password = ""
    //  content = "Hello, World....! This is working and working just fine!"
    //  destination = "C:/out.txt"
    //  destination = "C:\\something\\output.txt"   // ......... this doesn't work.. because dir something is not present
    source = "C:\\Users\\e062721\\junk\\output.txt"
    destination = "C:\\abc.txt"
    //  source = "C:\\Users\\e062721\\junk"
    //  destination = "C:\\junk\\junkie"
  }

  triggers = {
    "always" = uuid()
    // "x" = timestamp()
  }

}




```
