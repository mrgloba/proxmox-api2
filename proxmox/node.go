package proxmox

import (
	"errors"
	"fmt"
	"strconv"
)


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

	apitarget,err := n.parent.(*Proxmox).makeAPITarget(target)
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

	jsonErr := n.parent.(*Proxmox).dataUnmarshal(responseData,&taskID,nil)
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

	apitarget,err := n.parent.(*Proxmox).makeAPITarget(target)
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

	jsonErr := n.parent.(*Proxmox).dataUnmarshal(responseData,&taskID,nil)
	if jsonErr != nil {
		return nil,jsonErr
	}

	return &taskID, nil

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

func (n *Node) ScanUsb() ([]USBDevice, error) {
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