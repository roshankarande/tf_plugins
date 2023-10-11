package vcremote_exec

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/roshankarande/go-vsphere/vsphere"
	"github.com/sethvargo/go-retry"
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
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_HOST", "stl0vcsa01.corp.ctf.local"),
			},
			"vsphere_username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_USERNAME", "svc_vra@corp.ctf.local"),
			},
			"vsphere_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_PASSWORD", "Theflame@5"),
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
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_GUEST_USERNAME", "Administrator"),
			},
			"guest_password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VC_VSPHERE_GUEST_PASSWORD", "5NbUJxG6"),
			},
			"commands": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"script"},
			},
			"script" : &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"commands"},
			},
			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VC_RX_TIMEOUT", 300),
			},
			"interval": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VC_RX_INTERVAL", 15),
			},
		},

		ApplyFunc: applyFn,
	}
}

func applyFn(ctx context.Context) error {
	d := ctx.Value(schema.ProvConfigDataKey).(*schema.ResourceData)
	o := ctx.Value(schema.ProvOutputKey).(terraform.UIOutput)

	// Execute the command with env
	vSphereHost := fmt.Sprintf("https://%s/sdk", d.Get("vsphere_host").(string))
	vSphereUsername := d.Get("vsphere_username").(string)
	vSpherePassword := d.Get("vsphere_password").(string)
	vmName := d.Get("vm_name").(string)
	guestUsername := d.Get("guest_username").(string)
	guestPassword := d.Get("guest_password").(string)
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

	b, err := retry.NewConstant(10 * time.Second)
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

	if _, ok := d.GetOk("commands"); ok {
		commands := SliceInterfacesToStrings(d.Get("commands").([]interface{}))
		vsphere.InvokeCommands(ctx,c,vmName,guestUsername,guestPassword,commands,options)
	}

	if _, ok := d.GetOk("script"); ok {
		script := d.Get("script").(string)
		vsphere.InvokeScript(ctx,c,vmName,guestUsername,guestPassword,script,options)
	}

	return nil
}

func SliceInterfacesToStrings(s []interface{}) []string {
	var d []string
	for _, v := range s {
		if o, ok := v.(string); ok {
			d = append(d, o)
		}
	}
	return d
}
