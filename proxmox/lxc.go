package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
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
	Mp0 string `json:"mp0"`
	Mp1 string `json:"mp1"`
	Mp2 string `json:"mp2"`
	Mp3 string `json:"mp3"`
	Mp4 string `json:"mp4"`
	Mp5 string `json:"mp5"`
	Mp6 string `json:"mp6"`
	Mp7 string `json:"mp7"`
	Mp8 string `json:"mp8"`
	Mp9 string `json:"mp9"`

	Net0 string `json:"net0"`
	Net1 string `json:"net1"`
	Net2 string `json:"net2"`
	Net3 string `json:"net3"`
	Net4 string `json:"net4"`
	Net5 string `json:"net5"`
	Net6 string `json:"net6"`
	Net7 string `json:"net7"`
	Net8 string `json:"net8"`
	Net9 string `json:"net9"`

	Startup string`json:"startup"`

	Lxc []interface{} `json:"lxc"`

	BaseLxcConfig
}

type LxcConfig struct {

	MountPoints []MountPoint
	Networks []NetworkConfig
	Startup StartupConfig

	BaseLxcConfig
}


func (sc *StartupConfig) SetFromString(str string) {
	if len(str) > 0 {
		keypairs := parseKeyPairs(str)
		for _, kp := range keypairs {
			if len(kp) == 2 {
				switch kp[0] {
				case "order":
					sc.Order, _ = strconv.Atoi(kp[1])
				case "up":
					sc.UpDelay, _ = strconv.Atoi(kp[1])
				case "down":
					sc.DownDelay, _ = strconv.Atoi(kp[1])
				}
			}
		}
	}
}

func (sc *StartupConfig) String() (string) {
	var res []string

	if sc.Order>0 { res=append(res,fmt.Sprintf("order=%d",sc.Order)) }
	if sc.UpDelay>0 { res=append(res,fmt.Sprintf("up=%d",sc.UpDelay)) }
	if sc.DownDelay >0 { res=append(res,fmt.Sprintf("down=%d",sc.DownDelay)) }

	var str string
	for i,v := range res {
		str += v
		if len(res)-1 != i {
			str += ","
		}
	}

	return str
}

func (nc *NetworkConfig) SetFromString(idx int, str string) {
	nc.Index = idx
	if len(str)>0 {
		keypairs := parseKeyPairs(str)

		for _, kp := range keypairs {
			if len(kp) == 2 {
				switch kp[0] {
				case "name":
					nc.Name = kp[1]
				case "bridge":
					nc.Bridge = kp[1]
				case "firewall":
					f,_ := strconv.Atoi(kp[1])
					nc.Firewall = f == 1
				case "gw":
					nc.Gateway = kp[1]
				case "gw6":
					nc.GatewayV6 = kp[1]
				case "hwaddr":
					nc.HWAddr = kp[1]
				case "ip":
					nc.IPAddress = kp[1]
				case "ip6":
					nc.IPAddresV6 = kp[1]
				case "mtu":
					nc.MTU,_ = strconv.Atoi(kp[1])
				case "rate":
					nc.Rate,_ =strconv.Atoi(kp[1])
				case "tag":
					nc.Tag,_ = strconv.Atoi(kp[1])
				case "trunks":
					nc.Trunks = kp[1]
				case "type":
					nc.Type = kp[1]
				}
			}
		}
	}
}

func (nc *NetworkConfig) String() (string) {
	var res []string
	if len(nc.Name)>0 { res = append(res,fmt.Sprintf("name=%s",nc.Name)) }
	if len(nc.Bridge)>0 { res = append(res,fmt.Sprintf("bridge=%s",nc.Bridge)) }
	if nc.Firewall { res = append(res,fmt.Sprintf("firewall=%d",1)) }
	if len(nc.Gateway)>0 { res = append(res,fmt.Sprintf("gw=%s",nc.Gateway)) }
	if len(nc.GatewayV6)>0 { res = append(res,fmt.Sprintf("gw6=%s",nc.GatewayV6)) }
	if len(nc.HWAddr)>0 { res = append(res,fmt.Sprintf("hwaddr=%s",nc.HWAddr)) }
	if len(nc.IPAddress)>0 { res = append(res,fmt.Sprintf("ip=%s",nc.IPAddress)) }
	if len(nc.IPAddresV6)>0 { res = append(res,fmt.Sprintf("ip6=%s",nc.IPAddresV6)) }
	if nc.MTU > 0 { res = append(res,fmt.Sprintf("mtu=%d",nc.MTU)) }
	if nc.Rate > 0 { res = append(res,fmt.Sprintf("rate=%d",nc.Rate)) }
	if nc.Tag >0 { res = append(res,fmt.Sprintf("tag=%d",nc.Tag)) }
	if len(nc.Trunks) > 0 { res = append(res,fmt.Sprintf("trunks=%s",nc.Trunks)) }
	if len(nc.Type) >0 { res = append(res,fmt.Sprintf("type=%s",nc.Type)) }


	var str string
	for i,v := range res {
		str += v
		if len(res)-1 != i {
			str += ","
		}
	}

	return str
}

func (mp *MountPoint) SetFromString(idx int, str string) {
	mp.Index = idx
	if len(str)>0 {
		keypairs := parseKeyPairs(str)

		for _,kp := range keypairs {
			if len(kp) ==2 {
				switch kp[0] {

				case "volume":
					mp.Volume = kp[1]
				case "mp":
					mp.Path = kp[1]
				case "acl":
					a,_ := strconv.Atoi(kp[1])
					mp.ACL = a ==1
				case "quota":
					q,_ := strconv.Atoi(kp[1])
					mp.Quota = q==1
				case "backup":
					b,_ := strconv.Atoi(kp[1])
					mp.Backup = b ==1
				case "ro":
					ro,_ := strconv.Atoi(kp[1])
					mp.ReadOnly = ro ==1
				case "size":
					r := regexp.MustCompile("^(\\d+)G$")
					if r.MatchString(kp[1]) {
						sm := r.FindStringSubmatch(kp[1])
						if len(sm) == 2 {
							s,_:=strconv.Atoi(sm[1])
							mp.Size = s
						}
					}

				}
			}
			if len(kp) == 1 {
				mp.Volume = kp[0]
			}
		}
	}
}

func (mp *MountPoint) String() (string) {
	var res []string

	if len(mp.Volume) == 0 { return ""}

	res = append(res,mp.Volume)
	if len(mp.Path) > 0 { res = append(res,fmt.Sprintf("mp=%s",mp.Path)) }
	if mp.ACL { res = append(res,"acl=1") }
	if mp.Quota { res = append(res,"quota=1") }
	if mp.Backup { res = append(res,"backup=1") }
	if mp.ReadOnly { res = append(res,"ro=1") }
	if mp.Size > 0 { res = append(res,fmt.Sprintf("size=%dG",mp.Size)) }

	var str string
	for i,v := range res {
		str += v
		if len(res)-1 != i {
			str += ","
		}
	}

	return str
}

func (lcr *LxcConfigReceiver) Parse() (*LxcConfig){

	lxcConfig := LxcConfig{ BaseLxcConfig: lcr.BaseLxcConfig}

	lxcConfig.Startup.SetFromString(lcr.Startup)

	rlcrv := reflect.ValueOf(*lcr)
	rlcrt := reflect.TypeOf(*lcr)

	r := regexp.MustCompile("^(Mp|Net)(\\d+)$")

	for i:=0; i<rlcrv.NumField(); i++ {
		if r.MatchString(rlcrt.Field(i).Name) {
			sm := r.FindStringSubmatch(rlcrt.Field(i).Name)
			switch sm[1] {
			case "Mp":
				val := fmt.Sprintf("%v",rlcrv.Field(i))
				if len(val)>0 {
					mp := MountPoint{}
					idx, _ := strconv.Atoi(sm[2])
					mp.SetFromString(idx,val)
					lxcConfig.MountPoints = append(lxcConfig.MountPoints,mp)
				}
			case "Net":
				val := fmt.Sprintf("%v",rlcrv.Field(i))
				if len(val)>0 {
					nc := NetworkConfig{}
					idx, _ := strconv.Atoi(sm[2])
					nc.SetFromString(idx,val)
					lxcConfig.Networks = append(lxcConfig.Networks,nc)
				}
			}
		}
	}

	return &lxcConfig
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