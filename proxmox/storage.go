package proxmox


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

