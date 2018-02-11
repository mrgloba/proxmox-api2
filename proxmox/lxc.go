package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"strconv"
	"reflect"
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

type LxcBase struct {
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
	Status string	`json:"status"`
	Swap int64		`json:"swap"`
	Template string	`json:"template"`
	Type string		`json:"type"`
	Uptime int64	`json:"uptime"`
	VmId int64		`json:"vmid,string"`
}

type LxcStatus struct {
	HA map[string]interface{} `json:"ha"`
	LxcBase
}

type Lxc struct {
	Pid int			`json:"pid,string"`
	LxcBase
	BasicObject
}

type MountPoint struct {
	Index int
	Volume string
	Path string
	BaseStorageItem
}

type StartupConfig struct {
	Order int
	UpDelay int
	DownDelay int
}

type NetworkConfig struct {
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

type BaseLxcConfig struct {
	Arch string				`json:"arch"`
	CMode string			`json:"cmode"`
	Console bool			`json:"console"`
	Cores int				`json:"cores"`
	CpuLimit int			`json:"cpulimit"`
	CpuUnits int			`json:"cpuunits"`
	Description string		`json:"description"`
	Force bool				`json:"force"`
	Hostname string			`json:"hostname"`
	Lock string				`json:"lock"`
	Memory int				`json:"memory"`
	NameServer string		`json:"nameserver"`
	OnBoot bool				`json:"onboot"`
	OSTemplate string		`json:"ostemplate"`
	OSType string			`json:"ostype"`
	Password string			`json:"password"`
	Pool string				`json:"pool"`
	Protection bool			`json:"protection"`
	Restore bool			`json:"restore"`
	SearchDomain string		`json:"search_domain"`
	Storage string			`json:"storage"`
	Swap int				`json:"swap"`
	Template bool			`json:"template"`
	Tty int					`json:"tty"`
	Unprivileged bool		`json:"unprivileged"`
	VmId int64				`json:"vmid"`
	RootFS string 			`json:"rootfs"`
}

type LxcConfigReceiver struct {
	mp0 string `json:"mp0"`
	mp1 string `json:"mp1"`
	mp2 string `json:"mp2"`
	mp3 string `json:"mp3"`
	mp4 string `json:"mp4"`
	mp5 string `json:"mp5"`
	mp6 string `json:"mp6"`
	mp7 string `json:"mp7"`
	mp8 string `json:"mp8"`
	mp9 string `json:"mp9"`

	net0 string `json:"net0"`
	net1 string `json:"net1"`
	net2 string `json:"net2"`
	net3 string `json:"net3"`
	net4 string `json:"net4"`
	net5 string `json:"net5"`
	net6 string `json:"net6"`
	net7 string `json:"net7"`
	net8 string `json:"net8"`
	net9 string `json:"net9"`

	startup string`json:"startup"`

	lxc []interface{} `json:"lxc"`

	BaseLxcConfig
}

type LxcConfig struct {

	MountPoints []MountPoint
	Networks []NetworkConfig
	Startup StartupConfig

	BaseLxcConfig
}

func (lcr *LxcConfigReceiver) Parse() (*LxcConfig){

	lxcConfig := LxcConfig{ BaseLxcConfig: lcr.BaseLxcConfig}

	if len(lcr.startup) > 0 {

		lxcConfig.Startup = StartupConfig{}
		keyval := strings.Split(lcr.startup,",")

		for _, kvpair := range keyval {
			params := strings.Split(kvpair,"=")

			if len(params) == 2 {
				switch params[0] {
				case "order":
					lxcConfig.Startup.Order, _ = strconv.Atoi(params[1])
				case "up":
					lxcConfig.Startup.UpDelay, _ = strconv.Atoi(params[1])
				case "down":
					lxcConfig.Startup.DownDelay, _ = strconv.Atoi(params[1])
				}
			}
		}
	}

	rlcrv := reflect.ValueOf(*lcr)
	rlcrt := reflect.TypeOf(*lcr)
//todo: end
	for i:=0; i<rlcrv.NumField(); i++ {
		fmt.Printf("%v\n",rlcrt.Field(i))
	}

	return nil
}

func (clp *LxcConfig) Validate() error {
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

	if len(clp.RootFS) == 0 {
		return errors.New("RootFS could not be zero")
	}

	if clp.VmId < 0 {
		return errors.New("VmId has wrong value")
	}

	return nil
}

func (clp *LxcConfig) GetUrlDataValues() url.Values {

	reg,_:=regexp.Compile("([a-z]+)=")
	reg.ReplaceAllString("","\"$1\"")

	return url.Values{}
}