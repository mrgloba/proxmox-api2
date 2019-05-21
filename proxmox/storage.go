package proxmox

import (
	"errors"
	"fmt"
)

type BaseStorageItem struct {
	ACL bool
	Backup bool
	Quota bool
	ReadOnly bool
	Shared bool
	Size int
}

type Storage struct {
	Content string		`json:"content"`
	Digest string		`json:"digest"`
	Maxfiles int64		`json:"maxfiles"`
	Path string			`json:"path"`
	Shared int			`json:"shared"`
	Storage string		`json:"storage"`
	Type string			`json:"type"`

	BasicObject
}

type StorageContentItem struct {
	Content string		`json:"content"`
	Format string 		`json:"format"`
	Parent string 		`json:"parent"`
	Size int	  		`json:"size"`
	Used int	  		`json:"used"`
	VMid string	  		`json:"vmid"`
	Volid string  		`json:"volid"`
}

func (s *Storage) GetContents() ([]StorageContentItem, error) {


	var storageContent []StorageContentItem

	px := s.GetProxmox()
	n := s.GetNode()



	if n == nil {
		return nil, errors.New("Could not be called from cluster object")
	}

	target := "nodes/" + n.Node + "/storage/" + s.Storage + "/content"

	httpCode, err := px.APICall2("GET", target, nil, &storageContent, n)

	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}

	return storageContent, nil
}
