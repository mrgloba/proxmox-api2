package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)


type LxcStatus struct {
	HA map[string]interface{} `json:"ha"`
	LxcBase
}

type Lxc struct {
	Pid int			`json:"pid,string"`
	LxcBase
	BasicObject
}

type LxcConfig struct {

	MountPoints []MountPoint
	Networks []NetworkConfig
	Startup StartupConfig

	BaseLxcConfig
}


func (lxc *Lxc) Start() (*TaskID, error){
	target
}

func (lxc *Lxc) Stop() (*TaskID, error){}

func (lxc *Lxc) Shutdown() (*TaskID, error){}

func (lxc *Lxc) GetStatus() (*LxcStatus, error) {}

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
	var data url.Values

	data = make(url.Values)

	if len(clp.Arch) > 0 { data.Add("arch",clp.Arch) }
	if len(clp.CMode) > 0 { data.Add("cmode",clp.CMode) }
	if len(clp.Arch) > 0 { data.Add("arch",clp.Arch) }
	if clp.Console { data.Add("console", "1")}
	data.Add("cores",strconv.Itoa(clp.Cores))
	if clp.CpuLimit > 0 { data.Add("cpulimit",strconv.Itoa(clp.CpuLimit)) }
	if clp.CpuUnits > 0 { data.Add("cpuunits",strconv.Itoa(clp.CpuUnits)) }
	if len(clp.Description) > 0 { data.Add("description", clp.Description) }
	if clp.Force { data.Add("force", "1") }
	if len(clp.Hostname)>0 { data.Add("hostname",clp.Hostname) }
	if len(clp.Lock)>0 { data.Add("lock",clp.Lock) }
	if clp.Memory >0 {data.Add("memory",strconv.Itoa(clp.Memory))}
	if len(clp.NameServer)>0 { data.Add("nameserver",clp.NameServer) }
	if clp.OnBoot { data.Add("onboot", "1") }
	if len(clp.OSTemplate)>0 { data.Add("ostemplate", clp.OSTemplate) }
	if len(clp.OSType) > 0 { data.Add("ostype", clp.OSType) }
	if len(clp.Password) >0 { data.Add("password", clp.Password) }
	if len(clp.Pool) > 0 { data.Add("pool",clp.Pool) }
	if clp.Protection { data.Add("protection", "1") }
	if clp.Restore { data.Add("restore", "1") }
	if len(clp.SearchDomain)>0 { data.Add("searchdomain",clp.SearchDomain) }
	if len(clp.Storage)>0 { data.Add("storage",clp.Storage) }
	if clp.Swap >0 {data.Add("swap",strconv.Itoa(clp.Swap))}
	if clp.Template { data.Add("template", "1") }
	if clp.Tty >0 {data.Add("tty",strconv.Itoa(clp.Tty))}
	if clp.Unprivileged { data.Add("unprivileged", "1") }
	data.Add("vmid",strconv.Itoa(int(clp.VmId)))
	data.Add("rootfs",clp.RootFS)

	if len(clp.Startup.String()) > 0 { data.Add("startup",clp.Startup.String()) }
	if len(clp.Networks) > 0 {
		for _,n := range clp.Networks {
			if len(n.String()) >0 {
				data.Add("net" + strconv.Itoa(n.Index), n.String())
			}
		}
	}

	if len(clp.MountPoints) > 0 {
		for _,m := range clp.MountPoints {
			if len(m.String()) > 0 {
				data.Add("mp" + strconv.Itoa(m.Index), m.String())
			}
		}
	}
	return data
}

