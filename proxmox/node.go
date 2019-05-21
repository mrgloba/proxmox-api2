package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)


const (
	BACKUP_MODE_SNAPSHOT BackupMode = 0
	BACKUP_MODE_SUSPEND BackupMode = 1
	BACKUP_MODE_STOP BackupMode = 2
	BACKUP_COMP_LZO BackupComp = 0
	BACKUP_COMP_GZIP BackupComp = 1
)


type BackupComp int
type BackupMode int

type Node struct {
	Cpu float64			`json:"cpu"`
	Disk int64			`json:"disk"`
	Id string			`json:"id"`
	Level string		`json:"level"`
	MaxCpu	int			`json:"maxcpu"`
	MaxDisk int64		`json:"maxdisk"`
	MaxMem int64		`json:"maxmem"`
	Mem int64			`json:"mem"`
	Node string			`json:"node"`
	Type string			`json:"type"`
	Uptime int64		`json:"uptime"`

	BasicObject
}

type USBDevice struct {
	Busnum int			`json:"busnum"`
	Class int			`json:"class"`
	Devnum int			`json:"devnum"`
	Level int			`json:"level"`
	Manufacturer string `json:"manufacturer"`
	Port int			`json:"port"`
	Prodid string		`json:"prodid"`
	Product string		`json:"product"`
	Speed string		`json:"speed"`
	Usbpath string		`json:"usbpath"`
	Vendid string		`json:"vendid"`
}

type LVMVolumeGroup struct {
	Free int64 `json:"free"`
	Size int64 `json:"size"`
	VG string  `json:"vg"`
}

func (bm BackupMode) String() string {
	strvalue := [...]string{
		"snapshot",
		"suspend",
		"stop",
	}

	if bm < BACKUP_MODE_SNAPSHOT || bm > BACKUP_MODE_STOP {
		return "unknown"
	}

	return strvalue[bm]
}

func (bc BackupComp) String() string {
	strvalue := [...]string{
		"lzo",
		"gzip",
	}

	if bc < BACKUP_COMP_LZO || bc > BACKUP_COMP_GZIP {
		return "unknown"
	}

	return strvalue[bc]
}

func (n *Node) fillParent(v interface{}, parent interface{}) {
	n.parent.(*Proxmox).fillParent(v, n)
}

func (n *Node) GetStorageList() ([]Storage,error){


	target := "nodes/" + n.Node + "/storage"

	var storageList []Storage

	httpCode, err := n.parent.(*Proxmox).APICall2("GET", target, nil, &storageList,n)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	return storageList, nil
}

func (n *Node) GetLxcList() ([]Lxc,error) {
	target := "nodes/" + n.Node + "/lxc"

	var lxcList []Lxc

	httpCode, err := n.parent.(*Proxmox).APICall2("GET", target, nil, &lxcList, n)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	return lxcList, nil
}

func (n *Node) GetLxc(vmid int64) (*Lxc,error) {
	lxcList, err := n.GetLxcList()
	if err != nil {
		return nil, err
	}

	for _,v := range lxcList {
		if v.VmId == vmid {
			lxc := v
			return &lxc, nil

		}
	}

	return nil, errors.New("Lxc container VMID: " + strconv.Itoa(int(vmid)) + " not found.")
}

func (n *Node) RemoveLxc(vmid int64) (*TaskID, error) {

	target := "nodes/" + n.Node + "/lxc/" + strconv.Itoa(int(vmid))

	apitarget,err := n.GetProxmox().MakeAPITarget(target)
	if err != nil {
		return nil, err
	}

	responseData, httpCode, err := n.parent.(*Proxmox).APICall("DELETE", apitarget, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	var taskID TaskID

	jsonErr := n.GetProxmox().DataUnmarshal(responseData,&taskID,nil)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &taskID, nil
}

func (n *Node) CreateLxc(lxcParams LxcConfig) (*TaskID, error) {
	err := lxcParams.Validate()
	if err != nil {
		return nil,err
	}

	target := "nodes/" + n.Node + "/lxc"

	apitarget,err := n.GetProxmox().MakeAPITarget(target)
	if err != nil {
		return nil,err
	}

	data := lxcParams.GetUrlDataValues()

	responseData, httpCode, err := n.parent.(*Proxmox).APICall("POST", apitarget, data)
	if err != nil {
		return nil,err
	}
	if httpCode != 200 {
		return nil,errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	var taskID TaskID

	jsonErr := n.GetProxmox().DataUnmarshal(responseData,&taskID,nil)
	if jsonErr != nil {
		return nil,jsonErr
	}

	return &taskID, nil

}

func (n *Node) VZDump(vmid int64, storage Storage, mode BackupMode, comp BackupComp, remove bool) (*TaskID, error) {
	target := "nodes/" + n.Node + "/vzdump"

	apitarget,err := n.GetProxmox().MakeAPITarget(target)
	if err != nil {
		return nil, err
	}


	var data url.Values

	data = make(url.Values)
	data.Add("vmid", strconv.Itoa(int(vmid)))
	data.Add("storage", storage.Storage)
	data.Add("mode", mode.String())
	data.Add("compress", comp.String())

	if remove {
		data.Add("remove", "1")
	} else {
		data.Add("remove", "0")
	}


	responseData, httpCode, err := n.parent.(*Proxmox).APICall("POST", apitarget, data)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	var taskID TaskID

	jsonErr := n.GetProxmox().DataUnmarshal(responseData,&taskID,nil)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &taskID, nil
}

func (n *Node) RestoreLxc(vmid int64, storageContentItem StorageContentItem, storage string, force bool,  newLxcParams LxcConfig) (*TaskID, error) {
	newLxcParams.VmId = vmid
	newLxcParams.OSTemplate = storageContentItem.Volid
	newLxcParams.Force = force
	newLxcParams.Storage = storage
	newLxcParams.Restore = true
	return n.CreateLxc(newLxcParams)
}

func (n *Node) GetTasks() ([]Task, error) {
	target := "nodes/" + n.Node + "/tasks"

	var tasks []Task

	httpCode, err := n.parent.(*Proxmox).APICall2("GET", target, nil, &tasks, n)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	return tasks, nil
}

func (n *Node) ScanUSB() ([]USBDevice, error) {
	target := "nodes/" + n.Node + "/scan/usb"
	var devices []USBDevice

	httpCode, err := n.parent.(*Proxmox).APICall2("GET", target, nil, &devices, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}

	return devices,nil
}

func (n *Node) ScanLVM() ([]LVMVolumeGroup, error) {
	target := "nodes/" + n.Node + "/scan/lvm"
	var groups []LVMVolumeGroup

	httpCode, err := n.parent.(*Proxmox).APICall2("GET", target, nil, &groups, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}

	return groups,nil
}