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

func (n *Node) GetStorageList() ([]Storage,error){


	target := "nodes/" + n.Node + "/storage"

	var storageList []Storage

	httpCode, err := n.px.APICall2("GET", target, nil, &storageList)
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

	httpCode, err := n.px.APICall2("GET", target, nil, &lxcList)
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

	apitarget,err := n.px.makeAPITarget(target)
	if err != nil {
		return nil, err
	}

	responseData, httpCode, err := n.px.APICall("DELETE", apitarget, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	var taskID TaskID

	jsonErr := n.px.dataUnmarshal(responseData,&taskID)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &taskID, nil
}



/*
vmid:999
hostname:test
password:111111
ostemplate:local:vztmpl/debian-9.0-standard_9.0-2_amd64.tar.gz
rootfs:storage:8
cores:2
memory:1024
swap:512
net0:bridge=vmbr0,name=eth0,ip=10.10.10.0/24,gw=10.10.10.1,tag=99
searchdomain:test.loc
nameserver:8.8.8.8
 */

func (n *Node) CreateLxc(lxcParams CreateLxcParams) (*TaskID, error) {
	err := lxcParams.Validate()
	if err != nil {
		return nil,err
	}

	target := "nodes/" + n.Node + "/lxc"

	apitarget,err := n.px.makeAPITarget(target)
	if err != nil {
		return nil,err
	}

	data := lxcParams.GetUrlDataValues()

	responseData, httpCode, err := n.px.APICall("CREATE", apitarget, data)
	if err != nil {
		return nil,err
	}
	if httpCode != 200 {
		return nil,errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	var taskID TaskID

	jsonErr := n.px.dataUnmarshal(responseData,&taskID)
	if jsonErr != nil {
		return nil,jsonErr
	}

	return &taskID, nil

}

