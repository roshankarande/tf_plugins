package chef

import (
	"context"
	"github.com/roshankarande/go-vsphere/vsphere"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"io"
	"net/url"
)

type vcCommunicator struct {
	ctx           context.Context
	toolBoxClient *vsphere.ToolBoxClient
}

func NewVcCommunicator(ctx context.Context,opsmgr *guest.OperationsManager,guestUser, guestPassword string) (*vcCommunicator,error) {
	tboxClient, err := vsphere.NewToolBoxClient(ctx,opsmgr,guestUser,guestPassword,types.VirtualMachineGuestOsFamilyWindowsGuest)
	if err != nil {
		return nil, err
	}

	return &vcCommunicator{
		ctx:           ctx,
		toolBoxClient: tboxClient,
	},nil

}

//func NewVcCommunicator(ctx context.Context,toolBoxClient *govmomi.Client, vmName, guestUser, guestPassword string) (*vcCommunicator,error) {
//
//	vm, err := find.NewFinder(toolBoxClient.Client).VirtualMachine(ctx, vmName)
//
//	if err != nil {
//		return nil, fmt.Errorf("[vm] %s does not exist in [vc]", vmName)
//	}
//
//	pmgr, err := opsmgr.ProcessManager(ctx)
//
//	if err != nil {
//		return nil,err
//	}
//
//	fmgr, err := opsmgr.FileManager(ctx)
//
//	if err != nil {
//		return nil,err
//	}
//
//	return &vcCommunicator{
//		ctx: ctx,
//		toolBoxClient:  &toolbox.Client{
//			ProcessManager: pmgr,
//			FileManager:    fmgr,
//			Authentication: auth,
//			GuestFamily:    types.VirtualMachineGuestOsFamilyWindowsGuest,
//		},
//	},nil
//
//}

func NewClient(ctx context.Context, vSphereHost, vSphereUsername, vSpherePassword string) (*govmomi.Client, error) {

	u, err := soap.ParseURL(vSphereHost)
	if err != nil {
		return nil, err
	}

	u.User = url.UserPassword(vSphereUsername, vSpherePassword)

	return govmomi.NewClient(ctx, u, true)
}

func (comm *vcCommunicator) Upload(dst string, f io.Reader) error {
	return comm.toolBoxClient.UploadFile(comm.ctx,dst,f,"",false)
}
