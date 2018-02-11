package proxmox


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

