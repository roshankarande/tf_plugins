package vcfile

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/roshankarande/go-vsphere/vsphere"
	"github.com/sethvargo/go-retry"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// maxBufSize limits how much output we collect from a local
	// invocation. This is to prevent TF memory usage from growing
	// to an enormous amount due to a faulty process.
	maxBufSize = 8 * 1024
)

func Provisioner() terraform.ResourceProvisioner {
	return &schema.Provisioner{
		Schema: map[string]*schema.Schema{
			"vsphere_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_HOST","stl0vcsa01.corp.ctf.local"),
			},
			"vsphere_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_USERNAME","svc_vra@corp.ctf.local"),
			},
			"vsphere_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_PASSWORD","Theflame@5"),
			},
			"vm_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter" : &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Default: "STL_CTF",
			},
			"guest_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_GUEST_USERNAME","Administrator"),
			},
			"guest_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_GUEST_PASSWORD","5NbUJxG6"),
			},
			"source": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},
			"content": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
			},
			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "",
			},
			"timeout" : &schema.Schema{
			Type: schema.TypeInt,
			Optional: true,
			DefaultFunc: schema.EnvDefaultFunc("VC_FILE_TIMEOUT", 300),
			},
			"interval" : &schema.Schema{
			Type: schema.TypeInt,
			Optional: true,
			DefaultFunc: schema.EnvDefaultFunc("VC_FILE_INTERVAL", 20),
			},
		},

		ApplyFunc: applyFn,
		ValidateFunc:validateFn,
	}
}

func applyFn(ctx context.Context) error {
	d := ctx.Value(schema.ProvConfigDataKey).(*schema.ResourceData)
	o := ctx.Value(schema.ProvOutputKey).(terraform.UIOutput)

	// Execute the command with env
	vSphereHost := fmt.Sprintf("https://%s/sdk",d.Get("vsphere_host").(string))
	vSphereUsername := d.Get("vsphere_username").(string)
	vSpherePassword := d.Get("vsphere_password").(string)
	vmName := d.Get("vm_name").(string)
	guestUsername := d.Get("guest_username").(string)
	guestPassword := d.Get("guest_password").(string)
	dst := d.Get("destination").(string)
	datacenter := d.Get("datacenter").(string)


	options := map[string]interface{}{
		"timeout" : time.Duration(d.Get("timeout").(int)),
		"delay" : time.Duration(d.Get("interval").(int)),
		"output" : o,
		"datacenter" : datacenter,
	}

	c, err := vsphere.NewClient(ctx, vSphereHost, vSphereUsername, vSpherePassword)

	if err != nil {
		return err
	}

	vimClient := c.Client

	b, err := retry.NewConstant(10*time.Second)
	if err != nil {
		return err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(3, b), func(ctx context.Context) error {
		if !vimClient.IsVC() && vimClient.Valid() {
			o.Output("testing client connectivity to VC...")
			return retry.RetryableError(err)
		}
		o.Output("client connection to VC successful!")
		return nil
	})

	if err != nil {
		return err
	}

	o.Output("Starting to upload file!")

	if _, ok := d.GetOk("source"); ok {
		src := d.Get("source").(string)

		f, err := os.Open(src)
		fstat, err := f.Stat()

		if err != nil {
			return err
		}

		if fstat.IsDir(){
			tarfile := fmt.Sprintf("%s.tar",src)
			err = CreateTar(src, tarfile)

			if err != nil {
				return err
			}

			f, err = os.Open(tarfile)
			if err != nil {
				return err
			}
		}
		err = vsphere.Upload(ctx, c, vmName, guestUsername, guestPassword, f, filepath.Ext(src), dst, fstat.IsDir(), options)

	} else if _, ok := d.GetOk("content"); ok {
		content := d.Get("content").(string)
		err = vsphere.Upload(ctx, c, vmName, guestUsername, guestPassword, strings.NewReader(content),"", dst,false, options)
	}

	if err != nil {
		return err
	}

	o.Output("File uploaded!")

	return nil
}

func validateFn(c *terraform.ResourceConfig) (ws []string, es []error) {
	if !c.IsSet("source") && !c.IsSet("content") {
		es = append(es, fmt.Errorf("must provide one of 'source' or 'content'"))
	}
	return ws, es
}

