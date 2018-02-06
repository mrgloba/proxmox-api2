package proxmox

import (
	"fmt"
	"errors"
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

func (n *Node) GetStorages() ([]Storage,error){
	target,err := n.px.makeAPITarget("nodes/" + n.Node + "/storage")
	if err != nil {
		return nil, err
	}

	responseData, httpCode, err := n.px.APICall("GET", target, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}

	var storages []Storage

	jsonErr := n.px.dataUnmarshal(responseData, &storages)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return storages, nil
}

func (n *Node) GetLxcsis() ([]Lxc,error) {
	target,err := n.px.makeAPITarget("nodes/" + n.Node + "/lxc")
	if err != nil {
		return nil, err
	}

	responseData, httpCode, err := n.px.APICall("GET", target, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}

	var lxcs []Lxc

	jsonErr := n.px.dataUnmarshal(responseData, &lxcs)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return lxcs, nil
}