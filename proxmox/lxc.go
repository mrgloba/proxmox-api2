package proxmox

import (
	"net/url"
	"errors"
	"fmt"
)

const (
	LXC_CMODE_CONSOLE = "console"
	LXC_CMODE_SHELL = "shell"
	LXC_CMODE_TTY = "tty"

	LXC_ARCH_AMD64 = "amd64"
	LXC_ARCH_I386 = "i386"

	LXC_LOCK_BACKUP = "backup"
	LXC_LOCK_MIGRATE = "migrate"
	LXC_LOCK_ROLBACK = "rollback"
	LXC_LOCK_SNAPSHOT = "snapshot"

)

type Lxc struct {
	Cpu float64 	`json:"cpu"`
	Cpus int 		`json:"cpus,string"`
	Disk int64		`json:"disk"`
	DiskRead int64	`json:"diskread,string"`
	DiskWrite int64	`json:"diskwrite,string"`
	MaxDisk	int64	`json:"maxdisk"`
	MaxMem	int64	`json:"maxmem"`
	MaxSwap int64	`json:"maxswap"`
	Mem int64		`json:"mem"`
	Name string		`json:"name"`
	NetIn int64		`json:"netin"`
	NetOut int64	`json:"netout"`
	Pid int			`json:"pid,string"`
	Status string	`json:"status"`
	Swap int64		`json:"swap"`
	Template string	`json:"template"`
	Type string		`json:"type"`
	Uptime int64	`json:"uptime"`
	VmId int64		`json:"vmid,string"`

	BasicObject
}

type BaseStorageParams struct {
	ACL bool
	Backup bool
	Quota bool
	ReadOnly bool
	Shared bool
	Size int
}

type MountPointParams struct {
	Index int
	Volume string
	Path string
	BaseStorageParams
}

type RootFSParams struct {
	Volume string
	BaseStorageParams
}

type StartupParams struct {
	Order int
	UpDelay int
	DownDelay int
}

type NetworkParams {
	Index int
	Name string
	Bridge string
	Firewall bool
	Gateway string
	GatewayV6 string
	HWAddr string
	IPAddress string
	IPAddresV6 string
	MTU int
	Rate int
	Tag int
	Trunks string
	Type string
}

type CreateLxcParams struct {
	Arch string
	CMode string
	Console bool
	Cores int
	CpuLimit int
	CpuUnits int
	Description string
	Force bool
	Hostname string
	Lock string
	Memory int
	NameServer string
	OnBoot bool
	OSTemplate string
	OSType string
	Password string
	Pool string
	Protection bool
	Restore bool
	SearchDomain string
	Storage string
	Swap int
	Template bool
	Tty int
	Unprivileged bool
	VmId int64

	MountPoints []MountPointParams
	RootFS RootFSParams
	Networks []NetworkParams
	Startup StartupParams
}


func (clp *CreateLxcParams) Validate() error {
	if clp.Arch != "" && clp.Arch != LXC_ARCH_AMD64 && clp.Arch != LXC_ARCH_I386 {
		return errors.New(fmt.Sprintf("Arch has wrong value. Posible values is: %s or %s or empty", LXC_ARCH_I386, LXC_ARCH_AMD64))
	}

	if clp.CMode != "" && clp.CMode != LXC_CMODE_CONSOLE && clp.CMode != LXC_CMODE_SHELL && clp.CMode != LXC_CMODE_TTY {
		return errors.New(fmt.Sprintf("CMode has wrong value. Posible values is: [%s|%s|%s|empty]",LXC_CMODE_CONSOLE,LXC_CMODE_SHELL,LXC_CMODE_TTY))
	}

	if clp.Cores <1 || clp.Cores >128 {
		return errors.New("Cores has wrong value. It shuld be 1-128")
	}

	if clp.CpuLimit <0 || clp.CpuLimit >128 {
		return errors.New("CpuLimit has wrong value. It shuld be 0-128")
	}

	if clp.CpuUnits <0 || clp.CpuUnits > 500000 {
		return errors.New("CpuLimit has wrong value. It shuld be 0-500000")
	}

	switch clp.Lock {
	case LXC_LOCK_BACKUP:
	case LXC_LOCK_MIGRATE:
	case LXC_LOCK_ROLBACK:
	case LXC_LOCK_SNAPSHOT:
	case "":
	default:
		return errors.New(fmt.Sprintf("CMode has wrong value. Posible values is: [%s|%s|%s|%s|empty]",LXC_LOCK_SNAPSHOT,LXC_LOCK_ROLBACK,LXC_LOCK_MIGRATE,LXC_LOCK_BACKUP))
	}

	if clp.Memory < 16 {
		return errors.New("Memory has wrong value. It shuld be 16-N")
	}

	if len(clp.OSTemplate) == 0 {
		return errors.New("OSTemplate could not be zero")
	}

	if len(clp.Password) < 6 {
		return errors.New("Password length couldn ot be less than 6")
	}

	if clp.RootFS == (RootFSParams{}) {
		return errors.New("RootFS could not be zero")
	}

	if clp.VmId < 0 {
		return errors.New("VmId has wrong value")
	}

	return nil
}

func (clp *CreateLxcParams) GetUrlDataValues() url.Values {
	return url.Values{}
}