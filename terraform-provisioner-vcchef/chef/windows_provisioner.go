package chef

import (
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

// $url = "http://omnitruck.chef.io/%s/chef/download?p=windows&pv=$machine_os&m=$machine_arch&v=%s"
const installScript = `
$winver = [System.Environment]::OSVersion.Version | %% {"{0}.{1}" -f $_.Major,$_.Minor}

switch ($winver)
{
  "6.0" {$machine_os = "2008"}
  "6.1" {$machine_os = "2008r2"}
  "6.2" {$machine_os = "2012"}
  "6.3" {$machine_os = "2012"}
  default {$machine_os = "2008r2"}
}

if ([System.IntPtr]::Size -eq 4) {$machine_arch = "i686"} else {$machine_arch = "x86_64"}

$url = "%s"
$dest = [System.IO.Path]::GetTempFileName()
$dest = [System.IO.Path]::ChangeExtension($dest, ".msi")
$downloader = New-Object System.Net.WebClient

$http_proxy = '%s'
if ($http_proxy -ne '') {
	$no_proxy = '%s'
  if ($no_proxy -eq ''){
    $no_proxy = "127.0.0.1"
  }

  $proxy = New-Object System.Net.WebProxy($http_proxy, $true, ,$no_proxy.Split(','))
  $downloader.proxy = $proxy
}

Write-Host 'Downloading Chef Client...'
$downloader.DownloadFile($url, $dest)

Write-Host 'Installing Chef Client...'
Start-Process -FilePath msiexec -ArgumentList /qn, /i, $dest -Wait
`

func (p *provisioner) windowsInstallChefClient(o terraform.UIOutput, comm *vcCommunicator) error {
	script := path.Join(path.Dir("C:\\Windows\\Temp"), "ChefClient.ps1")
	//content := fmt.Sprintf(installScript, p.Channel, p.Version, p.HTTPProxy, strings.Join(p.NOProxy, ","))
	content := fmt.Sprintf(installScript,p.InstallerUrl, p.HTTPProxy, strings.Join(p.NOProxy, ","))

	// Copy the script to the new instance
	//o.Output(content)
	if err := comm.toolBoxClient.UploadFile(comm.ctx,script,strings.NewReader(content),"",false) ; err != nil {
		return fmt.Errorf("Uploading client.rb failed: %v", err)
	}


	// Execute the script to install Chef Client
	installCmd := fmt.Sprintf("powershell -NoProfile -ExecutionPolicy Bypass -File %s", script)
	//o.Output(installCmd)
	return p.runCommand(o, comm, installCmd)
}

func (p *provisioner) windowsCreateConfigFiles(o terraform.UIOutput, comm *vcCommunicator) error {
	// Make sure the config directory exists
	cmd := fmt.Sprintf("cmd /c if not exist %q mkdir %q", windowsConfDir, windowsConfDir)
	if err := p.runCommand(o, comm, cmd); err != nil {
		return err
	}

	if len(p.OhaiHints) > 0 {
		// Make sure the hits directory exists
		hintsDir := path.Join(windowsConfDir, "ohai/hints")
		cmd := fmt.Sprintf("cmd /c if not exist %q mkdir %q", hintsDir, hintsDir)
		if err := p.runCommand(o, comm, cmd); err != nil {
			return err
		}

		if err := p.deployOhaiHints(o, comm, hintsDir); err != nil {
			return err
		}
	}

	return p.deployConfigFiles(o, comm, windowsConfDir)
}
